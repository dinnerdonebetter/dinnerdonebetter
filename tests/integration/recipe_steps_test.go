package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types/fakes"

	"github.com/prixfixeco/api_server/pkg/types"
)

func checkRecipeStepEquality(t *testing.T, expected, actual *types.RecipeStep) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Index, actual.Index, "expected Index for recipe step %s to be %v, but it was %v", expected.ID, expected.Index, actual.Index)
	assert.Equal(t, expected.MinimumEstimatedTimeInSeconds, actual.MinimumEstimatedTimeInSeconds, "expected MinimumEstimatedTimeInSeconds for recipe step %s to be %v, but it was %v", expected.ID, expected.MinimumEstimatedTimeInSeconds, actual.MinimumEstimatedTimeInSeconds)
	assert.Equal(t, expected.MaximumEstimatedTimeInSeconds, actual.MaximumEstimatedTimeInSeconds, "expected MaximumEstimatedTimeInSeconds for recipe step %s to be %v, but it was %v", expected.ID, expected.MaximumEstimatedTimeInSeconds, actual.MaximumEstimatedTimeInSeconds)
	assert.Equal(t, expected.MinimumTemperatureInCelsius, actual.MinimumTemperatureInCelsius, "expected MinimumTemperatureInCelsius for recipe step %s to be %v, but it was %v", expected.ID, expected.MinimumTemperatureInCelsius, actual.MinimumTemperatureInCelsius)
	assert.Equal(t, expected.MaximumTemperatureInCelsius, actual.MaximumTemperatureInCelsius, "expected MaximumTemperatureInCelsius for recipe step %s to be %v, but it was %v", expected.ID, expected.MaximumTemperatureInCelsius, actual.MaximumTemperatureInCelsius)
	assert.Equal(t, expected.Notes, actual.Notes, "expected Notes for recipe step %s to be %v, but it was %v", expected.ID, expected.Notes, actual.Notes)
	assert.NotZero(t, actual.CreatedOn)
}

// convertRecipeStepToRecipeStepUpdateInput creates an RecipeStepUpdateRequestInput struct from a recipe step.
func convertRecipeStepToRecipeStepUpdateInput(x *types.RecipeStep) *types.RecipeStepUpdateRequestInput {
	return &types.RecipeStepUpdateRequestInput{
		Index:                         &x.Index,
		MinimumEstimatedTimeInSeconds: &x.MinimumEstimatedTimeInSeconds,
		MaximumEstimatedTimeInSeconds: &x.MaximumEstimatedTimeInSeconds,
		MinimumTemperatureInCelsius:   x.MinimumTemperatureInCelsius,
		MaximumTemperatureInCelsius:   x.MaximumTemperatureInCelsius,
		Notes:                         &x.Notes,
	}
}

func (s *TestSuite) TestRecipeSteps_CompleteLifecycle() {
	s.runForEachClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			createdValidIngredients, createdValidPreparation, createdRecipe := createRecipeForTest(ctx, t, testClients.admin, testClients.main, nil)

			var createdRecipeStep *types.RecipeStep
			for _, step := range createdRecipe.Steps {
				createdRecipeStep = step
				break
			}

			t.Log("changing recipe step")
			newRecipeStep := fakes.BuildFakeRecipeStep()
			newRecipeStep.BelongsToRecipe = createdRecipe.ID
			for j := range newRecipeStep.Ingredients {
				newRecipeStep.Ingredients[j].IngredientID = stringPointer(createdValidIngredients[j].ID)
			}

			updateInput := convertRecipeStepToRecipeStepUpdateInput(newRecipeStep)
			updateInput.Preparation = createdValidPreparation
			createdRecipeStep.Update(updateInput)
			assert.NoError(t, testClients.main.UpdateRecipeStep(ctx, createdRecipeStep))

			t.Log("fetching changed recipe step")
			actual, err := testClients.main.GetRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert recipe step equality
			checkRecipeStepEquality(t, newRecipeStep, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			t.Log("cleaning up recipe step")
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

			t.Log("cleaning up recipe")
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipeSteps_Listing() {
	s.runForEachClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			createdValidIngredients, createdValidPreparation, createdRecipe := createRecipeForTest(ctx, t, testClients.admin, testClients.main, nil)

			t.Log("creating recipe steps")
			var expected []*types.RecipeStep
			for i := 0; i < 5; i++ {
				exampleRecipeStep := fakes.BuildFakeRecipeStep()
				exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
				for j := range exampleRecipeStep.Ingredients {
					exampleRecipeStep.Ingredients[j].IngredientID = stringPointer(createdValidIngredients[j].ID)
				}

				exampleRecipeStepInput := fakes.BuildFakeRecipeStepCreationRequestInputFromRecipeStep(exampleRecipeStep)
				exampleRecipeStepInput.PreparationID = createdValidPreparation.ID
				createdRecipeStep, createdRecipeStepErr := testClients.main.CreateRecipeStep(ctx, exampleRecipeStepInput)
				require.NoError(t, createdRecipeStepErr)
				t.Logf("recipe step %q created", createdRecipeStep.ID)

				checkRecipeStepEquality(t, exampleRecipeStep, createdRecipeStep)

				createdRecipeStep, createdRecipeStepErr = testClients.main.GetRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID)
				requireNotNilAndNoProblems(t, createdRecipeStep, createdRecipeStepErr)
				require.Equal(t, createdRecipe.ID, createdRecipeStep.BelongsToRecipe)

				expected = append(expected, createdRecipeStep)
			}

			// assert recipe step list equality
			actual, err := testClients.main.GetRecipeSteps(ctx, createdRecipe.ID, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.RecipeSteps),
				"expected %d to be <= %d",
				len(expected),
				len(actual.RecipeSteps),
			)

			t.Log("cleaning up")
			for _, createdRecipeStep := range expected {
				assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))
			}

			t.Log("cleaning up recipe")
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}
