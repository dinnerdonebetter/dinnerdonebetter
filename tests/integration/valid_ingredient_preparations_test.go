package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prixfixeco/api_server/pkg/types"
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

/*

func (s *TestSuite) TestValidIngredientPreparations_CompleteLifecycle() {
	s.runForCookieClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			stopChan := make(chan bool, 1)
			notificationsChan, err := testClients.main.SubscribeToNotifications(ctx, stopChan)
			require.NotNil(t, notificationsChan)
			require.NoError(t, err)

			var n *types.DataChangeMessage

			t.Log("creating prerequisite valid preparation")
			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakes.BuildFakeValidPreparationCreationRequestInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparation, err := testClients.main.CreateValidPreparation(ctx, exampleValidPreparationInput)
			require.NoError(t, err)
			t.Logf("valid preparation %q created", createdValidPreparation.ID)

			n = <-notificationsChan
			assert.Equal(t, types.ValidPreparationDataType, n.DataType)
			require.NotNil(t, n.ValidPreparation)
			checkValidPreparationEquality(t, exampleValidPreparation, n.ValidPreparation)

			createdValidPreparation, err = testClients.main.GetValidPreparation(ctx, createdValidPreparation.ID)
			requireNotNilAndNoProblems(t, createdValidPreparation, err)
			checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

			t.Log("creating prerequisite valid ingredient")
			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakes.BuildFakeValidIngredientCreationRequestInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredient, err := testClients.main.CreateValidIngredient(ctx, exampleValidIngredientInput)
			require.NoError(t, err)

			n = <-notificationsChan
			assert.Equal(t, types.ValidIngredientDataType, n.DataType)
			require.NotNil(t, n.ValidIngredient)
			checkValidIngredientEquality(t, exampleValidIngredient, n.ValidIngredient)

			createdValidIngredient, err = testClients.main.GetValidIngredient(ctx, createdValidIngredient.ID)
			requireNotNilAndNoProblems(t, createdValidIngredient, err)
			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)
			t.Logf("valid ingredient %q created", createdValidIngredient.ID)

			t.Log("creating valid ingredient preparation")
			exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
			exampleValidIngredientPreparation.ValidIngredientID = createdValidIngredient.ID
			exampleValidIngredientPreparation.ValidPreparationID = createdValidPreparation.ID
			exampleValidIngredientPreparationInput := fakes.BuildFakeValidIngredientPreparationCreationRequestInputFromValidIngredientPreparation(exampleValidIngredientPreparation)
			createdValidIngredientPreparation, err := testClients.main.CreateValidIngredientPreparation(ctx, exampleValidIngredientPreparationInput)
			require.NoError(t, err)
			t.Logf("valid ingredient preparation %q created", createdValidIngredientPreparation.ID)

			n = <-notificationsChan
			assert.Equal(t, types.ValidIngredientPreparationDataType, n.DataType)
			require.NotNil(t, n.ValidIngredientPreparation)
			checkValidIngredientPreparationEquality(t, exampleValidIngredientPreparation, n.ValidIngredientPreparation)

			createdValidIngredientPreparation, err = testClients.main.GetValidIngredientPreparation(ctx, createdValidIngredientPreparation.ID)
			requireNotNilAndNoProblems(t, createdValidIngredientPreparation, err)

			checkValidIngredientPreparationEquality(t, exampleValidIngredientPreparation, createdValidIngredientPreparation)

			t.Log("changing valid ingredient preparation")
			newValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
			newValidIngredientPreparation.ValidIngredientID = createdValidIngredient.ID
			newValidIngredientPreparation.ValidPreparationID = createdValidPreparation.ID
			createdValidIngredientPreparation.Update(convertValidIngredientPreparationToValidIngredientPreparationUpdateInput(newValidIngredientPreparation))
			assert.NoError(t, testClients.main.UpdateValidIngredientPreparation(ctx, createdValidIngredientPreparation))

			n = <-notificationsChan
			assert.Equal(t, types.ValidIngredientPreparationDataType, n.DataType)

			t.Log("fetching changed valid ingredient preparation")
			actual, err := testClients.main.GetValidIngredientPreparation(ctx, createdValidIngredientPreparation.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert valid ingredient preparation equality
			checkValidIngredientPreparationEquality(t, newValidIngredientPreparation, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			t.Log("cleaning up valid ingredient preparation")
			assert.NoError(t, testClients.main.ArchiveValidIngredientPreparation(ctx, createdValidIngredientPreparation.ID))
		}
	})

	s.runForPASETOClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			var checkFunc func() bool
			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating valid preparation")
			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakes.BuildFakeValidPreparationCreationRequestInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparation, validPreparationCreationErr := testClients.main.CreateValidPreparation(ctx, exampleValidPreparationInput)
			require.NoError(t, validPreparationCreationErr)
			checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)
			t.Logf("valid preparation %q created", createdValidPreparation.ID)

			t.Log("creating valid ingredient")
			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakes.BuildFakeValidIngredientCreationRequestInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredient, validIngredientCreationErr := testClients.main.CreateValidIngredient(ctx, exampleValidIngredientInput)
			require.NoError(t, validIngredientCreationErr)
			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)
			t.Logf("valid ingredient %q created", createdValidIngredient.ID)

			t.Log("creating valid ingredient preparation")
			exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
			exampleValidIngredientPreparation.ValidPreparationID = createdValidPreparation.ID
			exampleValidIngredientPreparation.ValidIngredientID = createdValidIngredient.ID
			exampleValidIngredientPreparationInput := fakes.BuildFakeValidIngredientPreparationCreationRequestInputFromValidIngredientPreparation(exampleValidIngredientPreparation)
			createdValidIngredientPreparation, err := testClients.main.CreateValidIngredientPreparation(ctx, exampleValidIngredientPreparationInput)
			require.NoError(t, err)
			t.Logf("valid ingredient preparation %q created", createdValidIngredientPreparation.ID)
			checkValidIngredientPreparationEquality(t, exampleValidIngredientPreparation, createdValidIngredientPreparation)

			// assert valid ingredient preparation equality
			checkValidIngredientPreparationEquality(t, exampleValidIngredientPreparation, createdValidIngredientPreparation)

			// change valid ingredient preparation
			newValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
			newValidIngredientPreparation.ValidIngredientID = createdValidIngredient.ID
			newValidIngredientPreparation.ValidPreparationID = createdValidPreparation.ID
			createdValidIngredientPreparation.Update(convertValidIngredientPreparationToValidIngredientPreparationUpdateInput(newValidIngredientPreparation))
			assert.NoError(t, testClients.main.UpdateValidIngredientPreparation(ctx, createdValidIngredientPreparation))

			time.Sleep(2 * time.Second)

			// retrieve changed valid ingredient preparation
			var actual *types.ValidIngredientPreparation
			checkFunc = func() bool {
				actual, err = testClients.main.GetValidIngredientPreparation(ctx, createdValidIngredientPreparation.ID)
				return assert.NotNil(t, createdValidIngredientPreparation) && assert.NoError(t, err)
			}
			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)

			requireNotNilAndNoProblems(t, actual, err)

			// assert valid ingredient preparation equality
			checkValidIngredientPreparationEquality(t, newValidIngredientPreparation, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			t.Log("cleaning up valid ingredient preparation")
			assert.NoError(t, testClients.main.ArchiveValidIngredientPreparation(ctx, createdValidIngredientPreparation.ID))
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
			notificationsChan, err := testClients.main.SubscribeToNotifications(ctx, stopChan)
			require.NotNil(t, notificationsChan)
			require.NoError(t, err)

			var n *types.DataChangeMessage

			t.Log("creating valid ingredient preparations")
			var expected []*types.ValidIngredientPreparation
			for i := 0; i < 5; i++ {
				t.Log("creating prerequisite valid preparation")
				exampleValidPreparation := fakes.BuildFakeValidPreparation()
				exampleValidPreparationInput := fakes.BuildFakeValidPreparationCreationRequestInputFromValidPreparation(exampleValidPreparation)
				createdValidPreparation, err := testClients.main.CreateValidPreparation(ctx, exampleValidPreparationInput)
				require.NoError(t, err)
				t.Logf("valid preparation %q created", createdValidPreparation.ID)

				n = <-notificationsChan
				assert.Equal(t, types.ValidPreparationDataType, n.DataType)
				require.NotNil(t, n.ValidPreparation)
				checkValidPreparationEquality(t, exampleValidPreparation, n.ValidPreparation)

				createdValidPreparation, err = testClients.main.GetValidPreparation(ctx, createdValidPreparation.ID)
				requireNotNilAndNoProblems(t, createdValidPreparation, err)
				checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

				t.Log("creating prerequisite valid ingredient")
				exampleValidIngredient := fakes.BuildFakeValidIngredient()
				exampleValidIngredientInput := fakes.BuildFakeValidIngredientCreationRequestInputFromValidIngredient(exampleValidIngredient)
				createdValidIngredient, err := testClients.main.CreateValidIngredient(ctx, exampleValidIngredientInput)
				require.NoError(t, err)

				n = <-notificationsChan
				assert.Equal(t, types.ValidIngredientDataType, n.DataType)
				require.NotNil(t, n.ValidIngredient)
				checkValidIngredientEquality(t, exampleValidIngredient, n.ValidIngredient)

				createdValidIngredient, err = testClients.main.GetValidIngredient(ctx, createdValidIngredient.ID)
				requireNotNilAndNoProblems(t, createdValidIngredient, err)
				checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)
				t.Logf("valid ingredient %q created", createdValidIngredient.ID)

				exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
				exampleValidIngredientPreparation.ValidIngredientID = createdValidIngredient.ID
				exampleValidIngredientPreparation.ValidPreparationID = createdValidPreparation.ID
				exampleValidIngredientPreparationInput := fakes.BuildFakeValidIngredientPreparationCreationRequestInputFromValidIngredientPreparation(exampleValidIngredientPreparation)
				createdValidIngredientPreparation, createdValidIngredientPreparationErr := testClients.main.CreateValidIngredientPreparation(ctx, exampleValidIngredientPreparationInput)
				require.NoError(t, createdValidIngredientPreparationErr)

				n = <-notificationsChan
				assert.Equal(t, types.ValidIngredientPreparationDataType, n.DataType)
				require.NotNil(t, n.ValidIngredientPreparation)
				checkValidIngredientPreparationEquality(t, exampleValidIngredientPreparation, n.ValidIngredientPreparation)

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

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating valid ingredient preparations")
			var expected []*types.ValidIngredientPreparation
			for i := 0; i < 5; i++ {
				t.Log("creating valid preparation")
				exampleValidPreparation := fakes.BuildFakeValidPreparation()
				exampleValidPreparationInput := fakes.BuildFakeValidPreparationCreationRequestInputFromValidPreparation(exampleValidPreparation)
				createdValidPreparation, validPreparationCreationErr := testClients.main.CreateValidPreparation(ctx, exampleValidPreparationInput)
				require.NoError(t, validPreparationCreationErr)
				checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)
				t.Logf("valid preparation %q created", createdValidPreparation.ID)

				t.Log("creating valid ingredient")
				exampleValidIngredient := fakes.BuildFakeValidIngredient()
				exampleValidIngredientInput := fakes.BuildFakeValidIngredientCreationRequestInputFromValidIngredient(exampleValidIngredient)
				createdValidIngredient, validIngredientCreationErr := testClients.main.CreateValidIngredient(ctx, exampleValidIngredientInput)
				require.NoError(t, validIngredientCreationErr)
				checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)
				t.Logf("valid ingredient %q created", createdValidIngredient.ID)

				exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
				exampleValidIngredientPreparation.ValidPreparationID = createdValidPreparation.ID
				exampleValidIngredientPreparation.ValidIngredientID = createdValidIngredient.ID
				exampleValidIngredientPreparationInput := fakes.BuildFakeValidIngredientPreparationCreationRequestInputFromValidIngredientPreparation(exampleValidIngredientPreparation)
				createdValidIngredientPreparation, err := testClients.main.CreateValidIngredientPreparation(ctx, exampleValidIngredientPreparationInput)
				require.NoError(t, err)
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

*/
