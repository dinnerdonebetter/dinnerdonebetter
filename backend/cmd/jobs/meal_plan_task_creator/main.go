package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	analyticscfg "github.com/dinnerdonebetter/backend/internal/analytics/config"
	"github.com/dinnerdonebetter/backend/internal/business/recipeanalysis"
	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/database/postgres"
	msgconfig "github.com/dinnerdonebetter/backend/internal/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/workers"

	"go.opentelemetry.io/otel"
	_ "go.uber.org/automaxprocs"
)

func doTheThing() error {
	ctx := context.Background()

	if config.ShouldCeaseOperation() {
		slog.Info("CEASE_OPERATION is set to true, exiting")
		return nil
	}

	cfg, err := config.LoadConfigFromEnvironment[config.MealPlanTaskCreatorConfig]()
	if err != nil {
		return fmt.Errorf("error getting config: %w", err)
	}
	cfg.Database.RunMigrations = false

	logger, err := cfg.Observability.Logging.ProvideLogger(ctx)
	if err != nil {
		return fmt.Errorf("could not create logger: %w", err)
	}

	tracerProvider, initializeTracerErr := cfg.Observability.Tracing.ProvideTracerProvider(ctx, logger)
	if initializeTracerErr != nil {
		logger.Error("initializing tracer", initializeTracerErr)
	}
	otel.SetTracerProvider(tracerProvider)

	metricsProvider, err := cfg.Observability.Metrics.ProvideMetricsProvider(ctx, logger)
	if err != nil {
		logger.Error("initializing metrics", err)
	}

	ctx, span := tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer("meal_plan_task_creator_job")).StartSpan(ctx)
	defer span.End()

	analyticsEventReporter, err := analyticscfg.ProvideEventReporter(&cfg.Analytics, logger, tracerProvider, metricsProvider)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "configuring customer data collector")
	}

	defer analyticsEventReporter.Close()

	dataManager, err := postgres.ProvideDatabaseClient(ctx, logger, tracerProvider, &cfg.Database)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "establishing database connection")
	}

	defer dataManager.Close()

	publisherProvider, err := msgconfig.ProvidePublisherProvider(ctx, logger, tracerProvider, &cfg.Events)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "configuring queue manager")
	}

	defer publisherProvider.Close()

	dataChangesPublisher, err := publisherProvider.ProvidePublisher(cfg.Queues.DataChangesTopicName)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "configuring data changes publisher")
	}

	defer dataChangesPublisher.Stop()

	mealPlanTaskCreationEnsurerWorker := workers.ProvideMealPlanTaskCreationEnsurerWorker(
		logger,
		dataManager,
		recipeanalysis.NewRecipeAnalyzer(logger, tracerProvider),
		dataChangesPublisher,
		tracerProvider,
	)

	if err = mealPlanTaskCreationEnsurerWorker.CreateMealPlanTasksForFinalizedMealPlans(ctx, nil); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "handling message")
	}

	return nil
}

func main() {
	log.Println("doing the thing")
	if err := doTheThing(); err != nil {
		log.Fatal(err)
	}
	log.Println("the thing is done")
}
