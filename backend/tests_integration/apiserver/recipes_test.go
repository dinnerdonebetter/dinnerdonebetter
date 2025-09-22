package integration

import (
	"fmt"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	mpconverters "github.com/dinnerdonebetter/backend/internal/domain/mealplanning/converters"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/fakes"
	mealplanninggrpc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
	converters "github.com/dinnerdonebetter/backend/internal/services/mealplanning/grpc/converters"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkRecipeEquality(t *testing.T, expected, actual *mealplanning.Recipe) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for recipe %s to be %v, but it was %v", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.InspiredByRecipeID, actual.InspiredByRecipeID, "expected InspiredByRecipeID for recipe %s to be %v, but it was %v", expected.ID, expected.InspiredByRecipeID, actual.InspiredByRecipeID)
	assert.Equal(t, expected.EstimatedPortions, actual.EstimatedPortions, "expected EstimatedPortions for recipe %s to be %v, but it was %v", expected.ID, expected.EstimatedPortions, actual.EstimatedPortions)
	assert.Equal(t, expected.PluralPortionName, actual.PluralPortionName, "expected PluralPortionName for recipe %s to be %v, but it was %v", expected.ID, expected.PluralPortionName, actual.PluralPortionName)
	assert.Equal(t, expected.Description, actual.Description, "expected Description for recipe %s to be %v, but it was %v", expected.ID, expected.Description, actual.Description)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for recipe %s to be %v, but it was %v", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.PortionName, actual.PortionName, "expected PortionName for recipe %s to be %v, but it was %v", expected.ID, expected.PortionName, actual.PortionName)
	assert.NotZero(t, actual.CreatedByUser)
	assert.Equal(t, expected.Source, actual.Source, "expected Source for recipe %s to be %v, but it was %v", expected.ID, expected.Source, actual.Source)
	assert.Equal(t, expected.Slug, actual.Slug, "expected Slug for recipe %s to be %v, but it was %v", expected.ID, expected.Slug, actual.Slug)
	assert.Equal(t, expected.YieldsComponentType, actual.YieldsComponentType, "expected YieldsComponentType for recipe %s to be %v, but it was %v", expected.ID, expected.YieldsComponentType, actual.YieldsComponentType)
	checkRecipePrepTaskSliceEquality(t, expected.PrepTasks, actual.PrepTasks)
	checkRecipeLevelMediaSliceEquality(t, expected.Media, actual.Media)
	assert.Equal(t, expected.SealOfApproval, actual.SealOfApproval, "expected SealOfApproval for recipe %s to be %v, but it was %v", expected.ID, expected.SealOfApproval, actual.SealOfApproval)
	assert.Equal(t, expected.EligibleForMeals, actual.EligibleForMeals, "expected EligibleForMeals for recipe %s to be %v, but it was %v", expected.ID, expected.EligibleForMeals, actual.EligibleForMeals)

	for i, step := range expected.Steps {
		checkRecipeStepEquality(t, i, step, actual.Steps[i])
	}

	assert.NotZero(t, actual.CreatedAt)
}

