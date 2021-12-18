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

func checkValidInstrumentEquality(t *testing.T, expected, actual *types.ValidInstrument) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for valid instrument %s to be %v, but it was %v", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.Variant, actual.Variant, "expected Variant for valid instrument %s to be %v, but it was %v", expected.ID, expected.Variant, actual.Variant)
	assert.Equal(t, expected.Description, actual.Description, "expected Description for valid instrument %s to be %v, but it was %v", expected.ID, expected.Description, actual.Description)
	assert.Equal(t, expected.IconPath, actual.IconPath, "expected IconPath for valid instrument %s to be %v, but it was %v", expected.ID, expected.IconPath, actual.IconPath)
	assert.NotZero(t, actual.CreatedOn)
}

// convertValidInstrumentToValidInstrumentUpdateInput creates an ValidInstrumentUpdateRequestInput struct from a valid instrument.
func convertValidInstrumentToValidInstrumentUpdateInput(x *types.ValidInstrument) *types.ValidInstrumentUpdateRequestInput {
	return &types.ValidInstrumentUpdateRequestInput{
		Name:        x.Name,
		Variant:     x.Variant,
		Description: x.Description,
		IconPath:    x.IconPath,
	}
}

func (s *TestSuite) TestValidInstruments_CompleteLifecycle() {
	s.runForCookieClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			stopChan := make(chan bool, 1)
			notificationsChan, err := testClients.admin.SubscribeToNotifications(ctx, stopChan)
			require.NotNil(t, notificationsChan)
			require.NoError(t, err)

			var n *types.DataChangeMessage

			t.Log("creating valid instrument")
			exampleValidInstrument := fakes.BuildFakeValidInstrument()
			exampleValidInstrumentInput := fakes.BuildFakeValidInstrumentCreationRequestInputFromValidInstrument(exampleValidInstrument)
			createdValidInstrumentID, err := testClients.admin.CreateValidInstrument(ctx, exampleValidInstrumentInput)
			require.NoError(t, err)
			t.Logf("valid instrument %q created", createdValidInstrumentID)

			n = <-notificationsChan
			assert.Equal(t, types.ValidInstrumentDataType, n.DataType)
			require.NotNil(t, n.ValidInstrument)
			checkValidInstrumentEquality(t, exampleValidInstrument, n.ValidInstrument)

			createdValidInstrument, err := testClients.admin.GetValidInstrument(ctx, createdValidInstrumentID)
			requireNotNilAndNoProblems(t, createdValidInstrument, err)

			checkValidInstrumentEquality(t, exampleValidInstrument, createdValidInstrument)

			t.Log("changing valid instrument")
			newValidInstrument := fakes.BuildFakeValidInstrument()
			createdValidInstrument.Update(convertValidInstrumentToValidInstrumentUpdateInput(newValidInstrument))
			assert.NoError(t, testClients.admin.UpdateValidInstrument(ctx, createdValidInstrument))

			n = <-notificationsChan
			assert.Equal(t, types.ValidInstrumentDataType, n.DataType)

			t.Log("fetching changed valid instrument")
			actual, err := testClients.admin.GetValidInstrument(ctx, createdValidInstrumentID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert valid instrument equality
			checkValidInstrumentEquality(t, newValidInstrument, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			t.Log("cleaning up valid instrument")
			assert.NoError(t, testClients.admin.ArchiveValidInstrument(ctx, createdValidInstrumentID))
		}
	})

	s.runForPASETOClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			var checkFunc func() bool
			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating valid instrument")
			exampleValidInstrument := fakes.BuildFakeValidInstrument()
			exampleValidInstrumentInput := fakes.BuildFakeValidInstrumentCreationRequestInputFromValidInstrument(exampleValidInstrument)
			createdValidInstrumentID, err := testClients.admin.CreateValidInstrument(ctx, exampleValidInstrumentInput)
			require.NoError(t, err)
			t.Logf("valid instrument %q created", createdValidInstrumentID)

			var createdValidInstrument *types.ValidInstrument
			checkFunc = func() bool {
				createdValidInstrument, err = testClients.admin.GetValidInstrument(ctx, createdValidInstrumentID)
				return assert.NotNil(t, createdValidInstrument) && assert.NoError(t, err)
			}
			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
			checkValidInstrumentEquality(t, exampleValidInstrument, createdValidInstrument)

			// assert valid instrument equality
			checkValidInstrumentEquality(t, exampleValidInstrument, createdValidInstrument)

			// change valid instrument
			newValidInstrument := fakes.BuildFakeValidInstrument()
			createdValidInstrument.Update(convertValidInstrumentToValidInstrumentUpdateInput(newValidInstrument))
			assert.NoError(t, testClients.admin.UpdateValidInstrument(ctx, createdValidInstrument))

			time.Sleep(2 * time.Second)

			// retrieve changed valid instrument
			var actual *types.ValidInstrument
			checkFunc = func() bool {
				actual, err = testClients.admin.GetValidInstrument(ctx, createdValidInstrumentID)
				return assert.NotNil(t, createdValidInstrument) && assert.NoError(t, err)
			}
			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)

			requireNotNilAndNoProblems(t, actual, err)

			// assert valid instrument equality
			checkValidInstrumentEquality(t, newValidInstrument, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			t.Log("cleaning up valid instrument")
			assert.NoError(t, testClients.admin.ArchiveValidInstrument(ctx, createdValidInstrumentID))
		}
	})
}

