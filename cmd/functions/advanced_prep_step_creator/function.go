package mealplanfinalizerfunction

import (
	"context"
	"fmt"
	"log"

	_ "github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"go.opentelemetry.io/otel"

	"github.com/prixfixeco/api_server/internal/config"
	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/database/postgres"
	"github.com/prixfixeco/api_server/internal/messagequeue"
	msgconfig "github.com/prixfixeco/api_server/internal/messagequeue/config"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/logging/zerolog"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
)

const (
	dataChangesTopicName = "data_changes"
)

// PubSubMessage is the payload of a Pub/Sub event. See the documentation for more details:
// https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// CreateAdvancedPrepSteps is our cloud function entrypoint.
func CreateAdvancedPrepSteps(ctx context.Context, m PubSubMessage) error {
	logger := zerolog.NewZerologLogger()
	logger.SetLevel(logging.DebugLevel)

	cfg, err := config.GetAdvancedPrepStepCreatorWorkerConfigFromGoogleCloudSecretManager(ctx)
	if err != nil {
		return fmt.Errorf("error getting config: %w", err)
	}

	tracerProvider := tracing.NewNoopTracerProvider()
	otel.SetTracerProvider(tracerProvider)
	tracer := tracing.NewTracer(tracerProvider.Tracer("meal_plan_finalizer"))

	dataManager, err := postgres.ProvideDatabaseClient(ctx, logger, &cfg.Database, tracerProvider)
	if err != nil {
		return fmt.Errorf("error setting up database client: %w", err)
	}

	publisherProvider, err := msgconfig.ProvidePublisherProvider(logger, tracerProvider, &cfg.Events)
	if err != nil {
		log.Fatal(err)
	}

	dataChangesPublisher, err := publisherProvider.ProviderPublisher(dataChangesTopicName)
	if err != nil {
		log.Fatal(err)
	}

	if mainErr := ensureAdvancedPrepStepsAreCreatedForUpcomingMealPlans(ctx, tracer, dataManager, dataChangesPublisher); mainErr != nil {
		observability.AcknowledgeError(mainErr, logger, nil, "closing database connection")
		return mainErr
	}

	if closeErr := dataManager.DB().Close(); closeErr != nil {
		observability.AcknowledgeError(closeErr, logger, nil, "closing database connection")
		return closeErr
	}

	return nil
}

func ensureAdvancedPrepStepsAreCreatedForUpcomingMealPlans(ctx context.Context, tracer tracing.Tracer, dbmanager database.DataManager, publisher messagequeue.Publisher) error {
	_, span := tracer.StartSpan(ctx)
	defer span.End()

	// get all meal plans that are due to start in the coming week

	// iterate through the chosen meal plan options and for the chosen ones, check to see if they have any pure steps

	// for all pure steps, figure out which, if any, need to be created in the database

	return nil
}
