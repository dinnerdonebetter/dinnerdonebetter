package recipeanalysis

import (
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/fakes"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"

	"github.com/stretchr/testify/assert"
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

		ctx := t.Context()
		r := &mealplanning.Recipe{
			Steps: []*mealplanning.RecipeStep{
				{},
			},
		}

		actual, err := g.MakeGraphForRecipe(ctx, r)
		assert.NoError(t, err)
		assert.NotNil(t, actual)
	})
}

func TestRecipeGrapher_makeDAGForRecipe(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		g := newAnalyzerForTest(t)

		ctx := t.Context()
		r := &mealplanning.Recipe{
			Steps: []*mealplanning.RecipeStep{
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
		ctx := t.Context()

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
								Min: 900,
								Max: pointer.To(float32(900)),
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

		expected := []*mealplanning.MealPlanTaskDatabaseCreationInput{
			{
				CreationExplanation: buildThawStepCreationExplanation(1, 0),
				MealPlanOptionID:    exampleMealPlanOption.ID,
			},
		}

		actual, err := g.GenerateMealPlanTasksForRecipe(ctx, exampleMealPlanOption.ID, exampleRecipe)
		assert.NoError(t, err)

		for i := range expected {
			expected[i].ID = actual[i].ID
		}

		assert.Equal(t, expected, actual)
	})
}

func Test_recipeAnalyzer_RenderMermaidDiagramForRecipe(T *testing.T) {
	T.Parallel()

	T.Run("basic", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		g := newAnalyzerForTest(t)

		dice := fakes.BuildFakeValidPreparation()
		dice.Name = "dice"
		sautee := fakes.BuildFakeValidPreparation()
		sautee.Name = "sautee"

		step1ID := fakes.BuildFakeID()
		step2ID := fakes.BuildFakeID()
		step3ID := fakes.BuildFakeID()
		step4ID := fakes.BuildFakeID()
		dicedOnionRecipeStepProductID := fakes.BuildFakeID()
		dicedCarrotRecipeStepProductID := fakes.BuildFakeID()
		dicedCeleryRecipeStepProductID := fakes.BuildFakeID()

		recipe := &mealplanning.Recipe{
			Name: "example recipe",
			Steps: []*mealplanning.RecipeStep{
				{
					ID:          step1ID,
					Preparation: *dice,
					Ingredients: []*mealplanning.RecipeStepIngredient{
						{
							Ingredient: fakes.BuildFakeValidIngredient(),
							Name:       "onion",
						},
					},
					Products: []*mealplanning.RecipeStepProduct{
						{
							ID:   dicedOnionRecipeStepProductID,
							Name: "diced onion",
							Type: mealplanning.RecipeStepProductIngredientType,
						},
					},
					Notes: "first step",
					Index: 0,
				},
				{
					ID:          step2ID,
					Preparation: *dice,
					Ingredients: []*mealplanning.RecipeStepIngredient{
						{
							Ingredient: fakes.BuildFakeValidIngredient(),
							Name:       "carrot",
						},
					},
					Products: []*mealplanning.RecipeStepProduct{
						{
							ID:   dicedCarrotRecipeStepProductID,
							Name: "diced carrot",
							Type: mealplanning.RecipeStepProductIngredientType,
						},
					},
					Notes: "second step",
					Index: 1,
				},
				{
					ID:          step3ID,
					Preparation: *dice,
					Ingredients: []*mealplanning.RecipeStepIngredient{
						{
							Ingredient: fakes.BuildFakeValidIngredient(),
							Name:       "celery",
						},
					},
					Products: []*mealplanning.RecipeStepProduct{
						{
							ID:   dicedCeleryRecipeStepProductID,
							Name: "diced celery",
							Type: mealplanning.RecipeStepProductIngredientType,
						},
					},
					Notes: "third step",
					Index: 2,
				},
				{
					ID:          step4ID,
					Preparation: *sautee,
					Ingredients: []*mealplanning.RecipeStepIngredient{
						{
							Name:                "diced onion",
							RecipeStepProductID: pointer.To(dicedOnionRecipeStepProductID),
						},
						{
							Name:                "diced carrot",
							RecipeStepProductID: pointer.To(dicedCarrotRecipeStepProductID),
						},
						{
							Name:                "diced celery",
							RecipeStepProductID: pointer.To(dicedCeleryRecipeStepProductID),
						},
					},
					Products: []*mealplanning.RecipeStepProduct{
						{
							ID:   dicedOnionRecipeStepProductID,
							Name: "sauteed mire poix",
							Type: mealplanning.RecipeStepProductIngredientType,
						},
					},
					Notes: "fourth step",
					Index: 3,
				},
			},
		}

		expected := `flowchart TD;
	Step1["Step #1 (dice)"];
	Step2["Step #2 (dice)"];
	Step3["Step #3 (dice)"];
	Step4["Step #4 (sautee)"];
	Step1 -->|ingredient| Step4;
	Step2 -->|ingredient| Step4;
	Step3 -->|ingredient| Step4;
`
		actual := g.RenderMermaidDiagramForRecipe(ctx, recipe)

		assert.Equal(t, expected, actual)
	})
}