func (s *TestSuite) TestValidInstruments_Listing() {
	s.runForCookieClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			stopChan := make(chan bool, 1)
			notificationsChan, err := testClients.admin.SubscribeToNotifications(ctx, stopChan)
			require.NotNil(t, notificationsChan)
			require.NoError(t, err)

			var n *types.DataChangeMessage

			t.Log("creating valid instruments")
			var expected []*types.ValidInstrument
			for i := 0; i < 5; i++ {
				exampleValidInstrument := fakes.BuildFakeValidInstrument()
				exampleValidInstrumentInput := fakes.BuildFakeValidInstrumentCreationRequestInputFromValidInstrument(exampleValidInstrument)
				createdValidInstrumentID, createdValidInstrumentErr := testClients.admin.CreateValidInstrument(ctx, exampleValidInstrumentInput)
				require.NoError(t, createdValidInstrumentErr)
				t.Logf("valid instrument %q created", createdValidInstrumentID)

				n = <-notificationsChan
				assert.Equal(t, types.ValidInstrumentDataType, n.DataType)
				require.NotNil(t, n.ValidInstrument)
				checkValidInstrumentEquality(t, exampleValidInstrument, n.ValidInstrument)

				createdValidInstrument, createdValidInstrumentErr := testClients.admin.GetValidInstrument(ctx, createdValidInstrumentID)
				requireNotNilAndNoProblems(t, createdValidInstrument, createdValidInstrumentErr)

				expected = append(expected, createdValidInstrument)
			}

			// assert valid instrument list equality
			actual, err := testClients.admin.GetValidInstruments(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.ValidInstruments),
				"expected %d to be <= %d",
				len(expected),
				len(actual.ValidInstruments),
			)

			t.Log("cleaning up")
			for _, createdValidInstrument := range expected {
				assert.NoError(t, testClients.admin.ArchiveValidInstrument(ctx, createdValidInstrument.ID))
			}
		}
	})

	s.runForPASETOClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			var checkFunc func() bool
			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating valid instruments")
			var expected []*types.ValidInstrument
			for i := 0; i < 5; i++ {
				exampleValidInstrument := fakes.BuildFakeValidInstrument()
				exampleValidInstrumentInput := fakes.BuildFakeValidInstrumentCreationRequestInputFromValidInstrument(exampleValidInstrument)
				createdValidInstrumentID, err := testClients.admin.CreateValidInstrument(ctx, exampleValidInstrumentInput)
				require.NoError(t, err)

				var createdValidInstrument *types.ValidInstrument
				checkFunc = func() bool {
					createdValidInstrument, err = testClients.admin.GetValidInstrument(ctx, createdValidInstrumentID)
					return assert.NotNil(t, createdValidInstrument) && assert.NoError(t, err)
				}
				assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
				checkValidInstrumentEquality(t, exampleValidInstrument, createdValidInstrument)

				expected = append(expected, createdValidInstrument)
			}

			// assert valid instrument list equality
			actual, err := testClients.admin.GetValidInstruments(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.ValidInstruments),
				"expected %d to be <= %d",
				len(expected),
				len(actual.ValidInstruments),
			)

			t.Log("cleaning up")
			for _, createdValidInstrument := range expected {
				assert.NoError(t, testClients.admin.ArchiveValidInstrument(ctx, createdValidInstrument.ID))
			}
		}
	})
}

