package datachangesfunction

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/prixfixeco/backend/internal/analytics"
	analyticsconfig "github.com/prixfixeco/backend/internal/analytics/config"
	"github.com/prixfixeco/backend/internal/config"
	"github.com/prixfixeco/backend/internal/database"
	"github.com/prixfixeco/backend/internal/database/postgres"
	"github.com/prixfixeco/backend/internal/email"
	"github.com/prixfixeco/backend/internal/messagequeue"
	msgconfig "github.com/prixfixeco/backend/internal/messagequeue/config"
	"github.com/prixfixeco/backend/internal/observability"
	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/logging/zerolog"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/pkg/types"

	_ "github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"go.opentelemetry.io/otel"
	_ "go.uber.org/automaxprocs"
)

func init() {
	// Register a CloudEvent function with the Functions Framework
	functions.CloudEvent("ProcessDataChange", ProcessDataChange)
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

// ProcessDataChange handles a data change.
func ProcessDataChange(ctx context.Context, e event.Event) error {
	var msg MessagePublishedData
	if err := e.DataAs(&msg); err != nil {
		return fmt.Errorf("event.DataAs: %v", err)
	}

	logger := zerolog.NewZerologLogger(logging.DebugLevel)

	cfg, err := config.GetDataChangesWorkerConfigFromGoogleCloudSecretManager(ctx)
	if err != nil {
		return fmt.Errorf("error getting config: %w", err)
	}

	tracerProvider, initializeTracerErr := cfg.Observability.Tracing.ProvideTracerProvider(ctx, logger)
	if initializeTracerErr != nil {
		logger.Error(initializeTracerErr, "initializing tracer")
	}
	otel.SetTracerProvider(tracerProvider)

	tracer := tracing.NewTracer(tracerProvider.Tracer("data_changes_job"))

	ctx, span := tracer.StartSpan(ctx)
	defer span.End()

	analyticsEventReporter, err := analyticsconfig.ProvideEventReporter(&cfg.Analytics, logger, tracerProvider)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "error setting up customer data collector")
	}

	defer analyticsEventReporter.Close()

	publisherProvider, err := msgconfig.ProvidePublisherProvider(logger, tracerProvider, &cfg.Events)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "configuring queue manager")
	}

	defer publisherProvider.Close()

	outboundEmailsPublisher, err := publisherProvider.ProvidePublisher(os.Getenv("OUTBOUND_EMAILS_TOPIC_NAME"))
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "configuring data changes publisher")
	}

	defer outboundEmailsPublisher.Stop()

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

	var changeMessage types.DataChangeMessage
	if err = json.Unmarshal(msg.Message.Data, &changeMessage); err != nil {
		logger = logger.WithValue("raw_data", msg.Message.Data)
		return observability.PrepareAndLogError(err, logger, span, "unmarshalling data change message")
	}

	logger = logger.WithValue("event_type", changeMessage.EventType)

	if changeMessage.UserID != "" && changeMessage.EventType != "" {
		if err = analyticsEventReporter.EventOccurred(ctx, changeMessage.EventType, changeMessage.UserID, changeMessage.Context); err != nil {
			observability.AcknowledgeError(err, logger, span, "notifying customer data platform")
		}
	}

	if err = handleOutboundNotifications(ctx, logger, tracer, dataManager, outboundEmailsPublisher, analyticsEventReporter, changeMessage); err != nil {
		observability.AcknowledgeError(err, logger, span, "notifying customer(s)")
	}

	return nil
}

func handleOutboundNotifications(
	ctx context.Context,
	logger logging.Logger,
	tracer tracing.Tracer,
	dataManager database.DataManager,
	outboundEmailsPublisher messagequeue.Publisher,
	analyticsEventReporter analytics.EventReporter,
	changeMessage types.DataChangeMessage,
) error {
	ctx, span := tracer.StartSpan(ctx)
	defer span.End()

	var (
		emailType string
		edrs      []*email.DeliveryRequest
	)

	switch changeMessage.EventType {
	case types.UserSignedUpCustomerEventType:
		emailType = "user signup"
		if err := analyticsEventReporter.AddUser(ctx, changeMessage.UserID, changeMessage.Context); err != nil {
			observability.AcknowledgeError(err, logger, span, "notifying customer data platform")
		}

		edrs = append(edrs, &email.DeliveryRequest{
			UserID:                 changeMessage.UserID,
			Template:               email.TemplateTypeVerifyEmailAddress,
			EmailVerificationToken: changeMessage.EmailVerificationToken,
		})

		break
	case types.MealPlanCreatedCustomerEventType:
		emailType = "meal plan created"
		mealPlan := changeMessage.MealPlan
		if mealPlan == nil {
			return observability.PrepareError(fmt.Errorf("meal plan is nil"), span, "publishing meal plan created email")
		}

		household, err := dataManager.GetHousehold(ctx, mealPlan.BelongsToHousehold)
		if err != nil {
			return observability.PrepareError(err, span, "getting household")
		}

		for _, member := range household.Members {
			if member.BelongsToUser.EmailAddressVerifiedAt != nil {
				edrs = append(edrs, &email.DeliveryRequest{
					UserID:   member.BelongsToUser.ID,
					Template: email.TemplateTypeMealPlanCreated,
					MealPlan: mealPlan,
				})
			}
		}

		break
	case types.PasswordResetTokenCreatedEventType:
		emailType = "password reset request"
		if changeMessage.PasswordResetToken == nil {
			return observability.PrepareError(fmt.Errorf("password reset token is nil"), span, "publishing password reset token email")
		}

		edrs = append(edrs, &email.DeliveryRequest{
			UserID:             changeMessage.UserID,
			Template:           email.TemplateTypePasswordReset,
			PasswordResetToken: changeMessage.PasswordResetToken,
		})

		break

	case types.UsernameReminderRequestedEventType:
		emailType = "username reminder"
		edrs = append(edrs, &email.DeliveryRequest{
			UserID:   changeMessage.UserID,
			Template: email.TemplateTypeUsernameReminder,
		})

		break

	case types.PasswordResetTokenRedeemedEventType:
		emailType = "password reset token redeemed"
		edrs = append(edrs, &email.DeliveryRequest{
			UserID:   changeMessage.UserID,
			Template: email.TemplateTypePasswordResetTokenRedeemed,
		})

		break
	case types.HouseholdInvitationCreatedCustomerEventType:
		emailType = "household invitation created"
		if changeMessage.HouseholdInvitation == nil {
			return observability.PrepareError(fmt.Errorf("household invitation is nil"), span, "publishing password reset token redemption email")
		}

		edrs = append(edrs, &email.DeliveryRequest{
			UserID:     changeMessage.UserID,
			Template:   email.TemplateTypeInvite,
			Invitation: changeMessage.HouseholdInvitation,
		})

		break
	}

	for _, edr := range edrs {
		if err := outboundEmailsPublisher.Publish(ctx, edr); err != nil {
			observability.AcknowledgeError(err, logger, span, "publishing %s request email", emailType)
		}
	}

	return nil
}