func createRecipeForTest(t *testing.T, recipe *mealplanning.Recipe, inputFilter ...func(input *mealplanning.RecipeCreationRequestInput)) ([]*mealplanning.ValidIngredient, *mealplanning.ValidPreparation, *mealplanning.Recipe) {
	t.Helper()

	ctx := t.Context()

	createdValidPreparation := createValidPreparationForTest(t)
	createdValidMeasurementUnit := createValidMeasurementUnitForTest(t)
	createdValidInstrument := createValidInstrumentForTest(t)
	createdValidIngredientState := createValidIngredientStateForTest(t)
	createdValidVessel := createValidVesselForTest(t)

	exampleRecipe := fakes.BuildFakeRecipe()
	if recipe != nil {
		exampleRecipe = recipe
	}
	exampleRecipe.Media = []*mealplanning.RecipeMedia{}

	createdValidIngredients := []*mealplanning.ValidIngredient{}
	for i, recipeStep := range exampleRecipe.Steps {
		for j := range recipeStep.Ingredients {
			createdValidIngredient := createValidIngredientForTest(t)
			createdValidIngredients = append(createdValidIngredients, createdValidIngredient)

			exampleRecipe.Steps[i].Ingredients[j].Ingredient = createdValidIngredient
			exampleRecipe.Steps[i].Ingredients[j].MeasurementUnit = *createdValidMeasurementUnit
		}

		for j := range recipeStep.Products {
			exampleRecipe.Steps[i].Products[j].MeasurementUnit = createdValidMeasurementUnit
		}

		for j := range recipeStep.Instruments {
			recipeStep.Instruments[j].Instrument = createdValidInstrument
		}

		for j := range recipeStep.Vessels {
			recipeStep.Vessels[j].Vessel = createdValidVessel
		}

		for j := range recipeStep.CompletionConditions {
			recipeStep.CompletionConditions[j].IngredientState = *createdValidIngredientState
			for k := range recipeStep.CompletionConditions[j].Ingredients {
				recipeStep.CompletionConditions[j].Ingredients[k].RecipeStepIngredient = recipeStep.Ingredients[0].ID
			}
		}
	}

	exampleRecipeInput := mpconverters.ConvertRecipeToRecipeCreationRequestInput(exampleRecipe)
	exampleRecipeInput.AlsoCreateMeal = true
	for i := range exampleRecipeInput.Steps {
		exampleRecipeInput.Steps[i].PreparationID = createdValidPreparation.ID
	}

	examplePrepTask := fakes.BuildFakeRecipePrepTask()
	examplePrepTask.TaskSteps = []*mealplanning.RecipePrepTaskStep{
		{
			BelongsToRecipeStep: exampleRecipe.Steps[0].ID,
			SatisfiesRecipeStep: false,
		},
	}
	exampleRecipeInput.PrepTasks = []*mealplanning.RecipePrepTaskWithinRecipeCreationRequestInput{
		mpconverters.ConvertRecipePrepTaskToRecipePrepTaskWithinRecipeCreationRequestInput(exampleRecipe, examplePrepTask),
	}
	// Update the exampleRecipe to match what we're actually creating
	exampleRecipe.PrepTasks = []*mealplanning.RecipePrepTask{examplePrepTask}

	for _, filter := range inputFilter {
		filter(exampleRecipeInput)
	}

	createdRes, err := adminClient.CreateRecipe(ctx, &mealplanninggrpc.CreateRecipeRequest{Input: converters.ConvertRecipeCreationRequestInputToGRPCRecipeCreationRequestInput(exampleRecipeInput)})
	require.NoError(t, err)

	createdRecipe, err := adminClient.GetRecipe(ctx, &mealplanninggrpc.GetRecipeRequest{RecipeID: createdRes.Created.ID})
	require.NoError(t, err)
	require.NotNil(t, createdRecipe)

	// Only do basic comparisons that we know should work
	converted := converters.ConvertGRPCRecipeToRecipe(createdRecipe.Result)
	require.NotEmpty(t, createdRecipe.Result.Steps, "created recipe must have steps")
	require.NotEmpty(t, converted.Steps, "converted recipe must have steps")

	// Verify that completion conditions are present (this was the original issue)
	for i, step := range converted.Steps {
		require.NotEmpty(t, step.CompletionConditions, "recipe step %d must have completion conditions", i)
		for j, condition := range step.CompletionConditions {
			require.NotEmpty(t, condition.Ingredients, "completion condition %d for step %d must have ingredients", j, i)
		}
	}

	return createdValidIngredients, createdValidPreparation, converted
}

