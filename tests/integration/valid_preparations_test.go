package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
)

func checkValidPreparationEquality(t *testing.T, expected, actual *types.ValidPreparation) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for valid preparation %s to be %v, but it was %v", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.Description, actual.Description, "expected Description for valid preparation %s to be %v, but it was %v", expected.ID, expected.Description, actual.Description)
	assert.Equal(t, expected.IconPath, actual.IconPath, "expected IconPath for valid preparation %s to be %v, but it was %v", expected.ID, expected.IconPath, actual.IconPath)
	assert.NotZero(t, actual.CreatedOn)
}

// convertValidPreparationToValidPreparationUpdateInput creates an ValidPreparationUpdateRequestInput struct from a valid preparation.
func convertValidPreparationToValidPreparationUpdateInput(x *types.ValidPreparation) *types.ValidPreparationUpdateRequestInput {
	return &types.ValidPreparationUpdateRequestInput{
		Name:        x.Name,
		Description: x.Description,
		IconPath:    x.IconPath,
	}
}

func (s *TestSuite) TestValidPreparations_CompleteLifecycle() {
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

			t.Log("creating valid preparation")
			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakes.BuildFakeValidPreparationCreationRequestInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparationID, err := testClients.main.CreateValidPreparation(ctx, exampleValidPreparationInput)
			require.NoError(t, err)
			t.Logf("valid preparation %q created", createdValidPreparationID)

			n = <-notificationsChan
			assert.Equal(t, n.DataType, types.ValidPreparationDataType)
			require.NotNil(t, n.ValidPreparation)
			checkValidPreparationEquality(t, exampleValidPreparation, n.ValidPreparation)

			createdValidPreparation, err := testClients.main.GetValidPreparation(ctx, createdValidPreparationID)
			requireNotNilAndNoProblems(t, createdValidPreparation, err)
			checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

			t.Log("changing valid preparation")
			newValidPreparation := fakes.BuildFakeValidPreparation()
			createdValidPreparation.Update(convertValidPreparationToValidPreparationUpdateInput(newValidPreparation))
			assert.NoError(t, testClients.main.UpdateValidPreparation(ctx, createdValidPreparation))

			n = <-notificationsChan
			assert.Equal(t, n.DataType, types.ValidPreparationDataType)

			t.Log("fetching changed valid preparation")
			actual, err := testClients.main.GetValidPreparation(ctx, createdValidPreparationID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert valid preparation equality
			checkValidPreparationEquality(t, newValidPreparation, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			t.Log("cleaning up valid preparation")
			assert.NoError(t, testClients.main.ArchiveValidPreparation(ctx, createdValidPreparationID))
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
			createdValidPreparationID, err := testClients.main.CreateValidPreparation(ctx, exampleValidPreparationInput)
			require.NoError(t, err)
			t.Logf("valid preparation %q created", createdValidPreparationID)

			var createdValidPreparation *types.ValidPreparation
			checkFunc = func() bool {
				createdValidPreparation, err = testClients.main.GetValidPreparation(ctx, createdValidPreparationID)
				return assert.NotNil(t, createdValidPreparation) && assert.NoError(t, err)
			}
			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
			checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

			// change valid preparation
			newValidPreparation := fakes.BuildFakeValidPreparation()
			createdValidPreparation.Update(convertValidPreparationToValidPreparationUpdateInput(newValidPreparation))
			assert.NoError(t, testClients.main.UpdateValidPreparation(ctx, createdValidPreparation))

			time.Sleep(2 * time.Second)

			// retrieve changed valid preparation
			var actual *types.ValidPreparation
			checkFunc = func() bool {
				actual, err = testClients.main.GetValidPreparation(ctx, createdValidPreparationID)
				return assert.NotNil(t, createdValidPreparation) && assert.NoError(t, err)
			}
			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)

			requireNotNilAndNoProblems(t, actual, err)

			// assert valid preparation equality
			checkValidPreparationEquality(t, newValidPreparation, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			t.Log("cleaning up valid preparation")
			assert.NoError(t, testClients.main.ArchiveValidPreparation(ctx, createdValidPreparationID))
		}
	})
}

