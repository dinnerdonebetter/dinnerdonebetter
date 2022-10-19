package mealplanfinalizerfunction

import (
	"context"
	"fmt"
	"log"

	_ "github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"go.opentelemetry.io/otel"

	"github.com/prixfixeco/api_server/internal/config"
	customerdataconfig "github.com/prixfixeco/api_server/internal/customerdata/config"
	"github.com/prixfixeco/api_server/internal/database/postgres"
	"github.com/prixfixeco/api_server/internal/features/recipeanalysis"
	msgconfig "github.com/prixfixeco/api_server/internal/messagequeue/config"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/logging/zerolog"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/internal/workers"
)

const (
	dataChangesTopicName = "data_changes"
)

// PubSubMessage is the payload of a Pub/Sub event. See the documentation for more details:
// https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// CreateMealPlanTasks is our cloud function entrypoint.
func CreateMealPlanTasks(ctx context.Context, _ PubSubMessage) error {
	logger := zerolog.NewZerologLogger()
	logger.SetLevel(logging.DebugLevel)

	cfg, err := config.GetMealPlanTaskCreatorWorkerConfigFromGoogleCloudSecretManager(ctx)
	if err != nil {
		return fmt.Errorf("error getting config: %w", err)
	}
	cfg.Database.RunMigrations = false

	tracerProvider := tracing.NewNoopTracerProvider()
	otel.SetTracerProvider(tracerProvider)

	cdp, err := customerdataconfig.ProvideCollector(&cfg.CustomerData, logger)
	if err != nil {
		log.Fatal(err)
	}

	dataManager, err := postgres.ProvideDatabaseClient(ctx, logger, &cfg.Database, tracerProvider)
	if err != nil {
		log.Fatal(err)
	}

	if !dataManager.IsReady(ctx, 50) {
		log.Fatal("database is not ready")
	}

	publisherProvider, err := msgconfig.ProvidePublisherProvider(logger, tracerProvider, &cfg.Events)
	if err != nil {
		log.Fatal(err)
	}

	dataChangesPublisher, err := publisherProvider.ProviderPublisher(dataChangesTopicName)
	if err != nil {
		log.Fatal(err)
	}

	mealPlanTaskCreationEnsurerWorker := workers.ProvideMealPlanTaskCreationEnsurerWorker(
		logger,
		dataManager,
		recipeanalysis.NewRecipeAnalyzer(logger, tracerProvider),
		dataChangesPublisher,
		cdp,
		tracerProvider,
	)

	if err = mealPlanTaskCreationEnsurerWorker.HandleMessage(ctx, nil); err != nil {
		observability.AcknowledgeError(err, logger, nil, "handling message")
		return err
	}

	if closeErr := dataManager.DB().Close(); closeErr != nil {
		observability.AcknowledgeError(closeErr, logger, nil, "closing database connection")
		return closeErr
	}

	return nil
}
