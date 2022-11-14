package mealplanfinalizerfunction

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"

	"github.com/prixfixeco/backend/internal/config"
	customerdataconfig "github.com/prixfixeco/backend/internal/customerdata/config"
	"github.com/prixfixeco/backend/internal/database/postgres"
	emailconfig "github.com/prixfixeco/backend/internal/email/config"
	msgconfig "github.com/prixfixeco/backend/internal/messagequeue/config"
	"github.com/prixfixeco/backend/internal/observability"
	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/logging/zerolog"
	"github.com/prixfixeco/backend/internal/workers"
)

const (
	dataChangesTopicName = "data_changes"
)

func init() {
	// Register a CloudEvent function with the Functions Framework
	functions.CloudEvent("FinalizeMealPlans", FinalizeMealPlans)
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

// FinalizeMealPlans is our cloud function entrypoint.
func FinalizeMealPlans(ctx context.Context, _ event.Event) error {
	logger := zerolog.NewZerologLogger()
	logger.SetLevel(logging.DebugLevel)

	cfg, err := config.GetMealPlanFinalizerConfigFromGoogleCloudSecretManager(ctx)
	if err != nil {
		return fmt.Errorf("error getting config: %w", err)
	}

	tracerProvider, initializeTracerErr := cfg.Observability.Tracing.Initialize(ctx, logger)
	if initializeTracerErr != nil {
		logger.Error(initializeTracerErr, "initializing tracer")
	}

	cfg.Database.RunMigrations = false
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	emailer, err := emailconfig.ProvideEmailer(&cfg.Email, logger, client)
	if err != nil {
		log.Fatal(err)
	}

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

	mealPlanFinalizationWorker := workers.ProvideMealPlanFinalizationWorker(
		logger,
		dataManager,
		dataChangesPublisher,
		emailer,
		cdp,
		tracerProvider,
	)

	changedCount, mealPlanFinalizationErr := mealPlanFinalizationWorker.FinalizeExpiredMealPlans(ctx, nil)
	if mealPlanFinalizationErr != nil {
		return fmt.Errorf("finalizing meal plans: %w", mealPlanFinalizationErr)
	}

	if closeErr := dataManager.DB().Close(); closeErr != nil {
		observability.AcknowledgeError(closeErr, logger, nil, "closing database connection")
	}

	if changedCount > 0 {
		logger.WithValue("count", changedCount).Info("finalized meal plans")
	}

	return nil
}
