package mealplanfinalizerfunction

import (
	"context"
	"fmt"
	"os"
	"time"

	analyticsconfig "github.com/prixfixeco/backend/internal/analytics/config"
	"github.com/prixfixeco/backend/internal/config"
	"github.com/prixfixeco/backend/internal/database"
	"github.com/prixfixeco/backend/internal/database/postgres"
	msgconfig "github.com/prixfixeco/backend/internal/messagequeue/config"
	"github.com/prixfixeco/backend/internal/observability"
	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/logging/zerolog"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/internal/workers"

	_ "github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"go.opentelemetry.io/otel"
	_ "go.uber.org/automaxprocs"
)

func init() {
	// Register a CloudEvent function with the Functions Framework
	functions.CloudEvent("FinalizeMealPlans", FinalizeMealPlans)
}

// FinalizeMealPlans is our cloud function entrypoint.
func FinalizeMealPlans(ctx context.Context, _ event.Event) error {
	logger := zerolog.NewZerologLogger(logging.DebugLevel)

	cfg, err := config.GetMealPlanFinalizerConfigFromGoogleCloudSecretManager(ctx)
	if err != nil {
		return fmt.Errorf("error getting config: %w", err)
	}

	tracerProvider, initializeTracerErr := cfg.Observability.Tracing.ProvideTracerProvider(ctx, logger)
	if initializeTracerErr != nil {
		logger.Error(initializeTracerErr, "initializing tracer")
	}
	otel.SetTracerProvider(tracerProvider)

	ctx, span := tracing.NewTracer(tracerProvider.Tracer("meal_plan_finalizer_job")).StartSpan(ctx)
	defer span.End()

	analyticsEventReporter, err := analyticsconfig.ProvideEventReporter(&cfg.Analytics, logger, tracerProvider)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "configuring customer data collector")
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

	publisherProvider, err := msgconfig.ProvidePublisherProvider(logger, tracerProvider, &cfg.Events)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "configuring queue manager")
	}

	defer publisherProvider.Close()

	dataChangesPublisher, err := publisherProvider.ProvidePublisher(os.Getenv("DATA_CHANGES_TOPIC_NAME"))
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "configuring data changes publisher")
	}

	defer dataChangesPublisher.Stop()

	mealPlanFinalizationWorker := workers.ProvideMealPlanFinalizationWorker(
		logger,
		dataManager,
		dataChangesPublisher,
		analyticsEventReporter,
		tracerProvider,
	)

	changedCount, err := mealPlanFinalizationWorker.FinalizeExpiredMealPlans(ctx, nil)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "finalizing meal plans: %w")
	}

	if changedCount > 0 {
		logger.WithValue("count", changedCount).Info("finalized meal plans")
	}

	return nil
}
