package integration

import (
	"testing"

	types "github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	mpconverters "github.com/dinnerdonebetter/backend/internal/domain/mealplanning/converters"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/fakes"
	mealplanninggrpc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	converters "github.com/dinnerdonebetter/backend/internal/services/mealplanning/grpc/converters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecipeSteps_CompleteLifecycle(T *testing.T) {
	T.Parallel()

	T.Run("should update", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		createdValidIngredients, createdValidPreparation, createdRecipe := createRecipeForTest(t, nil)

		var (
			createdRecipeStep *types.RecipeStep
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

		_, err := adminClient.UpdateRecipeStep(ctx, &mealplanninggrpc.UpdateRecipeStepRequest{
			RecipeID:     createdRecipe.ID,
			RecipeStepID: createdRecipeStep.ID,
			Input:        converters.ConvertRecipeStepUpdateRequestInputToGRPCRecipeStepUpdateRequestInput(updateInput),
		})
		require.NoError(t, err)

		actual, err := adminClient.GetRecipeStep(ctx, &mealplanninggrpc.GetRecipeStepRequest{
			RecipeID:     createdRecipe.ID,
			RecipeStepID: createdRecipeStep.ID,
		})
		require.NoError(t, err)

		// assert recipe step equality
		checkRecipeStepEquality(t, stepIndex, newRecipeStep, converters.ConvertGRPCRecipeStepToRecipeStep(actual.Result))
		assert.NotNil(t, actual.Result.LastUpdatedAt)

		_, err = adminClient.ArchiveRecipeStep(ctx, &mealplanninggrpc.ArchiveRecipeStepRequest{
			RecipeID:     createdRecipe.ID,
			RecipeStepID: createdRecipeStep.ID,
		})
		assert.NoError(t, err)

		_, err = adminClient.ArchiveRecipe(ctx, &mealplanninggrpc.ArchiveRecipeRequest{RecipeID: createdRecipe.ID})
		assert.NoError(t, err)
	})
}

func TestRecipeSteps_Listing(T *testing.T) {
	T.SkipNow()

	T.Run("should be readable in paginated form", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		createdValidIngredients, createdValidPreparation, createdRecipe := createRecipeForTest(t, nil)
		createdValidMeasurementUnit := createValidMeasurementUnitForTest(t)
		createdValidInstrument := createValidInstrumentForTest(t)
		createdValidIngredientState := createValidIngredientStateForTest(t)
		createdValidVessel := createValidVesselForTest(t)

		var expected []*types.RecipeStep
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
					exampleRecipeStep.CompletionConditions[j].Ingredients[k].RecipeStepIngredient = createdValidIngredients[0].ID
				}
			}

			exampleRecipeStepInput := mpconverters.ConvertRecipeStepToRecipeStepCreationRequestInput(exampleRecipeStep)
			exampleRecipeStepInput.PreparationID = createdValidPreparation.ID

			createdRecipeStep, createdRecipeStepErr := adminClient.CreateRecipeStep(ctx, &mealplanninggrpc.CreateRecipeStepRequest{
				RecipeID: createdRecipe.ID,
				Input:    converters.ConvertRecipeStepCreationRequestInputToGRPCRecipeStepCreationRequestInput(exampleRecipeStepInput),
			})
			require.NoError(t, createdRecipeStepErr)

			checkRecipeStepEquality(t, -1, exampleRecipeStep, converters.ConvertGRPCRecipeStepToRecipeStep(createdRecipeStep.Created))

			recipeStep, err := adminClient.GetRecipeStep(ctx, &mealplanninggrpc.GetRecipeStepRequest{
				RecipeID:     createdRecipe.ID,
				RecipeStepID: createdRecipeStep.Created.ID,
			})
			require.NoError(t, err)
			require.Equal(t, createdRecipe.ID, recipeStep.Result.BelongsToRecipe)

			expected = append(expected, converters.ConvertGRPCRecipeStepToRecipeStep(recipeStep.Result))
			t.Logf("created step #%d", i)
		}

		// assert recipe step list equality
		actual, err := adminClient.GetRecipeSteps(ctx, &mealplanninggrpc.GetRecipeStepsRequest{RecipeID: createdRecipe.ID})
		require.NoError(t, err)
		assert.True(
			t,
			len(expected) <= len(actual.Results),
			"expected %d to be <= %d",
			len(expected),
			len(actual.Results),
		)

		for _, createdRecipeStep := range expected {
			_, err = adminClient.ArchiveRecipeStep(ctx, &mealplanninggrpc.ArchiveRecipeStepRequest{RecipeID: createdRecipe.ID, RecipeStepID: createdRecipeStep.ID})
			assert.NoError(t, err)
		}

		_, err = adminClient.ArchiveRecipe(ctx, &mealplanninggrpc.ArchiveRecipeRequest{RecipeID: createdRecipe.ID})
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
//			var createdRecipeStep *types.RecipeStep
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