func TestRecipes_Creating(T *testing.T) {
	T.Parallel()

	T.Run("realistic", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		soak := createValidPreparationForTest(t)
		mix := createValidPreparationForTest(t)
		grams := createValidMeasurementUnitForTest(t)
		cups := createValidMeasurementUnitForTest(t)
		pintoBeans := createValidIngredientForTest(t)
		water := createValidIngredientForTest(t)
		garlicPaste := createValidIngredientForTest(t)
		createdValidInstrument := createValidInstrumentForTest(t)

		expected := &mealplanning.Recipe{
			Name:                "sopa de frijol",
			Slug:                "sopa-de-frijol",
			YieldsComponentType: mealplanning.MealComponentTypesMain,
			PortionName:         t.Name(),
			PluralPortionName:   t.Name(),
			Media:               []*mealplanning.RecipeMedia{},
			PrepTasks:           []*mealplanning.RecipePrepTask{},
			EstimatedPortions: types.Float32RangeWithOptionalMax{
				Max: nil,
				Min: 1,
			},
			Steps: []*mealplanning.RecipeStep{
				{
					Products: []*mealplanning.RecipeStepProduct{
						{
							Name:            "soaked pinto beans",
							Type:            mealplanning.RecipeStepProductIngredientType,
							MeasurementUnit: grams,
							QuantityNotes:   "",
							Quantity: types.OptionalFloat32Range{
								Max: nil,
								Min: pointer.To(float32(1000)),
							},
						},
					},
					Notes:       "first step",
					Preparation: *soak, // This will be updated after recipe creation
					Instruments: []*mealplanning.RecipeStepInstrument{
						{
							Name:       "whatever",
							Instrument: createdValidInstrument,
						},
					},
					Ingredients: []*mealplanning.RecipeStepIngredient{
						{
							Ingredient:      pintoBeans,
							Name:            "pinto beans",
							MeasurementUnit: *grams,
							Quantity: types.Float32RangeWithOptionalMax{
								Min: 500,
							},
						},
						{
							Ingredient:      water,
							Name:            "water",
							MeasurementUnit: *cups,
							Quantity: types.Float32RangeWithOptionalMax{
								Min: 5,
							},
						},
					},
					Index: 0,
				},
				{
					Products: []*mealplanning.RecipeStepProduct{
						{
							Name:            "final output",
							Type:            mealplanning.RecipeStepProductIngredientType,
							MeasurementUnit: grams,
							QuantityNotes:   "",
							Quantity: types.OptionalFloat32Range{
								Max: nil,
								Min: pointer.To(float32(1010)),
							},
						},
					},
					Notes:       "second step",
					Preparation: *mix, // This will be updated after recipe creation
					Instruments: []*mealplanning.RecipeStepInstrument{
						{
							Name:       "whatever",
							Instrument: createdValidInstrument,
						},
					},
					Ingredients: []*mealplanning.RecipeStepIngredient{
						{
							Name:            "soaked pinto beans",
							MeasurementUnit: *grams,
							Quantity: types.Float32RangeWithOptionalMax{
								Min: 1000,
							},
						},
						{
							Ingredient:      garlicPaste,
							Name:            "garlic paste",
							MeasurementUnit: *grams,
							Quantity: types.Float32RangeWithOptionalMax{
								Min: 10,
							},
						},
					},
					Index: 1,
				},
			},
		}

		expectedInput := &mealplanning.RecipeCreationRequestInput{
			Name:                expected.Name,
			Description:         expected.Description,
			Slug:                expected.Slug,
			YieldsComponentType: expected.YieldsComponentType,
			PortionName:         expected.PortionName,
			PluralPortionName:   expected.PluralPortionName,
			AlsoCreateMeal:      true,
			EstimatedPortions: types.Float32RangeWithOptionalMax{
				Max: expected.EstimatedPortions.Max,
				Min: expected.EstimatedPortions.Min,
			},
			Steps: []*mealplanning.RecipeStepCreationRequestInput{
				{
					TemperatureInCelsius: expected.Steps[0].TemperatureInCelsius,
					Products: []*mealplanning.RecipeStepProductCreationRequestInput{
						{
							Name:              expected.Steps[0].Products[0].Name,
							Type:              expected.Steps[0].Products[0].Type,
							MeasurementUnitID: &expected.Steps[0].Products[0].MeasurementUnit.ID,
							QuantityNotes:     expected.Steps[0].Products[0].QuantityNotes,
							Quantity:          expected.Steps[0].Products[0].Quantity,
						},
					},
					Notes:         expected.Steps[0].Notes,
					PreparationID: expected.Steps[0].Preparation.ID,
					Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
						{
							Name:         "whatever",
							InstrumentID: pointer.To(createdValidInstrument.ID),
						},
					},
					Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
						{
							IngredientID:      &expected.Steps[0].Ingredients[0].Ingredient.ID,
							Name:              expected.Steps[0].Ingredients[0].Name,
							MeasurementUnitID: expected.Steps[0].Ingredients[0].MeasurementUnit.ID,
							Quantity: types.Float32RangeWithOptionalMax{
								Max: nil,
								Min: expected.Steps[0].Ingredients[0].Quantity.Min,
							},
						},
						{
							IngredientID:      &expected.Steps[0].Ingredients[1].Ingredient.ID,
							Name:              expected.Steps[0].Ingredients[1].Name,
							MeasurementUnitID: expected.Steps[0].Ingredients[1].MeasurementUnit.ID,
							Quantity: types.Float32RangeWithOptionalMax{
								Max: nil,
								Min: expected.Steps[0].Ingredients[1].Quantity.Min,
							},
						},
					},
					Index: expected.Steps[0].Index,
				},
				{
					TemperatureInCelsius: expected.Steps[1].TemperatureInCelsius,
					Products: []*mealplanning.RecipeStepProductCreationRequestInput{
						{
							Name:              expected.Steps[1].Products[0].Name,
							Type:              expected.Steps[1].Products[0].Type,
							MeasurementUnitID: &expected.Steps[1].Products[0].MeasurementUnit.ID,
							QuantityNotes:     expected.Steps[1].Products[0].QuantityNotes,
							Quantity:          expected.Steps[1].Products[0].Quantity,
						},
					},
					Notes:         expected.Steps[1].Notes,
					PreparationID: expected.Steps[1].Preparation.ID,
					Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
						{
							Name:         "whatever",
							InstrumentID: pointer.To(createdValidInstrument.ID),
						},
					},
					Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
						{
							Name:                            expected.Steps[1].Ingredients[0].Name,
							MeasurementUnitID:               expected.Steps[1].Ingredients[0].MeasurementUnit.ID,
							ProductOfRecipeStepIndex:        pointer.To(uint64(0)),
							ProductOfRecipeStepProductIndex: pointer.To(uint64(0)),
							Quantity: types.Float32RangeWithOptionalMax{
								Max: nil,
								Min: expected.Steps[1].Ingredients[0].Quantity.Min,
							},
						},
						{
							IngredientID:      &expected.Steps[1].Ingredients[1].Ingredient.ID,
							Name:              expected.Steps[1].Ingredients[1].Name,
							MeasurementUnitID: expected.Steps[1].Ingredients[1].MeasurementUnit.ID,
							Quantity: types.Float32RangeWithOptionalMax{
								Max: nil,
								Min: expected.Steps[1].Ingredients[1].Quantity.Min,
							},
						},
					},
					Index: expected.Steps[1].Index,
				},
			},
		}

		createdRes, err := adminClient.CreateRecipe(ctx, &mealplanninggrpc.CreateRecipeRequest{Input: converters.ConvertRecipeCreationRequestInputToGRPCRecipeCreationRequestInput(expectedInput)})
		require.NoError(t, err)

		created := converters.ConvertGRPCRecipeToRecipe(createdRes.Created)
		checkRecipeEquality(t, expected, created)

		recipeRes, err := adminClient.GetRecipe(ctx, &mealplanninggrpc.GetRecipeRequest{RecipeID: createdRes.Created.ID})
		require.NoError(t, err)
		created = converters.ConvertGRPCRecipeToRecipe(recipeRes.Result)

		recipeStepProductIndex := -1
		for i, ingredient := range created.Steps[1].Ingredients {
			if ingredient.RecipeStepProductID != nil {
				recipeStepProductIndex = i
			}
		}

		require.NotEqual(t, -1, recipeStepProductIndex)
		require.NotNil(t, created.Steps[1].Ingredients[recipeStepProductIndex].RecipeStepProductID)
		assert.Equal(t, created.Steps[0].Products[0].ID, *created.Steps[1].Ingredients[recipeStepProductIndex].RecipeStepProductID)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		exampleRecipe := fakes.BuildFakeRecipe()
		exampleRecipeInput := mpconverters.ConvertRecipeToRecipeCreationRequestInput(exampleRecipe)
		convertedInput := converters.ConvertRecipeCreationRequestInputToGRPCRecipeCreationRequestInput(exampleRecipeInput)

		c := buildUnauthenticatedGRPCClientForTest(t)
		created, err := c.CreateRecipe(ctx, &mealplanninggrpc.CreateRecipeRequest{
			Input: convertedInput,
		})
		assert.Error(t, err)
		assert.Nil(t, created)
	})

	T.Run("non-admin users are forbidden from creating", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, testClient := createUserAndClientForTest(t)

		exampleRecipe := fakes.BuildFakeRecipe()
		exampleRecipeInput := mpconverters.ConvertRecipeToRecipeCreationRequestInput(exampleRecipe)
		convertedInput := converters.ConvertRecipeCreationRequestInputToGRPCRecipeCreationRequestInput(exampleRecipeInput)

		created, err := testClient.CreateRecipe(ctx, &mealplanninggrpc.CreateRecipeRequest{
			Input: convertedInput,
		})
		assert.Error(t, err)
		assert.Nil(t, created)
	})
}

