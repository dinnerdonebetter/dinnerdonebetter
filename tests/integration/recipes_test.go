package integration

import (
	"context"
	"fmt"
	"github.com/prixfixeco/api_server/pkg/types/converters"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/internal/pointers"
	"github.com/prixfixeco/api_server/pkg/apiclient"
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
	assert.Equal(t, expected.YieldsPortions, actual.YieldsPortions, "expected YieldsPortions for recipe %s to be %v, but it was %v", expected.ID, expected.YieldsPortions, actual.YieldsPortions)
	assert.Equal(t, expected.SealOfApproval, actual.SealOfApproval, "expected SealOfApproval for recipe %s to be %v, but it was %v", expected.ID, expected.SealOfApproval, actual.SealOfApproval)
	assert.NotZero(t, actual.CreatedAt)
}

func createRecipeForTest(ctx context.Context, t *testing.T, adminClient, client *apiclient.Client, recipe *types.Recipe) ([]*types.ValidIngredient, *types.ValidPreparation, *types.Recipe) {
	t.Helper()

	t.Log("creating prerequisite valid preparation")
	exampleValidPreparation := fakes.BuildFakeValidPreparation()
	exampleValidPreparationInput := converters.ConvertValidPreparationToValidPreparationCreationRequestInput(exampleValidPreparation)
	createdValidPreparation, err := adminClient.CreateValidPreparation(ctx, exampleValidPreparationInput)
	require.NoError(t, err)
	t.Logf("valid preparation %q created", createdValidPreparation.ID)

	t.Log("creating valid measurement unit")
	exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()
	exampleValidMeasurementUnitInput := converters.ConvertValidMeasurementUnitToValidMeasurementUnitCreationRequestInput(exampleValidMeasurementUnit)
	createdValidMeasurementUnit, err := adminClient.CreateValidMeasurementUnit(ctx, exampleValidMeasurementUnitInput)
	require.NoError(t, err)
	t.Logf("valid measurement unit %q created", createdValidMeasurementUnit.ID)
	checkValidMeasurementUnitEquality(t, exampleValidMeasurementUnit, createdValidMeasurementUnit)

	t.Log("creating valid instrument")
	exampleValidInstrument := fakes.BuildFakeValidInstrument()
	exampleValidInstrumentInput := converters.ConvertValidInstrumentToValidInstrumentCreationRequestInput(exampleValidInstrument)
	createdValidInstrument, err := adminClient.CreateValidInstrument(ctx, exampleValidInstrumentInput)
	require.NoError(t, err)
	t.Logf("valid instrument %q created", createdValidInstrument.ID)
	checkValidInstrumentEquality(t, exampleValidInstrument, createdValidInstrument)

	createdValidMeasurementUnit, err = adminClient.GetValidMeasurementUnit(ctx, createdValidMeasurementUnit.ID)
	requireNotNilAndNoProblems(t, createdValidMeasurementUnit, err)
	checkValidMeasurementUnitEquality(t, exampleValidMeasurementUnit, createdValidMeasurementUnit)

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
			exampleValidIngredientInput := converters.ConvertValidIngredientToValidIngredientCreationRequestInput(exampleValidIngredient)
			createdValidIngredient, createdValidIngredientErr := adminClient.CreateValidIngredient(ctx, exampleValidIngredientInput)
			require.NoError(t, createdValidIngredientErr)

			createdValidIngredients = append(createdValidIngredients, createdValidIngredient)

			exampleRecipe.Steps[i].Ingredients[j].Ingredient = createdValidIngredient
			exampleRecipe.Steps[i].Ingredients[j].ProductOfRecipeStep = false
			exampleRecipe.Steps[i].Ingredients[j].MeasurementUnit = *createdValidMeasurementUnit
		}

		for j := range recipeStep.Products {
			exampleRecipe.Steps[i].Products[j].MeasurementUnit = *createdValidMeasurementUnit
		}

		for j := range recipeStep.Instruments {
			recipeStep.Instruments[j].Instrument = createdValidInstrument
		}
	}

	exampleRecipeInput := converters.ConvertRecipeToRecipeCreationRequestInputFromRecipe(exampleRecipe)
	for i := range exampleRecipeInput.Steps {
		exampleRecipeInput.Steps[i].PreparationID = createdValidPreparation.ID
	}

	examplePrepTask := fakes.BuildFakeRecipePrepTask()
	examplePrepTask.TaskSteps = []*types.RecipePrepTaskStep{
		{
			BelongsToRecipeStep: exampleRecipe.Steps[0].ID,
			SatisfiesRecipeStep: false,
		},
	}
	exampleRecipeInput.PrepTasks = []*types.RecipePrepTaskWithinRecipeCreationRequestInput{
		converters.ConvertRecipePrepTaskToRecipePrepTaskWithinRecipeCreationRequestInput(exampleRecipe, examplePrepTask),
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
			soakInput := converters.ConvertValidPreparationToValidPreparationCreationRequestInput(soakBase)
			soak, err := testClients.admin.CreateValidPreparation(ctx, soakInput)
			require.NoError(t, err)
			t.Logf("valid preparation %q created", soak.ID)

			t.Log("creating prerequisite valid preparation")
			mixBase := fakes.BuildFakeValidPreparation()
			mixInput := converters.ConvertValidPreparationToValidPreparationCreationRequestInput(mixBase)
			mix, err := testClients.admin.CreateValidPreparation(ctx, mixInput)
			require.NoError(t, err)
			t.Logf("valid preparation %q created", mix.ID)

			t.Log("creating valid measurement units")
			exampleGrams := fakes.BuildFakeValidMeasurementUnit()
			exampleGramsInput := converters.ConvertValidMeasurementUnitToValidMeasurementUnitCreationRequestInput(exampleGrams)
			grams, err := testClients.admin.CreateValidMeasurementUnit(ctx, exampleGramsInput)
			require.NoError(t, err)
			t.Logf("valid measurement unit %q created", grams.ID)
			checkValidMeasurementUnitEquality(t, exampleGrams, grams)

			grams, err = testClients.admin.GetValidMeasurementUnit(ctx, grams.ID)
			requireNotNilAndNoProblems(t, grams, err)
			checkValidMeasurementUnitEquality(t, exampleGrams, grams)

			exampleCups := fakes.BuildFakeValidMeasurementUnit()
			exampleCupsInput := converters.ConvertValidMeasurementUnitToValidMeasurementUnitCreationRequestInput(exampleCups)
			cups, err := testClients.admin.CreateValidMeasurementUnit(ctx, exampleCupsInput)
			require.NoError(t, err)
			t.Logf("valid measurement unit %q created", cups.ID)
			checkValidMeasurementUnitEquality(t, exampleCups, cups)

			cups, err = testClients.admin.GetValidMeasurementUnit(ctx, cups.ID)
			requireNotNilAndNoProblems(t, cups, err)
			checkValidMeasurementUnitEquality(t, exampleCups, cups)

			t.Log("creating prerequisite valid ingredient")
			pintoBeanBase := fakes.BuildFakeValidIngredient()
			pintoBeanInput := converters.ConvertValidIngredientToValidIngredientCreationRequestInput(pintoBeanBase)
			pintoBeans, createdValidIngredientErr := testClients.admin.CreateValidIngredient(ctx, pintoBeanInput)
			require.NoError(t, createdValidIngredientErr)

			t.Log("creating prerequisite valid ingredient")
			waterBase := fakes.BuildFakeValidIngredient()
			waterInput := converters.ConvertValidIngredientToValidIngredientCreationRequestInput(waterBase)
			water, createdValidIngredientErr := testClients.admin.CreateValidIngredient(ctx, waterInput)
			require.NoError(t, createdValidIngredientErr)

			t.Log("creating prerequisite valid ingredient")
			garlicPaste := fakes.BuildFakeValidIngredient()
			garlicPasteInput := converters.ConvertValidIngredientToValidIngredientCreationRequestInput(garlicPaste)
			garlicPaste, garlicPasteErr := testClients.admin.CreateValidIngredient(ctx, garlicPasteInput)
			require.NoError(t, garlicPasteErr)

			t.Log("creating recipe")
			expected := &types.Recipe{
				Name:        "sopa de frijol",
				Description: "",
				Steps: []*types.RecipeStep{
					{
						MinimumTemperatureInCelsius: nil,
						Products: []*types.RecipeStepProduct{
							{
								Name:            "soaked pinto beans",
								Type:            types.RecipeStepProductIngredientType,
								MeasurementUnit: *grams,
								QuantityNotes:   "",
								MinimumQuantity: 1000,
							},
						},
						Notes:       "first step",
						Preparation: *soak,
						Ingredients: []*types.RecipeStepIngredient{
							{
								RecipeStepProductID: nil,
								Ingredient:          pintoBeans,
								Name:                "pinto beans",
								MeasurementUnit:     *grams,
								MinimumQuantity:     500,
								ProductOfRecipeStep: false,
							},
							{
								RecipeStepProductID: nil,
								Ingredient:          water,
								Name:                "water",
								MeasurementUnit:     *cups,
								MinimumQuantity:     5,
								ProductOfRecipeStep: false,
							},
						},
						Index: 0,
					},
					{
						MinimumTemperatureInCelsius: nil,
						Products: []*types.RecipeStepProduct{
							{
								Name:            "final output",
								Type:            types.RecipeStepProductIngredientType,
								MeasurementUnit: *grams,
								QuantityNotes:   "",
								MinimumQuantity: 1010,
							},
						},
						Notes:       "second step",
						Preparation: *mix,
						Ingredients: []*types.RecipeStepIngredient{
							{
								Name:                "soaked pinto beans",
								MeasurementUnit:     *grams,
								MinimumQuantity:     1000,
								ProductOfRecipeStep: true,
							},
							{
								Ingredient:          garlicPaste,
								Name:                "garlic paste",
								MeasurementUnit:     *grams,
								MinimumQuantity:     10,
								ProductOfRecipeStep: false,
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
					MinimumTemperatureInCelsius:   step.MinimumTemperatureInCelsius,
					Notes:                         step.Notes,
					PreparationID:                 step.Preparation.ID,
					BelongsToRecipe:               step.BelongsToRecipe,
					ID:                            step.ID,
					Index:                         step.Index,
					MinimumEstimatedTimeInSeconds: step.MinimumEstimatedTimeInSeconds,
					MaximumEstimatedTimeInSeconds: step.MaximumEstimatedTimeInSeconds,
					Optional:                      step.Optional,
				}

				for _, ingredient := range step.Ingredients {
					newIngredient := &types.RecipeStepIngredientCreationRequestInput{
						ID:                  ingredient.ID,
						BelongsToRecipeStep: ingredient.BelongsToRecipeStep,
						Name:                ingredient.Name,
						MeasurementUnitID:   ingredient.MeasurementUnit.ID,
						QuantityNotes:       ingredient.QuantityNotes,
						IngredientNotes:     ingredient.IngredientNotes,
						MinimumQuantity:     ingredient.MinimumQuantity,
						ProductOfRecipeStep: ingredient.ProductOfRecipeStep,
					}

					if ingredient.Ingredient != nil {
						newIngredient.IngredientID = &ingredient.Ingredient.ID
					}

					newStep.Ingredients = append(newStep.Ingredients, newIngredient)
				}

				for _, product := range step.Products {
					newProduct := &types.RecipeStepProductCreationRequestInput{
						ID:                  product.ID,
						Name:                product.Name,
						Type:                product.Type,
						MeasurementUnitID:   product.MeasurementUnit.ID,
						QuantityNotes:       product.QuantityNotes,
						BelongsToRecipeStep: product.BelongsToRecipeStep,
						MinimumQuantity:     product.MinimumQuantity,
					}
					newStep.Products = append(newStep.Products, newProduct)
				}

				exampleRecipeInput.Steps = append(exampleRecipeInput.Steps, newStep)
			}

			created, err := testClients.user.CreateRecipe(ctx, exampleRecipeInput)
			require.NoError(t, err)
			t.Logf("recipe %q created", created.ID)
			checkRecipeEquality(t, expected, created)

			created, err = testClients.user.GetRecipe(ctx, created.ID)
			requireNotNilAndNoProblems(t, created, err)
			checkRecipeEquality(t, expected, created)

			recipeStepProductIndex := -1
			for i, ingredient := range created.Steps[1].Ingredients {
				if ingredient.RecipeStepProductID != nil {
					recipeStepProductIndex = i
				}
			}

			require.NotEqual(t, -1, recipeStepProductIndex)
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

			_, _, createdRecipe := createRecipeForTest(ctx, t, testClients.admin, testClients.user, nil)

			t.Log("changing recipe")
			newRecipe := fakes.BuildFakeRecipe()
			createdRecipe.Update(converters.ConvertRecipeToRecipeUpdateRequestInput(newRecipe))
			assert.NoError(t, testClients.user.UpdateRecipe(ctx, createdRecipe))

			t.Log("fetching changed recipe")
			actual, err := testClients.user.GetRecipe(ctx, createdRecipe.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert recipe equality
			checkRecipeEquality(t, newRecipe, actual)
			assert.NotNil(t, actual.LastUpdatedAt)

			t.Log("cleaning up recipe")
			assert.NoError(t, testClients.user.ArchiveRecipe(ctx, createdRecipe.ID))
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
			exampleValidPreparationInput := converters.ConvertValidPreparationToValidPreparationCreationRequestInput(exampleValidPreparation)
			createdValidPreparation, err := testClients.admin.CreateValidPreparation(ctx, exampleValidPreparationInput)
			require.NoError(t, err)

			t.Log("creating valid instrument")
			exampleValidInstrument := fakes.BuildFakeValidInstrument()
			exampleValidInstrumentInput := converters.ConvertValidInstrumentToValidInstrumentCreationRequestInput(exampleValidInstrument)
			createdValidInstrument, err := testClients.admin.CreateValidInstrument(ctx, exampleValidInstrumentInput)
			require.NoError(t, err)
			t.Logf("valid instrument %q created", createdValidInstrument.ID)
			checkValidInstrumentEquality(t, exampleValidInstrument, createdValidInstrument)

			t.Log("creating valid measurement unit")
			exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()
			exampleValidMeasurementUnitInput := converters.ConvertValidMeasurementUnitToValidMeasurementUnitCreationRequestInput(exampleValidMeasurementUnit)
			createdValidMeasurementUnit, err := testClients.admin.CreateValidMeasurementUnit(ctx, exampleValidMeasurementUnitInput)
			require.NoError(t, err)
			t.Logf("valid measurement unit %q created", createdValidMeasurementUnit.ID)
			checkValidMeasurementUnitEquality(t, exampleValidMeasurementUnit, createdValidMeasurementUnit)

			createdValidMeasurementUnit, err = testClients.admin.GetValidMeasurementUnit(ctx, createdValidMeasurementUnit.ID)
			requireNotNilAndNoProblems(t, createdValidMeasurementUnit, err)
			checkValidMeasurementUnitEquality(t, exampleValidMeasurementUnit, createdValidMeasurementUnit)

			exampleRecipe := fakes.BuildFakeRecipe()
			createdValidIngredients := []*types.ValidIngredient{}
			for i, recipeStep := range exampleRecipe.Steps {
				for j := range recipeStep.Ingredients {
					t.Log("creating prerequisite valid ingredient")
					exampleValidIngredient := fakes.BuildFakeValidIngredient()
					exampleValidIngredientInput := converters.ConvertValidIngredientToValidIngredientCreationRequestInput(exampleValidIngredient)
					createdValidIngredient, createdValidIngredientErr := testClients.admin.CreateValidIngredient(ctx, exampleValidIngredientInput)
					require.NoError(t, createdValidIngredientErr)

					createdValidIngredients = append(createdValidIngredients, createdValidIngredient)

					exampleRecipe.Steps[i].Ingredients[j].Ingredient = createdValidIngredient
					exampleRecipe.Steps[i].Ingredients[j].ProductOfRecipeStep = false
					exampleRecipe.Steps[i].Ingredients[j].MeasurementUnit = *createdValidMeasurementUnit
				}

				for j := range recipeStep.Products {
					exampleRecipe.Steps[i].Products[j].MeasurementUnit = *createdValidMeasurementUnit
				}

				for j := range recipeStep.Instruments {
					recipeStep.Instruments[j].Instrument = createdValidInstrument
				}
			}

			exampleRecipeInput := converters.ConvertRecipeToRecipeCreationRequestInputFromRecipe(exampleRecipe)
			for i := range exampleRecipeInput.Steps {
				exampleRecipeInput.Steps[i].PreparationID = createdValidPreparation.ID
			}
			exampleRecipeInput.AlsoCreateMeal = true

			examplePrepTask := fakes.BuildFakeRecipePrepTask()
			examplePrepTask.TaskSteps = []*types.RecipePrepTaskStep{}
			exampleRecipeInput.PrepTasks = []*types.RecipePrepTaskWithinRecipeCreationRequestInput{
				converters.ConvertRecipePrepTaskToRecipePrepTaskWithinRecipeCreationRequestInput(exampleRecipe, examplePrepTask),
			}

			createdRecipe, err := testClients.user.CreateRecipe(ctx, exampleRecipeInput)
			require.NoError(t, err)
			t.Logf("recipe %q created", createdRecipe.ID)
			checkRecipeEquality(t, exampleRecipe, createdRecipe)

			mealResults, err := testClients.user.SearchForMeals(ctx, createdRecipe.Name, nil)
			requireNotNilAndNoProblems(t, mealResults, err)

			foundMealID := ""
			for _, m := range mealResults.Meals {
				meal, mealFetchErr := testClients.user.GetMeal(ctx, m.ID)
				requireNotNilAndNoProblems(t, meal, mealFetchErr)

				for _, r := range meal.Recipes {
					if r.ID == createdRecipe.ID {
						foundMealID = r.ID
					}
				}
			}

			require.NotEmpty(t, foundMealID)

			t.Log("cleaning up recipe")
			assert.NoError(t, testClients.user.ArchiveRecipe(ctx, createdRecipe.ID))
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
			exampleValidIngredientInput := converters.ConvertValidIngredientToValidIngredientCreationRequestInput(exampleValidIngredient)
			createdValidIngredient, err := testClients.admin.CreateValidIngredient(ctx, exampleValidIngredientInput)
			require.NoError(t, err)
			t.Logf("valid ingredient %q created", createdValidIngredient.ID)

			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

			createdValidIngredient, err = testClients.user.GetValidIngredient(ctx, createdValidIngredient.ID)
			requireNotNilAndNoProblems(t, createdValidIngredient, err)
			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

			t.Log("creating prerequisite valid preparation")
			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparationInput := converters.ConvertValidPreparationToValidPreparationCreationRequestInput(exampleValidPreparation)
			createdValidPreparation, err := testClients.admin.CreateValidPreparation(ctx, exampleValidPreparationInput)
			require.NoError(t, err)
			t.Logf("valid preparation %q created", createdValidPreparation.ID)

			checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

			createdValidPreparation, err = testClients.user.GetValidPreparation(ctx, createdValidPreparation.ID)
			requireNotNilAndNoProblems(t, createdValidPreparation, err)
			checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

			t.Log("creating recipes")
			var expected []*types.Recipe
			for i := 0; i < 5; i++ {
				_, _, createdRecipe := createRecipeForTest(ctx, t, testClients.admin, testClients.user, nil)

				expected = append(expected, createdRecipe)
			}

			// assert recipe list equality
			actual, err := testClients.user.GetRecipes(ctx, nil)
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
				assert.NoError(t, testClients.user.ArchiveRecipe(ctx, createdRecipe.ID))
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
			exampleValidIngredientInput := converters.ConvertValidIngredientToValidIngredientCreationRequestInput(exampleValidIngredient)
			createdValidIngredient, err := testClients.admin.CreateValidIngredient(ctx, exampleValidIngredientInput)
			require.NoError(t, err)
			t.Logf("valid ingredient %q created", createdValidIngredient.ID)

			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

			createdValidIngredient, err = testClients.user.GetValidIngredient(ctx, createdValidIngredient.ID)
			requireNotNilAndNoProblems(t, createdValidIngredient, err)
			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

			t.Log("creating prerequisite valid preparation")
			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparationInput := converters.ConvertValidPreparationToValidPreparationCreationRequestInput(exampleValidPreparation)
			createdValidPreparation, err := testClients.admin.CreateValidPreparation(ctx, exampleValidPreparationInput)
			require.NoError(t, err)
			t.Logf("valid preparation %q created", createdValidPreparation.ID)

			checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

			createdValidPreparation, err = testClients.user.GetValidPreparation(ctx, createdValidPreparation.ID)
			requireNotNilAndNoProblems(t, createdValidPreparation, err)
			checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

			exampleRecipe := fakes.BuildFakeRecipe()

			t.Log("creating recipes")
			var expected []*types.Recipe
			for i := 0; i < 5; i++ {
				exampleRecipe.Name = fmt.Sprintf("example%d", i)
				_, _, createdRecipe := createRecipeForTest(ctx, t, testClients.admin, testClients.user, exampleRecipe)

				expected = append(expected, createdRecipe)
			}

			// assert recipe list equality
			actual, err := testClients.user.SearchForRecipes(ctx, "example", nil)
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
				assert.NoError(t, testClients.user.ArchiveRecipe(ctx, createdRecipe.ID))
			}
		}
	})
}

