package outboundemailerfunction

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	analyticsconfig "github.com/prixfixeco/backend/internal/analytics/config"
	"github.com/prixfixeco/backend/internal/config"
	"github.com/prixfixeco/backend/internal/database"
	"github.com/prixfixeco/backend/internal/database/postgres"
	"github.com/prixfixeco/backend/internal/email"
	emailconfig "github.com/prixfixeco/backend/internal/email/config"
	"github.com/prixfixeco/backend/internal/observability"
	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/logging/zerolog"
	"github.com/prixfixeco/backend/internal/observability/tracing"

	_ "github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	_ "go.uber.org/automaxprocs"
)

func init() {
	// Register a CloudEvent function with the Functions Framework
	functions.CloudEvent("SendEmail", SendEmail)
}

// MessagePublishedData contains the full Pub/Sub message
// See the documentation for more details:
// https://cloud.google.com/eventarc/docs/cloudevents#pubsub
type MessagePublishedData struct {
	Message PubSubMessage
}

// PubSubMessage is the payload of a Pub/Sub event.
// See the documentation for more details:
// https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// SendEmail handles a data change.
func SendEmail(ctx context.Context, e event.Event) error {
	var msg MessagePublishedData
	if err := e.DataAs(&msg); err != nil {
		return fmt.Errorf("event.DataAs: %v", err)
	}

	logger := zerolog.NewZerologLogger(logging.DebugLevel)

	envCfg := email.GetConfigForEnvironment(os.Getenv("PF_ENVIRONMENT"))
	if envCfg == nil {
		return observability.PrepareAndLogError(email.ErrMissingEnvCfg, logger, nil, "getting environment config")
	}

	cfg, err := config.GetDataChangesWorkerConfigFromGoogleCloudSecretManager(ctx)
	if err != nil {
		return fmt.Errorf("error getting config: %w", err)
	}

	tracerProvider, err := cfg.Observability.Tracing.ProvideTracerProvider(ctx, logger)
	if err != nil {
		logger.Error(err, "initializing tracer")
	}
	otel.SetTracerProvider(tracerProvider)

	ctx, span := tracing.NewTracer(tracerProvider.Tracer("outbound_emailer_job")).StartSpan(ctx)
	defer span.End()

	analyticsEventReporter, err := analyticsconfig.ProvideEventReporter(&cfg.Analytics, logger, tracerProvider)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "error setting up customer data collector")
	}

	defer analyticsEventReporter.Close()

	// manual db timeout until I find out what's wrong
	dbConnectionContext, cancel := context.WithTimeout(ctx, 15*time.Second)
	dataManager, err := postgres.ProvideDatabaseClient(dbConnectionContext, logger, &cfg.Database, tracerProvider)
	if err != nil {
		cancel()
		return observability.PrepareAndLogError(err, logger, span, "establishing database connection")
	}

	cancel()
	defer dataManager.Close()

	if !dataManager.IsReady(ctx, 50) {
		return observability.PrepareAndLogError(database.ErrDatabaseNotReady, logger, span, "pinging database")
	}

	emailer, err := emailconfig.ProvideEmailer(&cfg.Email, logger, tracerProvider, otelhttp.DefaultClient)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "configuring outbound emailer")
	}

	var emailDeliveryRequest email.DeliveryRequest
	if err = json.Unmarshal(msg.Message.Data, &emailDeliveryRequest); err != nil {
		logger = logger.WithValue("raw_data", msg.Message.Data)
		return observability.PrepareAndLogError(err, logger, span, "unmarshalling data change message")
	}

	logger = logger.WithValue("template", emailDeliveryRequest.Template)

	user, err := dataManager.GetUser(ctx, emailDeliveryRequest.UserID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "getting user")
	}

	if user.EmailAddressVerifiedAt == nil {
		logger.Info("user email address not verified, skipping email delivery")
		return nil
	}

	var (
		mail      *email.OutboundEmailMessage
		emailType string
	)

	switch emailDeliveryRequest.Template {
	case email.TemplateTypeInvite:
		if emailDeliveryRequest.Invitation == nil {
			return observability.PrepareAndLogError(err, logger, span, "missing household invitation")
		}

		mail, err = email.BuildInviteMemberEmail(emailDeliveryRequest.Invitation, envCfg)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "building email message")
		}
		emailType = "invite"

		break
	case email.TemplateTypeUsernameReminder:
		mail, err = email.BuildUsernameReminderEmail(user.EmailAddress, user.Username, envCfg)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "building username reminder email")
		}
		emailType = "username reminder"

		break
	case email.TemplateTypePasswordReset:
		if emailDeliveryRequest.PasswordResetToken == nil {
			return observability.PrepareAndLogError(err, logger, span, "missing password reset token")
		}

		mail, err = email.BuildGeneratedPasswordResetTokenEmail(user.EmailAddress, emailDeliveryRequest.PasswordResetToken, envCfg)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "building password reset token email")
		}
		emailType = "password reset token"

		break
	case email.TemplateTypePasswordResetTokenRedeemed:
		mail, err = email.BuildPasswordResetTokenRedeemedEmail(user.EmailAddress, envCfg)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "building password reset token redemption email")
		}
		emailType = "password reset token redemption"

		break
	case email.TemplateTypeMealPlanCreated:
		mail, err = email.BuildMealPlanCreatedEmail(user.EmailAddress, emailDeliveryRequest.MealPlan, envCfg)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "building meal plan created email")
		}
		emailType = "meal plan created"

		break
	}

	if err = emailer.SendEmail(ctx, mail); err != nil {
		observability.AcknowledgeError(err, logger, span, "sending %s email", emailType)
	}

	if err = analyticsEventReporter.EventOccurred(ctx, email.SentEventType, emailDeliveryRequest.UserID, emailDeliveryRequest.TemplateParams); err != nil {
		observability.AcknowledgeError(err, logger, span, "notifying customer data platform")
	}

	return nil
}
