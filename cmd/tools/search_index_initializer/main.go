package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/dinnerdonebetter/backend/internal/database"
	dbconfig "github.com/dinnerdonebetter/backend/internal/database/config"
	"github.com/dinnerdonebetter/backend/internal/database/postgres"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/logging/zap"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/pointers"
	"github.com/dinnerdonebetter/backend/internal/search"
	"github.com/dinnerdonebetter/backend/internal/search/algolia"
	searchcfg "github.com/dinnerdonebetter/backend/internal/search/config"
	"github.com/dinnerdonebetter/backend/internal/search/indexing"
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
	logger := zap.NewZapLogger(logging.DebugLevel)
	tracerProvider := tracing.NewNoopTracerProvider()

	cfg := &searchcfg.Config{
		Provider: searchcfg.AlgoliaProvider,
		Algolia: &algolia.Config{
			AppID:  os.Getenv("ALGOLIA_APP_ID"),
			APIKey: os.Getenv("ALGOLIA_API_KEY"),
		},
	}

	dbConfig := &dbconfig.Config{
		ConnectionDetails: database.ConnectionDetails(os.Getenv("DATABASE_URL")),
	}

	dataManager, err := postgres.ProvideDatabaseClient(ctx, logger, dbConfig, tracing.NewNoopTracerProvider())
	if dataManager != nil {
		defer dataManager.Close()
	}

	if err != nil {
		log.Println(fmt.Errorf("initializing database client: %w", err))
		return
	}

	var (
		im               search.IndexManager
		indexRequestChan = make(chan *indexing.IndexRequest)
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

			if err = indexing.HandleIndexRequest(ctx, logger, tracerProvider, cfg, dataManager, x); err != nil {
				observability.AcknowledgeError(err, logger, nil, "indexing row")
			}

			waitGroup.Done()
		}
	}()

	for i, index := range indices {
		if i > 0 {
			waitGroup.Wait()
		}

		filter := types.DefaultQueryFilter()
		filter.Limit = pointers.Pointer(uint8(50))
		thresholdMet := false

		switch index {
		case indexing.IndexTypeRecipes:
			im, err = searchcfg.ProvideIndexManager[types.RecipeSearchSubset](ctx, logger, tracerProvider, cfg, index)
			if err != nil {
				observability.AcknowledgeError(err, logger, nil, "initializing index manager")
				return
			}

			for !thresholdMet {
				var data *types.QueryFilteredResult[types.Recipe]
				data, err = dataManager.GetRecipes(ctx, filter)
				if err != nil {
					log.Println(fmt.Errorf("getting Recipe data: %w", err))
					return
				}

				for _, x := range data.Data {
					indexRequestChan <- &indexing.IndexRequest{
						RowID:     x.ID,
						IndexType: indexing.IndexTypeRecipes,
					}
					waitGroup.Add(1)
				}

				thresholdMet = len(data.Data) == 0
				*filter.Page++
			}
		case indexing.IndexTypeMeals:
			im, err = searchcfg.ProvideIndexManager[types.MealSearchSubset](ctx, logger, tracerProvider, cfg, index)
			if err != nil {
				observability.AcknowledgeError(err, logger, nil, "initializing index manager")
				return
			}

			for !thresholdMet {
				var data *types.QueryFilteredResult[types.Meal]
				data, err = dataManager.GetMeals(ctx, filter)
				if err != nil {
					log.Println(fmt.Errorf("getting Meal data: %w", err))
					return
				}

				for _, x := range data.Data {
					indexRequestChan <- &indexing.IndexRequest{
						RowID:     x.ID,
						IndexType: indexing.IndexTypeMeals,
					}
					waitGroup.Add(1)
				}

				thresholdMet = len(data.Data) == 0
				*filter.Page++
			}
		case indexing.IndexTypeValidIngredients:
			im, err = searchcfg.ProvideIndexManager[types.ValidIngredient](ctx, logger, tracerProvider, cfg, index)
			if err != nil {
				observability.AcknowledgeError(err, logger, nil, "initializing index manager")
				return
			}

			for !thresholdMet {
				var data *types.QueryFilteredResult[types.ValidIngredient]
				data, err = dataManager.GetValidIngredients(ctx, filter)
				if err != nil {
					log.Println(fmt.Errorf("getting ValidIngredient data: %w", err))
					return
				}

				for _, x := range data.Data {
					indexRequestChan <- &indexing.IndexRequest{
						RowID:     x.ID,
						IndexType: indexing.IndexTypeValidIngredients,
					}
					waitGroup.Add(1)
				}

				thresholdMet = len(data.Data) == 0
				*filter.Page++
			}
		case indexing.IndexTypeValidInstruments:
			im, err = searchcfg.ProvideIndexManager[types.ValidInstrument](ctx, logger, tracerProvider, cfg, index)
			if err != nil {
				observability.AcknowledgeError(err, logger, nil, "initializing index manager")
				return
			}

			for !thresholdMet {
				var data *types.QueryFilteredResult[types.ValidInstrument]
				data, err = dataManager.GetValidInstruments(ctx, filter)
				if err != nil {
					log.Println(fmt.Errorf("getting ValidInstrument data: %w", err))
					return
				}

				for _, x := range data.Data {
					indexRequestChan <- &indexing.IndexRequest{
						RowID:     x.ID,
						IndexType: indexing.IndexTypeValidInstruments,
					}
					waitGroup.Add(1)
				}

				thresholdMet = len(data.Data) == 0
				*filter.Page++
			}
		case indexing.IndexTypeValidMeasurementUnits:
			im, err = searchcfg.ProvideIndexManager[types.ValidMeasurementUnit](ctx, logger, tracerProvider, cfg, index)
			if err != nil {
				observability.AcknowledgeError(err, logger, nil, "initializing index manager")
				return
			}

			for !thresholdMet {
				var data *types.QueryFilteredResult[types.ValidMeasurementUnit]
				data, err = dataManager.GetValidMeasurementUnits(ctx, filter)
				if err != nil {
					log.Println(fmt.Errorf("getting ValidMeasurementUnit data: %w", err))
					return
				}

				for _, x := range data.Data {
					indexRequestChan <- &indexing.IndexRequest{
						RowID:     x.ID,
						IndexType: indexing.IndexTypeValidMeasurementUnits,
					}
					waitGroup.Add(1)
				}

				thresholdMet = len(data.Data) == 0
				*filter.Page++
			}
		case indexing.IndexTypeValidPreparations:
			im, err = searchcfg.ProvideIndexManager[types.ValidPreparation](ctx, logger, tracerProvider, cfg, index)
			if err != nil {
				observability.AcknowledgeError(err, logger, nil, "initializing index manager")
				return
			}

			for !thresholdMet {
				var data *types.QueryFilteredResult[types.ValidPreparation]
				data, err = dataManager.GetValidPreparations(ctx, filter)
				if err != nil {
					log.Println(fmt.Errorf("getting ValidPreparation data: %w", err))
					return
				}

				for _, x := range data.Data {
					indexRequestChan <- &indexing.IndexRequest{
						RowID:     x.ID,
						IndexType: indexing.IndexTypeValidPreparations,
					}
					waitGroup.Add(1)
				}

				thresholdMet = len(data.Data) == 0
				*filter.Page++
			}
		case indexing.IndexTypeValidIngredientStates:
			im, err = searchcfg.ProvideIndexManager[types.ValidIngredientState](ctx, logger, tracerProvider, cfg, index)
			if err != nil {
				observability.AcknowledgeError(err, logger, nil, "initializing index manager")
				return
			}

			for !thresholdMet {
				var data *types.QueryFilteredResult[types.ValidIngredientState]
				data, err = dataManager.GetValidIngredientStates(ctx, filter)
				if err != nil {
					log.Println(fmt.Errorf("getting ValidIngredientState data: %w", err))
					return
				}

				for _, x := range data.Data {
					indexRequestChan <- &indexing.IndexRequest{
						RowID:     x.ID,
						IndexType: indexing.IndexTypeValidIngredientStates,
					}
					waitGroup.Add(1)
				}

				thresholdMet = len(data.Data) == 0
				*filter.Page++
			}
		case indexing.IndexTypeValidIngredientMeasurementUnits:
			im, err = searchcfg.ProvideIndexManager[types.ValidIngredientMeasurementUnit](ctx, logger, tracerProvider, cfg, index)
			if err != nil {
				observability.AcknowledgeError(err, logger, nil, "initializing index manager")
				return
			}

			for !thresholdMet {
				var data *types.QueryFilteredResult[types.ValidIngredientMeasurementUnit]
				data, err = dataManager.GetValidIngredientMeasurementUnits(ctx, filter)
				if err != nil {
					log.Println(fmt.Errorf("getting ValidIngredientMeasurementUnit data: %w", err))
					return
				}

				for _, x := range data.Data {
					indexRequestChan <- &indexing.IndexRequest{
						RowID:     x.ID,
						IndexType: indexing.IndexTypeValidIngredientMeasurementUnits,
					}
					waitGroup.Add(1)
				}

				thresholdMet = len(data.Data) == 0
				*filter.Page++
			}
		case indexing.IndexTypeValidMeasurementUnitConversions:
			im, err = searchcfg.ProvideIndexManager[types.ValidMeasurementUnitConversion](ctx, logger, tracerProvider, cfg, index)
			if err != nil {
				observability.AcknowledgeError(err, logger, nil, "initializing index manager")
				return
			}

			for !thresholdMet {
				var measurementUnits *types.QueryFilteredResult[types.ValidMeasurementUnit]
				measurementUnits, err = dataManager.GetValidMeasurementUnits(ctx, filter)
				if err != nil {
					log.Println(fmt.Errorf("getting ValidMeasurementUnit data: %w", err))
					return
				}

				for _, y := range measurementUnits.Data {
					var fromUnit []*types.ValidMeasurementUnitConversion
					fromUnit, err = dataManager.GetValidMeasurementUnitConversionsFromUnit(ctx, y.ID)
					if err != nil {
						log.Println(fmt.Errorf("getting ValidMeasurementUnitConversion data: %w", err))
						return
					}

					for _, z := range fromUnit {
						indexRequestChan <- &indexing.IndexRequest{
							RowID:     z.ID,
							IndexType: indexing.IndexTypeValidMeasurementUnitConversions,
						}
						waitGroup.Add(1)
					}

					var toUnit []*types.ValidMeasurementUnitConversion
					toUnit, err = dataManager.GetValidMeasurementUnitConversionsToUnit(ctx, y.ID)
					if err != nil {
						log.Println(fmt.Errorf("getting ValidMeasurementUnitConversion data: %w", err))
						return
					}

					for _, z := range toUnit {
						indexRequestChan <- &indexing.IndexRequest{
							RowID:     z.ID,
							IndexType: indexing.IndexTypeValidMeasurementUnitConversions,
						}
						waitGroup.Add(1)
					}
				}

				thresholdMet = len(measurementUnits.Data) == 0
				*filter.Page++
			}
		case indexing.IndexTypeValidPreparationInstruments:
			im, err = searchcfg.ProvideIndexManager[types.ValidPreparationInstrument](ctx, logger, tracerProvider, cfg, index)
			if err != nil {
				observability.AcknowledgeError(err, logger, nil, "initializing index manager")
				return
			}

			for !thresholdMet {
				var data *types.QueryFilteredResult[types.ValidPreparationInstrument]
				data, err = dataManager.GetValidPreparationInstruments(ctx, filter)
				if err != nil {
					log.Println(fmt.Errorf("getting ValidPreparationInstrument data: %w", err))
					return
				}

				for _, x := range data.Data {
					indexRequestChan <- &indexing.IndexRequest{
						RowID:     x.ID,
						IndexType: indexing.IndexTypeValidPreparationInstruments,
					}
					waitGroup.Add(1)
				}

				thresholdMet = len(data.Data) == 0
				*filter.Page++
			}
		case indexing.IndexTypeValidIngredientPreparations:
			im, err = searchcfg.ProvideIndexManager[types.ValidIngredientPreparation](ctx, logger, tracerProvider, cfg, index)
			if err != nil {
				observability.AcknowledgeError(err, logger, nil, "initializing index manager")
				return
			}

			for !thresholdMet {
				var data *types.QueryFilteredResult[types.ValidIngredientPreparation]
				data, err = dataManager.GetValidIngredientPreparations(ctx, filter)
				if err != nil {
					log.Println(fmt.Errorf("getting ValidIngredientPreparation data: %w", err))
					return
				}

				for _, x := range data.Data {
					indexRequestChan <- &indexing.IndexRequest{
						RowID:     x.ID,
						IndexType: indexing.IndexTypeValidIngredientPreparations,
					}
					waitGroup.Add(1)
				}

				thresholdMet = len(data.Data) == 0
				*filter.Page++
			}
		}
	}
}
