package integration

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkRecipeStepEquality(t *testing.T, expected, actual *types.RecipeStep) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Index, actual.Index, "expected Index for recipe step %s to be %v, but it was %v", expected.ID, expected.Index, actual.Index)
	assert.Equal(t, expected.EstimatedTimeInSeconds, actual.EstimatedTimeInSeconds, "expected EstimatedTimeInSeconds for recipe step %s to be %v, but it was %v", expected.ID, expected.EstimatedTimeInSeconds, actual.EstimatedTimeInSeconds)
	assert.Equal(t, expected.TemperatureInCelsius, actual.TemperatureInCelsius, "expected TemperatureInCelsius for recipe step %s to be %v, but it was %v", expected.ID, expected.TemperatureInCelsius, actual.TemperatureInCelsius)
	assert.Equal(t, expected.Notes, actual.Notes, "expected StatusExplanation for recipe step %s to be %v, but it was %v", expected.ID, expected.Notes, actual.Notes)
	assert.Equal(t, expected.ConditionExpression, actual.ConditionExpression, "expected StatusExplanation for recipe step %s to be %v, but it was %v", expected.ID, expected.ConditionExpression, actual.ConditionExpression)
	assert.Equal(t, expected.ExplicitInstructions, actual.ExplicitInstructions, "expected ExplicitInstructions for recipe step %s to be %v, but it was %v", expected.ID, expected.ExplicitInstructions, actual.ExplicitInstructions)
	assert.Equal(t, expected.Optional, actual.Optional, "expected Optional for recipe step %s to be %v, but it was %v", expected.ID, expected.Optional, actual.Optional)
	assert.Equal(t, expected.StartTimerAutomatically, actual.StartTimerAutomatically, "expected StartTimerAutomatically for recipe step %s to be %v, but it was %v", expected.ID, expected.StartTimerAutomatically, actual.StartTimerAutomatically)
	assert.NotZero(t, actual.CreatedAt)
}

func (s *TestSuite) TestRecipeSteps_CompleteLifecycle() {
	s.runTest("should CRUD", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			createdValidIngredients, createdValidPreparation, createdRecipe := createRecipeForTest(ctx, t, testClients.adminClient, testClients.userClient, nil)

			var createdRecipeStep *types.RecipeStep
			for _, step := range createdRecipe.Steps {
				createdRecipeStep = step
				break
			}

			require.NotNil(t, createdRecipeStep)

			newRecipeStep := fakes.BuildFakeRecipeStep()
			newRecipeStep.BelongsToRecipe = createdRecipe.ID
			for j := range newRecipeStep.Ingredients {
				newRecipeStep.Ingredients[j].Ingredient = createdValidIngredients[j]
			}

			updateInput := converters.ConvertRecipeStepToRecipeStepUpdateRequestInput(newRecipeStep)
			updateInput.Preparation = createdValidPreparation
			createdRecipeStep.Update(updateInput)
			assert.NoError(t, testClients.adminClient.UpdateRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID, updateInput))

			actual, err := testClients.userClient.GetRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert recipe step equality
			checkRecipeStepEquality(t, newRecipeStep, actual)
			assert.NotNil(t, actual.LastUpdatedAt)

			assert.NoError(t, testClients.userClient.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

			assert.NoError(t, testClients.adminClient.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

// TODO: uncomment me.
//func (s *TestSuite) TestRecipeSteps_ContentUploading() {
//	s.runTest("should be able to upload content for a recipe step", func(testClients *testClientWrapper) func() {
//		return func() {
//			t := s.T()
//
//			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
//			defer span.End()
//
//			_, _, createdRecipe := createRecipeForTest(ctx, t, testClients.adminClient, testClients.userClient, nil)
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
//			require.NoError(t, testClients.userClient.UploadMediaForRecipeStep(ctx, files, createdRecipe.ID, createdRecipeStep.ID))
//
//			assert.NoError(t, testClients.adminClient.ArchiveRecipe(ctx, createdRecipe.ID))
//		}
//	})
//}

func (s *TestSuite) TestRecipeSteps_Listing() {
	s.runTest("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			createdValidIngredients, createdValidPreparation, createdRecipe := createRecipeForTest(ctx, t, testClients.adminClient, testClients.userClient, nil)

			createdValidMeasurementUnit := createValidMeasurementUnitForTest(t, ctx, testClients.adminClient)
			createdValidInstrument := createValidInstrumentForTest(t, ctx, testClients.adminClient)
			createdValidIngredientState := createValidIngredientStateForTest(t, ctx, testClients.adminClient)
			createdValidVessel := createValidVesselForTest(t, ctx, nil, testClients.adminClient)

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

				exampleRecipeStepInput := converters.ConvertRecipeStepToRecipeStepCreationRequestInput(exampleRecipeStep)
				exampleRecipeStepInput.PreparationID = createdValidPreparation.ID

				createdRecipeStep, createdRecipeStepErr := testClients.adminClient.CreateRecipeStep(ctx, createdRecipe.ID, exampleRecipeStepInput)
				require.NoError(t, createdRecipeStepErr)

				checkRecipeStepEquality(t, exampleRecipeStep, createdRecipeStep)

				createdRecipeStep, createdRecipeStepErr = testClients.userClient.GetRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID)
				requireNotNilAndNoProblems(t, createdRecipeStep, createdRecipeStepErr)
				require.Equal(t, createdRecipe.ID, createdRecipeStep.BelongsToRecipe)

				expected = append(expected, createdRecipeStep)
				t.Logf("created step #%d", i)
			}

			// assert recipe step list equality
			actual, err := testClients.userClient.GetRecipeSteps(ctx, createdRecipe.ID, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Data),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Data),
			)

			for _, createdRecipeStep := range expected {
				assert.NoError(t, testClients.userClient.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))
			}

			assert.NoError(t, testClients.adminClient.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}