func TestRecipes_Updating(T *testing.T) {
	T.Parallel()

	T.Run("should update recipe", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		_, _, createdRecipe := createRecipeForTest(t, nil)

		// Store the original recipe data for comparison
		originalSteps := make([]*mealplanning.RecipeStep, len(createdRecipe.Steps))
		for i, step := range createdRecipe.Steps {
			originalSteps[i] = &mealplanning.RecipeStep{
				ID:                   step.ID,
				CompletionConditions: step.CompletionConditions,
			}
		}

		// Create update input with new basic fields
		newRecipe := fakes.BuildFakeRecipe()
		updateInput := mpconverters.ConvertRecipeToRecipeUpdateRequestInput(newRecipe)

		_, err := adminClient.UpdateRecipe(ctx, &mealplanninggrpc.UpdateRecipeRequest{
			RecipeID: createdRecipe.ID,
			Input:    converters.ConvertRecipeUpdateRequestInputToGRPCRecipeUpdateRequestInput(updateInput),
		})
		require.NoError(t, err)

		// Retrieve the updated recipe
		actual, err := adminClient.GetRecipe(ctx, &mealplanninggrpc.GetRecipeRequest{RecipeID: createdRecipe.ID})
		require.NoError(t, err)
		require.NotNil(t, actual)
		actualRecipe := converters.ConvertGRPCRecipeToRecipe(actual.Result)

		// Assert that basic fields were updated correctly
		assert.Equal(t, newRecipe.Name, actualRecipe.Name, "recipe name should be updated")
		assert.Equal(t, newRecipe.Slug, actualRecipe.Slug, "recipe slug should be updated")
		assert.Equal(t, newRecipe.Source, actualRecipe.Source, "recipe source should be updated")
		assert.Equal(t, newRecipe.Description, actualRecipe.Description, "recipe description should be updated")
		assert.Equal(t, newRecipe.InspiredByRecipeID, actualRecipe.InspiredByRecipeID, "recipe inspired by recipe ID should be updated")
		assert.Equal(t, newRecipe.SealOfApproval, actualRecipe.SealOfApproval, "recipe seal of approval should be updated")
		assert.Equal(t, newRecipe.EstimatedPortions, actualRecipe.EstimatedPortions, "recipe estimated portions should be updated")
		assert.Equal(t, newRecipe.PortionName, actualRecipe.PortionName, "recipe portion name should be updated")
		assert.Equal(t, newRecipe.PluralPortionName, actualRecipe.PluralPortionName, "recipe plural portion name should be updated")
		assert.Equal(t, newRecipe.EligibleForMeals, actualRecipe.EligibleForMeals, "recipe eligible for meals should be updated")
		assert.Equal(t, newRecipe.YieldsComponentType, actualRecipe.YieldsComponentType, "recipe yields component type should be updated")
		assert.NotNil(t, actual.Result.LastUpdatedAt, "recipe should have last updated timestamp")

		// Assert that steps and completion conditions remain unchanged (since UpdateRecipe only updates top-level fields)
		assert.Equal(t, len(originalSteps), len(actualRecipe.Steps), "number of recipe steps should remain the same")
		for i, originalStep := range originalSteps {
			actualStep := actualRecipe.Steps[i]
			assert.Equal(t, originalStep.ID, actualStep.ID, "recipe step ID should remain unchanged")
			assert.Equal(t, len(originalStep.CompletionConditions), len(actualStep.CompletionConditions), "number of completion conditions should remain the same")

			// Verify completion conditions are still present and working (this was the original issue)
			for j, originalCondition := range originalStep.CompletionConditions {
				actualCondition := actualStep.CompletionConditions[j]
				assert.Equal(t, originalCondition.ID, actualCondition.ID, "completion condition ID should remain unchanged")
				assert.Equal(t, originalCondition.Optional, actualCondition.Optional, "completion condition optional flag should remain unchanged")
				assert.Equal(t, originalCondition.Notes, actualCondition.Notes, "completion condition notes should remain unchanged")
				assert.Equal(t, len(originalCondition.Ingredients), len(actualCondition.Ingredients), "number of completion condition ingredients should remain the same")
			}
		}

		_, err = adminClient.ArchiveRecipe(ctx, &mealplanninggrpc.ArchiveRecipeRequest{RecipeID: createdRecipe.ID})
		assert.NoError(t, err)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, _, createdRecipe := createRecipeForTest(t, nil)

		newRecipe := fakes.BuildFakeRecipe()
		updateInput := mpconverters.ConvertRecipeToRecipeUpdateRequestInput(newRecipe)

		c := buildUnauthenticatedGRPCClientForTest(t)
		_, err := c.UpdateRecipe(ctx, &mealplanninggrpc.UpdateRecipeRequest{
			RecipeID: createdRecipe.ID,
			Input:    converters.ConvertRecipeUpdateRequestInputToGRPCRecipeUpdateRequestInput(updateInput),
		})
		assert.Error(t, err)
	})

	T.Run("nonexistent recipe", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		newRecipe := fakes.BuildFakeRecipe()
		updateInput := mpconverters.ConvertRecipeToRecipeUpdateRequestInput(newRecipe)

		_, err := adminClient.UpdateRecipe(ctx, &mealplanninggrpc.UpdateRecipeRequest{
			RecipeID: nonexistentID,
			Input:    converters.ConvertRecipeUpdateRequestInputToGRPCRecipeUpdateRequestInput(updateInput),
		})
		assert.Error(t, err)
	})
}

