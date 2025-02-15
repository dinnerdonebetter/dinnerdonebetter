package mealplantaskcreator

import (
	"context"
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/internal/database"
	msgconfig "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/config"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/lib/pointer"
	"github.com/dinnerdonebetter/backend/internal/lib/testutils"
	"github.com/dinnerdonebetter/backend/internal/services/eating/businesslogic/recipeanalysis"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func buildNewMealPlanTaskCreatorForTest(t *testing.T) *Worker {
	t.Helper()

	cfg := &msgconfig.QueuesConfig{DataChangesTopicName: "data_changes"}

	pp := &mockpublishers.PublisherProvider{}
	pp.On("ProvidePublisher", cfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

	x, err := NewMealPlanTaskCreator(
		logging.NewNoopLogger(),
		tracing.NewNoopTracerProvider(),
		&recipeanalysis.MockRecipeAnalyzer{},
		database.NewMockDatabase(),
		pp,
		metrics.NewNoopMetricsProvider(),
		cfg,
	)
	require.NoError(t, err)

	return x
}

func TestWorker_Work(T *testing.T) {
	T.Parallel()

	T.Run("with nothing to do", func(t *testing.T) {
		t.Parallel()

		w := buildNewMealPlanTaskCreatorForTest(t)
		assert.NotNil(t, w)

		ctx := context.Background()

		mdm := database.NewMockDatabase()
		mdm.MealPlanDataManagerMock.On("GetFinalizedMealPlanIDsForTheNextWeek", testutils.ContextMatcher).Return([]*types.FinalizedMealPlanDatabaseResult{}, nil)
		w.dataManager = mdm

		assert.NoError(t, w.Work(ctx))

		mock.AssertExpectationsForObjects(t, mdm)
	})

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		w := buildNewMealPlanTaskCreatorForTest(t)
		assert.NotNil(t, w)

		ctx := context.Background()

		exampleMeal := fakes.BuildFakeMeal()
		exampleMealPlan := fakes.BuildFakeMealPlan()

		exampleMealPlanEvent := fakes.BuildFakeMealPlanEvent()
		exampleMealPlanEvent.BelongsToMealPlan = exampleMealPlan.ID
		now := time.Now().Add(0).Truncate(time.Second).UTC()
		inThreeDays := now.Add((time.Hour * 24) * 3).Add(0).Truncate(time.Second).UTC()
		inOneWeek := now.Add((time.Hour * 24) * 7).Add(0).Truncate(time.Second).UTC()
		exampleMealPlanEvent.StartsAt = inThreeDays
		exampleMealPlanEvent.EndsAt = inOneWeek

		exampleMealPlanOption := fakes.BuildFakeMealPlanOption()
		exampleMealPlanOption.BelongsToMealPlanEvent = exampleMealPlanEvent.ID
		exampleMealPlanOption.Meal = *exampleMeal

		recipeStepID := fakes.BuildFakeID()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipe := &types.Recipe{
			Name: "Recipe 1",
			ID:   exampleRecipeID,
			Steps: []*types.RecipeStep{
				{
					BelongsToRecipe: exampleRecipeID,
					ID:              recipeStepID,
					Preparation:     types.ValidPreparation{Name: "dice"},
					Ingredients: []*types.RecipeStepIngredient{
						{
							Ingredient: &types.ValidIngredient{
								StorageTemperatureInCelsius: types.OptionalFloat32Range{
									Min: pointer.To(float32(2.5)),
								},
								PluralName:          "chicken breasts",
								StorageInstructions: "keep frozen",
								Name:                "chicken breast",
								ID:                  fakes.BuildFakeID(),
							},
							Name:                "chicken breast",
							ID:                  fakes.BuildFakeID(),
							BelongsToRecipeStep: recipeStepID,
							MeasurementUnit:     types.ValidMeasurementUnit{Name: "gram", PluralName: "grams"},
							Quantity: types.Float32RangeWithOptionalMax{
								Max: pointer.To(float32(900)),
								Min: 900,
							},
						},
					},
					Products: []*types.RecipeStepProduct{
						{
							Name:                "diced chicken breast",
							Type:                types.RecipeStepProductIngredientType,
							BelongsToRecipeStep: recipeStepID,
							ID:                  fakes.BuildFakeID(),
							MeasurementUnit:     &types.ValidMeasurementUnit{},
						},
					},
				},
			},
		}

		recipeMap := map[string]*types.Recipe{
			exampleRecipe.ID: exampleRecipe,
		}

		exampleFinalizedMealPlanResult := &types.FinalizedMealPlanDatabaseResult{
			MealPlanID:       exampleMealPlan.ID,
			MealPlanEventID:  exampleMealPlanEvent.ID,
			MealPlanOptionID: exampleMealPlanOption.ID,
			MealID:           exampleMeal.ID,
			RecipeIDs: []string{
				exampleRecipe.ID,
			},
		}

		exampleFinalizedMealPlanResults := []*types.FinalizedMealPlanDatabaseResult{exampleFinalizedMealPlanResult}

		createdMealPlanTasks := fakes.BuildFakeMealPlanTasksList().Data

		mdm := database.NewMockDatabase()
		mdm.MealPlanTaskDataManagerMock.On("CreateMealPlanTasksForMealPlanOption", testutils.ContextMatcher, testutils.MatchType[[]*types.MealPlanTaskDatabaseCreationInput]()).Return(createdMealPlanTasks, nil)
		mdm.MealPlanDataManagerMock.On("GetFinalizedMealPlanIDsForTheNextWeek", testutils.ContextMatcher).Return(exampleFinalizedMealPlanResults, nil)
		mdm.MealPlanTaskDataManagerMock.On("MarkMealPlanAsHavingTasksCreated", testutils.ContextMatcher, testutils.MatchType[string]()).Return(nil)

		expectedReturnResults := []*types.MealPlanTaskDatabaseCreationInput{
			{
				CreationExplanation: t.Name(),
				MealPlanOptionID:    exampleFinalizedMealPlanResult.MealPlanOptionID,
			},
		}

		mockAnalyzer := &recipeanalysis.MockRecipeAnalyzer{}
		for _, result := range exampleFinalizedMealPlanResults {
			mdm.MealPlanEventDataManagerMock.On("GetMealPlanEvent", testutils.ContextMatcher, result.MealPlanID, result.MealPlanEventID).Return(exampleMealPlanEvent, nil)

			for _, recipeID := range result.RecipeIDs {
				mdm.RecipeDataManagerMock.On("GetRecipe", testutils.ContextMatcher, recipeID).Return(recipeMap[recipeID], nil)

				mockAnalyzer.On(
					"GenerateMealPlanTasksForRecipe",
					testutils.ContextMatcher,
					result.MealPlanOptionID,
					recipeMap[recipeID],
				).Return(expectedReturnResults, nil)
			}
		}
		w.analyzer = mockAnalyzer

		mmp := &mockpublishers.Publisher{}
		mmp.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.MatchType[*types.DataChangeMessage](),
		).Return(nil)
		w.postUpdatesPublisher = mmp

		w.dataManager = mdm

		assert.NoError(t, w.Work(ctx))

		mock.AssertExpectationsForObjects(t, mdm, mockAnalyzer, mmp)
	})
}
