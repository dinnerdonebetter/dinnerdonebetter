package integration

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	mpconverters "github.com/dinnerdonebetter/backend/internal/domain/mealplanning/converters"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/fakes"
	mealplanninggrpc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	converters "github.com/dinnerdonebetter/backend/internal/services/mealplanning/grpc/converters"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkRecipeStepEquality(t *testing.T, index int, expected, actual *mealplanning.RecipeStep) {
	t.Helper()

	assert.NotZero(t, actual.CreatedAt, "expected recipe step %d", index)
	assert.Equal(t, expected.EstimatedTimeInSeconds, actual.EstimatedTimeInSeconds, "expected recipe step %d", index)
	assert.Equal(t, expected.TemperatureInCelsius, actual.TemperatureInCelsius, "expected recipe step %d", index)
	assert.NotEmpty(t, actual.BelongsToRecipe, "expected recipe step %d", index)
	assert.Equal(t, expected.ConditionExpression, actual.ConditionExpression, "expected recipe step %d", index)
	assert.NotEmpty(t, actual.ID, "expected recipe step %d", index)
	assert.Equal(t, expected.Notes, actual.Notes, "expected recipe step %d", index)
	assert.Equal(t, expected.ExplicitInstructions, actual.ExplicitInstructions, "expected recipe step %d", index)
	checkRecipeMediaSliceEquality(t, index, expected.Media, actual.Media)
	checkRecipeStepProductSliceEquality(t, index, expected.Products, actual.Products)
	checkRecipeStepInstrumentSliceEquality(t, index, expected.Instruments, actual.Instruments)
	checkRecipeStepVesselSliceEquality(t, index, expected.Vessels, actual.Vessels)
	checkRecipeStepCompletionConditionSliceEquality(t, index, expected.CompletionConditions, actual.CompletionConditions)
	checkRecipeStepIngredientSliceEquality(t, index, expected.Ingredients, actual.Ingredients)
	checkValidPreparationEquality(t, index, &expected.Preparation, &actual.Preparation)
	assert.Equal(t, expected.Index, actual.Index, "expected recipe step %d", index)
	assert.Equal(t, expected.Optional, actual.Optional, "expected recipe step %d", index)
	assert.Equal(t, expected.StartTimerAutomatically, actual.StartTimerAutomatically, "expected recipe step %d", index)
}