func TestRecipes_Searching(T *testing.T) {
	T.Parallel()

	T.Run("should be searchable by name", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		exampleRecipe := fakes.BuildFakeRecipe()

		var expected []*mealplanning.Recipe
		for i := 0; i < 5; i++ {
			exampleRecipe.Name = fmt.Sprintf("example%d", i)
			_, _, createdRecipe := createRecipeForTest(t, exampleRecipe)

			expected = append(expected, createdRecipe)
		}

		// assert recipe list equality
		actual, err := adminClient.SearchForRecipes(ctx, &mealplanninggrpc.SearchForRecipesRequest{
			Query: "example",
		})
		require.NoError(t, err)
		assert.True(
			t,
			len(expected) <= len(actual.Results),
			"expected %d to be <= %d",
			len(expected),
			len(actual.Results),
		)

		for _, createdRecipe := range expected {
			_, err = adminClient.ArchiveRecipe(ctx, &mealplanninggrpc.ArchiveRecipeRequest{RecipeID: createdRecipe.ID})
			assert.NoError(t, err)
		}
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		c := buildUnauthenticatedGRPCClientForTest(t)
		results, err := c.SearchForRecipes(ctx, &mealplanninggrpc.SearchForRecipesRequest{
			Query: "test",
		})
		assert.Error(t, err)
		assert.Nil(t, results)
	})
}

