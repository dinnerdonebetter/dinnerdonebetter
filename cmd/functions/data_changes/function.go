package datachangesfunction

import (
	"context"
	"encoding/json"
	"fmt"

	analyticsconfig "github.com/prixfixeco/backend/internal/analytics/config"
	"github.com/prixfixeco/backend/internal/config"
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

	ctx, span := tracing.NewTracer(tracerProvider.Tracer("data_changes_job")).StartSpan(ctx)
	defer span.End()

	analyticsEventReporter, err := analyticsconfig.ProvideEventReporter(&cfg.Analytics, logger, tracerProvider)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "error setting up customer data collector")
	}

	defer analyticsEventReporter.Close()

	var changeMessage types.DataChangeMessage
	if err = json.Unmarshal(msg.Message.Data, &changeMessage); err != nil {
		logger = logger.WithValue("raw_data", msg.Message.Data)
		return observability.PrepareAndLogError(err, logger, span, "unmarshalling data change message")
	}

	logger = logger.WithValue("event_type", changeMessage.EventType)

	if changeMessage.EventType == types.UserSignedUpCustomerEventType {
		if err = analyticsEventReporter.AddUser(ctx, changeMessage.UserID, changeMessage.Context); err != nil {
			return observability.PrepareError(err, span, "notifying customer data platform")
		}

		return nil
	}

	if dataCollectionErr := analyticsEventReporter.EventOccurred(ctx, changeMessage.EventType, changeMessage.UserID, changeMessage.Context); dataCollectionErr != nil {
		observability.AcknowledgeError(dataCollectionErr, logger, span, "notifying customer data platform")
	}

	return nil
}