func TestRecipeSteps_CompleteLifecycle(T *testing.T) {
	T.Parallel()

	T.Run("should update", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		createdValidIngredients, createdValidPreparation, createdRecipe := createRecipeForTest(t, nil)

		var (
			createdRecipeStep *mealplanning.RecipeStep
			stepIndex         int
		)
		for i, step := range createdRecipe.Steps {
			createdRecipeStep = step
			stepIndex = i
			break
		}

		require.NotNil(t, createdRecipeStep)

		newRecipeStep := fakes.BuildFakeRecipeStep()
		newRecipeStep.BelongsToRecipe = createdRecipe.ID
		newRecipeStep.Media = createdRecipeStep.Media
		newRecipeStep.Products = createdRecipeStep.Products
		newRecipeStep.Instruments = createdRecipeStep.Instruments
		newRecipeStep.Vessels = createdRecipeStep.Vessels
		newRecipeStep.CompletionConditions = createdRecipeStep.CompletionConditions
		newRecipeStep.Ingredients = createdRecipeStep.Ingredients
		for j := range newRecipeStep.Ingredients {
			newRecipeStep.Ingredients[j].Ingredient = createdValidIngredients[j]
		}

		updateInput := mpconverters.ConvertRecipeStepToRecipeStepUpdateRequestInput(newRecipeStep)
		updateInput.Preparation = createdValidPreparation

		createdRecipeStep.Update(updateInput)

		updateResponse, err := adminClient.UpdateRecipeStep(ctx, &mealplanninggrpc.UpdateRecipeStepRequest{
			RecipeId:     createdRecipe.ID,
			RecipeStepId: createdRecipeStep.ID,
			Input:        converters.ConvertRecipeStepUpdateRequestInputToGRPCRecipeStepUpdateRequestInput(updateInput),
		})
		require.NoError(t, err)

		// Test the response from UpdateRecipeStep first
		checkRecipeStepEquality(t, stepIndex, createdRecipeStep, converters.ConvertGRPCRecipeStepToRecipeStep(updateResponse.Updated))
		assert.NotNil(t, updateResponse.Updated.LastUpdatedAt)

		// Also test the separate GetRecipeStep call
		actual, err := adminClient.GetRecipeStep(ctx, &mealplanninggrpc.GetRecipeStepRequest{
			RecipeId:     createdRecipe.ID,
			RecipeStepId: createdRecipeStep.ID,
		})
		require.NoError(t, err)

		// assert recipe step equality for the separate get call
		checkRecipeStepEquality(t, stepIndex, createdRecipeStep, converters.ConvertGRPCRecipeStepToRecipeStep(actual.Result))
		assert.NotNil(t, actual.Result.LastUpdatedAt)

		_, err = adminClient.ArchiveRecipeStep(ctx, &mealplanninggrpc.ArchiveRecipeStepRequest{
			RecipeId:     createdRecipe.ID,
			RecipeStepId: createdRecipeStep.ID,
		})
		assert.NoError(t, err)

		_, err = adminClient.ArchiveRecipe(ctx, &mealplanninggrpc.ArchiveRecipeRequest{RecipeId: createdRecipe.ID})
		assert.NoError(t, err)
	})
}