func (s *TestSuite) TestRecipes_GetMealPlanTasksForRecipe() {
	s.runForEachClient("meal plan tasks with frozen chicken breast", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating prerequisite valid preparation")
			diceBase := fakes.BuildFakeValidPreparation()
			diceInput := converters.ConvertValidPreparationToValidPreparationCreationRequestInput(diceBase)
			dice, err := testClients.admin.CreateValidPreparation(ctx, diceInput)
			require.NoError(t, err)
			t.Logf("valid preparation %q created", dice.ID)

			t.Log("creating valid measurement units")
			exampleGrams := fakes.BuildFakeValidMeasurementUnit()
			exampleGramsInput := converters.ConvertValidMeasurementUnitToValidMeasurementUnitCreationRequestInput(exampleGrams)
			grams, err := testClients.admin.CreateValidMeasurementUnit(ctx, exampleGramsInput)
			require.NoError(t, err)
			t.Logf("valid measurement unit %q created", grams.ID)
			checkValidMeasurementUnitEquality(t, exampleGrams, grams)

			grams, err = testClients.admin.GetValidMeasurementUnit(ctx, grams.ID)
			requireNotNilAndNoProblems(t, grams, err)
			checkValidMeasurementUnitEquality(t, exampleGrams, grams)

			t.Log("creating prerequisite valid ingredient")
			chickenBreastBase := fakes.BuildFakeValidIngredient()
			chickenBreastInput := converters.ConvertValidIngredientToValidIngredientCreationRequestInput(chickenBreastBase)
			chickenBreastInput.MinimumIdealStorageTemperatureInCelsius = pointers.Float32Pointer(2.5)
			chickenBreast, createdValidIngredientErr := testClients.admin.CreateValidIngredient(ctx, chickenBreastInput)
			require.NoError(t, createdValidIngredientErr)

			t.Log("creating prerequisite valid preparation")
			sauteeBase := fakes.BuildFakeValidPreparation()
			sauteeInput := converters.ConvertValidPreparationToValidPreparationCreationRequestInput(sauteeBase)
			sautee, err := testClients.admin.CreateValidPreparation(ctx, sauteeInput)
			require.NoError(t, err)
			t.Logf("valid preparation %q created", sautee.ID)

			t.Log("creating recipe")
			expected := &types.Recipe{
				Name:        "sopa de frijol",
				Description: "",
				Steps: []*types.RecipeStep{
					{
						MinimumTemperatureInCelsius: nil,
						Products: []*types.RecipeStepProduct{
							{
								Name:            "diced chicken breast",
								Type:            types.RecipeStepProductIngredientType,
								MeasurementUnit: *grams,
								QuantityNotes:   "",
								MinimumQuantity: 1000,
							},
						},
						Notes:       "first step",
						Preparation: *dice,
						Ingredients: []*types.RecipeStepIngredient{
							{
								RecipeStepProductID: nil,
								Ingredient:          chickenBreast,
								Name:                "pinto beans",
								MeasurementUnit:     *grams,
								MinimumQuantity:     500,
								ProductOfRecipeStep: false,
							},
						},
						Index: 0,
					},
					{
						MinimumTemperatureInCelsius: nil,
						Products: []*types.RecipeStepProduct{
							{
								Name:            "final output",
								Type:            types.RecipeStepProductIngredientType,
								MeasurementUnit: *grams,
								QuantityNotes:   "",
								MinimumQuantity: 1010,
							},
						},
						Notes:       "second step",
						Preparation: *sautee,
						Ingredients: []*types.RecipeStepIngredient{
							{
								Name:                "diced chicken breast",
								MeasurementUnit:     *grams,
								MinimumQuantity:     1000,
								ProductOfRecipeStep: true,
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
					MinimumTemperatureInCelsius:   step.MinimumTemperatureInCelsius,
					Notes:                         step.Notes,
					PreparationID:                 step.Preparation.ID,
					BelongsToRecipe:               step.BelongsToRecipe,
					ID:                            step.ID,
					Index:                         step.Index,
					MinimumEstimatedTimeInSeconds: step.MinimumEstimatedTimeInSeconds,
					MaximumEstimatedTimeInSeconds: step.MaximumEstimatedTimeInSeconds,
					Optional:                      step.Optional,
				}

				for _, ingredient := range step.Ingredients {
					newIngredient := &types.RecipeStepIngredientCreationRequestInput{
						ID:                  ingredient.ID,
						BelongsToRecipeStep: ingredient.BelongsToRecipeStep,
						Name:                ingredient.Name,
						MeasurementUnitID:   ingredient.MeasurementUnit.ID,
						QuantityNotes:       ingredient.QuantityNotes,
						IngredientNotes:     ingredient.IngredientNotes,
						MinimumQuantity:     ingredient.MinimumQuantity,
						ProductOfRecipeStep: ingredient.ProductOfRecipeStep,
					}

					if ingredient.Ingredient != nil {
						newIngredient.IngredientID = &ingredient.Ingredient.ID
					}

					newStep.Ingredients = append(newStep.Ingredients, newIngredient)
				}

				for _, product := range step.Products {
					newProduct := &types.RecipeStepProductCreationRequestInput{
						ID:                  product.ID,
						Name:                product.Name,
						Type:                product.Type,
						MeasurementUnitID:   product.MeasurementUnit.ID,
						QuantityNotes:       product.QuantityNotes,
						BelongsToRecipeStep: product.BelongsToRecipeStep,
						MinimumQuantity:     product.MinimumQuantity,
					}
					newStep.Products = append(newStep.Products, newProduct)
				}

				exampleRecipeInput.Steps = append(exampleRecipeInput.Steps, newStep)
			}

			created, err := testClients.user.CreateRecipe(ctx, exampleRecipeInput)
			require.NoError(t, err)
			t.Logf("recipe %q created", created.ID)
			checkRecipeEquality(t, expected, created)

			steps, err := testClients.user.GetMealPlanTasksForRecipe(ctx, created.ID)
			requireNotNilAndNoProblems(t, created, err)

			require.NotEmpty(t, steps)
		}
	})
}

func (s *TestSuite) TestRecipes_DAGGeneration() {
	s.runForEachClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			_, _, createdRecipe := createRecipeForTest(ctx, t, testClients.admin, testClients.user, nil)

			t.Log("fetching changed recipe")
			actual, err := testClients.user.GetRecipeDAG(ctx, createdRecipe.ID)
			requireNotNilAndNoProblems(t, actual, err)

			t.Log("cleaning up recipe")
			assert.NoError(t, testClients.user.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}
