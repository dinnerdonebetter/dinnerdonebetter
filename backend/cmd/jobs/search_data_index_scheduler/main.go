package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/database/postgres"
	msgconfig "github.com/dinnerdonebetter/backend/internal/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/random"
	"github.com/dinnerdonebetter/backend/internal/search/text"
	"github.com/dinnerdonebetter/backend/internal/search/text/indexing"

	"github.com/hashicorp/go-multierror"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	_ "go.uber.org/automaxprocs"
)

func doTheThing() error {
	ctx := context.Background()

	if config.ShouldCeaseOperation() {
		slog.Info("CEASE_OPERATION is set to true, exiting")
		return nil
	}

	cfg, err := config.LoadConfigFromEnvironment[config.SearchDataIndexSchedulerConfig]()
	if err != nil {
		return fmt.Errorf("error getting config: %w", err)
	}
	cfg.Database.RunMigrations = false

	logger, tracerProvider, metricsProvider, err := cfg.Observability.ProvideThreePillars(ctx)
	if err != nil {
		return fmt.Errorf("could not establish observability pillars: %w", err)
	}

	handledRecordsCounter, err := metricsProvider.NewInt64Counter("search_data_index_scheduler.handled_records")
	if err != nil {
		return fmt.Errorf("could not create metrics counter: %w", err)
	}

	tracer := tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer("search_indexer_cloud_function"))

	ctx, span := tracer.StartSpan(ctx)
	defer span.End()

	dbConnectionContext, cancel := context.WithTimeout(ctx, 15*time.Second)
	dataManager, err := postgres.ProvideDatabaseClient(dbConnectionContext, logger, tracerProvider, &cfg.Database)
	if err != nil {
		cancel()
		return observability.PrepareError(err, span, "establishing database connection")
	}

	if err = dataManager.DB().PingContext(ctx); err != nil {
		cancel()
		return observability.PrepareError(err, span, "pinging database")
	}
	defer dataManager.Close()

	cancel()

	publisherProvider, err := msgconfig.ProvidePublisherProvider(ctx, logger, tracerProvider, &cfg.Events)
	if err != nil {
		return observability.PrepareError(err, span, "configuring queue manager")
	}
	defer publisherProvider.Close()

	searchDataIndexPublisher, err := publisherProvider.ProvidePublisher(cfg.Queues.SearchIndexRequestsTopicName)
	if err != nil {
		return observability.PrepareError(err, span, "configuring search indexing publisher")
	}
	defer searchDataIndexPublisher.Stop()

	// figure out what records to join
	chosenIndex := random.Element(indexing.AllIndexTypes)

	logger = logger.WithValue("chosen_index_type", chosenIndex)
	logger.Info("index type chosen")

	var actionFunc func(context.Context) ([]string, error)
	switch chosenIndex {
	case textsearch.IndexTypeValidPreparations:
		actionFunc = dataManager.GetValidPreparationIDsThatNeedSearchIndexing
	case textsearch.IndexTypeRecipes:
		actionFunc = dataManager.GetRecipeIDsThatNeedSearchIndexing
	case textsearch.IndexTypeMeals:
		actionFunc = dataManager.GetMealIDsThatNeedSearchIndexing
	case textsearch.IndexTypeValidIngredients:
		actionFunc = dataManager.GetValidIngredientIDsThatNeedSearchIndexing
	case textsearch.IndexTypeValidInstruments:
		actionFunc = dataManager.GetValidInstrumentIDsThatNeedSearchIndexing
	case textsearch.IndexTypeValidMeasurementUnits:
		actionFunc = dataManager.GetValidMeasurementUnitIDsThatNeedSearchIndexing
	case textsearch.IndexTypeValidIngredientStates:
		actionFunc = dataManager.GetValidIngredientStateIDsThatNeedSearchIndexing
	case textsearch.IndexTypeValidVessels:
		actionFunc = dataManager.GetValidVesselIDsThatNeedSearchIndexing
	case textsearch.IndexTypeUsers:
		actionFunc = dataManager.GetUserIDsThatNeedSearchIndexing
	default:
		logger.Info("unhandled index type chosen, exiting")
		return nil
	}

	var ids []string
	ids, err = actionFunc(ctx)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			observability.AcknowledgeError(err, logger, span, "getting %s IDs that need search indexing", chosenIndex)
			return err
		}
		return nil
	}

	if len(ids) > 0 {
		logger.WithValue("count", len(ids)).Info("publishing search index requests")
	}

	publishedIDCount := int64(0)
	errs := &multierror.Error{}
	for _, id := range ids {
		indexReq := &indexing.IndexRequest{
			RowID:     id,
			IndexType: chosenIndex,
		}
		if err = searchDataIndexPublisher.Publish(ctx, indexReq); err != nil {
			errs = multierror.Append(errs, err)
		} else {
			publishedIDCount++
		}
	}

	handledRecordsCounter.Add(ctx, publishedIDCount, metric.WithAttributes(
		attribute.KeyValue{
			Key:   "record.type",
			Value: attribute.StringValue(chosenIndex),
		},
	))

	return errs.ErrorOrNil()
}

func main() {
	log.Println("doing the thing")
	if err := doTheThing(); err != nil {
		log.Fatal(err)
	}
	log.Println("the thing is done")
}
