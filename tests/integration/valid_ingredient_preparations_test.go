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

func checkValidIngredientPreparationEquality(t *testing.T, expected, actual *types.ValidIngredientPreparation) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Notes, actual.Notes, "expected Notes for valid ingredient preparation %s to be %v, but it was %v", expected.ID, expected.Notes, actual.Notes)
	assert.Equal(t, expected.ValidPreparationID, actual.ValidPreparationID, "expected ValidPreparationID for valid ingredient preparation %s to be %v, but it was %v", expected.ID, expected.ValidPreparationID, actual.ValidPreparationID)
	assert.Equal(t, expected.ValidIngredientID, actual.ValidIngredientID, "expected ValidIngredientID for valid ingredient preparation %s to be %v, but it was %v", expected.ID, expected.ValidIngredientID, actual.ValidIngredientID)
	assert.NotZero(t, actual.CreatedOn)
}

// convertValidIngredientPreparationToValidIngredientPreparationUpdateInput creates an ValidIngredientPreparationUpdateRequestInput struct from a valid ingredient preparation.
func convertValidIngredientPreparationToValidIngredientPreparationUpdateInput(x *types.ValidIngredientPreparation) *types.ValidIngredientPreparationUpdateRequestInput {
	return &types.ValidIngredientPreparationUpdateRequestInput{
		Notes:              x.Notes,
		ValidPreparationID: x.ValidPreparationID,
		ValidIngredientID:  x.ValidIngredientID,
	}
}

func (s *TestSuite) TestValidIngredientPreparations_CompleteLifecycle() {
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

			t.Log("creating valid ingredient preparation")
			exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
			exampleValidIngredientPreparationInput := fakes.BuildFakeValidIngredientPreparationCreationRequestInputFromValidIngredientPreparation(exampleValidIngredientPreparation)
			createdValidIngredientPreparationID, err := testClients.main.CreateValidIngredientPreparation(ctx, exampleValidIngredientPreparationInput)
			require.NoError(t, err)
			t.Logf("valid ingredient preparation %q created", createdValidIngredientPreparationID)

			n = <-notificationsChan
			assert.Equal(t, n.DataType, types.ValidIngredientPreparationDataType)
			require.NotNil(t, n.ValidIngredientPreparation)
			checkValidIngredientPreparationEquality(t, exampleValidIngredientPreparation, n.ValidIngredientPreparation)

			createdValidIngredientPreparation, err := testClients.main.GetValidIngredientPreparation(ctx, createdValidIngredientPreparationID)
			requireNotNilAndNoProblems(t, createdValidIngredientPreparation, err)

			checkValidIngredientPreparationEquality(t, exampleValidIngredientPreparation, createdValidIngredientPreparation)

			t.Log("changing valid ingredient preparation")
			newValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
			createdValidIngredientPreparation.Update(convertValidIngredientPreparationToValidIngredientPreparationUpdateInput(newValidIngredientPreparation))
			assert.NoError(t, testClients.main.UpdateValidIngredientPreparation(ctx, createdValidIngredientPreparation))

			n = <-notificationsChan
			assert.Equal(t, n.DataType, types.ValidIngredientPreparationDataType)

			t.Log("fetching changed valid ingredient preparation")
			actual, err := testClients.main.GetValidIngredientPreparation(ctx, createdValidIngredientPreparationID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert valid ingredient preparation equality
			checkValidIngredientPreparationEquality(t, newValidIngredientPreparation, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			t.Log("cleaning up valid ingredient preparation")
			assert.NoError(t, testClients.main.ArchiveValidIngredientPreparation(ctx, createdValidIngredientPreparationID))
		}
	})

	s.runForPASETOClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			var checkFunc func() bool
			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating valid ingredient preparation")
			exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
			exampleValidIngredientPreparationInput := fakes.BuildFakeValidIngredientPreparationCreationRequestInputFromValidIngredientPreparation(exampleValidIngredientPreparation)
			createdValidIngredientPreparationID, err := testClients.main.CreateValidIngredientPreparation(ctx, exampleValidIngredientPreparationInput)
			require.NoError(t, err)
			t.Logf("valid ingredient preparation %q created", createdValidIngredientPreparationID)

			var createdValidIngredientPreparation *types.ValidIngredientPreparation
			checkFunc = func() bool {
				createdValidIngredientPreparation, err = testClients.main.GetValidIngredientPreparation(ctx, createdValidIngredientPreparationID)
				return assert.NotNil(t, createdValidIngredientPreparation) && assert.NoError(t, err)
			}
			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
			checkValidIngredientPreparationEquality(t, exampleValidIngredientPreparation, createdValidIngredientPreparation)

			// assert valid ingredient preparation equality
			checkValidIngredientPreparationEquality(t, exampleValidIngredientPreparation, createdValidIngredientPreparation)

			// change valid ingredient preparation
			newValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
			createdValidIngredientPreparation.Update(convertValidIngredientPreparationToValidIngredientPreparationUpdateInput(newValidIngredientPreparation))
			assert.NoError(t, testClients.main.UpdateValidIngredientPreparation(ctx, createdValidIngredientPreparation))

			time.Sleep(2 * time.Second)

			// retrieve changed valid ingredient preparation
			var actual *types.ValidIngredientPreparation
			checkFunc = func() bool {
				actual, err = testClients.main.GetValidIngredientPreparation(ctx, createdValidIngredientPreparationID)
				return assert.NotNil(t, createdValidIngredientPreparation) && assert.NoError(t, err)
			}
			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)

			requireNotNilAndNoProblems(t, actual, err)

			// assert valid ingredient preparation equality
			checkValidIngredientPreparationEquality(t, newValidIngredientPreparation, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			t.Log("cleaning up valid ingredient preparation")
			assert.NoError(t, testClients.main.ArchiveValidIngredientPreparation(ctx, createdValidIngredientPreparationID))
		}
	})
}

