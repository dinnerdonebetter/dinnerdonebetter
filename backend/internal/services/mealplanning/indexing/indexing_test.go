package indexing

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/fakes"
	mealplanningmock "github.com/dinnerdonebetter/backend/internal/domain/mealplanning/mocks"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/reflection"
	textsearch "github.com/dinnerdonebetter/backend/internal/platform/search/text"
	mocksearch "github.com/dinnerdonebetter/backend/internal/platform/search/text/mock"
	"github.com/dinnerdonebetter/backend/internal/platform/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleIndexRequest(T *testing.T) {
	T.Parallel()

	T.Run("recipe index type", func(t *testing.T) {
		t.Parallel()

		exampleRecipe := fakes.BuildFakeRecipe()

		ctx := t.Context()
		logger := logging.NewNoopLogger()

		mealPlanningRepo := &mealplanningmock.Repository{}
		mealPlanningRepo.On(reflection.GetMethodName(mealPlanningRepo.GetRecipe), testutils.ContextMatcher, exampleRecipe.ID).Return(exampleRecipe, nil)
		mealPlanningRepo.On(reflection.GetMethodName(mealPlanningRepo.MarkRecipeAsIndexed), testutils.ContextMatcher, exampleRecipe.ID).Return(nil)

		rim := &mocksearch.IndexManager[RecipeSearchSubset]{}
		ss := ConvertRecipeToRecipeSearchSubset(exampleRecipe)
		rim.On(reflection.GetMethodName(rim.Index), testutils.ContextMatcher, exampleRecipe.ID, ss).Return(nil)

		mim := &mocksearch.IndexManager[MealSearchSubset]{}
		vinm := &mocksearch.IndexManager[ValidIngredientSearchSubset]{}
		vism := &mocksearch.IndexManager[ValidInstrumentSearchSubset]{}
		vmuim := &mocksearch.IndexManager[ValidMeasurementUnitSearchSubset]{}
		vpim := &mocksearch.IndexManager[ValidPreparationSearchSubset]{}
		visim := &mocksearch.IndexManager[ValidIngredientStateSearchSubset]{}
		vvim := &mocksearch.IndexManager[ValidVesselSearchSubset]{}

		cdi := NewMealPlanningDataIndexer(
			logger,
			tracing.NewNoopTracerProvider(),
			mealPlanningRepo,
			rim,
			mim,
			vinm,
			vism,
			vmuim,
			vpim,
			visim,
			vvim,
		)

		indexReq := &textsearch.IndexRequest{
			RowID:     exampleRecipe.ID,
			IndexType: IndexTypeRecipes,
			Delete:    false,
		}

		assert.NoError(t, cdi.HandleIndexRequest(ctx, indexReq))

		mock.AssertExpectationsForObjects(t,
			rim,
			mim,
			vinm,
			vism,
			vmuim,
			vpim,
			visim,
			vvim,
		)
	})

	T.Run("meal index type", func(t *testing.T) {
		t.Parallel()

		exampleMeal := fakes.BuildFakeMeal()

		ctx := t.Context()
		logger := logging.NewNoopLogger()

		mealPlanningRepo := &mealplanningmock.Repository{}
		mealPlanningRepo.On(reflection.GetMethodName(mealPlanningRepo.GetMeal), testutils.ContextMatcher, exampleMeal.ID).Return(exampleMeal, nil)
		mealPlanningRepo.On(reflection.GetMethodName(mealPlanningRepo.MarkMealAsIndexed), testutils.ContextMatcher, exampleMeal.ID).Return(nil)

		rim := &mocksearch.IndexManager[RecipeSearchSubset]{}

		mim := &mocksearch.IndexManager[MealSearchSubset]{}
		ss := ConvertMealToMealSearchSubset(exampleMeal)
		mim.On(reflection.GetMethodName(mim.Index), testutils.ContextMatcher, exampleMeal.ID, ss).Return(nil)

		vinm := &mocksearch.IndexManager[ValidIngredientSearchSubset]{}
		vism := &mocksearch.IndexManager[ValidInstrumentSearchSubset]{}
		vmuim := &mocksearch.IndexManager[ValidMeasurementUnitSearchSubset]{}
		vpim := &mocksearch.IndexManager[ValidPreparationSearchSubset]{}
		visim := &mocksearch.IndexManager[ValidIngredientStateSearchSubset]{}
		vvim := &mocksearch.IndexManager[ValidVesselSearchSubset]{}

		cdi := NewMealPlanningDataIndexer(
			logger,
			tracing.NewNoopTracerProvider(),
			mealPlanningRepo,
			rim,
			mim,
			vinm,
			vism,
			vmuim,
			vpim,
			visim,
			vvim,
		)

		indexReq := &textsearch.IndexRequest{
			RowID:     exampleMeal.ID,
			IndexType: IndexTypeMeals,
			Delete:    false,
		}

		assert.NoError(t, cdi.HandleIndexRequest(ctx, indexReq))

		mock.AssertExpectationsForObjects(t,
			rim,
			mim,
			vinm,
			vism,
			vmuim,
			vpim,
			visim,
			vvim,
		)
	})

	T.Run("valid vessel index type", func(t *testing.T) {
		t.Parallel()

		exampleValidVessel := fakes.BuildFakeValidVessel()

		ctx := t.Context()
		logger := logging.NewNoopLogger()

		mealPlanningRepo := &mealplanningmock.Repository{}
		mealPlanningRepo.On(reflection.GetMethodName(mealPlanningRepo.GetValidVessel), testutils.ContextMatcher, exampleValidVessel.ID).Return(exampleValidVessel, nil)
		mealPlanningRepo.On(reflection.GetMethodName(mealPlanningRepo.MarkValidVesselAsIndexed), testutils.ContextMatcher, exampleValidVessel.ID).Return(nil)

		rim := &mocksearch.IndexManager[RecipeSearchSubset]{}
		mim := &mocksearch.IndexManager[MealSearchSubset]{}
		vinm := &mocksearch.IndexManager[ValidIngredientSearchSubset]{}
		vism := &mocksearch.IndexManager[ValidInstrumentSearchSubset]{}
		vmuim := &mocksearch.IndexManager[ValidMeasurementUnitSearchSubset]{}
		vpim := &mocksearch.IndexManager[ValidPreparationSearchSubset]{}
		visim := &mocksearch.IndexManager[ValidIngredientStateSearchSubset]{}

		vvim := &mocksearch.IndexManager[ValidVesselSearchSubset]{}
		ss := ConvertValidVesselToValidVesselSearchSubset(exampleValidVessel)
		vvim.On(reflection.GetMethodName(vvim.Index), testutils.ContextMatcher, exampleValidVessel.ID, ss).Return(nil)

		cdi := NewMealPlanningDataIndexer(
			logger,
			tracing.NewNoopTracerProvider(),
			mealPlanningRepo,
			rim,
			mim,
			vinm,
			vism,
			vmuim,
			vpim,
			visim,
			vvim,
		)

		indexReq := &textsearch.IndexRequest{
			RowID:     exampleValidVessel.ID,
			IndexType: IndexTypeValidVessels,
			Delete:    false,
		}

		assert.NoError(t, cdi.HandleIndexRequest(ctx, indexReq))

		mock.AssertExpectationsForObjects(t,
			rim,
			mim,
			vinm,
			vism,
			vmuim,
			vpim,
			visim,
			vvim,
		)
	})

	T.Run("valid ingredient index type", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredient := fakes.BuildFakeValidIngredient()

		ctx := t.Context()
		logger := logging.NewNoopLogger()

		mealPlanningRepo := &mealplanningmock.Repository{}
		mealPlanningRepo.On(reflection.GetMethodName(mealPlanningRepo.GetValidIngredient), testutils.ContextMatcher, exampleValidIngredient.ID).Return(exampleValidIngredient, nil)
		mealPlanningRepo.On(reflection.GetMethodName(mealPlanningRepo.MarkValidIngredientAsIndexed), testutils.ContextMatcher, exampleValidIngredient.ID).Return(nil)

		rim := &mocksearch.IndexManager[RecipeSearchSubset]{}
		mim := &mocksearch.IndexManager[MealSearchSubset]{}
		vism := &mocksearch.IndexManager[ValidInstrumentSearchSubset]{}
		vmuim := &mocksearch.IndexManager[ValidMeasurementUnitSearchSubset]{}
		vpim := &mocksearch.IndexManager[ValidPreparationSearchSubset]{}
		visim := &mocksearch.IndexManager[ValidIngredientStateSearchSubset]{}
		vvim := &mocksearch.IndexManager[ValidVesselSearchSubset]{}

		vinm := &mocksearch.IndexManager[ValidIngredientSearchSubset]{}
		ss := ConvertValidIngredientToValidIngredientSearchSubset(exampleValidIngredient)
		vinm.On(reflection.GetMethodName(vinm.Index), testutils.ContextMatcher, exampleValidIngredient.ID, ss).Return(nil)

		cdi := NewMealPlanningDataIndexer(
			logger,
			tracing.NewNoopTracerProvider(),
			mealPlanningRepo,
			rim,
			mim,
			vinm,
			vism,
			vmuim,
			vpim,
			visim,
			vvim,
		)

		indexReq := &textsearch.IndexRequest{
			RowID:     exampleValidIngredient.ID,
			IndexType: IndexTypeValidIngredients,
			Delete:    false,
		}

		assert.NoError(t, cdi.HandleIndexRequest(ctx, indexReq))

		mock.AssertExpectationsForObjects(t,
			rim,
			mim,
			vinm,
			vism,
			vmuim,
			vpim,
			visim,
			vvim,
		)
	})

	T.Run("valid instrument index type", func(t *testing.T) {
		t.Parallel()

		exampleValidInstrument := fakes.BuildFakeValidInstrument()

		ctx := t.Context()
		logger := logging.NewNoopLogger()

		mealPlanningRepo := &mealplanningmock.Repository{}
		mealPlanningRepo.On(reflection.GetMethodName(mealPlanningRepo.GetValidInstrument), testutils.ContextMatcher, exampleValidInstrument.ID).Return(exampleValidInstrument, nil)
		mealPlanningRepo.On(reflection.GetMethodName(mealPlanningRepo.MarkValidInstrumentAsIndexed), testutils.ContextMatcher, exampleValidInstrument.ID).Return(nil)

		rim := &mocksearch.IndexManager[RecipeSearchSubset]{}
		mim := &mocksearch.IndexManager[MealSearchSubset]{}
		vinm := &mocksearch.IndexManager[ValidIngredientSearchSubset]{}

		vism := &mocksearch.IndexManager[ValidInstrumentSearchSubset]{}
		ss := ConvertValidInstrumentToValidInstrumentSearchSubset(exampleValidInstrument)
		vism.On(reflection.GetMethodName(vism.Index), testutils.ContextMatcher, exampleValidInstrument.ID, ss).Return(nil)

		vmuim := &mocksearch.IndexManager[ValidMeasurementUnitSearchSubset]{}
		vpim := &mocksearch.IndexManager[ValidPreparationSearchSubset]{}
		visim := &mocksearch.IndexManager[ValidIngredientStateSearchSubset]{}
		vvim := &mocksearch.IndexManager[ValidVesselSearchSubset]{}

		cdi := NewMealPlanningDataIndexer(
			logger,
			tracing.NewNoopTracerProvider(),
			mealPlanningRepo,
			rim,
			mim,
			vinm,
			vism,
			vmuim,
			vpim,
			visim,
			vvim,
		)

		indexReq := &textsearch.IndexRequest{
			RowID:     exampleValidInstrument.ID,
			IndexType: IndexTypeValidInstruments,
			Delete:    false,
		}

		assert.NoError(t, cdi.HandleIndexRequest(ctx, indexReq))

		mock.AssertExpectationsForObjects(t,
			rim,
			mim,
			vinm,
			vism,
			vmuim,
			vpim,
			visim,
			vvim,
		)
	})

	T.Run("valid preparation index type", func(t *testing.T) {
		t.Parallel()

		exampleValidPreparation := fakes.BuildFakeValidPreparation()

		ctx := t.Context()
		logger := logging.NewNoopLogger()

		mealPlanningRepo := &mealplanningmock.Repository{}
		mealPlanningRepo.On(reflection.GetMethodName(mealPlanningRepo.GetValidPreparation), testutils.ContextMatcher, exampleValidPreparation.ID).Return(exampleValidPreparation, nil)
		mealPlanningRepo.On(reflection.GetMethodName(mealPlanningRepo.MarkValidPreparationAsIndexed), testutils.ContextMatcher, exampleValidPreparation.ID).Return(nil)

		rim := &mocksearch.IndexManager[RecipeSearchSubset]{}
		mim := &mocksearch.IndexManager[MealSearchSubset]{}
		vinm := &mocksearch.IndexManager[ValidIngredientSearchSubset]{}
		vism := &mocksearch.IndexManager[ValidInstrumentSearchSubset]{}
		vmuim := &mocksearch.IndexManager[ValidMeasurementUnitSearchSubset]{}

		vpim := &mocksearch.IndexManager[ValidPreparationSearchSubset]{}
		ss := ConvertValidPreparationToValidPreparationSearchSubset(exampleValidPreparation)
		vpim.On(reflection.GetMethodName(vpim.Index), testutils.ContextMatcher, exampleValidPreparation.ID, ss).Return(nil)

		visim := &mocksearch.IndexManager[ValidIngredientStateSearchSubset]{}
		vvim := &mocksearch.IndexManager[ValidVesselSearchSubset]{}

		cdi := NewMealPlanningDataIndexer(
			logger,
			tracing.NewNoopTracerProvider(),
			mealPlanningRepo,
			rim,
			mim,
			vinm,
			vism,
			vmuim,
			vpim,
			visim,
			vvim,
		)

		indexReq := &textsearch.IndexRequest{
			RowID:     exampleValidPreparation.ID,
			IndexType: IndexTypeValidPreparations,
			Delete:    false,
		}

		assert.NoError(t, cdi.HandleIndexRequest(ctx, indexReq))

		mock.AssertExpectationsForObjects(t,
			rim,
			mim,
			vinm,
			vism,
			vmuim,
			vpim,
			visim,
			vvim,
		)
	})

	T.Run("valid measurement unit index type", func(t *testing.T) {
		t.Parallel()

		exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()

		ctx := t.Context()
		logger := logging.NewNoopLogger()

		mealPlanningRepo := &mealplanningmock.Repository{}
		mealPlanningRepo.On(reflection.GetMethodName(mealPlanningRepo.GetValidMeasurementUnit), testutils.ContextMatcher, exampleValidMeasurementUnit.ID).Return(exampleValidMeasurementUnit, nil)
		mealPlanningRepo.On(reflection.GetMethodName(mealPlanningRepo.MarkValidMeasurementUnitAsIndexed), testutils.ContextMatcher, exampleValidMeasurementUnit.ID).Return(nil)

		rim := &mocksearch.IndexManager[RecipeSearchSubset]{}
		mim := &mocksearch.IndexManager[MealSearchSubset]{}
		vinm := &mocksearch.IndexManager[ValidIngredientSearchSubset]{}
		vism := &mocksearch.IndexManager[ValidInstrumentSearchSubset]{}

		vmuim := &mocksearch.IndexManager[ValidMeasurementUnitSearchSubset]{}
		ss := ConvertValidMeasurementUnitToValidMeasurementUnitSearchSubset(exampleValidMeasurementUnit)
		vmuim.On(reflection.GetMethodName(vmuim.Index), testutils.ContextMatcher, exampleValidMeasurementUnit.ID, ss).Return(nil)

		vpim := &mocksearch.IndexManager[ValidPreparationSearchSubset]{}
		visim := &mocksearch.IndexManager[ValidIngredientStateSearchSubset]{}
		vvim := &mocksearch.IndexManager[ValidVesselSearchSubset]{}

		cdi := NewMealPlanningDataIndexer(
			logger,
			tracing.NewNoopTracerProvider(),
			mealPlanningRepo,
			rim,
			mim,
			vinm,
			vism,
			vmuim,
			vpim,
			visim,
			vvim,
		)

		indexReq := &textsearch.IndexRequest{
			RowID:     exampleValidMeasurementUnit.ID,
			IndexType: IndexTypeValidMeasurementUnits,
			Delete:    false,
		}

		assert.NoError(t, cdi.HandleIndexRequest(ctx, indexReq))

		mock.AssertExpectationsForObjects(t,
			rim,
			mim,
			vinm,
			vism,
			vmuim,
			vpim,
			visim,
			vvim,
		)
	})

	T.Run("valid ingredient state index type", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredientState := fakes.BuildFakeValidIngredientState()

		ctx := t.Context()
		logger := logging.NewNoopLogger()

		mealPlanningRepo := &mealplanningmock.Repository{}
		mealPlanningRepo.On(reflection.GetMethodName(mealPlanningRepo.GetValidIngredientState), testutils.ContextMatcher, exampleValidIngredientState.ID).Return(exampleValidIngredientState, nil)
		mealPlanningRepo.On(reflection.GetMethodName(mealPlanningRepo.MarkValidIngredientStateAsIndexed), testutils.ContextMatcher, exampleValidIngredientState.ID).Return(nil)

		rim := &mocksearch.IndexManager[RecipeSearchSubset]{}
		mim := &mocksearch.IndexManager[MealSearchSubset]{}
		vinm := &mocksearch.IndexManager[ValidIngredientSearchSubset]{}
		vism := &mocksearch.IndexManager[ValidInstrumentSearchSubset]{}
		vmuim := &mocksearch.IndexManager[ValidMeasurementUnitSearchSubset]{}
		vpim := &mocksearch.IndexManager[ValidPreparationSearchSubset]{}

		visim := &mocksearch.IndexManager[ValidIngredientStateSearchSubset]{}
		ss := ConvertValidIngredientStateToValidIngredientStateSearchSubset(exampleValidIngredientState)
		visim.On(reflection.GetMethodName(visim.Index), testutils.ContextMatcher, exampleValidIngredientState.ID, ss).Return(nil)

		vvim := &mocksearch.IndexManager[ValidVesselSearchSubset]{}

		cdi := NewMealPlanningDataIndexer(
			logger,
			tracing.NewNoopTracerProvider(),
			mealPlanningRepo,
			rim,
			mim,
			vinm,
			vism,
			vmuim,
			vpim,
			visim,
			vvim,
		)

		indexReq := &textsearch.IndexRequest{
			RowID:     exampleValidIngredientState.ID,
			IndexType: IndexTypeValidIngredientStates,
			Delete:    false,
		}

		assert.NoError(t, cdi.HandleIndexRequest(ctx, indexReq))

		mock.AssertExpectationsForObjects(t,
			rim,
			mim,
			vinm,
			vism,
			vmuim,
			vpim,
			visim,
			vvim,
		)
	})
}
