package integration

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"
)

func checkRecipeStepInstrumentEquality(t *testing.T, expected, actual *types.RecipeStepInstrument) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.InstrumentID, actual.InstrumentID, "expected InstrumentID for recipe step instrument %s to be %v, but it was %v", expected.ID, expected.InstrumentID, actual.InstrumentID)
	assert.Equal(t, expected.RecipeStepID, actual.RecipeStepID, "expected RecipeStepID for recipe step instrument %s to be %v, but it was %v", expected.ID, expected.RecipeStepID, actual.RecipeStepID)
	assert.Equal(t, expected.Notes, actual.Notes, "expected Notes for recipe step instrument %s to be %v, but it was %v", expected.ID, expected.Notes, actual.Notes)
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
	s.runForCookieClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			stopChan := make(chan bool, 1)
			notificationsChan, err := testClients.main.SubscribeToDataChangeNotifications(ctx, stopChan)
			require.NotNil(t, notificationsChan)
			require.NoError(t, err)

			var n *types.DataChangeMessage

			t.Log("creating prerequisite recipe")
			exampleRecipe := fakes.BuildFakeRecipe()
			exampleRecipeInput := fakes.BuildFakeRecipeCreationRequestInputFromRecipe(exampleRecipe)
			createdRecipeID, err := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
			require.NoError(t, err)
			t.Logf("recipe %q created", createdRecipeID)

			n = <-notificationsChan
			assert.Equal(t, n.DataType, types.RecipeDataType)
			require.NotNil(t, n.Recipe)
			checkRecipeEquality(t, exampleRecipe, n.Recipe)

			createdRecipe, err := testClients.main.GetRecipe(ctx, createdRecipeID)
			requireNotNilAndNoProblems(t, createdRecipe, err)

			t.Log("creating prerequisite recipe step")
			exampleRecipeStep := fakes.BuildFakeRecipeStep()
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeStepInput := fakes.BuildFakeRecipeStepCreationRequestInputFromRecipeStep(exampleRecipeStep)
			createdRecipeStepID, err := testClients.main.CreateRecipeStep(ctx, exampleRecipeStepInput)
			require.NoError(t, err)
			t.Logf("recipe step %q created", createdRecipeStepID)

			n = <-notificationsChan
			assert.Equal(t, n.DataType, types.RecipeStepDataType)
			require.NotNil(t, n.RecipeStep)
			checkRecipeStepEquality(t, exampleRecipeStep, n.RecipeStep)

			createdRecipeStep, err := testClients.main.GetRecipeStep(ctx, createdRecipe.ID, createdRecipeStepID)
			requireNotNilAndNoProblems(t, createdRecipeStep, err)
			require.Equal(t, createdRecipe.ID, createdRecipeStep.BelongsToRecipe)

			t.Log("creating recipe step instrument")
			exampleRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()
			exampleRecipeStepInstrument.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepInstrumentInput := fakes.BuildFakeRecipeStepInstrumentCreationRequestInputFromRecipeStepInstrument(exampleRecipeStepInstrument)
			createdRecipeStepInstrumentID, err := testClients.main.CreateRecipeStepInstrument(ctx, createdRecipe.ID, exampleRecipeStepInstrumentInput)
			require.NoError(t, err)
			t.Logf("recipe step instrument %q created", createdRecipeStepInstrumentID)

			n = <-notificationsChan
			assert.Equal(t, n.DataType, types.RecipeStepInstrumentDataType)
			require.NotNil(t, n.RecipeStepInstrument)
			checkRecipeStepInstrumentEquality(t, exampleRecipeStepInstrument, n.RecipeStepInstrument)

			createdRecipeStepInstrument, err := testClients.main.GetRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepInstrumentID)
			requireNotNilAndNoProblems(t, createdRecipeStepInstrument, err)
			require.Equal(t, createdRecipeStep.ID, createdRecipeStepInstrument.BelongsToRecipeStep)

			checkRecipeStepInstrumentEquality(t, exampleRecipeStepInstrument, createdRecipeStepInstrument)

			t.Log("changing recipe step instrument")
			newRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()
			createdRecipeStepInstrument.Update(convertRecipeStepInstrumentToRecipeStepInstrumentUpdateInput(newRecipeStepInstrument))
			assert.NoError(t, testClients.main.UpdateRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepInstrument))

			n = <-notificationsChan
			assert.Equal(t, n.DataType, types.RecipeStepInstrumentDataType)

			t.Log("fetching changed recipe step instrument")
			actual, err := testClients.main.GetRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepInstrumentID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert recipe step instrument equality
			checkRecipeStepInstrumentEquality(t, newRecipeStepInstrument, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			t.Log("cleaning up recipe step instrument")
			assert.NoError(t, testClients.main.ArchiveRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepInstrumentID))

			t.Log("cleaning up recipe step")
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStepID))

			t.Log("cleaning up recipe")
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipeID))
		}
	})

	s.runForPASETOClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			var checkFunc func() bool
			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating prerequisite recipe")
			exampleRecipe := fakes.BuildFakeRecipe()
			exampleRecipeInput := fakes.BuildFakeRecipeCreationRequestInputFromRecipe(exampleRecipe)
			createdRecipeID, err := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
			require.NoError(t, err)
			t.Logf("recipe %q created", createdRecipeID)

			var createdRecipe *types.Recipe
			checkFunc = func() bool {
				createdRecipe, err = testClients.main.GetRecipe(ctx, createdRecipeID)
				return assert.NotNil(t, createdRecipe) && assert.NoError(t, err)
			}
			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
			checkRecipeEquality(t, exampleRecipe, createdRecipe)

			t.Log("creating prerequisite recipe step")
			exampleRecipeStep := fakes.BuildFakeRecipeStep()
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeStepInput := fakes.BuildFakeRecipeStepCreationRequestInputFromRecipeStep(exampleRecipeStep)
			createdRecipeStepID, err := testClients.main.CreateRecipeStep(ctx, exampleRecipeStepInput)
			require.NoError(t, err)
			t.Logf("recipe step %q created", createdRecipeStepID)

			var createdRecipeStep *types.RecipeStep
			checkFunc = func() bool {
				createdRecipeStep, err = testClients.main.GetRecipeStep(ctx, createdRecipe.ID, createdRecipeStepID)
				return assert.NotNil(t, createdRecipeStep) && assert.NoError(t, err)
			}
			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
			require.Equal(t, createdRecipe.ID, createdRecipeStep.BelongsToRecipe)
			checkRecipeStepEquality(t, exampleRecipeStep, createdRecipeStep)

			t.Log("creating recipe step instrument")
			exampleRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()
			exampleRecipeStepInstrument.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepInstrumentInput := fakes.BuildFakeRecipeStepInstrumentCreationRequestInputFromRecipeStepInstrument(exampleRecipeStepInstrument)
			createdRecipeStepInstrumentID, err := testClients.main.CreateRecipeStepInstrument(ctx, createdRecipe.ID, exampleRecipeStepInstrumentInput)
			require.NoError(t, err)
			t.Logf("recipe step instrument %q created", createdRecipeStepInstrumentID)

			var createdRecipeStepInstrument *types.RecipeStepInstrument
			checkFunc = func() bool {
				createdRecipeStepInstrument, err = testClients.main.GetRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepInstrumentID)
				return assert.NotNil(t, createdRecipeStepInstrument) && assert.NoError(t, err)
			}
			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
			require.Equal(t, createdRecipeStep.ID, createdRecipeStepInstrument.BelongsToRecipeStep)
			checkRecipeStepInstrumentEquality(t, exampleRecipeStepInstrument, createdRecipeStepInstrument)

			// assert recipe step instrument equality
			checkRecipeStepInstrumentEquality(t, exampleRecipeStepInstrument, createdRecipeStepInstrument)

			// change recipe step instrument
			newRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()
			createdRecipeStepInstrument.Update(convertRecipeStepInstrumentToRecipeStepInstrumentUpdateInput(newRecipeStepInstrument))
			assert.NoError(t, testClients.main.UpdateRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepInstrument))

			time.Sleep(time.Second)

			// retrieve changed recipe step instrument
			var actual *types.RecipeStepInstrument
			checkFunc = func() bool {
				actual, err = testClients.main.GetRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepInstrumentID)
				return assert.NotNil(t, createdRecipeStepInstrument) && assert.NoError(t, err)
			}
			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)

			requireNotNilAndNoProblems(t, actual, err)

			// assert recipe step instrument equality
			checkRecipeStepInstrumentEquality(t, newRecipeStepInstrument, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			t.Log("cleaning up recipe step instrument")
			assert.NoError(t, testClients.main.ArchiveRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepInstrumentID))

			t.Log("cleaning up recipe step")
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStepID))

			t.Log("cleaning up recipe")
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipeID))
		}
	})
}