func TestRecipeSteps_Listing(T *testing.T) {
	T.Parallel()

	T.Run("should be readable in paginated form", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		createdValidIngredients, createdValidPreparation, createdRecipe := createRecipeForTest(t, nil)
		createdValidMeasurementUnit := createValidMeasurementUnitForTest(t)
		createdValidInstrument := createValidInstrumentForTest(t)
		createdValidIngredientState := createValidIngredientStateForTest(t)
		createdValidVessel := createValidVesselForTest(t)

		var expected []*mealplanning.RecipeStep
		for i := 0; i < 5; i++ {
			t.Logf("creating recipe step #%d", i+1)

			exampleRecipeStep := fakes.BuildFakeRecipeStep()
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			for j := range exampleRecipeStep.Ingredients {
				exampleRecipeStep.Ingredients[j].Ingredient = createdValidIngredients[j]
				exampleRecipeStep.Ingredients[j].MeasurementUnit = *createdValidMeasurementUnit
			}

			for j := range exampleRecipeStep.Products {
				exampleRecipeStep.Products[j].MeasurementUnit = createdValidMeasurementUnit
			}

			for j := range exampleRecipeStep.Instruments {
				exampleRecipeStep.Instruments[j].Instrument = createdValidInstrument
			}

			for j := range exampleRecipeStep.Vessels {
				exampleRecipeStep.Vessels[j].Vessel = createdValidVessel
			}

			for j := range exampleRecipeStep.CompletionConditions {
				exampleRecipeStep.CompletionConditions[j].IngredientState = *createdValidIngredientState
				for k := range exampleRecipeStep.CompletionConditions[j].Ingredients {
					// Reference the first recipe step ingredient that will be created
					exampleRecipeStep.CompletionConditions[j].Ingredients[k].RecipeStepIngredient = exampleRecipeStep.Ingredients[0].ID
				}
			}

			exampleRecipeStepInput := mpconverters.ConvertRecipeStepToRecipeStepCreationRequestInput(exampleRecipeStep)
			exampleRecipeStepInput.PreparationID = createdValidPreparation.ID

			// Set the preparation on the example recipe step for comparison
			exampleRecipeStep.Preparation = *createdValidPreparation

			createdRecipeStep, createdRecipeStepErr := adminClient.CreateRecipeStep(ctx, &mealplanninggrpc.CreateRecipeStepRequest{
				RecipeId: createdRecipe.ID,
				Input:    converters.ConvertRecipeStepCreationRequestInputToGRPCRecipeStepCreationRequestInput(exampleRecipeStepInput),
			})
			require.NoError(t, createdRecipeStepErr)

			// Update the completion condition ingredient references to match the actual created recipe step
			createdRecipeStepConverted := converters.ConvertGRPCRecipeStepToRecipeStep(createdRecipeStep.Created)
			for j := range exampleRecipeStep.CompletionConditions {
				for k := range exampleRecipeStep.CompletionConditions[j].Ingredients {
					// Use the actual created recipe step ingredient ID
					exampleRecipeStep.CompletionConditions[j].Ingredients[k].RecipeStepIngredient = createdRecipeStepConverted.Ingredients[k].ID
				}
			}

			checkRecipeStepEquality(t, -1, exampleRecipeStep, createdRecipeStepConverted)

			recipeStep, err := adminClient.GetRecipeStep(ctx, &mealplanninggrpc.GetRecipeStepRequest{
				RecipeId:     createdRecipe.ID,
				RecipeStepId: createdRecipeStep.Created.Id,
			})
			require.NoError(t, err)
			require.Equal(t, createdRecipe.ID, recipeStep.Result.BelongsToRecipe)

			expected = append(expected, converters.ConvertGRPCRecipeStepToRecipeStep(recipeStep.Result))
			t.Logf("created step #%d", i)
		}

		// assert recipe step list equality
		actual, err := adminClient.GetRecipeSteps(ctx, &mealplanninggrpc.GetRecipeStepsRequest{RecipeId: createdRecipe.ID})
		require.NoError(t, err)
		assert.True(
			t,
			len(expected) <= len(actual.Results),
			"expected %d to be <= %d",
			len(expected),
			len(actual.Results),
		)

		for _, createdRecipeStep := range expected {
			_, err = adminClient.ArchiveRecipeStep(ctx, &mealplanninggrpc.ArchiveRecipeStepRequest{RecipeId: createdRecipe.ID, RecipeStepId: createdRecipeStep.ID})
			assert.NoError(t, err)
		}

		_, err = adminClient.ArchiveRecipe(ctx, &mealplanninggrpc.ArchiveRecipeRequest{RecipeId: createdRecipe.ID})
		assert.NoError(t, err)
	})
}

/*

// TODO: uncomment me.
//func (s *TestSuite) TestRecipeSteps_ContentUploading() {
//	s.runTest("should be able to upload content for a recipe step", func(testClients *testClientWrapper) func() {
//		return func() {
//			t := s.T()
//
//			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
//			defer span.End()
//
//			_, _, createdRecipe := createRecipeForTest(ctx, t, adminClient, adminClient, nil)
//
//			var createdRecipeStep *mealplanning.RecipeStep
//			for _, step := range createdRecipe.Steps {
//				createdRecipeStep = step
//				break
//			}
//
//			require.NotNil(t, createdRecipeStep)
//
//			_, img1Bytes := testutils.BuildArbitraryImagePNGBytes(200)
//			_, img2Bytes := testutils.BuildArbitraryImagePNGBytes(250)
//			_, img3Bytes := testutils.BuildArbitraryImagePNGBytes(300)
//
//			files := map[string][]byte{
//				"image_1.png": img1Bytes,
//				"image_2.png": img2Bytes,
//				"image_3.png": img3Bytes,
//			}
//
//			require.NoError(t, adminClient.UploadMediaForRecipeStep(ctx, files, createdRecipe.ID, createdRecipeStep.ID))
//
//			assert.NoError(t, adminClient.ArchiveRecipe(ctx, createdRecipe.ID))
//		}
//	})
//}


*/
