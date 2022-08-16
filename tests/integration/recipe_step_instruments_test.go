package integration

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types/fakes"

	"github.com/prixfixeco/api_server/pkg/types"
)

func checkRecipeStepInstrumentEquality(t *testing.T, expected, actual *types.RecipeStepInstrument, checkInstrument bool) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	if checkInstrument {
		checkValidInstrumentEquality(t, expected.Instrument, actual.Instrument)
	} else {
		assert.Equal(t, expected.Instrument.ID, actual.Instrument.ID, "expected Instrument.ID for recipe step instrument %s to be %v, but it was %v", expected.ID, expected.Instrument.ID, actual.Instrument.ID)
	}
	assert.Equal(t, expected.Name, actual.Name, "expected Name for recipe step instrument %s to be %v, but it was %v", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.ProductOfRecipeStep, actual.ProductOfRecipeStep, "expected ProductOfRecipeStep for recipe step instrument %s to be %v, but it was %v", expected.ID, expected.ProductOfRecipeStep, actual.ProductOfRecipeStep)
	assert.Equal(t, expected.RecipeStepProductID, actual.RecipeStepProductID, "expected RecipeStepProductID for recipe step instrument %s to be %v, but it was %v", expected.ID, expected.RecipeStepProductID, actual.RecipeStepProductID)
	assert.Equal(t, expected.Notes, actual.Notes, "expected Notes for recipe step instrument %s to be %v, but it was %v", expected.ID, expected.Notes, actual.Notes)
	assert.Equal(t, expected.PreferenceRank, actual.PreferenceRank, "expected PreferenceRank for recipe step instrument %s to be %v, but it was %v", expected.ID, expected.PreferenceRank, actual.PreferenceRank)
	assert.NotZero(t, actual.CreatedOn)
}

// convertRecipeStepInstrumentToRecipeStepInstrumentUpdateInput creates an RecipeStepInstrumentUpdateRequestInput struct from a recipe step instrument.
func convertRecipeStepInstrumentToRecipeStepInstrumentUpdateInput(x *types.RecipeStepInstrument) *types.RecipeStepInstrumentUpdateRequestInput {
	return &types.RecipeStepInstrumentUpdateRequestInput{
		InstrumentID:        &x.Instrument.ID,
		RecipeStepProductID: x.RecipeStepProductID,
		ProductOfRecipeStep: &x.ProductOfRecipeStep,
		Notes:               &x.Notes,
		PreferenceRank:      &x.PreferenceRank,
		BelongsToRecipeStep: &x.BelongsToRecipeStep,
		Name:                &x.Name,
	}
}

