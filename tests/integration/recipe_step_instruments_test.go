package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types/fakes"

	"github.com/prixfixeco/api_server/pkg/types"
)

func checkRecipeStepInstrumentEquality(t *testing.T, expected, actual *types.RecipeStepInstrument) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.InstrumentID, actual.InstrumentID, "expected InstrumentID for recipe step instrument %s to be %v, but it was %v", expected.ID, expected.InstrumentID, actual.InstrumentID)
	assert.Equal(t, expected.RecipeStepID, actual.RecipeStepID, "expected RecipeStepID for recipe step instrument %s to be %v, but it was %v", expected.ID, expected.RecipeStepID, actual.RecipeStepID)
	assert.Equal(t, expected.Notes, actual.Notes, "expected IngredientNotes for recipe step instrument %s to be %v, but it was %v", expected.ID, expected.Notes, actual.Notes)
	assert.NotZero(t, actual.CreatedOn)
}

// convertRecipeStepInstrumentToRecipeStepInstrumentUpdateInput creates an RecipeStepInstrumentUpdateRequestInput struct from a recipe step instrument.
func convertRecipeStepInstrumentToRecipeStepInstrumentUpdateInput(x *types.RecipeStepInstrument) *types.RecipeStepInstrumentUpdateRequestInput {
	return &types.RecipeStepInstrumentUpdateRequestInput{
		InstrumentID: x.InstrumentID,
		RecipeStepID: x.RecipeStepID,
		Notes:        x.Notes,
	}
}

func (s *TestSuite) TestRecipeStepInstruments_CompleteLifecycle() {
	s.runForEachClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			_, _, createdRecipe := createRecipeForTest(ctx, t, testClients.admin, testClients.main, nil)

			var createdRecipeStepID string
			for _, step := range createdRecipe.Steps {
				createdRecipeStepID = step.ID
				break
			}

			t.Log("creating recipe step instrument")
			exampleRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()
			exampleRecipeStepInstrument.BelongsToRecipeStep = createdRecipeStepID
			exampleRecipeStepInstrumentInput := fakes.BuildFakeRecipeStepInstrumentCreationRequestInputFromRecipeStepInstrument(exampleRecipeStepInstrument)
			createdRecipeStepInstrument, err := testClients.main.CreateRecipeStepInstrument(ctx, createdRecipe.ID, exampleRecipeStepInstrumentInput)
			require.NoError(t, err)
			t.Logf("recipe step instrument %q created", createdRecipeStepInstrument.ID)

			checkRecipeStepInstrumentEquality(t, exampleRecipeStepInstrument, createdRecipeStepInstrument)

			createdRecipeStepInstrument, err = testClients.main.GetRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepInstrument.ID)
			requireNotNilAndNoProblems(t, createdRecipeStepInstrument, err)
			require.Equal(t, createdRecipeStepID, createdRecipeStepInstrument.BelongsToRecipeStep)

			checkRecipeStepInstrumentEquality(t, exampleRecipeStepInstrument, createdRecipeStepInstrument)

			t.Log("changing recipe step instrument")
			newRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()
			newRecipeStepInstrument.BelongsToRecipeStep = createdRecipeStepID
			createdRecipeStepInstrument.Update(convertRecipeStepInstrumentToRecipeStepInstrumentUpdateInput(newRecipeStepInstrument))
			assert.NoError(t, testClients.main.UpdateRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepInstrument))

			t.Log("fetching changed recipe step instrument")
			actual, err := testClients.main.GetRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepInstrument.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert recipe step instrument equality
			checkRecipeStepInstrumentEquality(t, newRecipeStepInstrument, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			t.Log("cleaning up recipe step instrument")
			assert.NoError(t, testClients.main.ArchiveRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepInstrument.ID))

			t.Log("cleaning up recipe step")
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStepID))

			t.Log("cleaning up recipe")
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipeStepInstruments_Listing() {
	s.runForEachClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			_, _, createdRecipe := createRecipeForTest(ctx, t, testClients.admin, testClients.main, nil)

			var createdRecipeStepID string
			for _, step := range createdRecipe.Steps {
				createdRecipeStepID = step.ID
				break
			}

			t.Log("creating recipe step instruments")
			var expected []*types.RecipeStepInstrument
			for i := 0; i < 5; i++ {
				exampleRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()
				exampleRecipeStepInstrument.BelongsToRecipeStep = createdRecipeStepID
				exampleRecipeStepInstrumentInput := fakes.BuildFakeRecipeStepInstrumentCreationRequestInputFromRecipeStepInstrument(exampleRecipeStepInstrument)
				createdRecipeStepInstrument, createdRecipeStepInstrumentErr := testClients.main.CreateRecipeStepInstrument(ctx, createdRecipe.ID, exampleRecipeStepInstrumentInput)
				require.NoError(t, createdRecipeStepInstrumentErr)
				t.Logf("recipe step instrument %q created", createdRecipeStepInstrument.ID)
				checkRecipeStepInstrumentEquality(t, exampleRecipeStepInstrument, createdRecipeStepInstrument)

				createdRecipeStepInstrument, createdRecipeStepInstrumentErr = testClients.main.GetRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepInstrument.ID)
				requireNotNilAndNoProblems(t, createdRecipeStepInstrument, createdRecipeStepInstrumentErr)
				require.Equal(t, createdRecipeStepID, createdRecipeStepInstrument.BelongsToRecipeStep)

				expected = append(expected, createdRecipeStepInstrument)
			}

			// assert recipe step instrument list equality
			actual, err := testClients.main.GetRecipeStepInstruments(ctx, createdRecipe.ID, createdRecipeStepID, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.RecipeStepInstruments),
				"expected %d to be <= %d",
				len(expected),
				len(actual.RecipeStepInstruments),
			)

			t.Log("cleaning up")
			for _, createdRecipeStepInstrument := range expected {
				assert.NoError(t, testClients.main.ArchiveRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepInstrument.ID))
			}

			t.Log("cleaning up recipe step")
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStepID))

			t.Log("cleaning up recipe")
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}
