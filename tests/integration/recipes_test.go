package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/client/httpclient"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
)

func checkRecipeEquality(t *testing.T, expected, actual *types.Recipe) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for recipe %s to be %v, but it was %v", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.Source, actual.Source, "expected Source for recipe %s to be %v, but it was %v", expected.ID, expected.Source, actual.Source)
	assert.Equal(t, expected.Description, actual.Description, "expected Description for recipe %s to be %v, but it was %v", expected.ID, expected.Description, actual.Description)
	assert.Equal(t, expected.InspiredByRecipeID, actual.InspiredByRecipeID, "expected InspiredByRecipeID for recipe %s to be %v, but it was %v", expected.ID, expected.InspiredByRecipeID, actual.InspiredByRecipeID)
	assert.NotZero(t, actual.CreatedOn)
}

// convertRecipeToRecipeUpdateInput creates an RecipeUpdateRequestInput struct from a recipe.
func convertRecipeToRecipeUpdateInput(x *types.Recipe) *types.RecipeUpdateRequestInput {
	return &types.RecipeUpdateRequestInput{
		Name:               &x.Name,
		Source:             &x.Source,
		Description:        &x.Description,
		InspiredByRecipeID: x.InspiredByRecipeID,
	}
}

func createRecipeForTest(ctx context.Context, t *testing.T, adminClient, client *httpclient.Client, recipe *types.Recipe) ([]*types.ValidIngredient, *types.ValidPreparation, *types.Recipe) {
	t.Helper()

	t.Log("creating prerequisite valid preparation")
	exampleValidPreparation := fakes.BuildFakeValidPreparation()
	exampleValidPreparationInput := fakes.BuildFakeValidPreparationCreationRequestInputFromValidPreparation(exampleValidPreparation)
	createdValidPreparation, err := adminClient.CreateValidPreparation(ctx, exampleValidPreparationInput)
	require.NoError(t, err)
	t.Logf("valid preparation %q created", createdValidPreparation.ID)

	t.Log("creating recipe")
	exampleRecipe := fakes.BuildFakeRecipe()
	if recipe != nil {
		exampleRecipe = recipe
	}

	createdValidIngredients := []*types.ValidIngredient{}
	for i, recipeStep := range exampleRecipe.Steps {
		for j := range recipeStep.Ingredients {
			t.Log("creating prerequisite valid ingredient")
			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakes.BuildFakeValidIngredientCreationRequestInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredient, createdValidIngredientErr := adminClient.CreateValidIngredient(ctx, exampleValidIngredientInput)
			require.NoError(t, createdValidIngredientErr)

			createdValidIngredients = append(createdValidIngredients, createdValidIngredient)

			exampleRecipe.Steps[i].Ingredients[j].IngredientID = stringPointer(createdValidIngredient.ID)
			exampleRecipe.Steps[i].Ingredients[j].ProductOfRecipeStep = false
		}
	}

	exampleRecipeInput := fakes.BuildFakeRecipeCreationRequestInputFromRecipe(exampleRecipe)
	for i := range exampleRecipeInput.Steps {
		exampleRecipeInput.Steps[i].PreparationID = createdValidPreparation.ID
	}

	createdRecipe, err := client.CreateRecipe(ctx, exampleRecipeInput)
	require.NoError(t, err)
	t.Logf("recipe %q created", createdRecipe.ID)
	checkRecipeEquality(t, exampleRecipe, createdRecipe)

	createdRecipe, err = client.GetRecipe(ctx, createdRecipe.ID)
	requireNotNilAndNoProblems(t, createdRecipe, err)
	checkRecipeEquality(t, exampleRecipe, createdRecipe)

	require.NotEmpty(t, createdRecipe.Steps, "created recipe must have steps")

	return createdValidIngredients, createdValidPreparation, createdRecipe
}

