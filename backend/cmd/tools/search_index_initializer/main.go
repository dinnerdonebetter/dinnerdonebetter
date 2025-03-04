package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	databasecfg "github.com/dinnerdonebetter/backend/internal/database/config"
	"github.com/dinnerdonebetter/backend/internal/database/postgres"
	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	loggingcfg "github.com/dinnerdonebetter/backend/internal/lib/observability/logging/config"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/lib/pointer"
	textsearch "github.com/dinnerdonebetter/backend/internal/lib/search/text"
	"github.com/dinnerdonebetter/backend/internal/lib/search/text/algolia"
	textsearchcfg "github.com/dinnerdonebetter/backend/internal/lib/search/text/config"
	coreindexing "github.com/dinnerdonebetter/backend/internal/services/core/indexing"
	eatingindexing "github.com/dinnerdonebetter/backend/internal/services/eating/indexing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

func main() {
	indicesPtr := flag.String("indices", "", "indices to initialize")
	wipePtr := flag.Bool("wipe", false, "whether to wipe the indices or not")

	flag.Parse()

	indices := strings.Split(*indicesPtr, ",")
	if len(indices) == 0 {
		log.Fatal("indices are required")
	}

	ctx := context.Background()

	logger, err := (&loggingcfg.Config{Level: logging.DebugLevel, Provider: loggingcfg.ProviderSlog}).ProvideLogger(ctx)
	if err != nil {
		log.Fatalf("could not create logger: %v", err)
	}

	tracerProvider := tracing.NewNoopTracerProvider()
	metricsProvider := metrics.NewNoopMetricsProvider()

	cfg := &textsearchcfg.Config{
		Provider: textsearchcfg.AlgoliaProvider,
		Algolia: &algolia.Config{
			AppID:  os.Getenv("ALGOLIA_APP_ID"),
			APIKey: os.Getenv("ALGOLIA_API_KEY"),
		},
	}

	dbConfig := &databasecfg.Config{}
	if err = dbConfig.LoadConnectionDetailsFromURL(os.Getenv("DATABASE_URL")); err != nil {
		log.Fatal(err)
	}

	dataManager, err := postgres.ProvideDatabaseClient(ctx, logger, tracerProvider, dbConfig)
	if dataManager != nil {
		defer dataManager.Close()
	}

	if err != nil {
		log.Println(fmt.Errorf("initializing database client: %w", err))
		return
	}

	var (
		im               textsearch.IndexManager
		indexRequestChan = make(chan *textsearch.IndexRequest)
		wipeOnce         sync.Once
		waitGroup        sync.WaitGroup
	)

	go func() {
		for x := range indexRequestChan {
			wipeOnce.Do(func() {
				if *wipePtr {
					log.Println("wiping index")
					if err = im.Wipe(ctx); err != nil {
						log.Println(fmt.Errorf("wiping index: %w", err))
						return
					}
					log.Println("wiped index")
				}
			})

			if err = eatingindexing.HandleIndexRequest(ctx, logger, tracerProvider, metricsProvider, cfg, dataManager, x); err != nil {
				observability.AcknowledgeError(err, logger, nil, "indexing row")
			}

			waitGroup.Done()
		}
	}()

	for i, index := range indices {
		if i > 0 {
			waitGroup.Wait()
		}

		filter := filtering.DefaultQueryFilter()
		filter.Limit = pointer.To(uint8(50))
		thresholdMet := false

		switch index {
		case eatingindexing.IndexTypeRecipes:
			im, err = textsearchcfg.ProvideIndex[eatingindexing.RecipeSearchSubset](ctx, logger, tracerProvider, metricsProvider, cfg, index)
			if err != nil {
				observability.AcknowledgeError(err, logger, nil, "initializing index manager")
				return
			}

			for !thresholdMet {
				var data *filtering.QueryFilteredResult[types.Recipe]
				data, err = dataManager.GetRecipes(ctx, filter)
				if err != nil {
					log.Println(fmt.Errorf("getting Recipe data: %w", err))
					return
				}

				for _, x := range data.Data {
					indexRequestChan <- &textsearch.IndexRequest{
						RowID:     x.ID,
						IndexType: eatingindexing.IndexTypeRecipes,
					}
					waitGroup.Add(1)
				}

				thresholdMet = len(data.Data) == 0
				*filter.Page++
			}
		case eatingindexing.IndexTypeMeals:
			im, err = textsearchcfg.ProvideIndex[eatingindexing.MealSearchSubset](ctx, logger, tracerProvider, metricsProvider, cfg, index)
			if err != nil {
				observability.AcknowledgeError(err, logger, nil, "initializing index manager")
				return
			}

			for !thresholdMet {
				var data *filtering.QueryFilteredResult[types.Meal]
				data, err = dataManager.GetMeals(ctx, filter)
				if err != nil {
					log.Println(fmt.Errorf("getting Meal data: %w", err))
					return
				}

				for _, x := range data.Data {
					indexRequestChan <- &textsearch.IndexRequest{
						RowID:     x.ID,
						IndexType: eatingindexing.IndexTypeMeals,
					}
					waitGroup.Add(1)
				}

				thresholdMet = len(data.Data) == 0
				*filter.Page++
			}
		case eatingindexing.IndexTypeValidIngredients:
			im, err = textsearchcfg.ProvideIndex[eatingindexing.ValidIngredientSearchSubset](ctx, logger, tracerProvider, metricsProvider, cfg, index)
			if err != nil {
				observability.AcknowledgeError(err, logger, nil, "initializing index manager")
				return
			}

			for !thresholdMet {
				var data *filtering.QueryFilteredResult[types.ValidIngredient]
				data, err = dataManager.GetValidIngredients(ctx, filter)
				if err != nil {
					log.Println(fmt.Errorf("getting ValidIngredient data: %w", err))
					return
				}

				for _, x := range data.Data {
					indexRequestChan <- &textsearch.IndexRequest{
						RowID:     x.ID,
						IndexType: eatingindexing.IndexTypeValidIngredients,
					}
					waitGroup.Add(1)
				}

				thresholdMet = len(data.Data) == 0
				*filter.Page++
			}
		case eatingindexing.IndexTypeValidInstruments:
			im, err = textsearchcfg.ProvideIndex[eatingindexing.ValidInstrumentSearchSubset](ctx, logger, tracerProvider, metricsProvider, cfg, index)
			if err != nil {
				observability.AcknowledgeError(err, logger, nil, "initializing index manager")
				return
			}

			for !thresholdMet {
				var data *filtering.QueryFilteredResult[types.ValidInstrument]
				data, err = dataManager.GetValidInstruments(ctx, filter)
				if err != nil {
					log.Println(fmt.Errorf("getting ValidInstrument data: %w", err))
					return
				}

				for _, x := range data.Data {
					indexRequestChan <- &textsearch.IndexRequest{
						RowID:     x.ID,
						IndexType: eatingindexing.IndexTypeValidInstruments,
					}
					waitGroup.Add(1)
				}

				thresholdMet = len(data.Data) == 0
				*filter.Page++
			}
		case eatingindexing.IndexTypeValidMeasurementUnits:
			im, err = textsearchcfg.ProvideIndex[eatingindexing.ValidMeasurementUnitSearchSubset](ctx, logger, tracerProvider, metricsProvider, cfg, index)
			if err != nil {
				observability.AcknowledgeError(err, logger, nil, "initializing index manager")
				return
			}

			for !thresholdMet {
				var data *filtering.QueryFilteredResult[types.ValidMeasurementUnit]
				data, err = dataManager.GetValidMeasurementUnits(ctx, filter)
				if err != nil {
					log.Println(fmt.Errorf("getting ValidMeasurementUnit data: %w", err))
					return
				}

				for _, x := range data.Data {
					indexRequestChan <- &textsearch.IndexRequest{
						RowID:     x.ID,
						IndexType: eatingindexing.IndexTypeValidMeasurementUnits,
					}
					waitGroup.Add(1)
				}

				thresholdMet = len(data.Data) == 0
				*filter.Page++
			}
		case eatingindexing.IndexTypeValidPreparations:
			im, err = textsearchcfg.ProvideIndex[eatingindexing.ValidPreparationSearchSubset](ctx, logger, tracerProvider, metricsProvider, cfg, index)
			if err != nil {
				observability.AcknowledgeError(err, logger, nil, "initializing index manager")
				return
			}

			for !thresholdMet {
				var data *filtering.QueryFilteredResult[types.ValidPreparation]
				data, err = dataManager.GetValidPreparations(ctx, filter)
				if err != nil {
					log.Println(fmt.Errorf("getting ValidPreparation data: %w", err))
					return
				}

				for _, x := range data.Data {
					indexRequestChan <- &textsearch.IndexRequest{
						RowID:     x.ID,
						IndexType: eatingindexing.IndexTypeValidPreparations,
					}
					waitGroup.Add(1)
				}

				thresholdMet = len(data.Data) == 0
				*filter.Page++
			}
		case eatingindexing.IndexTypeValidIngredientStates:
			im, err = textsearchcfg.ProvideIndex[eatingindexing.ValidIngredientStateSearchSubset](ctx, logger, tracerProvider, metricsProvider, cfg, index)
			if err != nil {
				observability.AcknowledgeError(err, logger, nil, "initializing index manager")
				return
			}

			for !thresholdMet {
				var data *filtering.QueryFilteredResult[types.ValidIngredientState]
				data, err = dataManager.GetValidIngredientStates(ctx, filter)
				if err != nil {
					log.Println(fmt.Errorf("getting ValidIngredientState data: %w", err))
					return
				}

				for _, x := range data.Data {
					indexRequestChan <- &textsearch.IndexRequest{
						RowID:     x.ID,
						IndexType: eatingindexing.IndexTypeValidIngredientStates,
					}
					waitGroup.Add(1)
				}

				thresholdMet = len(data.Data) == 0
				*filter.Page++
			}
		case eatingindexing.IndexTypeValidVessels:
			im, err = textsearchcfg.ProvideIndex[eatingindexing.ValidVesselSearchSubset](ctx, logger, tracerProvider, metricsProvider, cfg, index)
			if err != nil {
				observability.AcknowledgeError(err, logger, nil, "initializing index manager")
				return
			}

			for !thresholdMet {
				var data *filtering.QueryFilteredResult[types.ValidVessel]
				data, err = dataManager.GetValidVessels(ctx, filter)
				if err != nil {
					log.Println(fmt.Errorf("getting ValidVessel data: %w", err))
					return
				}

				for _, x := range data.Data {
					indexRequestChan <- &textsearch.IndexRequest{
						RowID:     x.ID,
						IndexType: eatingindexing.IndexTypeValidVessels,
					}
					waitGroup.Add(1)
				}

				thresholdMet = len(data.Data) == 0
				*filter.Page++
			}
		case coreindexing.IndexTypeUsers:
			im, err = textsearchcfg.ProvideIndex[coreindexing.UserSearchSubset](ctx, logger, tracerProvider, metricsProvider, cfg, index)
			if err != nil {
				observability.AcknowledgeError(err, logger, nil, "initializing index manager")
				return
			}

			for !thresholdMet {
				var data *filtering.QueryFilteredResult[types.User]
				data, err = dataManager.GetUsers(ctx, filter)
				if err != nil {
					log.Println(fmt.Errorf("getting user data: %w", err))
					return
				}

				for _, x := range data.Data {
					indexRequestChan <- &textsearch.IndexRequest{
						RowID:     x.ID,
						IndexType: coreindexing.IndexTypeUsers,
					}
					waitGroup.Add(1)
				}

				thresholdMet = len(data.Data) == 0
				*filter.Page++
			}
		}
	}
}