func (s *TestSuite) TestValidIngredientPreparations_Listing() {
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

			t.Log("creating valid ingredient preparations")
			var expected []*types.ValidIngredientPreparation
			for i := 0; i < 5; i++ {
				exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
				exampleValidIngredientPreparationInput := fakes.BuildFakeValidIngredientPreparationCreationRequestInputFromValidIngredientPreparation(exampleValidIngredientPreparation)
				createdValidIngredientPreparationID, createdValidIngredientPreparationErr := testClients.main.CreateValidIngredientPreparation(ctx, exampleValidIngredientPreparationInput)
				require.NoError(t, createdValidIngredientPreparationErr)
				t.Logf("valid ingredient preparation %q created", createdValidIngredientPreparationID)

				n = <-notificationsChan
				assert.Equal(t, n.DataType, types.ValidIngredientPreparationDataType)
				require.NotNil(t, n.ValidIngredientPreparation)
				checkValidIngredientPreparationEquality(t, exampleValidIngredientPreparation, n.ValidIngredientPreparation)

				createdValidIngredientPreparation, createdValidIngredientPreparationErr := testClients.main.GetValidIngredientPreparation(ctx, createdValidIngredientPreparationID)
				requireNotNilAndNoProblems(t, createdValidIngredientPreparation, createdValidIngredientPreparationErr)

				expected = append(expected, createdValidIngredientPreparation)
			}

			// assert valid ingredient preparation list equality
			actual, err := testClients.main.GetValidIngredientPreparations(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.ValidIngredientPreparations),
				"expected %d to be <= %d",
				len(expected),
				len(actual.ValidIngredientPreparations),
			)

			t.Log("cleaning up")
			for _, createdValidIngredientPreparation := range expected {
				assert.NoError(t, testClients.main.ArchiveValidIngredientPreparation(ctx, createdValidIngredientPreparation.ID))
			}
		}
	})

	s.runForPASETOClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			var checkFunc func() bool
			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating valid ingredient preparations")
			var expected []*types.ValidIngredientPreparation
			for i := 0; i < 5; i++ {
				exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
				exampleValidIngredientPreparationInput := fakes.BuildFakeValidIngredientPreparationCreationRequestInputFromValidIngredientPreparation(exampleValidIngredientPreparation)
				createdValidIngredientPreparationID, err := testClients.main.CreateValidIngredientPreparation(ctx, exampleValidIngredientPreparationInput)
				require.NoError(t, err)

				var createdValidIngredientPreparation *types.ValidIngredientPreparation
				checkFunc = func() bool {
					createdValidIngredientPreparation, err = testClients.main.GetValidIngredientPreparation(ctx, createdValidIngredientPreparationID)
					return assert.NotNil(t, createdValidIngredientPreparation) && assert.NoError(t, err)
				}
				assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
				checkValidIngredientPreparationEquality(t, exampleValidIngredientPreparation, createdValidIngredientPreparation)

				expected = append(expected, createdValidIngredientPreparation)
			}

			// assert valid ingredient preparation list equality
			actual, err := testClients.main.GetValidIngredientPreparations(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.ValidIngredientPreparations),
				"expected %d to be <= %d",
				len(expected),
				len(actual.ValidIngredientPreparations),
			)

			t.Log("cleaning up")
			for _, createdValidIngredientPreparation := range expected {
				assert.NoError(t, testClients.main.ArchiveValidIngredientPreparation(ctx, createdValidIngredientPreparation.ID))
			}
		}
	})
}