func (s *TestSuite) TestValidInstruments_Searching() {
	s.runForCookieClient("should be able to be search for valid instruments", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			stopChan := make(chan bool, 1)
			notificationsChan, err := testClients.admin.SubscribeToNotifications(ctx, stopChan)
			require.NotNil(t, notificationsChan)
			require.NoError(t, err)

			var n *types.DataChangeMessage

			t.Log("creating valid instruments")
			var expected []*types.ValidInstrument
			exampleValidInstrument := fakes.BuildFakeValidInstrument()
			exampleValidInstrument.Name = "example"
			searchQuery := exampleValidInstrument.Name
			for i := 0; i < 5; i++ {
				exampleValidInstrument.Name = fmt.Sprintf("%s %d", searchQuery, i)
				exampleValidInstrumentInput := fakes.BuildFakeValidInstrumentCreationRequestInputFromValidInstrument(exampleValidInstrument)
				createdValidInstrumentID, createdValidInstrumentErr := testClients.admin.CreateValidInstrument(ctx, exampleValidInstrumentInput)
				require.NoError(t, createdValidInstrumentErr)
				t.Logf("valid instrument %q created", createdValidInstrumentID)

				n = <-notificationsChan
				assert.Equal(t, types.ValidInstrumentDataType, n.DataType)
				require.NotNil(t, n.ValidInstrument)
				checkValidInstrumentEquality(t, exampleValidInstrument, n.ValidInstrument)

				createdValidInstrument, createdValidInstrumentErr := testClients.admin.GetValidInstrument(ctx, createdValidInstrumentID)
				requireNotNilAndNoProblems(t, createdValidInstrument, createdValidInstrumentErr)

				expected = append(expected, createdValidInstrument)
			}

			exampleLimit := uint8(20)

			// give the index a moment
			time.Sleep(3 * time.Second)

			// assert valid instrument list equality
			actual, err := testClients.admin.SearchValidInstruments(ctx, searchQuery, exampleLimit)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual),
				"expected %d to be <= %d",
				len(expected),
				len(actual),
			)

			t.Log("cleaning up")
			for _, createdValidInstrument := range expected {
				assert.NoError(t, testClients.admin.ArchiveValidInstrument(ctx, createdValidInstrument.ID))
			}
		}
	})

	s.runForPASETOClient("should be able to be search for valid instruments", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			var checkFunc func() bool
			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating valid instruments")
			var expected []*types.ValidInstrument
			exampleValidInstrument := fakes.BuildFakeValidInstrument()
			exampleValidInstrument.Name = "example"
			searchQuery := exampleValidInstrument.Name
			for i := 0; i < 5; i++ {
				exampleValidInstrument.Name = fmt.Sprintf("%s %d", searchQuery, i)
				exampleValidInstrumentInput := fakes.BuildFakeValidInstrumentCreationRequestInputFromValidInstrument(exampleValidInstrument)
				createdValidInstrumentID, err := testClients.admin.CreateValidInstrument(ctx, exampleValidInstrumentInput)
				require.NoError(t, err)
				t.Logf("valid instrument %q created", createdValidInstrumentID)

				var createdValidInstrument *types.ValidInstrument
				checkFunc = func() bool {
					createdValidInstrument, err = testClients.admin.GetValidInstrument(ctx, createdValidInstrumentID)
					return assert.NotNil(t, createdValidInstrument) && assert.NoError(t, err)
				}
				assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
				requireNotNilAndNoProblems(t, createdValidInstrument, err)

				expected = append(expected, createdValidInstrument)
			}

			exampleLimit := uint8(20)
			time.Sleep(2 * time.Second) // give the index a moment

			// assert valid instrument list equality
			actual, err := testClients.admin.SearchValidInstruments(ctx, searchQuery, exampleLimit)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual),
				"expected %d to be <= %d",
				len(expected),
				len(actual),
			)

			t.Log("cleaning up")
			for _, createdValidInstrument := range expected {
				assert.NoError(t, testClients.admin.ArchiveValidInstrument(ctx, createdValidInstrument.ID))
			}
		}
	})
}
