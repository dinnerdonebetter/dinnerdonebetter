package mealplanfinalizerfunction

import (
	"context"
	"fmt"
	"github.com/prixfixeco/backend/internal/features/grocerylistpreparation"
	"log"

	_ "github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"go.opentelemetry.io/otel"

	"github.com/prixfixeco/backend/internal/config"
	customerdataconfig "github.com/prixfixeco/backend/internal/customerdata/config"
	"github.com/prixfixeco/backend/internal/database/postgres"
	"github.com/prixfixeco/backend/internal/features/recipeanalysis"
	msgconfig "github.com/prixfixeco/backend/internal/messagequeue/config"
	"github.com/prixfixeco/backend/internal/observability"
	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/logging/zerolog"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/internal/workers"
)

const (
	dataChangesTopicName = "data_changes"
)

// PubSubMessage is the payload of a Pub/Sub event. See the documentation for more details:
// https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// InitializeGroceryListsItemsForMealPlans is our cloud function entrypoint.
func InitializeGroceryListsItemsForMealPlans(ctx context.Context, _ PubSubMessage) error {
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

	mealPlanGroceryListInitializationWorker := workers.ProvideMealPlanGroceryListInitializer(
		logger,
		dataManager,
		recipeanalysis.NewRecipeAnalyzer(logger, tracerProvider),
		dataChangesPublisher,
		cdp,
		tracerProvider,
		grocerylistpreparation.NewGroceryListCreator(logger, tracerProvider),
	)

	if err = mealPlanGroceryListInitializationWorker.HandleMessage(ctx, nil); err != nil {
		observability.AcknowledgeError(err, logger, nil, "handling message")
		return err
	}

	if closeErr := dataManager.DB().Close(); closeErr != nil {
		observability.AcknowledgeError(closeErr, logger, nil, "closing database connection")
		return closeErr
	}

	return nil
}