func (s *TestSuite) TestRecipeStepInstruments_CompleteLifecycle() {
	s.runForEachClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			_, _, createdRecipe := createRecipeForTest(ctx, t, testClients.admin, testClients.user, nil)

			var createdRecipeStepID string
			for _, step := range createdRecipe.Steps {
				createdRecipeStepID = step.ID
				break
			}

			t.Log("creating valid instrument")
			exampleValidInstrument := fakes.BuildFakeValidInstrument()
			exampleValidInstrumentInput := fakes.BuildFakeValidInstrumentCreationRequestInputFromValidInstrument(exampleValidInstrument)
			createdValidInstrument, err := testClients.admin.CreateValidInstrument(ctx, exampleValidInstrumentInput)
			require.NoError(t, err)
			t.Logf("valid instrument %q created", createdValidInstrument.ID)
			checkValidInstrumentEquality(t, exampleValidInstrument, createdValidInstrument)

			t.Log("creating recipe step instrument")
			exampleRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()
			exampleRecipeStepInstrument.BelongsToRecipeStep = createdRecipeStepID
			exampleRecipeStepInstrument.Instrument = &types.ValidInstrument{ID: createdValidInstrument.ID}
			exampleRecipeStepInstrumentInput := fakes.BuildFakeRecipeStepInstrumentCreationRequestInputFromRecipeStepInstrument(exampleRecipeStepInstrument)
			createdRecipeStepInstrument, err := testClients.user.CreateRecipeStepInstrument(ctx, createdRecipe.ID, exampleRecipeStepInstrumentInput)
			require.NoError(t, err)
			t.Logf("recipe step instrument %q created", createdRecipeStepInstrument.ID)

			checkRecipeStepInstrumentEquality(t, exampleRecipeStepInstrument, createdRecipeStepInstrument, false)

			createdRecipeStepInstrument, err = testClients.user.GetRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepInstrument.ID)
			requireNotNilAndNoProblems(t, createdRecipeStepInstrument, err)
			require.Equal(t, createdRecipeStepID, createdRecipeStepInstrument.BelongsToRecipeStep)
			exampleRecipeStepInstrument.Instrument = createdValidInstrument
			exampleRecipeStepInstrument.Instrument.CreatedOn = createdRecipeStepInstrument.Instrument.CreatedOn

			checkRecipeStepInstrumentEquality(t, exampleRecipeStepInstrument, createdRecipeStepInstrument, false)

			t.Log("creating valid instrument")
			newExampleValidInstrument := fakes.BuildFakeValidInstrument()
			newExampleValidInstrumentInput := fakes.BuildFakeValidInstrumentCreationRequestInputFromValidInstrument(newExampleValidInstrument)
			newValidInstrument, err := testClients.admin.CreateValidInstrument(ctx, newExampleValidInstrumentInput)
			require.NoError(t, err)
			t.Logf("valid instrument %q created", createdValidInstrument.ID)
			checkValidInstrumentEquality(t, newExampleValidInstrument, newValidInstrument)

			t.Log("changing recipe step instrument")
			newRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()
			newRecipeStepInstrument.BelongsToRecipeStep = createdRecipeStepID
			newRecipeStepInstrument.Instrument = newValidInstrument
			createdRecipeStepInstrument.Update(convertRecipeStepInstrumentToRecipeStepInstrumentUpdateInput(newRecipeStepInstrument))
			assert.NoError(t, testClients.user.UpdateRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepInstrument))

			t.Log("fetching changed recipe step instrument")
			actual, err := testClients.user.GetRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepInstrument.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert recipe step instrument equality
			checkRecipeStepInstrumentEquality(t, newRecipeStepInstrument, actual, false)
			assert.NotNil(t, actual.LastUpdatedOn)

			t.Log("cleaning up recipe step instrument")
			assert.NoError(t, testClients.user.ArchiveRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepInstrument.ID))

			t.Log("cleaning up recipe step")
			assert.NoError(t, testClients.user.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStepID))

			t.Log("cleaning up recipe")
			assert.NoError(t, testClients.user.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

//func (s *TestSuite) TestRecipeStepInstruments_AsRecipeStepProducts() {
//	s.runForEachClient("should be able to use a recipe step instrument that was the product of a prior recipe step", func(testClients *testClientWrapper) func() {
//		return func() {
//			t := s.T()
//
//			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
//			defer span.End()
//
//			_, _, createdRecipe := createRecipeForTest(ctx, t, testClients.admin, testClients.user, nil)
//
//			var createdRecipeStepID string
//			for _, step := range createdRecipe.Steps {
//				createdRecipeStepID = step.ID
//				break
//			}
//
//			t.Log("creating valid instrument")
//			exampleValidInstrument := fakes.BuildFakeValidInstrument()
//			exampleValidInstrumentInput := fakes.BuildFakeValidInstrumentCreationRequestInputFromValidInstrument(exampleValidInstrument)
//			createdValidInstrument, err := testClients.admin.CreateValidInstrument(ctx, exampleValidInstrumentInput)
//			require.NoError(t, err)
//			t.Logf("valid instrument %q created", createdValidInstrument.ID)
//			checkValidInstrumentEquality(t, exampleValidInstrument, createdValidInstrument)
//
//			time.Sleep(10 * time.Second)
//
//			t.Log("creating recipe step instrument")
//			exampleRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()
//			exampleRecipeStepInstrument.BelongsToRecipeStep = createdRecipeStepID
//			exampleRecipeStepInstrument.Instrument = &types.ValidInstrument{ID: createdValidInstrument.ID}
//			exampleRecipeStepInstrumentInput := fakes.BuildFakeRecipeStepInstrumentCreationRequestInputFromRecipeStepInstrument(exampleRecipeStepInstrument)
//			createdRecipeStepInstrument, err := testClients.user.CreateRecipeStepInstrument(ctx, createdRecipe.ID, exampleRecipeStepInstrumentInput)
//			require.NoError(t, err)
//			t.Logf("recipe step instrument %q created", createdRecipeStepInstrument.ID)
//
//			checkRecipeStepInstrumentEquality(t, exampleRecipeStepInstrument, createdRecipeStepInstrument, false)
//
//			createdRecipeStepInstrument, err = testClients.user.GetRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepInstrument.ID)
//			requireNotNilAndNoProblems(t, createdRecipeStepInstrument, err)
//			require.Equal(t, createdRecipeStepID, createdRecipeStepInstrument.BelongsToRecipeStep)
//
//			checkRecipeStepInstrumentEquality(t, exampleRecipeStepInstrument, createdRecipeStepInstrument, false)
//
//			t.Log("creating valid instrument")
//			newExampleValidInstrument := fakes.BuildFakeValidInstrument()
//			newExampleValidInstrumentInput := fakes.BuildFakeValidInstrumentCreationRequestInputFromValidInstrument(newExampleValidInstrument)
//			newValidInstrument, err := testClients.admin.CreateValidInstrument(ctx, newExampleValidInstrumentInput)
//			require.NoError(t, err)
//			t.Logf("valid instrument %q created", createdValidInstrument.ID)
//			checkValidInstrumentEquality(t, newExampleValidInstrument, newValidInstrument)
//
//			t.Log("changing recipe step instrument")
//			newRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()
//			newRecipeStepInstrument.BelongsToRecipeStep = createdRecipeStepID
//			newRecipeStepInstrument.Instrument = newValidInstrument
//			createdRecipeStepInstrument.Update(convertRecipeStepInstrumentToRecipeStepInstrumentUpdateInput(newRecipeStepInstrument))
//			assert.NoError(t, testClients.user.UpdateRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepInstrument))
//
//			t.Log("fetching changed recipe step instrument")
//			actual, err := testClients.user.GetRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepInstrument.ID)
//			requireNotNilAndNoProblems(t, actual, err)
//
//			// assert recipe step instrument equality
//			checkRecipeStepInstrumentEquality(t, newRecipeStepInstrument, actual, false)
//			assert.NotNil(t, actual.LastUpdatedOn)
//
//			t.Log("cleaning up recipe step instrument")
//			assert.NoError(t, testClients.user.ArchiveRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepInstrument.ID))
//
//			t.Log("cleaning up recipe step")
//			assert.NoError(t, testClients.user.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStepID))
//
//			t.Log("cleaning up recipe")
//			assert.NoError(t, testClients.user.ArchiveRecipe(ctx, createdRecipe.ID))
//		}
//	})
//}

func (s *TestSuite) TestRecipeStepInstruments_Listing() {
	s.runForEachClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			_, _, createdRecipe := createRecipeForTest(ctx, t, testClients.admin, testClients.user, nil)

			var createdRecipeStepID string
			for _, step := range createdRecipe.Steps {
				createdRecipeStepID = step.ID
				break
			}

			t.Log("creating valid instrument")
			exampleValidInstrument := fakes.BuildFakeValidInstrument()
			exampleValidInstrumentInput := fakes.BuildFakeValidInstrumentCreationRequestInputFromValidInstrument(exampleValidInstrument)
			createdValidInstrument, err := testClients.admin.CreateValidInstrument(ctx, exampleValidInstrumentInput)
			require.NoError(t, err)
			t.Logf("valid instrument %q created", createdValidInstrument.ID)
			checkValidInstrumentEquality(t, exampleValidInstrument, createdValidInstrument)

			t.Log("creating recipe step instruments")
			var expected []*types.RecipeStepInstrument
			for i := 0; i < 5; i++ {
				exampleRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()
				exampleRecipeStepInstrument.BelongsToRecipeStep = createdRecipeStepID
				exampleRecipeStepInstrument.Instrument = &types.ValidInstrument{ID: createdValidInstrument.ID}
				exampleRecipeStepInstrumentInput := fakes.BuildFakeRecipeStepInstrumentCreationRequestInputFromRecipeStepInstrument(exampleRecipeStepInstrument)
				createdRecipeStepInstrument, createdRecipeStepInstrumentErr := testClients.user.CreateRecipeStepInstrument(ctx, createdRecipe.ID, exampleRecipeStepInstrumentInput)
				require.NoError(t, createdRecipeStepInstrumentErr)
				t.Logf("recipe step instrument %q created", createdRecipeStepInstrument.ID)
				checkRecipeStepInstrumentEquality(t, exampleRecipeStepInstrument, createdRecipeStepInstrument, false)

				createdRecipeStepInstrument, createdRecipeStepInstrumentErr = testClients.user.GetRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepInstrument.ID)
				requireNotNilAndNoProblems(t, createdRecipeStepInstrument, createdRecipeStepInstrumentErr)
				require.Equal(t, createdRecipeStepID, createdRecipeStepInstrument.BelongsToRecipeStep)

				expected = append(expected, createdRecipeStepInstrument)
			}

			// assert recipe step instrument list equality
			actual, err := testClients.user.GetRecipeStepInstruments(ctx, createdRecipe.ID, createdRecipeStepID, nil)
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
				assert.NoError(t, testClients.user.ArchiveRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepInstrument.ID))
			}

			t.Log("cleaning up recipe step")
			assert.NoError(t, testClients.user.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStepID))

			t.Log("cleaning up recipe")
			assert.NoError(t, testClients.user.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}
