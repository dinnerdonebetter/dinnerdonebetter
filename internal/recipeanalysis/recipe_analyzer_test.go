package recipeanalysis

import (
	"context"
	"github.com/prixfixeco/api_server/internal/pointers"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/graph/simple"
	"testing"
	"time"

	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

func TestRecipeGrapher_makeGraphForRecipe(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.SkipNow()

		g := &recipeAnalyzer{
			tracer: tracing.NewTracerForTest(t.Name()),
		}

		ctx := context.Background()
		r := &types.Recipe{
			Steps: []*types.RecipeStep{
				{},
			},
		}

		expected := &simple.DirectedGraph{}

		actual, err := g.makeGraphForRecipe(ctx, r)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestRecipeAnalyzer_GenerateAdvancedStepCreationForRecipe(T *testing.T) {
	T.Parallel()

	T.Run("creates frozen thawing steps", func(t *testing.T) {
		t.Parallel()

		g := &recipeAnalyzer{
			tracer: tracing.NewTracerForTest(t.Name()),
		}
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

		exampleFinalizedMealPlanResult := &types.FinalizedMealPlanDatabaseResult{
			MealPlanID:       exampleMealPlan.ID,
			MealPlanEventID:  exampleMealPlanEvent.ID,
			MealPlanOptionID: exampleMealPlanOption.ID,
			MealID:           exampleMeal.ID,
			RecipeIDs: []string{
				exampleRecipe.ID,
			},
		}

		expected := []*types.AdvancedPrepStepDatabaseCreationInput{
			{
				CannotCompleteBefore: time.Now(),
				CannotCompleteAfter:  time.Now(),
				CompletedAt:          nil,
				Status:               types.AdvancedPrepStepStatusUnfinished,
				CreationExplanation:  t.Name(),
				MealPlanOptionID:     exampleFinalizedMealPlanResult.MealPlanOptionID,
				RecipeStepID:         exampleRecipe.ID,
			},
		}

		actual, err := g.GenerateAdvancedStepCreationForRecipe(ctx, exampleMealPlanEvent, exampleMealPlanOption.ID, exampleRecipe)
		assert.NoError(t, err)

		assert.Equal(t, expected, actual)
	})

	T.Run("creates step that can be done in advance and ignores later steps", func(t *testing.T) {
		t.Parallel()

		g := &recipeAnalyzer{
			tracer: tracing.NewTracerForTest(t.Name()),
		}

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

		exampleFinalizedMealPlanResult := &types.FinalizedMealPlanDatabaseResult{
			MealPlanID:       exampleMealPlan.ID,
			MealPlanEventID:  exampleMealPlanEvent.ID,
			MealPlanOptionID: exampleMealPlanOption.ID,
			MealID:           exampleMeal.ID,
			RecipeIDs: []string{
				exampleRecipe.ID,
			},
		}

		expected := []*types.AdvancedPrepStepDatabaseCreationInput{
			{
				CannotCompleteBefore: time.Now(),
				CannotCompleteAfter:  time.Now(),
				CompletedAt:          nil,
				Status:               types.AdvancedPrepStepStatusUnfinished,
				CreationExplanation:  t.Name(),
				MealPlanOptionID:     exampleFinalizedMealPlanResult.MealPlanOptionID,
				RecipeStepID:         recipeStep1ID,
			},
		}

		actual, err := g.GenerateAdvancedStepCreationForRecipe(ctx, exampleMealPlanEvent, exampleMealPlanOption.ID, exampleRecipe)
		assert.NoError(t, err)

		assert.Equal(t, expected, actual)
	})
}