func (s *TestSuite) TestRecipes_Realistic() {
	s.runForEachClient("sopa de frijol", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating prerequisite valid preparation")
			soakBase := fakes.BuildFakeValidPreparation()
			soakInput := fakes.BuildFakeValidPreparationCreationRequestInputFromValidPreparation(soakBase)
			soak, err := testClients.admin.CreateValidPreparation(ctx, soakInput)
			require.NoError(t, err)
			t.Logf("valid preparation %q created", soak.ID)

			t.Log("creating prerequisite valid ingredient")
			pintoBeanBase := fakes.BuildFakeValidIngredient()
			pintoBeanInput := fakes.BuildFakeValidIngredientCreationRequestInputFromValidIngredient(pintoBeanBase)
			pintoBeans, createdValidIngredientErr := testClients.admin.CreateValidIngredient(ctx, pintoBeanInput)
			require.NoError(t, createdValidIngredientErr)

			t.Log("creating prerequisite valid ingredient")
			waterBase := fakes.BuildFakeValidIngredient()
			waterInput := fakes.BuildFakeValidIngredientCreationRequestInputFromValidIngredient(waterBase)
			water, createdValidIngredientErr := testClients.admin.CreateValidIngredient(ctx, waterInput)
			require.NoError(t, createdValidIngredientErr)

			t.Log("creating prerequisite valid ingredient")
			garlicPaste := fakes.BuildFakeValidIngredient()
			garlicPasteInput := fakes.BuildFakeValidIngredientCreationRequestInputFromValidIngredient(garlicPaste)
			garlicPaste, garlicPasteErr := testClients.admin.CreateValidIngredient(ctx, garlicPasteInput)
			require.NoError(t, garlicPasteErr)

			t.Log("creating recipe")
			expected := &types.Recipe{
				Name:        "sopa de frijol",
				Description: "",
				Steps: []*types.RecipeStep{
					{
						MinTemperatureInCelsius: nil,
						Products: []*types.RecipeStepProduct{
							{
								Name:          "soaked pinto beans",
								QuantityType:  "grams",
								QuantityNotes: "",
								QuantityValue: 1000,
							},
						},
						Notes:       "first step",
						Preparation: *soak,
						Ingredients: []*types.RecipeStepIngredient{
							{
								RecipeStepProductID:  nil,
								IngredientID:         &pintoBeans.ID,
								Name:                 "pinto beans",
								QuantityType:         "grams",
								MinimumQuantityValue: 500,
								ProductOfRecipeStep:  false,
							},
							{
								RecipeStepProductID:  nil,
								IngredientID:         &water.ID,
								Name:                 "water",
								QuantityType:         "cups",
								MinimumQuantityValue: 5,
								ProductOfRecipeStep:  false,
							},
						},
						Index: 0,
					},
					{
						MinTemperatureInCelsius: nil,
						Products: []*types.RecipeStepProduct{
							{
								Name:          "final output",
								QuantityType:  "grams",
								QuantityNotes: "",
								QuantityValue: 1010,
							},
						},
						Notes:       "first step",
						Preparation: *soak,
						Ingredients: []*types.RecipeStepIngredient{
							{
								Name:                 "soaked pinto beans",
								QuantityType:         "grams",
								MinimumQuantityValue: 1000,
								ProductOfRecipeStep:  true,
							},
							{
								IngredientID:         &garlicPaste.ID,
								Name:                 "garlic paste",
								QuantityType:         "grams",
								MinimumQuantityValue: 10,
								ProductOfRecipeStep:  false,
							},
						},
						Index: 1,
					},
				},
			}

			exampleRecipeInput := &types.RecipeCreationRequestInput{
				Name:        expected.Name,
				Description: expected.Description,
			}
			for _, step := range expected.Steps {
				newStep := &types.RecipeStepCreationRequestInput{
					TemperatureInCelsius:      step.MinTemperatureInCelsius,
					Notes:                     step.Notes,
					PreparationID:             step.Preparation.ID,
					BelongsToRecipe:           step.BelongsToRecipe,
					ID:                        step.ID,
					Index:                     step.Index,
					MinEstimatedTimeInSeconds: step.MinEstimatedTimeInSeconds,
					MaxEstimatedTimeInSeconds: step.MaxEstimatedTimeInSeconds,
					Optional:                  step.Optional,
				}

				for _, ingredient := range step.Ingredients {
					newIngredient := &types.RecipeStepIngredientCreationRequestInput{
						IngredientID:        ingredient.IngredientID,
						ID:                  ingredient.ID,
						BelongsToRecipeStep: ingredient.BelongsToRecipeStep,
						Name:                ingredient.Name,
						QuantityType:        ingredient.QuantityType,
						QuantityNotes:       ingredient.QuantityNotes,
						IngredientNotes:     ingredient.IngredientNotes,
						QuantityValue:       ingredient.MinimumQuantityValue,
						ProductOfRecipeStep: ingredient.ProductOfRecipeStep,
					}
					newStep.Ingredients = append(newStep.Ingredients, newIngredient)
				}

				for _, product := range step.Products {
					newProduct := &types.RecipeStepProductCreationRequestInput{
						ID:                  product.ID,
						Name:                product.Name,
						QuantityType:        product.QuantityType,
						QuantityNotes:       product.QuantityNotes,
						BelongsToRecipeStep: product.BelongsToRecipeStep,
						QuantityValue:       product.QuantityValue,
					}
					newStep.Products = append(newStep.Products, newProduct)
				}

				exampleRecipeInput.Steps = append(exampleRecipeInput.Steps, newStep)
			}

			created, err := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
			require.NoError(t, err)
			t.Logf("recipe %q created", created.ID)
			checkRecipeEquality(t, expected, created)

			created, err = testClients.main.GetRecipe(ctx, created.ID)
			requireNotNilAndNoProblems(t, created, err)
			checkRecipeEquality(t, expected, created)

			createdJSON, _ := json.Marshal(created)
			t.Log(string(createdJSON))

			recipeStepProductIndex := -1
			for i, ingredient := range created.Steps[1].Ingredients {
				if ingredient.RecipeStepProductID != nil {
					recipeStepProductIndex = i
				}
			}

			assert.NotEqual(t, -1, recipeStepProductIndex)
			require.NotNil(t, created.Steps[1].Ingredients[recipeStepProductIndex].RecipeStepProductID)
			assert.Equal(t, created.Steps[0].Products[0].ID, *created.Steps[1].Ingredients[recipeStepProductIndex].RecipeStepProductID)
		}
	})
}