func (s *TestSuite) TestValidPreparations_Listing() {
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

			t.Log("creating valid preparations")
			var expected []*types.ValidPreparation
			for i := 0; i < 5; i++ {
				exampleValidPreparation := fakes.BuildFakeValidPreparation()
				exampleValidPreparationInput := fakes.BuildFakeValidPreparationCreationRequestInputFromValidPreparation(exampleValidPreparation)
				createdValidPreparationID, createdValidPreparationErr := testClients.main.CreateValidPreparation(ctx, exampleValidPreparationInput)
				require.NoError(t, createdValidPreparationErr)
				t.Logf("valid preparation %q created", createdValidPreparationID)

				n = <-notificationsChan
				assert.Equal(t, n.DataType, types.ValidPreparationDataType)
				require.NotNil(t, n.ValidPreparation)
				checkValidPreparationEquality(t, exampleValidPreparation, n.ValidPreparation)

				createdValidPreparation, createdValidPreparationErr := testClients.main.GetValidPreparation(ctx, createdValidPreparationID)
				requireNotNilAndNoProblems(t, createdValidPreparation, createdValidPreparationErr)

				expected = append(expected, createdValidPreparation)
			}

			// assert valid preparation list equality
			actual, err := testClients.main.GetValidPreparations(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.ValidPreparations),
				"expected %d to be <= %d",
				len(expected),
				len(actual.ValidPreparations),
			)

			t.Log("cleaning up")
			for _, createdValidPreparation := range expected {
				assert.NoError(t, testClients.main.ArchiveValidPreparation(ctx, createdValidPreparation.ID))
			}
		}
	})

	s.runForPASETOClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			var checkFunc func() bool
			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating valid preparations")
			var expected []*types.ValidPreparation
			for i := 0; i < 5; i++ {
				exampleValidPreparation := fakes.BuildFakeValidPreparation()
				exampleValidPreparationInput := fakes.BuildFakeValidPreparationCreationRequestInputFromValidPreparation(exampleValidPreparation)
				createdValidPreparationID, err := testClients.main.CreateValidPreparation(ctx, exampleValidPreparationInput)
				require.NoError(t, err)

				var createdValidPreparation *types.ValidPreparation
				checkFunc = func() bool {
					createdValidPreparation, err = testClients.main.GetValidPreparation(ctx, createdValidPreparationID)
					return assert.NotNil(t, createdValidPreparation) && assert.NoError(t, err)
				}
				assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
				checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

				expected = append(expected, createdValidPreparation)
			}

			// assert valid preparation list equality
			actual, err := testClients.main.GetValidPreparations(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.ValidPreparations),
				"expected %d to be <= %d",
				len(expected),
				len(actual.ValidPreparations),
			)

			t.Log("cleaning up")
			for _, createdValidPreparation := range expected {
				assert.NoError(t, testClients.main.ArchiveValidPreparation(ctx, createdValidPreparation.ID))
			}
		}
	})
}

func (s *TestSuite) TestValidPreparations_Searching() {
	s.runForCookieClient("should be able to be search for valid preparations", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			stopChan := make(chan bool, 1)
			notificationsChan, err := testClients.main.SubscribeToDataChangeNotifications(ctx, stopChan)
			require.NotNil(t, notificationsChan)
			require.NoError(t, err)

			var n *types.DataChangeMessage

			t.Log("creating valid preparations")
			var expected []*types.ValidPreparation
			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			searchQuery := exampleValidPreparation.Name
			for i := 0; i < 5; i++ {
				exampleValidPreparation.Name = fmt.Sprintf("%s %d", searchQuery, i)
				exampleValidPreparationInput := fakes.BuildFakeValidPreparationCreationRequestInputFromValidPreparation(exampleValidPreparation)
				createdValidPreparationID, createdValidPreparationErr := testClients.main.CreateValidPreparation(ctx, exampleValidPreparationInput)
				require.NoError(t, createdValidPreparationErr)
				t.Logf("valid preparation %q created", createdValidPreparationID)

				n = <-notificationsChan
				assert.Equal(t, n.DataType, types.ValidPreparationDataType)
				require.NotNil(t, n.ValidPreparation)
				checkValidPreparationEquality(t, exampleValidPreparation, n.ValidPreparation)

				createdValidPreparation, createdValidPreparationErr := testClients.main.GetValidPreparation(ctx, createdValidPreparationID)
				requireNotNilAndNoProblems(t, createdValidPreparation, createdValidPreparationErr)

				expected = append(expected, createdValidPreparation)
			}

			exampleLimit := uint8(20)

			// give the index a moment
			time.Sleep(3 * time.Second)

			// assert valid preparation list equality
			actual, err := testClients.main.SearchValidPreparations(ctx, searchQuery, exampleLimit)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual),
				"expected %d to be <= %d",
				len(expected),
				len(actual),
			)

			t.Log("cleaning up")
			for _, createdValidPreparation := range expected {
				assert.NoError(t, testClients.main.ArchiveValidPreparation(ctx, createdValidPreparation.ID))
			}
		}
	})

	s.runForPASETOClient("should be able to be search for valid preparations", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			var checkFunc func() bool
			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating valid preparations")
			var expected []*types.ValidPreparation
			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			searchQuery := exampleValidPreparation.Name
			for i := 0; i < 5; i++ {
				exampleValidPreparation.Name = fmt.Sprintf("%s %d", searchQuery, i)
				exampleValidPreparationInput := fakes.BuildFakeValidPreparationCreationRequestInputFromValidPreparation(exampleValidPreparation)
				createdValidPreparationID, err := testClients.main.CreateValidPreparation(ctx, exampleValidPreparationInput)
				require.NoError(t, err)
				t.Logf("valid preparation %q created", createdValidPreparationID)

				var createdValidPreparation *types.ValidPreparation
				checkFunc = func() bool {
					createdValidPreparation, err = testClients.main.GetValidPreparation(ctx, createdValidPreparationID)
					return assert.NotNil(t, createdValidPreparation) && assert.NoError(t, err)
				}
				assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
				requireNotNilAndNoProblems(t, createdValidPreparation, err)

				expected = append(expected, createdValidPreparation)
			}

			exampleLimit := uint8(20)
			time.Sleep(2 * time.Second) // give the index a moment

			// assert valid preparation list equality
			actual, err := testClients.main.SearchValidPreparations(ctx, searchQuery, exampleLimit)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual),
				"expected %d to be <= %d",
				len(expected),
				len(actual),
			)

			t.Log("cleaning up")
			for _, createdValidPreparation := range expected {
				assert.NoError(t, testClients.main.ArchiveValidPreparation(ctx, createdValidPreparation.ID))
			}
		}
	})
}
