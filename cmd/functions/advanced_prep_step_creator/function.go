package mealplanfinalizerfunction

import (
	"context"
	"fmt"
	customerdataconfig "github.com/prixfixeco/api_server/internal/customerdata/config"
	"github.com/prixfixeco/api_server/internal/recipeanalysis"
	"github.com/prixfixeco/api_server/internal/workers"
	"log"

	_ "github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/prixfixeco/api_server/internal/config"
	"github.com/prixfixeco/api_server/internal/database/postgres"
	msgconfig "github.com/prixfixeco/api_server/internal/messagequeue/config"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/logging/zerolog"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"go.opentelemetry.io/otel"
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
func CreateAdvancedPrepSteps(ctx context.Context, _ PubSubMessage) error {
	logger := zerolog.NewZerologLogger()
	logger.SetLevel(logging.DebugLevel)

	cfg, err := config.GetAdvancedPrepStepCreatorWorkerConfigFromGoogleCloudSecretManager(ctx)
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

	advancedPrepStepCreationEnsurerWorker := workers.ProvideAdvancedPrepStepCreationEnsurerWorker(
		logger,
		dataManager,
		recipeanalysis.NewRecipeAnalyzer(tracerProvider),
		dataChangesPublisher,
		cdp,
		tracerProvider,
	)

	if err = advancedPrepStepCreationEnsurerWorker.HandleMessage(ctx, nil); err != nil {
		observability.AcknowledgeError(err, logger, nil, "closing database connection")
		return err
	}

	if closeErr := dataManager.DB().Close(); closeErr != nil {
		observability.AcknowledgeError(closeErr, logger, nil, "closing database connection")
		return closeErr
	}

	return nil
}