func (s *TestSuite) TestRecipes_CompleteLifecycle() {
	s.runForEachClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			_, _, createdRecipe := createRecipeForTest(ctx, t, testClients.admin, testClients.main, nil)

			t.Log("changing recipe")
			newRecipe := fakes.BuildFakeRecipe()
			createdRecipe.Update(convertRecipeToRecipeUpdateInput(newRecipe))
			assert.NoError(t, testClients.main.UpdateRecipe(ctx, createdRecipe))

			t.Log("fetching changed recipe")
			actual, err := testClients.main.GetRecipe(ctx, createdRecipe.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert recipe equality
			checkRecipeEquality(t, newRecipe, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			t.Log("cleaning up recipe")
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipes_AlsoCreateMeal() {
	s.runForEachClient("should be able to create a meal and a recipe", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakes.BuildFakeValidPreparationCreationRequestInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparation, err := testClients.admin.CreateValidPreparation(ctx, exampleValidPreparationInput)
			require.NoError(t, err)

			exampleRecipe := fakes.BuildFakeRecipe()
			createdValidIngredients := []*types.ValidIngredient{}
			for i, recipeStep := range exampleRecipe.Steps {
				for j := range recipeStep.Ingredients {
					t.Log("creating prerequisite valid ingredient")
					exampleValidIngredient := fakes.BuildFakeValidIngredient()
					exampleValidIngredientInput := fakes.BuildFakeValidIngredientCreationRequestInputFromValidIngredient(exampleValidIngredient)
					createdValidIngredient, createdValidIngredientErr := testClients.admin.CreateValidIngredient(ctx, exampleValidIngredientInput)
					require.NoError(t, createdValidIngredientErr)

					createdValidIngredients = append(createdValidIngredients, createdValidIngredient)

					exampleRecipe.Steps[i].Ingredients[j].IngredientID = stringPointer(createdValidIngredient.ID)
					exampleRecipe.Steps[i].Ingredients[j].ProductOfRecipeStep = false
				}
			}

			exampleRecipeInput := fakes.BuildFakeRecipeCreationRequestInputFromRecipe(exampleRecipe)
			for i := range exampleRecipeInput.Steps {
				exampleRecipeInput.Steps[i].PreparationID = createdValidPreparation.ID
			}

			exampleRecipeInput.AlsoCreateMeal = true

			createdRecipe, err := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
			require.NoError(t, err)
			t.Logf("recipe %q created", createdRecipe.ID)
			checkRecipeEquality(t, exampleRecipe, createdRecipe)

			mealResults, err := testClients.main.SearchForMeals(ctx, createdRecipe.Name, nil)
			requireNotNilAndNoProblems(t, mealResults, err)

			foundMealID := ""
			for _, m := range mealResults.Meals {
				meal, mealFetchErr := testClients.main.GetMeal(ctx, m.ID)
				requireNotNilAndNoProblems(t, meal, mealFetchErr)

				for _, r := range meal.Recipes {
					if r.ID == createdRecipe.ID {
						foundMealID = r.ID
					}
				}
			}

			require.NotEmpty(t, foundMealID)

			t.Log("cleaning up recipe")
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipes_Listing() {
	s.runForEachClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating prerequisite valid ingredient")
			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakes.BuildFakeValidIngredientCreationRequestInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredient, err := testClients.admin.CreateValidIngredient(ctx, exampleValidIngredientInput)
			require.NoError(t, err)
			t.Logf("valid ingredient %q created", createdValidIngredient.ID)

			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

			createdValidIngredient, err = testClients.main.GetValidIngredient(ctx, createdValidIngredient.ID)
			requireNotNilAndNoProblems(t, createdValidIngredient, err)
			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

			t.Log("creating prerequisite valid preparation")
			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakes.BuildFakeValidPreparationCreationRequestInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparation, err := testClients.admin.CreateValidPreparation(ctx, exampleValidPreparationInput)
			require.NoError(t, err)
			t.Logf("valid preparation %q created", createdValidPreparation.ID)

			checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

			createdValidPreparation, err = testClients.main.GetValidPreparation(ctx, createdValidPreparation.ID)
			requireNotNilAndNoProblems(t, createdValidPreparation, err)
			checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

			t.Log("creating recipes")
			var expected []*types.Recipe
			for i := 0; i < 5; i++ {
				_, _, createdRecipe := createRecipeForTest(ctx, t, testClients.admin, testClients.main, nil)

				expected = append(expected, createdRecipe)
			}

			// assert recipe list equality
			actual, err := testClients.main.GetRecipes(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Recipes),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Recipes),
			)

			t.Log("cleaning up")
			for _, createdRecipe := range expected {
				assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
			}
		}
	})
}

