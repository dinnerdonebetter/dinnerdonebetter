package recipeanalysis

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/internal/pointers"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
)

func newAnalyzerForTest(t *testing.T) *recipeAnalyzer {
	t.Helper()

	return &recipeAnalyzer{
		tracer: tracing.NewTracerForTest(t.Name()),
		logger: logging.NewNoopLogger(),
	}
}

func TestRecipeGrapher_makeGraphForRecipe(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		g := newAnalyzerForTest(t)

		ctx := context.Background()
		r := &types.Recipe{
			Steps: []*types.RecipeStep{
				{},
			},
		}

		actual, err := g.makeGraphForRecipe(ctx, r)
		assert.NoError(t, err)
		assert.NotNil(t, actual)
	})
}

func TestRecipeGrapher_makeDAGForRecipe(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		g := newAnalyzerForTest(t)

		ctx := context.Background()
		r := &types.Recipe{
			Steps: []*types.RecipeStep{
				{},
			},
		}

		actual, err := g.makeDAGForRecipe(ctx, r)
		assert.NoError(t, err)
		assert.NotNil(t, actual)
	})
}

func TestRecipeAnalyzer_GenerateMealPlanTasksForRecipe(T *testing.T) {
	T.Parallel()

	T.Run("creates frozen thawing steps", func(t *testing.T) {
		t.Parallel()

		g := newAnalyzerForTest(t)
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
							MaximumStorageDurationInSeconds:    nil,
							MaximumQuantity:                    0,
							MinimumQuantity:                    0,
							Compostable:                        false,
						},
					},
					Instruments: nil,
				},
			},
		}

		expected := []*types.MealPlanTaskDatabaseCreationInput{
			{
				CannotCompleteBefore: time.Now(),
				CannotCompleteAfter:  time.Now(),
				CreationExplanation:  buildThawStepCreationExplanation(1, 0),
				MealPlanOptionID:     exampleMealPlanOption.ID,
				RecipeSteps: []*types.MealPlanTaskRecipeStepDatabaseCreationInput{
					{
						AppliesToRecipeStep: recipeStepID,
						SatisfiesRecipeStep: false,
					},
				},
			},
		}

		actual, err := g.GenerateMealPlanTasksForRecipe(ctx, exampleMealPlanEvent.StartsAt, exampleMealPlanOption.ID, exampleRecipe)
		assert.NoError(t, err)

		for i := range expected {
			expected[i].ID = actual[i].ID
			expected[i].CannotCompleteBefore = actual[i].CannotCompleteBefore
			expected[i].CannotCompleteAfter = actual[i].CannotCompleteAfter

			for j := range expected[i].RecipeSteps {
				expected[i].RecipeSteps[j].BelongsToMealPlanTask = actual[i].RecipeSteps[j].BelongsToMealPlanTask
				expected[i].RecipeSteps[j].ID = actual[i].RecipeSteps[j].ID
			}
		}

		assert.Equal(t, expected, actual)
	})

	T.Run("creates step that can be done in advance and ignores later steps", func(t *testing.T) {
		t.Parallel()

		g := newAnalyzerForTest(t)

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
		recipeStep3ID := fakes.BuildFakeID()
		recipeStepProduct1ID := fakes.BuildFakeID()
		recipeStepProduct2ID := fakes.BuildFakeID()

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
							MeasurementUnit: types.ValidMeasurementUnit{
								Name:       "gram",
								PluralName: "grams",
							},
							MinimumQuantity: 500,
							MaximumQuantity: 1000,
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
							ID:                                 recipeStepProduct1ID,
							QuantityNotes:                      "",
							MeasurementUnit: types.ValidMeasurementUnit{
								Name:       "gram",
								PluralName: "grams",
							},
							MaximumStorageDurationInSeconds: pointers.Uint32Pointer(259200),
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
					ID:                            recipeStep2ID,
					Preparation:                   types.ValidPreparation{Name: "chop"},
					Ingredients: []*types.RecipeStepIngredient{
						{
							RecipeStepProductID: nil,
							Ingredient: &types.ValidIngredient{
								MaximumIdealStorageTemperatureInCelsius: nil,
								MinimumIdealStorageTemperatureInCelsius: nil,
								PluralName:                              "cherry tomatoes",
								StorageInstructions:                     "",
								Name:                                    "cherry tomato",
								ID:                                      fakes.BuildFakeID(),
							},
							Name:                "cherry tomato",
							ID:                  fakes.BuildFakeID(),
							BelongsToRecipeStep: recipeStep2ID,
							MeasurementUnit: types.ValidMeasurementUnit{
								Name:       "gram",
								PluralName: "grams",
							},
							MinimumQuantity: 500,
							MaximumQuantity: 1000,
						},
					},
					Products: []*types.RecipeStepProduct{
						{
							MinimumStorageTemperatureInCelsius: nil,
							MaximumStorageTemperatureInCelsius: nil,
							StorageInstructions:                "",
							Name:                               "chopped cherry tomatoes",
							Type:                               types.RecipeStepProductIngredientType,
							BelongsToRecipeStep:                recipeStep2ID,
							ID:                                 recipeStepProduct2ID,
							QuantityNotes:                      "",
							MeasurementUnit: types.ValidMeasurementUnit{
								Name:       "gram",
								PluralName: "grams",
							},
							MaximumStorageDurationInSeconds: nil,
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
					Index:                         2,
					BelongsToRecipe:               exampleRecipeID,
					ID:                            recipeStep3ID,
					Preparation:                   types.ValidPreparation{Name: "sautee"},
					Ingredients: []*types.RecipeStepIngredient{
						{
							RecipeStepProductID: pointers.StringPointer(recipeStepProduct1ID),
							Ingredient:          nil,
							Name:                "massaged kale",
							ID:                  fakes.BuildFakeID(),
							ProductOfRecipeStep: true,
							BelongsToRecipeStep: recipeStep3ID,
							MeasurementUnit: types.ValidMeasurementUnit{
								Name:       "gram",
								PluralName: "grams",
							},
							MinimumQuantity: 500,
							MaximumQuantity: 1000,
						},
						{
							RecipeStepProductID: pointers.StringPointer(recipeStepProduct2ID),
							Ingredient:          nil,
							Name:                "massaged kale",
							ID:                  fakes.BuildFakeID(),
							ProductOfRecipeStep: true,
							BelongsToRecipeStep: recipeStep3ID,
							MeasurementUnit: types.ValidMeasurementUnit{
								Name:       "gram",
								PluralName: "grams",
							},
							MinimumQuantity: 500,
							MaximumQuantity: 1000,
						},
					},
					Products: []*types.RecipeStepProduct{
						{
							MinimumStorageTemperatureInCelsius: nil,
							MaximumStorageTemperatureInCelsius: nil,
							StorageInstructions:                "",
							Name:                               "cooked kale",
							Type:                               types.RecipeStepProductIngredientType,
							BelongsToRecipeStep:                recipeStep3ID,
							ID:                                 fakes.BuildFakeID(),
							QuantityNotes:                      "",
							MeasurementUnit: types.ValidMeasurementUnit{
								Name:       "gram",
								PluralName: "gram",
							},
							MaximumStorageDurationInSeconds: nil,
							MaximumQuantity:                 0,
							MinimumQuantity:                 0,
							Compostable:                     false,
						},
					},
					Instruments: nil,
				},
			},
		}

		expected := []*types.MealPlanTaskDatabaseCreationInput{
			{
				CannotCompleteBefore: time.Now(),
				CannotCompleteAfter:  time.Now(),
				CreationExplanation:  storagePrepCreationExplanation,
				MealPlanOptionID:     exampleMealPlanOption.ID,
				RecipeSteps: []*types.MealPlanTaskRecipeStepDatabaseCreationInput{
					{
						AppliesToRecipeStep: recipeStep1ID,
						SatisfiesRecipeStep: true,
					},
					{
						AppliesToRecipeStep: recipeStep2ID,
						SatisfiesRecipeStep: true,
					},
				},
			},
		}

		actual, err := g.GenerateMealPlanTasksForRecipe(ctx, exampleMealPlanEvent.StartsAt, exampleMealPlanOption.ID, exampleRecipe)
		assert.NoError(t, err)

		require.Equal(t, len(actual), len(expected))

		for i := range expected {
			expected[i].ID = actual[i].ID
			expected[i].CannotCompleteBefore = actual[i].CannotCompleteBefore
			expected[i].CannotCompleteAfter = actual[i].CannotCompleteAfter

			for j := range expected[i].RecipeSteps {
				expected[i].RecipeSteps[j].BelongsToMealPlanTask = actual[i].RecipeSteps[j].BelongsToMealPlanTask
				expected[i].RecipeSteps[j].ID = actual[i].RecipeSteps[j].ID
			}
		}

		assert.Equal(t, expected, actual)
	})
}
