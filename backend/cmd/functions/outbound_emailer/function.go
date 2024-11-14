package outboundemailerfunction

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/dinnerdonebetter/backend/cmd/functions/outbound_emailer/logic"
	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/email"
	"github.com/dinnerdonebetter/backend/internal/observability"

	_ "github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
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

// SendEmail handles sending an email.
func SendEmail(ctx context.Context, e event.Event) error {
	if config.ShouldCeaseOperation() {
		slog.Info("CEASE_OPERATION is set to true, exiting")
		return nil
	}

	cfg, err := config.GetOutboundEmailerConfigFromGoogleCloudSecretManager(ctx)
	if err != nil {
		return fmt.Errorf("error getting config: %w", err)
	}
	cfg.Database.RunMigrations = false

	logger := cfg.Observability.Logging.ProvideLogger()

	tracerProvider, err := cfg.Observability.Tracing.ProvideTracerProvider(ctx, logger)
	if err != nil {
		logger.Error(err, "initializing tracer")
	}
	otel.SetTracerProvider(tracerProvider)

	var msg MessagePublishedData
	if err = e.DataAs(&msg); err != nil {
		return fmt.Errorf("event.DataAs: %w", err)
	}

	var emailDeliveryRequest email.DeliveryRequest
	if err = json.Unmarshal(msg.Message.Data, &emailDeliveryRequest); err != nil {
		logger = logger.WithValue("raw_data", msg.Message.Data)
		return observability.PrepareAndLogError(err, logger, nil, "unmarshalling delivery request message")
	}

	return logic.HandleSendEmailRequest(ctx, logger, tracerProvider, cfg, &emailDeliveryRequest)
}