func (s *TestSuite) TestRecipes_Searching() {
	s.runForEachClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating prerequisite valid ingredient")
			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakes.BuildFakeValidIngredientCreationRequestInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredient, err := testClients.admin.CreateValidIngredient(ctx, exampleValidIngredientInput)
			require.NoError(t, err)
			t.Logf("valid ingredient %q created", createdValidIngredient.ID)

			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

			createdValidIngredient, err = testClients.main.GetValidIngredient(ctx, createdValidIngredient.ID)
			requireNotNilAndNoProblems(t, createdValidIngredient, err)
			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

			t.Log("creating prerequisite valid preparation")
			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakes.BuildFakeValidPreparationCreationRequestInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparation, err := testClients.admin.CreateValidPreparation(ctx, exampleValidPreparationInput)
			require.NoError(t, err)
			t.Logf("valid preparation %q created", createdValidPreparation.ID)

			checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

			createdValidPreparation, err = testClients.main.GetValidPreparation(ctx, createdValidPreparation.ID)
			requireNotNilAndNoProblems(t, createdValidPreparation, err)
			checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

			exampleRecipe := fakes.BuildFakeRecipe()

			t.Log("creating recipes")
			var expected []*types.Recipe
			for i := 0; i < 5; i++ {
				exampleRecipe.Name = fmt.Sprintf("example%d", i)
				_, _, createdRecipe := createRecipeForTest(ctx, t, testClients.admin, testClients.main, exampleRecipe)

				expected = append(expected, createdRecipe)
			}

			// assert recipe list equality
			actual, err := testClients.main.SearchForRecipes(ctx, "example", nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Recipes),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Recipes),
			)

			t.Log("cleaning up")
			for _, createdRecipe := range expected {
				assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
			}
		}
	})
}
