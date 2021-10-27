package integration

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
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

			_, _, createdRecipe := createRecipeWithNotificationChannel(ctx, t, notificationsChan, testClients.main)

			var createdRecipeStepID string
			for _, step := range createdRecipe.Steps {
				createdRecipeStepID = step.ID
				break
			}

			t.Log("creating recipe step instrument")
			exampleRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()
			exampleRecipeStepInstrument.BelongsToRecipeStep = createdRecipeStepID
			exampleRecipeStepInstrumentInput := fakes.BuildFakeRecipeStepInstrumentCreationRequestInputFromRecipeStepInstrument(exampleRecipeStepInstrument)
			createdRecipeStepInstrumentID, err := testClients.main.CreateRecipeStepInstrument(ctx, createdRecipe.ID, exampleRecipeStepInstrumentInput)
			require.NoError(t, err)
			t.Logf("recipe step instrument %q created", createdRecipeStepInstrumentID)

			n = <-notificationsChan
			assert.Equal(t, n.DataType, types.RecipeStepInstrumentDataType)
			require.NotNil(t, n.RecipeStepInstrument)
			checkRecipeStepInstrumentEquality(t, exampleRecipeStepInstrument, n.RecipeStepInstrument)

			createdRecipeStepInstrument, err := testClients.main.GetRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepInstrumentID)
			requireNotNilAndNoProblems(t, createdRecipeStepInstrument, err)
			require.Equal(t, createdRecipeStepID, createdRecipeStepInstrument.BelongsToRecipeStep)

			checkRecipeStepInstrumentEquality(t, exampleRecipeStepInstrument, createdRecipeStepInstrument)

			t.Log("changing recipe step instrument")
			newRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()
			newRecipeStepInstrument.BelongsToRecipeStep = createdRecipeStepID
			createdRecipeStepInstrument.Update(convertRecipeStepInstrumentToRecipeStepInstrumentUpdateInput(newRecipeStepInstrument))
			assert.NoError(t, testClients.main.UpdateRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepInstrument))

			n = <-notificationsChan
			assert.Equal(t, n.DataType, types.RecipeStepInstrumentDataType)

			t.Log("fetching changed recipe step instrument")
			actual, err := testClients.main.GetRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepInstrumentID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert recipe step instrument equality
			checkRecipeStepInstrumentEquality(t, newRecipeStepInstrument, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			t.Log("cleaning up recipe step instrument")
			assert.NoError(t, testClients.main.ArchiveRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepInstrumentID))

			t.Log("cleaning up recipe step")
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStepID))

			t.Log("cleaning up recipe")
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})

	s.runForPASETOClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			var checkFunc func() bool
			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			_, _, createdRecipe := createRecipeWhilePolling(ctx, t, testClients.main)

			var createdRecipeStepID string
			for _, step := range createdRecipe.Steps {
				createdRecipeStepID = step.ID
				break
			}

			t.Log("creating recipe step instrument")
			exampleRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()
			exampleRecipeStepInstrument.BelongsToRecipeStep = createdRecipeStepID
			exampleRecipeStepInstrumentInput := fakes.BuildFakeRecipeStepInstrumentCreationRequestInputFromRecipeStepInstrument(exampleRecipeStepInstrument)
			createdRecipeStepInstrumentID, err := testClients.main.CreateRecipeStepInstrument(ctx, createdRecipe.ID, exampleRecipeStepInstrumentInput)
			require.NoError(t, err)
			t.Logf("recipe step instrument %q created", createdRecipeStepInstrumentID)

			var createdRecipeStepInstrument *types.RecipeStepInstrument
			checkFunc = func() bool {
				createdRecipeStepInstrument, err = testClients.main.GetRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepInstrumentID)
				return assert.NotNil(t, createdRecipeStepInstrument) && assert.NoError(t, err)
			}
			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
			require.Equal(t, createdRecipeStepID, createdRecipeStepInstrument.BelongsToRecipeStep)
			checkRecipeStepInstrumentEquality(t, exampleRecipeStepInstrument, createdRecipeStepInstrument)

			// assert recipe step instrument equality
			checkRecipeStepInstrumentEquality(t, exampleRecipeStepInstrument, createdRecipeStepInstrument)

			// change recipe step instrument
			newRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()
			newRecipeStepInstrument.BelongsToRecipeStep = createdRecipeStepID
			createdRecipeStepInstrument.Update(convertRecipeStepInstrumentToRecipeStepInstrumentUpdateInput(newRecipeStepInstrument))
			assert.NoError(t, testClients.main.UpdateRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepInstrument))

			time.Sleep(2 * time.Second)

			// retrieve changed recipe step instrument
			var actual *types.RecipeStepInstrument
			checkFunc = func() bool {
				actual, err = testClients.main.GetRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepInstrumentID)
				return assert.NotNil(t, createdRecipeStepInstrument) && assert.NoError(t, err)
			}
			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)

			requireNotNilAndNoProblems(t, actual, err)

			// assert recipe step instrument equality
			checkRecipeStepInstrumentEquality(t, newRecipeStepInstrument, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			t.Log("cleaning up recipe step instrument")
			assert.NoError(t, testClients.main.ArchiveRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepInstrumentID))

			t.Log("cleaning up recipe step")
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStepID))

			t.Log("cleaning up recipe")
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
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

			_, _, createdRecipe := createRecipeWithNotificationChannel(ctx, t, notificationsChan, testClients.main)

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
				createdRecipeStepInstrumentID, createdRecipeStepInstrumentErr := testClients.main.CreateRecipeStepInstrument(ctx, createdRecipe.ID, exampleRecipeStepInstrumentInput)
				require.NoError(t, createdRecipeStepInstrumentErr)
				t.Logf("recipe step instrument %q created", createdRecipeStepInstrumentID)

				n = <-notificationsChan
				assert.Equal(t, n.DataType, types.RecipeStepInstrumentDataType)
				require.NotNil(t, n.RecipeStepInstrument)
				checkRecipeStepInstrumentEquality(t, exampleRecipeStepInstrument, n.RecipeStepInstrument)

				createdRecipeStepInstrument, createdRecipeStepInstrumentErr := testClients.main.GetRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepInstrumentID)
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

	s.runForPASETOClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			var checkFunc func() bool
			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			_, _, createdRecipe := createRecipeWhilePolling(ctx, t, testClients.main)

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
				createdRecipeStepInstrumentID, EcreatedRecipeStepInstrumentrr := testClients.main.CreateRecipeStepInstrument(ctx, createdRecipe.ID, exampleRecipeStepInstrumentInput)
				require.NoError(t, EcreatedRecipeStepInstrumentrr)

				var createdRecipeStepInstrument *types.RecipeStepInstrument
				checkFunc = func() bool {
					createdRecipeStepInstrument, EcreatedRecipeStepInstrumentrr = testClients.main.GetRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepInstrumentID)
					return assert.NotNil(t, createdRecipeStepInstrument) && assert.NoError(t, EcreatedRecipeStepInstrumentrr)
				}
				assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
				checkRecipeStepInstrumentEquality(t, exampleRecipeStepInstrument, createdRecipeStepInstrument)

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