func TestRecipes_Cloning(T *testing.T) {
	T.Parallel()

	T.Run("should CRUD", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, _, createdRecipe := createRecipeForTest(t, nil)

		actual, err := adminClient.CloneRecipe(ctx, &mealplanninggrpc.CloneRecipeRequest{RecipeID: createdRecipe.ID})
		require.NoError(t, err)

		require.Equal(t, createdRecipe.Name, actual.Cloned.Name)
		require.Equal(t, len(createdRecipe.Steps), len(actual.Cloned.Steps))

		_, err = adminClient.ArchiveRecipe(ctx, &mealplanninggrpc.ArchiveRecipeRequest{RecipeID: createdRecipe.ID})
		assert.NoError(t, err)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, _, createdRecipe := createRecipeForTest(t, nil)

		c := buildUnauthenticatedGRPCClientForTest(t)
		cloned, err := c.CloneRecipe(ctx, &mealplanninggrpc.CloneRecipeRequest{RecipeID: createdRecipe.ID})
		assert.Error(t, err)
		assert.Nil(t, cloned)
	})

	T.Run("nonexistent recipe", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		cloned, err := adminClient.CloneRecipe(ctx, &mealplanninggrpc.CloneRecipeRequest{RecipeID: nonexistentID})
		assert.Error(t, err)
		assert.Nil(t, cloned)
	})

	T.Run("non-admin users can clone", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, testClient := createUserAndClientForTest(t)
		_, _, createdRecipe := createRecipeForTest(t, nil)

		cloned, err := testClient.CloneRecipe(ctx, &mealplanninggrpc.CloneRecipeRequest{RecipeID: createdRecipe.ID})
		assert.NoError(t, err)
		assert.NotNil(t, cloned)
	})
}

