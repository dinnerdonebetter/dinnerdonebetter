package indexing

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	textsearch "github.com/dinnerdonebetter/backend/internal/lib/search/text"
	textsearchcfg "github.com/dinnerdonebetter/backend/internal/lib/search/text/config"
	"github.com/dinnerdonebetter/backend/internal/lib/testutils"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
)

func TestHandleIndexRequest(T *testing.T) {
	T.Parallel()

	T.Run("recipe index type", func(t *testing.T) {
		t.Parallel()

		exampleRecipe := fakes.BuildFakeRecipe()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		searchConfig := &textsearchcfg.Config{}

		dataManager := database.NewMockDatabase()
		dataManager.RecipeDataManagerMock.On("GetRecipe", testutils.ContextMatcher, exampleRecipe.ID).Return(exampleRecipe, nil)
		dataManager.RecipeDataManagerMock.On("MarkRecipeAsIndexed", testutils.ContextMatcher, exampleRecipe.ID).Return(nil)

		indexReq := &textsearch.IndexRequest{
			RowID:     exampleRecipe.ID,
			IndexType: IndexTypeRecipes,
			Delete:    false,
		}

		assert.NoError(t, HandleIndexRequest(ctx, logger, tracing.NewNoopTracerProvider(), metrics.NewNoopMetricsProvider(), searchConfig, dataManager, indexReq))
	})

	T.Run("meal index type", func(t *testing.T) {
		t.Parallel()

		exampleMeal := fakes.BuildFakeMeal()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		searchConfig := &textsearchcfg.Config{}

		dataManager := database.NewMockDatabase()
		dataManager.MealDataManagerMock.On("GetMeal", testutils.ContextMatcher, exampleMeal.ID).Return(exampleMeal, nil)
		dataManager.MealDataManagerMock.On("MarkMealAsIndexed", testutils.ContextMatcher, exampleMeal.ID).Return(nil)

		indexReq := &textsearch.IndexRequest{
			RowID:     exampleMeal.ID,
			IndexType: IndexTypeMeals,
			Delete:    false,
		}

		assert.NoError(t, HandleIndexRequest(ctx, logger, tracing.NewNoopTracerProvider(), metrics.NewNoopMetricsProvider(), searchConfig, dataManager, indexReq))
	})

	T.Run("valid vessel index type", func(t *testing.T) {
		t.Parallel()

		exampleValidVessel := fakes.BuildFakeValidVessel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		searchConfig := &textsearchcfg.Config{}

		dataManager := database.NewMockDatabase()
		dataManager.ValidVesselDataManagerMock.On("GetValidVessel", testutils.ContextMatcher, exampleValidVessel.ID).Return(exampleValidVessel, nil)
		dataManager.ValidVesselDataManagerMock.On("MarkValidVesselAsIndexed", testutils.ContextMatcher, exampleValidVessel.ID).Return(nil)

		indexReq := &textsearch.IndexRequest{
			RowID:     exampleValidVessel.ID,
			IndexType: IndexTypeValidVessels,
			Delete:    false,
		}

		assert.NoError(t, HandleIndexRequest(ctx, logger, tracing.NewNoopTracerProvider(), metrics.NewNoopMetricsProvider(), searchConfig, dataManager, indexReq))
	})

	T.Run("valid ingredient index type", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredient := fakes.BuildFakeValidIngredient()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		searchConfig := &textsearchcfg.Config{}

		dataManager := database.NewMockDatabase()
		dataManager.ValidIngredientDataManagerMock.On("GetValidIngredient", testutils.ContextMatcher, exampleValidIngredient.ID).Return(exampleValidIngredient, nil)
		dataManager.ValidIngredientDataManagerMock.On("MarkValidIngredientAsIndexed", testutils.ContextMatcher, exampleValidIngredient.ID).Return(nil)

		indexReq := &textsearch.IndexRequest{
			RowID:     exampleValidIngredient.ID,
			IndexType: IndexTypeValidIngredients,
			Delete:    false,
		}

		assert.NoError(t, HandleIndexRequest(ctx, logger, tracing.NewNoopTracerProvider(), metrics.NewNoopMetricsProvider(), searchConfig, dataManager, indexReq))
	})

	T.Run("valid instrument index type", func(t *testing.T) {
		t.Parallel()

		exampleValidInstrument := fakes.BuildFakeValidInstrument()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		searchConfig := &textsearchcfg.Config{}

		dataManager := database.NewMockDatabase()
		dataManager.ValidInstrumentDataManagerMock.On("GetValidInstrument", testutils.ContextMatcher, exampleValidInstrument.ID).Return(exampleValidInstrument, nil)
		dataManager.ValidInstrumentDataManagerMock.On("MarkValidInstrumentAsIndexed", testutils.ContextMatcher, exampleValidInstrument.ID).Return(nil)

		indexReq := &textsearch.IndexRequest{
			RowID:     exampleValidInstrument.ID,
			IndexType: IndexTypeValidInstruments,
			Delete:    false,
		}

		assert.NoError(t, HandleIndexRequest(ctx, logger, tracing.NewNoopTracerProvider(), metrics.NewNoopMetricsProvider(), searchConfig, dataManager, indexReq))
	})

	T.Run("valid preparation index type", func(t *testing.T) {
		t.Parallel()

		exampleValidPreparation := fakes.BuildFakeValidPreparation()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		searchConfig := &textsearchcfg.Config{}

		dataManager := database.NewMockDatabase()
		dataManager.ValidPreparationDataManagerMock.On("GetValidPreparation", testutils.ContextMatcher, exampleValidPreparation.ID).Return(exampleValidPreparation, nil)
		dataManager.ValidPreparationDataManagerMock.On("MarkValidPreparationAsIndexed", testutils.ContextMatcher, exampleValidPreparation.ID).Return(nil)

		indexReq := &textsearch.IndexRequest{
			RowID:     exampleValidPreparation.ID,
			IndexType: IndexTypeValidPreparations,
			Delete:    false,
		}

		assert.NoError(t, HandleIndexRequest(ctx, logger, tracing.NewNoopTracerProvider(), metrics.NewNoopMetricsProvider(), searchConfig, dataManager, indexReq))
	})

	T.Run("valid measurement unit index type", func(t *testing.T) {
		t.Parallel()

		exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		searchConfig := &textsearchcfg.Config{}

		dataManager := database.NewMockDatabase()
		dataManager.ValidMeasurementUnitDataManagerMock.On("GetValidMeasurementUnit", testutils.ContextMatcher, exampleValidMeasurementUnit.ID).Return(exampleValidMeasurementUnit, nil)
		dataManager.ValidMeasurementUnitDataManagerMock.On("MarkValidMeasurementUnitAsIndexed", testutils.ContextMatcher, exampleValidMeasurementUnit.ID).Return(nil)

		indexReq := &textsearch.IndexRequest{
			RowID:     exampleValidMeasurementUnit.ID,
			IndexType: IndexTypeValidMeasurementUnits,
			Delete:    false,
		}

		assert.NoError(t, HandleIndexRequest(ctx, logger, tracing.NewNoopTracerProvider(), metrics.NewNoopMetricsProvider(), searchConfig, dataManager, indexReq))
	})

	T.Run("valid ingredient state index type", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredientState := fakes.BuildFakeValidIngredientState()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		searchConfig := &textsearchcfg.Config{}

		dataManager := database.NewMockDatabase()
		dataManager.ValidIngredientStateDataManagerMock.On("GetValidIngredientState", testutils.ContextMatcher, exampleValidIngredientState.ID).Return(exampleValidIngredientState, nil)
		dataManager.ValidIngredientStateDataManagerMock.On("MarkValidIngredientStateAsIndexed", testutils.ContextMatcher, exampleValidIngredientState.ID).Return(nil)

		indexReq := &textsearch.IndexRequest{
			RowID:     exampleValidIngredientState.ID,
			IndexType: IndexTypeValidIngredientStates,
			Delete:    false,
		}

		assert.NoError(t, HandleIndexRequest(ctx, logger, tracing.NewNoopTracerProvider(), metrics.NewNoopMetricsProvider(), searchConfig, dataManager, indexReq))
	})
}
