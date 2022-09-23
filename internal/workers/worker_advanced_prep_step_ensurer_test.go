package workers

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/internal/customerdata"
	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/graphing"
	mockpublishers "github.com/prixfixeco/api_server/internal/messagequeue/mock"
	"github.com/prixfixeco/api_server/internal/observability/logging/zerolog"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
	"github.com/prixfixeco/api_server/pkg/utils/pointers"
	testutils "github.com/prixfixeco/api_server/tests/utils"
)

func TestProvideAdvancedPrepStepCreationEnsurerWorker(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		actual := ProvideAdvancedPrepStepCreationEnsurerWorker(
			zerolog.NewZerologLogger(),
			&database.MockDatabase{},
			&graphing.MockRecipeGrapher{},
			&mockpublishers.Publisher{},
			&customerdata.MockCollector{},
			tracing.NewNoopTracerProvider(),
		)
		assert.NotNil(t, actual)
	})
}

func TestAdvancedPrepStepCreationEnsurerWorker_HandleMessage(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		w := ProvideAdvancedPrepStepCreationEnsurerWorker(
			zerolog.NewZerologLogger(),
			&database.MockDatabase{},
			&graphing.MockRecipeGrapher{},
			&mockpublishers.Publisher{},
			&customerdata.MockCollector{},
			tracing.NewNoopTracerProvider(),
		)
		assert.NotNil(t, w)

		ctx := context.Background()

		mdm := database.NewMockDatabase()
		mdm.MealPlanDataManager.On("GetFinalizedMealPlanIDsForTheNextWeek", testutils.ContextMatcher).Return([]*types.FinalizedMealPlanDatabaseResult{}, nil)
		w.dataManager = mdm

		err := w.HandleMessage(ctx, []byte("{}"))
		assert.NoError(t, err)
	})
}