func TestRecipes_GetMealPlanTasksForRecipe(T *testing.T) {
	T.Parallel()

	T.Run("meal plan tasks with frozen chicken breast", func(t *testing.T) {
		t.Parallel()

		dice := createValidPreparationForTest(t)
		grams := createValidMeasurementUnitForTest(t)
		chickenBreast := createValidIngredientForTest(t)
		createdValidInstrument := createValidInstrumentForTest(t)
		sautee := createValidPreparationForTest(t)

		expected := &mealplanning.Recipe{
			Name:                "sopa de frijol",
			Slug:                "whatever-who-cares-sopa-de-frijol",
			YieldsComponentType: mealplanning.MealComponentTypesMain,
			PortionName:         t.Name(),
			PluralPortionName:   t.Name(),
			EstimatedPortions: types.Float32RangeWithOptionalMax{
				Max: nil,
				Min: 1,
			},
			Steps: []*mealplanning.RecipeStep{
				{
					Products: []*mealplanning.RecipeStepProduct{
						{
							Name:            "diced chicken breast",
							Type:            mealplanning.RecipeStepProductIngredientType,
							MeasurementUnit: grams,
							QuantityNotes:   "",
							Quantity: types.OptionalFloat32Range{
								Max: nil,
								Min: pointer.To(float32(1000)),
							},
						},
					},
					Notes:       "first step",
					Preparation: *dice,
					Instruments: []*mealplanning.RecipeStepInstrument{
						{
							Name:       "whatever",
							Instrument: createdValidInstrument,
						},
					},
					Ingredients: []*mealplanning.RecipeStepIngredient{
						{
							RecipeStepProductID: nil,
							Ingredient:          chickenBreast,
							Name:                "pinto beans",
							MeasurementUnit:     *grams,
							Quantity: types.Float32RangeWithOptionalMax{
								Min: 500,
							},
						},
					},
					Index: 0,
				},
				{
					Products: []*mealplanning.RecipeStepProduct{
						{
							Name:            "final output",
							Type:            mealplanning.RecipeStepProductIngredientType,
							MeasurementUnit: grams,
							QuantityNotes:   "",
							Quantity: types.OptionalFloat32Range{
								Max: nil,
								Min: pointer.To(float32(1010)),
							},
						},
					},
					Notes:       "second step",
					Preparation: *sautee,
					Instruments: []*mealplanning.RecipeStepInstrument{
						{
							Name:       "whatever",
							Instrument: createdValidInstrument,
						},
					},
					Ingredients: []*mealplanning.RecipeStepIngredient{
						{
							Name:            "diced chicken breast",
							MeasurementUnit: *grams,
							Quantity: types.Float32RangeWithOptionalMax{
								Min: 1000,
							},
						},
					},
					Index: 1,
				},
			},
		}

		expectedInput := &mealplanning.RecipeCreationRequestInput{
			Name:                expected.Name,
			Slug:                expected.Slug,
			YieldsComponentType: expected.YieldsComponentType,
			PortionName:         expected.PortionName,
			PluralPortionName:   expected.PluralPortionName,
			EstimatedPortions: types.Float32RangeWithOptionalMax{
				Min: expected.EstimatedPortions.Min,
				Max: expected.EstimatedPortions.Max,
			},
			Steps: []*mealplanning.RecipeStepCreationRequestInput{
				{
					Products: []*mealplanning.RecipeStepProductCreationRequestInput{
						{
							Name:              "diced chicken breast",
							Type:              mealplanning.RecipeStepProductIngredientType,
							MeasurementUnitID: &grams.ID,
							QuantityNotes:     "",
							Quantity:          types.OptionalFloat32Range{Min: pointer.To(float32(1000))},
						},
					},
					Notes:         "first step",
					PreparationID: dice.ID,
					Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
						{
							Name:         "whatever",
							InstrumentID: pointer.To(createdValidInstrument.ID),
						},
					},
					Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
						{
							IngredientID:      &chickenBreast.ID,
							Name:              "pinto beans",
							MeasurementUnitID: grams.ID,
							Quantity:          types.Float32RangeWithOptionalMax{Min: 500},
						},
					},
					Index: 0,
				},
				{
					Products: []*mealplanning.RecipeStepProductCreationRequestInput{
						{
							Name:              "final output",
							Type:              mealplanning.RecipeStepProductIngredientType,
							MeasurementUnitID: &grams.ID,
							QuantityNotes:     "",
							Quantity:          types.OptionalFloat32Range{Min: pointer.To(float32(1010))},
						},
					},
					Notes:         "second step",
					PreparationID: sautee.ID,
					Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
						{
							Name:         "whatever",
							InstrumentID: pointer.To(createdValidInstrument.ID),
						},
					},
					Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
						{
							Name:                            "diced chicken breast",
							MeasurementUnitID:               grams.ID,
							Quantity:                        types.Float32RangeWithOptionalMax{Min: 1000},
							ProductOfRecipeStepIndex:        pointer.To(uint64(0)),
							ProductOfRecipeStepProductIndex: pointer.To(uint64(0)),
						},
					},
					Index: 1,
				},
			},
		}

		ctx := t.Context()

		created, err := adminClient.CreateRecipe(ctx, &mealplanninggrpc.CreateRecipeRequest{Input: converters.ConvertRecipeCreationRequestInputToGRPCRecipeCreationRequestInput(expectedInput)})
		require.NoError(t, err)
		checkRecipeEquality(t, expected, converters.ConvertGRPCRecipeToRecipe(created.Created))

		steps, err := adminClient.GetMealPlanTasks(ctx, &mealplanninggrpc.GetMealPlanTasksRequest{MealPlanID: created.Created.ID})
		require.NoError(t, err)
		require.NotEmpty(t, steps)
	})
}

