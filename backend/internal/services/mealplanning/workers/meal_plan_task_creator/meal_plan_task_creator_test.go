package mealplantaskcreator

import (
	"context"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/recipeanalysis"
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/internal/domain/audit"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/fakes"
	mealplanningmock "github.com/dinnerdonebetter/backend/internal/domain/mealplanning/mocks"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/testutils"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
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
		&mealplanningmock.Repository{},
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

		mdm := &mealplanningmock.Repository{}
		mdm.On("GetFinalizedMealPlanIDsForTheNextWeek", testutils.ContextMatcher).Return([]*mealplanning.FinalizedMealPlanDatabaseResult{}, nil)
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
		exampleRecipe := &mealplanning.Recipe{
			Name: "Recipe 1",
			ID:   exampleRecipeID,
			Steps: []*mealplanning.RecipeStep{
				{
					BelongsToRecipe: exampleRecipeID,
					ID:              recipeStepID,
					Preparation:     mealplanning.ValidPreparation{Name: "dice"},
					Ingredients: []*mealplanning.RecipeStepIngredient{
						{
							Ingredient: &mealplanning.ValidIngredient{
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
							MeasurementUnit:     mealplanning.ValidMeasurementUnit{Name: "gram", PluralName: "grams"},
							Quantity: types.Float32RangeWithOptionalMax{
								Max: pointer.To(float32(900)),
								Min: 900,
							},
						},
					},
					Products: []*mealplanning.RecipeStepProduct{
						{
							Name:                "diced chicken breast",
							Type:                mealplanning.RecipeStepProductIngredientType,
							BelongsToRecipeStep: recipeStepID,
							ID:                  fakes.BuildFakeID(),
							MeasurementUnit:     &mealplanning.ValidMeasurementUnit{},
						},
					},
				},
			},
		}

		recipeMap := map[string]*mealplanning.Recipe{
			exampleRecipe.ID: exampleRecipe,
		}

		exampleFinalizedMealPlanResult := &mealplanning.FinalizedMealPlanDatabaseResult{
			MealPlanID:       exampleMealPlan.ID,
			MealPlanEventID:  exampleMealPlanEvent.ID,
			MealPlanOptionID: exampleMealPlanOption.ID,
			MealID:           exampleMeal.ID,
			RecipeIDs: []string{
				exampleRecipe.ID,
			},
		}

		exampleFinalizedMealPlanResults := []*mealplanning.FinalizedMealPlanDatabaseResult{exampleFinalizedMealPlanResult}

		createdMealPlanTasks := fakes.BuildFakeMealPlanTasksList().Data

		mdm := &mealplanningmock.Repository{}
		mdm.On("CreateMealPlanTasksForMealPlanOption", testutils.ContextMatcher, testutils.MatchType[[]*mealplanning.MealPlanTaskDatabaseCreationInput]()).Return(createdMealPlanTasks, nil)
		mdm.On("GetFinalizedMealPlanIDsForTheNextWeek", testutils.ContextMatcher).Return(exampleFinalizedMealPlanResults, nil)
		mdm.On("MarkMealPlanAsHavingTasksCreated", testutils.ContextMatcher, testutils.MatchType[string]()).Return(nil)

		expectedReturnResults := []*mealplanning.MealPlanTaskDatabaseCreationInput{
			{
				CreationExplanation: t.Name(),
				MealPlanOptionID:    exampleFinalizedMealPlanResult.MealPlanOptionID,
			},
		}

		mockAnalyzer := &recipeanalysis.MockRecipeAnalyzer{}
		for _, result := range exampleFinalizedMealPlanResults {
			for _, recipeID := range result.RecipeIDs {
				mdm.On("GetRecipe", testutils.ContextMatcher, recipeID).Return(recipeMap[recipeID], nil)

				mockAnalyzer.On(
					"GenerateMealPlanTasksForRecipe",
					testutils.ContextMatcher,
					result.MealPlanOptionID,
					recipeMap[recipeID],
				).Return(expectedReturnResults, nil)
			}
		}
		w.analyzer = mockAnalyzer
		w.dataManager = mdm

		mmp := &mockpublishers.Publisher{}
		mmp.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.MatchType[*audit.DataChangeMessage](),
		).Return(nil)
		w.postUpdatesPublisher = mmp

		assert.NoError(t, w.Work(ctx))

		mock.AssertExpectationsForObjects(t, mdm, mockAnalyzer, mmp)
	})
}
