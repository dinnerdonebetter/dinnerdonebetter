package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/database/postgres"
	msgconfig "github.com/dinnerdonebetter/backend/internal/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/logging/zerolog"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/search"
	"github.com/dinnerdonebetter/backend/internal/search/indexing"

	"go.opentelemetry.io/otel"
	_ "go.uber.org/automaxprocs"
)

func main() {
	ctx := context.Background()
	logger := zerolog.NewZerologLogger(logging.DebugLevel)

	if strings.TrimSpace(strings.ToLower(os.Getenv("CEASE_OPERATION"))) == "true" {
		logger.Info("CEASE_OPERATION is set to true, exiting")
	}

	cfg, err := config.GetSearchDataIndexSchedulerConfigFromGoogleCloudSecretManager(ctx)
	if err != nil {
		log.Fatal(fmt.Errorf("error getting config: %w", err))
	}

	logger = logger.WithValue("commit", cfg.Commit())

	tracerProvider, err := cfg.Observability.Tracing.ProvideTracerProvider(ctx, logger)
	if err != nil {
		logger.Error(err, "initializing tracer")
	}
	otel.SetTracerProvider(tracerProvider)

	tracer := tracing.NewTracer(tracerProvider.Tracer("search_indexer_cloud_function"))
	ctx, span := tracer.StartSpan(ctx)
	defer span.End()

	// manual db timeout until I find out what's wrong
	dbConnectionContext, cancel := context.WithTimeout(ctx, 15*time.Second)
	dataManager, err := postgres.ProvideDatabaseClient(dbConnectionContext, logger, &cfg.Database, tracerProvider)
	if err != nil {
		cancel()
		log.Fatal(observability.PrepareError(err, span, "establishing database connection"))
	}

	if err = dataManager.DB().PingContext(ctx); err != nil {
		cancel()
		log.Fatal(observability.PrepareError(err, span, "pinging database at %s", cfg.Database.ConnectionDetails))
	}

	cancel()
	defer dataManager.Close()

	publisherProvider, err := msgconfig.ProvidePublisherProvider(logger, tracerProvider, &cfg.Events)
	if err != nil {
		log.Fatal(observability.PrepareError(err, span, "configuring queue manager"))
	}

	defer publisherProvider.Close()

	searchDataIndexPublisher, err := publisherProvider.ProvidePublisher(os.Getenv("SEARCH_INDEXING_TOPIC_NAME"))
	if err != nil {
		log.Fatal(observability.PrepareError(err, span, "configuring search indexing publisher"))
	}

	defer searchDataIndexPublisher.Stop()

	var ids []string

	// figure out what records to join
	chosenIndex := indexing.AllIndexTypes[rand.Intn(len(indexing.AllIndexTypes))]
	logger = logger.WithValue("chosen_index_type", chosenIndex)

	logger.Info("index type chosen")

	var actionFunc func(context.Context) ([]string, error)
	switch chosenIndex {
	case search.IndexTypeValidPreparations:
		actionFunc = dataManager.GetValidPreparationIDsThatNeedSearchIndexing
	case search.IndexTypeRecipes:
		actionFunc = dataManager.GetRecipeIDsThatNeedSearchIndexing
	case search.IndexTypeMeals:
		actionFunc = dataManager.GetMealIDsThatNeedSearchIndexing
	case search.IndexTypeValidIngredients:
		actionFunc = dataManager.GetValidIngredientIDsThatNeedSearchIndexing
	case search.IndexTypeValidInstruments:
		actionFunc = dataManager.GetValidInstrumentIDsThatNeedSearchIndexing
	case search.IndexTypeValidMeasurementUnits:
		actionFunc = dataManager.GetValidMeasurementUnitIDsThatNeedSearchIndexing
	case search.IndexTypeValidIngredientStates:
		actionFunc = dataManager.GetValidIngredientStateIDsThatNeedSearchIndexing
	default:
		logger.Info("unhandled index type chosen, exiting")
		return
	}

	if actionFunc != nil {
		ids, err = actionFunc(ctx)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				log.Fatal(observability.PrepareError(err, span, "getting %s IDs that need search indexing", chosenIndex))
			}
			return
		}
	} else {
		logger.Info("unspecified action function, exiting")
		return
	}

	if len(ids) > 0 {
		logger.WithValue("count", len(ids)).Info("publishing search index requests")
	}

	for _, id := range ids {
		indexReq := &indexing.IndexRequest{
			RowID:     id,
			IndexType: chosenIndex,
		}
		if err = searchDataIndexPublisher.Publish(ctx, indexReq); err != nil {
			observability.AcknowledgeError(err, logger, span, "publishing search index request")
		}
	}
}