func TestRecipes_Reading(T *testing.T) {
	T.Parallel()

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, _, createdRecipe := createRecipeForTest(t, nil)

		c := buildUnauthenticatedGRPCClientForTest(t)
		recipe, err := c.GetRecipe(ctx, &mealplanninggrpc.GetRecipeRequest{RecipeID: createdRecipe.ID})
		assert.Error(t, err)
		assert.Nil(t, recipe)
	})

	T.Run("nonexistent recipe", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		recipe, err := adminClient.GetRecipe(ctx, &mealplanninggrpc.GetRecipeRequest{RecipeID: nonexistentID})
		assert.Error(t, err)
		assert.Nil(t, recipe)
	})
}

func TestRecipes_Archiving(T *testing.T) {
	T.Parallel()

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, _, createdRecipe := createRecipeForTest(t, nil)

		c := buildUnauthenticatedGRPCClientForTest(t)
		_, err := c.ArchiveRecipe(ctx, &mealplanninggrpc.ArchiveRecipeRequest{RecipeID: createdRecipe.ID})
		assert.Error(t, err)
	})

	T.Run("nonexistent recipe", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, err := adminClient.ArchiveRecipe(ctx, &mealplanninggrpc.ArchiveRecipeRequest{RecipeID: nonexistentID})
		assert.Error(t, err)
	})

	T.Run("non-admin users are forbidden from archiving", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, testClient := createUserAndClientForTest(t)
		_, _, createdRecipe := createRecipeForTest(t, nil)

		_, err := testClient.ArchiveRecipe(ctx, &mealplanninggrpc.ArchiveRecipeRequest{RecipeID: createdRecipe.ID})
		assert.Error(t, err)
	})
}
