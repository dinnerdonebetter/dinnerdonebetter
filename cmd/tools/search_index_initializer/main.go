package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
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
	indexPtr := flag.String("index", "", "index to initialize")
	wipePtr := flag.Bool("wipe", false, "whether to wipe the index or not")

	flag.Parse()

	index := *indexPtr
	if index == "" {
		log.Fatal("index is required")
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
	if err != nil {
		log.Fatal(fmt.Errorf("initializing database client: %w", err))
	}

	filter := types.DefaultQueryFilter()
	filter.Limit = pointers.Pointer(uint8(50))

	var (
		im               search.IndexManager
		thresholdMet     bool
		indexRequestChan = make(chan *indexing.IndexRequest)
		wipeOnce         sync.Once
	)

	go func() {
		for x := range indexRequestChan {
			wipeOnce.Do(func() {
				if *wipePtr {
					log.Println("wiping index")
					if err = im.Wipe(ctx); err != nil {
						log.Fatal(fmt.Errorf("wiping index: %w", err))
					}
					log.Println("wiped index")
				}
			})

			if err = indexing.HandleIndexRequest(ctx, logger, tracerProvider, cfg, dataManager, x); err != nil {
				observability.AcknowledgeError(err, logger, nil, "indexing row")
			}
		}
	}()

	switch index {
	case indexing.IndexTypeRecipes:
		im, err = searchcfg.ProvideIndexManager[search.RecipeSearchSubset](ctx, logger, tracerProvider, cfg, index)
		if err != nil {
			log.Fatal(fmt.Errorf("initializing index manager: %w", err))
		}

		for !thresholdMet {
			var data *types.QueryFilteredResult[types.Recipe]
			data, err = dataManager.GetRecipes(ctx, filter)
			if err != nil {
				log.Fatal(fmt.Errorf("getting Recipe data: %w", err))
			}

			for _, x := range data.Data {
				indexRequestChan <- &indexing.IndexRequest{
					RowID:     x.ID,
					IndexType: indexing.IndexTypeRecipes,
				}
			}

			thresholdMet = len(data.Data) == 0
			*filter.Page++
		}
		break
	case indexing.IndexTypeMeals:
		im, err = searchcfg.ProvideIndexManager[search.MealSearchSubset](ctx, logger, tracerProvider, cfg, index)
		if err != nil {
			log.Fatal(fmt.Errorf("initializing index manager: %w", err))
		}

		for !thresholdMet {
			var data *types.QueryFilteredResult[types.Meal]
			data, err = dataManager.GetMeals(ctx, filter)
			if err != nil {
				log.Fatal(fmt.Errorf("getting Meal data: %w", err))
			}

			for _, x := range data.Data {
				indexRequestChan <- &indexing.IndexRequest{
					RowID:     x.ID,
					IndexType: indexing.IndexTypeMeals,
				}
			}

			thresholdMet = len(data.Data) == 0
			*filter.Page++
		}
		break
	case indexing.IndexTypeValidIngredients:
		im, err = searchcfg.ProvideIndexManager[types.ValidIngredient](ctx, logger, tracerProvider, cfg, index)
		if err != nil {
			log.Fatal(fmt.Errorf("initializing index manager: %w", err))
		}

		for !thresholdMet {
			var data *types.QueryFilteredResult[types.ValidIngredient]
			data, err = dataManager.GetValidIngredients(ctx, filter)
			if err != nil {
				log.Fatal(fmt.Errorf("getting ValidIngredient data: %w", err))
			}

			for _, x := range data.Data {
				indexRequestChan <- &indexing.IndexRequest{
					RowID:     x.ID,
					IndexType: indexing.IndexTypeValidIngredients,
				}
			}

			thresholdMet = len(data.Data) == 0
			*filter.Page++
		}
		break
	case indexing.IndexTypeValidInstruments:
		im, err = searchcfg.ProvideIndexManager[types.ValidInstrument](ctx, logger, tracerProvider, cfg, index)
		if err != nil {
			log.Fatal(fmt.Errorf("initializing index manager: %w", err))
		}

		for !thresholdMet {
			var data *types.QueryFilteredResult[types.ValidInstrument]
			data, err = dataManager.GetValidInstruments(ctx, filter)
			if err != nil {
				log.Fatal(fmt.Errorf("getting ValidInstrument data: %w", err))
			}

			for _, x := range data.Data {
				indexRequestChan <- &indexing.IndexRequest{
					RowID:     x.ID,
					IndexType: indexing.IndexTypeValidInstruments,
				}
			}

			thresholdMet = len(data.Data) == 0
			*filter.Page++
		}
		break
	case indexing.IndexTypeValidMeasurementUnits:
		im, err = searchcfg.ProvideIndexManager[types.ValidMeasurementUnit](ctx, logger, tracerProvider, cfg, index)
		if err != nil {
			log.Fatal(fmt.Errorf("initializing index manager: %w", err))
		}

		for !thresholdMet {
			var data *types.QueryFilteredResult[types.ValidMeasurementUnit]
			data, err = dataManager.GetValidMeasurementUnits(ctx, filter)
			if err != nil {
				log.Fatal(fmt.Errorf("getting ValidMeasurementUnit data: %w", err))
			}

			for _, x := range data.Data {
				indexRequestChan <- &indexing.IndexRequest{
					RowID:     x.ID,
					IndexType: indexing.IndexTypeValidMeasurementUnits,
				}
			}

			thresholdMet = len(data.Data) == 0
			*filter.Page++
		}
		break
	case indexing.IndexTypeValidPreparations:
		im, err = searchcfg.ProvideIndexManager[types.ValidPreparation](ctx, logger, tracerProvider, cfg, index)
		if err != nil {
			log.Fatal(fmt.Errorf("initializing index manager: %w", err))
		}

		for !thresholdMet {
			var data *types.QueryFilteredResult[types.ValidPreparation]
			data, err = dataManager.GetValidPreparations(ctx, filter)
			if err != nil {
				log.Fatal(fmt.Errorf("getting ValidPreparation data: %w", err))
			}

			for _, x := range data.Data {
				indexRequestChan <- &indexing.IndexRequest{
					RowID:     x.ID,
					IndexType: indexing.IndexTypeValidPreparations,
				}
			}

			thresholdMet = len(data.Data) == 0
			*filter.Page++
		}
		break
	case indexing.IndexTypeValidIngredientStates:
		im, err = searchcfg.ProvideIndexManager[types.ValidIngredientState](ctx, logger, tracerProvider, cfg, index)
		if err != nil {
			log.Fatal(fmt.Errorf("initializing index manager: %w", err))
		}

		for !thresholdMet {
			var data *types.QueryFilteredResult[types.ValidIngredientState]
			data, err = dataManager.GetValidIngredientStates(ctx, filter)
			if err != nil {
				log.Fatal(fmt.Errorf("getting ValidIngredientState data: %w", err))
			}

			for _, x := range data.Data {
				indexRequestChan <- &indexing.IndexRequest{
					RowID:     x.ID,
					IndexType: indexing.IndexTypeValidIngredientStates,
				}
			}

			thresholdMet = len(data.Data) == 0
			*filter.Page++
		}
		break
	case indexing.IndexTypeValidIngredientMeasurementUnits:
		im, err = searchcfg.ProvideIndexManager[types.ValidIngredientMeasurementUnit](ctx, logger, tracerProvider, cfg, index)
		if err != nil {
			log.Fatal(fmt.Errorf("initializing index manager: %w", err))
		}

		for !thresholdMet {
			var data *types.QueryFilteredResult[types.ValidIngredientMeasurementUnit]
			data, err = dataManager.GetValidIngredientMeasurementUnits(ctx, filter)
			if err != nil {
				log.Fatal(fmt.Errorf("getting ValidIngredientMeasurementUnit data: %w", err))
			}

			for _, x := range data.Data {
				indexRequestChan <- &indexing.IndexRequest{
					RowID:     x.ID,
					IndexType: indexing.IndexTypeValidIngredientMeasurementUnits,
				}
			}

			thresholdMet = len(data.Data) == 0
			*filter.Page++
		}
		break
	//case indexing.IndexTypeValidMeasurementUnitConversions:
	//	im, err = searchcfg.ProvideIndexManager[types.ValidMeasurementUnitConversion](ctx, logger, tracerProvider, cfg, index)
	//	if err != nil {
	//		log.Fatal(fmt.Errorf("initializing index manager: %w", err))
	//	}
	//
	//	for !thresholdMet {
	//		var data *types.QueryFilteredResult[types.ValidMeasurementUnitConversion]
	//		data, err = dataManager.GetValidMeasur(ctx, filter)
	//		if err != nil {
	//			log.Fatal(fmt.Errorf("getting ValidMeasurementUnitConversion data: %w", err))
	//		}
	//
	//		for _, x := range data.Data {
	//			indexRequestChan <- &indexing.IndexRequest{
	//				RowID:     x.ID,
	//				IndexType: indexing.IndexTypeValidMeasurementUnitConversions,
	//			}
	//		}
	//
	//		thresholdMet = len(data.Data) == 0
	//		*filter.Page++
	//	}
	//	break
	case indexing.IndexTypeValidPreparationInstruments:
		im, err = searchcfg.ProvideIndexManager[types.ValidPreparationInstrument](ctx, logger, tracerProvider, cfg, index)
		if err != nil {
			log.Fatal(fmt.Errorf("initializing index manager: %w", err))
		}

		for !thresholdMet {
			var data *types.QueryFilteredResult[types.ValidPreparationInstrument]
			data, err = dataManager.GetValidPreparationInstruments(ctx, filter)
			if err != nil {
				log.Fatal(fmt.Errorf("getting ValidPreparationInstrument data: %w", err))
			}

			for _, x := range data.Data {
				indexRequestChan <- &indexing.IndexRequest{
					RowID:     x.ID,
					IndexType: indexing.IndexTypeValidPreparationInstruments,
				}
			}

			thresholdMet = len(data.Data) == 0
			*filter.Page++
		}
		break
	case indexing.IndexTypeValidIngredientPreparations:
		im, err = searchcfg.ProvideIndexManager[types.ValidIngredientPreparation](ctx, logger, tracerProvider, cfg, index)
		if err != nil {
			log.Fatal(fmt.Errorf("initializing index manager: %w", err))
		}

		for !thresholdMet {
			var data *types.QueryFilteredResult[types.ValidIngredientPreparation]
			data, err = dataManager.GetValidIngredientPreparations(ctx, filter)
			if err != nil {
				log.Fatal(fmt.Errorf("getting ValidIngredientPreparation data: %w", err))
			}

			for _, x := range data.Data {
				indexRequestChan <- &indexing.IndexRequest{
					RowID:     x.ID,
					IndexType: indexing.IndexTypeValidIngredientPreparations,
				}
			}

			thresholdMet = len(data.Data) == 0
			*filter.Page++
		}
		break
	}
}