func (s *TestSuite) TestRecipeStepInstruments_Listing() {
	s.runForCookieClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			stopChan := make(chan bool, 1)
			notificationsChan, err := testClients.main.SubscribeToDataChangeNotifications(ctx, stopChan)
			require.NotNil(t, notificationsChan)
			require.NoError(t, err)

			var n *types.DataChangeMessage

			t.Log("creating prerequisite recipe")
			exampleRecipe := fakes.BuildFakeRecipe()
			exampleRecipeInput := fakes.BuildFakeRecipeCreationRequestInputFromRecipe(exampleRecipe)
			createdRecipeID, err := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
			require.NoError(t, err)
			t.Logf("recipe %q created", createdRecipeID)

			n = <-notificationsChan
			assert.Equal(t, n.DataType, types.RecipeDataType)
			require.NotNil(t, n.Recipe)
			checkRecipeEquality(t, exampleRecipe, n.Recipe)

			createdRecipe, err := testClients.main.GetRecipe(ctx, createdRecipeID)
			requireNotNilAndNoProblems(t, createdRecipe, err)

			t.Log("creating prerequisite recipe step")
			exampleRecipeStep := fakes.BuildFakeRecipeStep()
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeStepInput := fakes.BuildFakeRecipeStepCreationRequestInputFromRecipeStep(exampleRecipeStep)
			createdRecipeStepID, err := testClients.main.CreateRecipeStep(ctx, exampleRecipeStepInput)
			require.NoError(t, err)
			t.Logf("recipe step %q created", createdRecipeStepID)

			n = <-notificationsChan
			assert.Equal(t, n.DataType, types.RecipeStepDataType)
			require.NotNil(t, n.RecipeStep)
			checkRecipeStepEquality(t, exampleRecipeStep, n.RecipeStep)

			createdRecipeStep, err := testClients.main.GetRecipeStep(ctx, createdRecipe.ID, createdRecipeStepID)
			requireNotNilAndNoProblems(t, createdRecipeStep, err)
			require.Equal(t, createdRecipe.ID, createdRecipeStep.BelongsToRecipe)

			t.Log("creating recipe step instruments")
			var expected []*types.RecipeStepInstrument
			for i := 0; i < 5; i++ {
				exampleRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()
				exampleRecipeStepInstrument.BelongsToRecipeStep = createdRecipeStep.ID
				exampleRecipeStepInstrumentInput := fakes.BuildFakeRecipeStepInstrumentCreationRequestInputFromRecipeStepInstrument(exampleRecipeStepInstrument)
				createdRecipeStepInstrumentID, createdRecipeStepInstrumentErr := testClients.main.CreateRecipeStepInstrument(ctx, createdRecipe.ID, exampleRecipeStepInstrumentInput)
				require.NoError(t, createdRecipeStepInstrumentErr)
				t.Logf("recipe step instrument %q created", createdRecipeStepInstrumentID)

				n = <-notificationsChan
				assert.Equal(t, n.DataType, types.RecipeStepInstrumentDataType)
				require.NotNil(t, n.RecipeStepInstrument)
				checkRecipeStepInstrumentEquality(t, exampleRecipeStepInstrument, n.RecipeStepInstrument)

				createdRecipeStepInstrument, createdRecipeStepInstrumentErr := testClients.main.GetRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepInstrumentID)
				requireNotNilAndNoProblems(t, createdRecipeStepInstrument, createdRecipeStepInstrumentErr)
				require.Equal(t, createdRecipeStep.ID, createdRecipeStepInstrument.BelongsToRecipeStep)

				expected = append(expected, createdRecipeStepInstrument)
			}

			// assert recipe step instrument list equality
			actual, err := testClients.main.GetRecipeStepInstruments(ctx, createdRecipe.ID, createdRecipeStep.ID, nil)
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
				assert.NoError(t, testClients.main.ArchiveRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepInstrument.ID))
			}

			t.Log("cleaning up recipe step")
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStepID))

			t.Log("cleaning up recipe")
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipeID))
		}
	})

	s.runForPASETOClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			var checkFunc func() bool
			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating prerequisite recipe")
			exampleRecipe := fakes.BuildFakeRecipe()
			exampleRecipeInput := fakes.BuildFakeRecipeCreationRequestInputFromRecipe(exampleRecipe)
			createdRecipeID, err := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
			require.NoError(t, err)
			t.Logf("recipe %q created", createdRecipeID)

			var createdRecipe *types.Recipe
			checkFunc = func() bool {
				createdRecipe, err = testClients.main.GetRecipe(ctx, createdRecipeID)
				return assert.NotNil(t, createdRecipe) && assert.NoError(t, err)
			}
			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
			checkRecipeEquality(t, exampleRecipe, createdRecipe)

			t.Log("creating prerequisite recipe step")
			exampleRecipeStep := fakes.BuildFakeRecipeStep()
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeStepInput := fakes.BuildFakeRecipeStepCreationRequestInputFromRecipeStep(exampleRecipeStep)
			createdRecipeStepID, err := testClients.main.CreateRecipeStep(ctx, exampleRecipeStepInput)
			require.NoError(t, err)
			t.Logf("recipe step %q created", createdRecipeStepID)

			var createdRecipeStep *types.RecipeStep
			checkFunc = func() bool {
				createdRecipeStep, err = testClients.main.GetRecipeStep(ctx, createdRecipe.ID, createdRecipeStepID)
				return assert.NotNil(t, createdRecipeStep) && assert.NoError(t, err)
			}
			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
			require.Equal(t, createdRecipe.ID, createdRecipeStep.BelongsToRecipe)
			checkRecipeStepEquality(t, exampleRecipeStep, createdRecipeStep)

			t.Log("creating recipe step instruments")
			var expected []*types.RecipeStepInstrument
			for i := 0; i < 5; i++ {
				exampleRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()
				exampleRecipeStepInstrument.BelongsToRecipeStep = createdRecipeStep.ID
				exampleRecipeStepInstrumentInput := fakes.BuildFakeRecipeStepInstrumentCreationRequestInputFromRecipeStepInstrument(exampleRecipeStepInstrument)
				createdRecipeStepInstrumentID, EcreatedRecipeStepInstrumentrr := testClients.main.CreateRecipeStepInstrument(ctx, createdRecipe.ID, exampleRecipeStepInstrumentInput)
				require.NoError(t, EcreatedRecipeStepInstrumentrr)

				var createdRecipeStepInstrument *types.RecipeStepInstrument
				checkFunc = func() bool {
					createdRecipeStepInstrument, EcreatedRecipeStepInstrumentrr = testClients.main.GetRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepInstrumentID)
					return assert.NotNil(t, createdRecipeStepInstrument) && assert.NoError(t, EcreatedRecipeStepInstrumentrr)
				}
				assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
				checkRecipeStepInstrumentEquality(t, exampleRecipeStepInstrument, createdRecipeStepInstrument)

				expected = append(expected, createdRecipeStepInstrument)
			}

			// assert recipe step instrument list equality
			actual, err := testClients.main.GetRecipeStepInstruments(ctx, createdRecipe.ID, createdRecipeStep.ID, nil)
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
				assert.NoError(t, testClients.main.ArchiveRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepInstrument.ID))
			}

			t.Log("cleaning up recipe step")
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStepID))

			t.Log("cleaning up recipe")
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipeID))
		}
	})
}