func TestAdvancedPrepStepCreationEnsurerWorker_DetermineCreatableSteps(T *testing.T) {
	T.Parallel()

	T.Run("with nothing to do", func(t *testing.T) {
		t.Parallel()

		w := ProvideAdvancedPrepStepCreationEnsurerWorker(
			zerolog.NewZerologLogger(),
			&database.MockDatabase{},
			&graphing.MockRecipeGrapher{},
			&mockpublishers.Publisher{},
			&customerdata.MockCollector{},
			tracing.NewNoopTracerProvider(),
		)
		assert.NotNil(t, w)

		ctx := context.Background()
		expected := map[string][]*types.AdvancedPrepStepDatabaseCreationInput{}

		mdm := database.NewMockDatabase()
		mdm.MealPlanDataManager.On("GetFinalizedMealPlanIDsForTheNextWeek", testutils.ContextMatcher).Return([]*types.FinalizedMealPlanDatabaseResult{}, nil)
		w.dataManager = mdm

		actual, err := w.DetermineCreatableSteps(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, mdm)
	})

	T.Run("creates frozen thawing steps", func(t *testing.T) {
		t.Parallel()

		w := ProvideAdvancedPrepStepCreationEnsurerWorker(
			zerolog.NewZerologLogger(),
			&database.MockDatabase{},
			&graphing.MockRecipeGrapher{},
			&mockpublishers.Publisher{},
			&customerdata.MockCollector{},
			tracing.NewNoopTracerProvider(),
		)
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
					MaximumEstimatedTimeInSeconds: nil,
					MinimumTemperatureInCelsius:   nil,
					MinimumEstimatedTimeInSeconds: nil,
					MaximumTemperatureInCelsius:   nil,
					BelongsToRecipe:               exampleRecipeID,
					ID:                            recipeStepID,
					Preparation:                   types.ValidPreparation{Name: "dice"},
					Ingredients: []*types.RecipeStepIngredient{
						{
							RecipeStepProductID: nil,
							Ingredient: &types.ValidIngredient{
								MaximumIdealStorageTemperatureInCelsius: nil,
								MinimumIdealStorageTemperatureInCelsius: pointers.Float32Pointer(2.5),
								PluralName:                              "chicken breasts",
								StorageInstructions:                     "keep frozen",
								Name:                                    "chicken breast",
								ID:                                      fakes.BuildFakeID(),
							},
							Name:                "chicken breast",
							ID:                  fakes.BuildFakeID(),
							BelongsToRecipeStep: recipeStepID,
							MeasurementUnit:     types.ValidMeasurementUnit{Name: "gram", PluralName: "grams"},
							MinimumQuantity:     900,
							MaximumQuantity:     900,
							Optional:            false,
							ProductOfRecipeStep: false,
						},
					},
					Products: []*types.RecipeStepProduct{
						{
							MinimumStorageTemperatureInCelsius: nil,
							MaximumStorageTemperatureInCelsius: nil,
							StorageInstructions:                "",
							Name:                               "diced chicken breast",
							Type:                               types.RecipeStepProductIngredientType,
							BelongsToRecipeStep:                recipeStepID,
							ID:                                 fakes.BuildFakeID(),
							QuantityNotes:                      "",
							MeasurementUnit:                    types.ValidMeasurementUnit{},
							MaximumStorageDurationInSeconds:    0,
							MaximumQuantity:                    0,
							MinimumQuantity:                    0,
							Compostable:                        false,
						},
					},
					Instruments: nil,
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

		mdm := database.NewMockDatabase()
		mdm.MealPlanDataManager.On("GetFinalizedMealPlanIDsForTheNextWeek", testutils.ContextMatcher).Return(exampleFinalizedMealPlanResults, nil)

		mockGrapher := &graphing.MockRecipeGrapher{}
		for _, result := range exampleFinalizedMealPlanResults {
			mdm.MealPlanEventDataManager.On("GetMealPlanEvent", testutils.ContextMatcher, result.MealPlanID, result.MealPlanEventID).Return(exampleMealPlanEvent, nil)

			for _, recipeID := range result.RecipeIDs {
				mdm.RecipeDataManager.On("GetRecipe", testutils.ContextMatcher, recipeID).Return(recipeMap[recipeID], nil)
				mockGrapher.On("FindStepsEligibleForAdvancedCreation", testutils.ContextMatcher, recipeMap[recipeID]).Return([]*types.RecipeStep{}, nil)
			}
		}
		w.grapher = mockGrapher
		w.dataManager = mdm

		expected := map[string][]*types.AdvancedPrepStepDatabaseCreationInput{
			exampleFinalizedMealPlanResult.MealPlanOptionID: {
				{
					CannotCompleteBefore: exampleMealPlanEvent.StartsAt.Add(2 * -time.Hour * 24),
					CannotCompleteAfter:  exampleMealPlanEvent.StartsAt.Add(-time.Hour * 24),
					CompletedAt:          nil,
					Status:               types.AdvancedPrepStepStatusUnfinished,
					CreationExplanation:  buildThawStepCreationExplanation([]int{0}),
					MealPlanOptionID:     exampleFinalizedMealPlanResult.MealPlanOptionID,
					RecipeStepID:         recipeStepID,
				},
			},
		}

		actual, err := w.DetermineCreatableSteps(ctx)
		assert.NoError(t, err)

		// because we can't guarantee what this will be without too much effort
		for k, v := range actual {
			for j := range v {
				actual[k][j].ID = ""
			}
		}
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, mdm, mockGrapher)
	})

	T.Run("creates step that can be done in advance and ignores later steps", func(t *testing.T) {
		t.Parallel()

		w := ProvideAdvancedPrepStepCreationEnsurerWorker(
			zerolog.NewZerologLogger(),
			&database.MockDatabase{},
			graphing.NewRecipeGrapher(tracing.NewNoopTracerProvider()),
			&mockpublishers.Publisher{},
			&customerdata.MockCollector{},
			tracing.NewNoopTracerProvider(),
		)
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

		recipeStep1ID := fakes.BuildFakeID()
		recipeStep2ID := fakes.BuildFakeID()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipe := &types.Recipe{
			Name: "Recipe 1",
			ID:   exampleRecipeID,
			Steps: []*types.RecipeStep{
				{
					MaximumEstimatedTimeInSeconds: nil,
					MinimumTemperatureInCelsius:   nil,
					MinimumEstimatedTimeInSeconds: nil,
					MaximumTemperatureInCelsius:   nil,
					Index:                         0,
					BelongsToRecipe:               exampleRecipeID,
					ID:                            recipeStep1ID,
					Preparation:                   types.ValidPreparation{Name: "massage"},
					Ingredients: []*types.RecipeStepIngredient{
						{
							RecipeStepProductID: nil,
							Ingredient: &types.ValidIngredient{
								MaximumIdealStorageTemperatureInCelsius: nil,
								MinimumIdealStorageTemperatureInCelsius: nil,
								PluralName:                              "kale",
								StorageInstructions:                     "",
								Name:                                    "kale",
								ID:                                      fakes.BuildFakeID(),
							},
							Name:                "kale",
							ID:                  fakes.BuildFakeID(),
							BelongsToRecipeStep: recipeStep1ID,
							MeasurementUnit:     types.ValidMeasurementUnit{Name: "gram", PluralName: "grams"},
							MinimumQuantity:     500,
							MaximumQuantity:     1000,
						},
					},
					Products: []*types.RecipeStepProduct{
						{
							MinimumStorageTemperatureInCelsius: nil,
							MaximumStorageTemperatureInCelsius: nil,
							StorageInstructions:                "store in an airtight container",
							Name:                               "massaged kale",
							Type:                               types.RecipeStepProductIngredientType,
							BelongsToRecipeStep:                recipeStep1ID,
							ID:                                 fakes.BuildFakeID(),
							QuantityNotes:                      "",
							MeasurementUnit: types.ValidMeasurementUnit{
								Name: "gram", PluralName: "gram",
							},
							MaximumStorageDurationInSeconds: 259200,
							MaximumQuantity:                 0,
							MinimumQuantity:                 0,
							Compostable:                     false,
						},
					},
					Instruments: nil,
				},
				{
					MaximumEstimatedTimeInSeconds: nil,
					MinimumTemperatureInCelsius:   nil,
					MinimumEstimatedTimeInSeconds: nil,
					MaximumTemperatureInCelsius:   nil,
					Index:                         1,
					BelongsToRecipe:               exampleRecipeID,
					ID:                            recipeStep1ID,
					Preparation:                   types.ValidPreparation{Name: "sautee"},
					Ingredients: []*types.RecipeStepIngredient{
						{
							RecipeStepProductID: pointers.StringPointer(recipeStep2ID),
							Ingredient:          nil,
							Name:                "massaged kale",
							ID:                  fakes.BuildFakeID(),
							BelongsToRecipeStep: recipeStep1ID,
							MeasurementUnit:     types.ValidMeasurementUnit{Name: "gram", PluralName: "grams"},
							MinimumQuantity:     500,
							MaximumQuantity:     1000,
						},
					},
					Products: []*types.RecipeStepProduct{
						{
							MinimumStorageTemperatureInCelsius: nil,
							MaximumStorageTemperatureInCelsius: nil,
							StorageInstructions:                "",
							Name:                               "cooked kale",
							Type:                               types.RecipeStepProductIngredientType,
							BelongsToRecipeStep:                recipeStep1ID,
							ID:                                 fakes.BuildFakeID(),
							QuantityNotes:                      "",
							MeasurementUnit: types.ValidMeasurementUnit{
								Name: "gram", PluralName: "gram",
							},
							MaximumStorageDurationInSeconds: 0,
							MaximumQuantity:                 0,
							MinimumQuantity:                 0,
							Compostable:                     false,
						},
					},
					Instruments: nil,
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

		mdm := database.NewMockDatabase()
		mdm.MealPlanDataManager.On("GetFinalizedMealPlanIDsForTheNextWeek", testutils.ContextMatcher).Return(exampleFinalizedMealPlanResults, nil)

		for _, result := range exampleFinalizedMealPlanResults {
			mdm.MealPlanEventDataManager.On("GetMealPlanEvent", testutils.ContextMatcher, result.MealPlanID, result.MealPlanEventID).Return(exampleMealPlanEvent, nil)

			for _, recipeID := range result.RecipeIDs {
				mdm.RecipeDataManager.On("GetRecipe", testutils.ContextMatcher, recipeID).Return(recipeMap[recipeID], nil)
			}
		}
		w.dataManager = mdm

		cannotCompleteBefore, cannotCompleteAfter := determineCreationMinAndMaxTimesForRecipeStep(exampleRecipe.Steps[0], exampleMealPlanEvent)
		expected := map[string][]*types.AdvancedPrepStepDatabaseCreationInput{
			exampleFinalizedMealPlanResult.MealPlanOptionID: {
				{
					CannotCompleteBefore: cannotCompleteBefore,
					CannotCompleteAfter:  cannotCompleteAfter,
					CompletedAt:          nil,
					Status:               types.AdvancedPrepStepStatusUnfinished,
					CreationExplanation:  advancedStepCreationExplanation,
					MealPlanOptionID:     exampleFinalizedMealPlanResult.MealPlanOptionID,
					RecipeStepID:         recipeStep1ID,
				},
			},
		}

		actual, err := w.DetermineCreatableSteps(ctx)
		assert.NoError(t, err)

		// because we can't guarantee what this will be without too much effort
		for k, v := range actual {
			for j := range v {
				actual[k][j].ID = ""
				assert.WithinDuration(t, expected[k][j].CannotCompleteBefore, actual[k][j].CannotCompleteBefore, time.Second*30)
				actual[k][j].CannotCompleteBefore = expected[k][j].CannotCompleteBefore
			}
		}
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, mdm)
	})
}
