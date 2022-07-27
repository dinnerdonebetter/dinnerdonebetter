package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
)

func checkValidPreparationInstrumentEquality(t *testing.T, expected, actual *types.ValidPreparationInstrument) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Notes, actual.Notes, "expected Notes for valid preparation instrument %s to be %v, but it was %v", expected.ID, expected.Notes, actual.Notes)
	assert.Equal(t, expected.ValidPreparationID, actual.ValidPreparationID, "expected ValidPreparationID for valid preparation instrument %s to be %v, but it was %v", expected.ID, expected.ValidPreparationID, actual.ValidPreparationID)
	assert.Equal(t, expected.ValidInstrumentID, actual.ValidInstrumentID, "expected ValidInstrumentID for valid preparation instrument %s to be %v, but it was %v", expected.ID, expected.ValidInstrumentID, actual.ValidInstrumentID)
	assert.NotZero(t, actual.CreatedOn)
}

// convertValidPreparationInstrumentToValidPreparationInstrumentUpdateInput creates an ValidPreparationInstrumentUpdateRequestInput struct from a valid preparation instrument.
func convertValidPreparationInstrumentToValidPreparationInstrumentUpdateInput(x *types.ValidPreparationInstrument) *types.ValidPreparationInstrumentUpdateRequestInput {
	return &types.ValidPreparationInstrumentUpdateRequestInput{
		Notes:              &x.Notes,
		ValidPreparationID: &x.ValidPreparationID,
		ValidInstrumentID:  &x.ValidInstrumentID,
	}
}

func (s *TestSuite) TestValidPreparationInstruments_CompleteLifecycle() {
	s.runForEachClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating prerequisite valid preparation")
			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakes.BuildFakeValidPreparationCreationRequestInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparation, err := testClients.admin.CreateValidPreparation(ctx, exampleValidPreparationInput)
			require.NoError(t, err)
			t.Logf("valid preparation %q created", createdValidPreparation.ID)

			checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

			createdValidPreparation, err = testClients.main.GetValidPreparation(ctx, createdValidPreparation.ID)
			requireNotNilAndNoProblems(t, createdValidPreparation, err)
			checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

			t.Log("creating prerequisite valid instrument")
			exampleValidInstrument := fakes.BuildFakeValidInstrument()
			exampleValidInstrumentInput := fakes.BuildFakeValidInstrumentCreationRequestInputFromValidInstrument(exampleValidInstrument)
			createdValidInstrument, err := testClients.admin.CreateValidInstrument(ctx, exampleValidInstrumentInput)
			require.NoError(t, err)

			checkValidInstrumentEquality(t, exampleValidInstrument, createdValidInstrument)

			createdValidInstrument, err = testClients.main.GetValidInstrument(ctx, createdValidInstrument.ID)
			requireNotNilAndNoProblems(t, createdValidInstrument, err)
			checkValidInstrumentEquality(t, exampleValidInstrument, createdValidInstrument)
			t.Logf("valid instrument %q created", createdValidInstrument.ID)

			t.Log("creating valid preparation instrument")
			exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
			exampleValidPreparationInstrument.ValidInstrumentID = createdValidInstrument.ID
			exampleValidPreparationInstrument.ValidPreparationID = createdValidPreparation.ID
			exampleValidPreparationInstrumentInput := fakes.BuildFakeValidPreparationInstrumentCreationRequestInputFromValidPreparationInstrument(exampleValidPreparationInstrument)
			createdValidPreparationInstrument, err := testClients.admin.CreateValidPreparationInstrument(ctx, exampleValidPreparationInstrumentInput)
			require.NoError(t, err)
			t.Logf("valid preparation instrument %q created", createdValidPreparationInstrument.ID)

			checkValidPreparationInstrumentEquality(t, exampleValidPreparationInstrument, createdValidPreparationInstrument)

			createdValidPreparationInstrument, err = testClients.main.GetValidPreparationInstrument(ctx, createdValidPreparationInstrument.ID)
			requireNotNilAndNoProblems(t, createdValidPreparationInstrument, err)

			checkValidPreparationInstrumentEquality(t, exampleValidPreparationInstrument, createdValidPreparationInstrument)

			t.Log("changing valid preparation instrument")
			newValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
			newValidPreparationInstrument.ValidInstrumentID = createdValidInstrument.ID
			newValidPreparationInstrument.ValidPreparationID = createdValidPreparation.ID
			createdValidPreparationInstrument.Update(convertValidPreparationInstrumentToValidPreparationInstrumentUpdateInput(newValidPreparationInstrument))
			assert.NoError(t, testClients.admin.UpdateValidPreparationInstrument(ctx, createdValidPreparationInstrument))

			t.Log("fetching changed valid preparation instrument")
			actual, err := testClients.main.GetValidPreparationInstrument(ctx, createdValidPreparationInstrument.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert valid preparation instrument equality
			checkValidPreparationInstrumentEquality(t, newValidPreparationInstrument, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			t.Log("cleaning up valid preparation instrument")
			assert.NoError(t, testClients.admin.ArchiveValidPreparationInstrument(ctx, createdValidPreparationInstrument.ID))
		}
	})
}

func (s *TestSuite) TestValidPreparationInstruments_Listing() {
	s.runForEachClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating valid preparation instruments")
			var expected []*types.ValidPreparationInstrument
			for i := 0; i < 5; i++ {
				t.Log("creating prerequisite valid preparation")
				exampleValidPreparation := fakes.BuildFakeValidPreparation()
				exampleValidPreparationInput := fakes.BuildFakeValidPreparationCreationRequestInputFromValidPreparation(exampleValidPreparation)
				createdValidPreparation, err := testClients.admin.CreateValidPreparation(ctx, exampleValidPreparationInput)
				require.NoError(t, err)
				t.Logf("valid preparation %q created", createdValidPreparation.ID)

				checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

				createdValidPreparation, err = testClients.main.GetValidPreparation(ctx, createdValidPreparation.ID)
				requireNotNilAndNoProblems(t, createdValidPreparation, err)
				checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

				t.Log("creating prerequisite valid instrument")
				exampleValidInstrument := fakes.BuildFakeValidInstrument()
				exampleValidInstrumentInput := fakes.BuildFakeValidInstrumentCreationRequestInputFromValidInstrument(exampleValidInstrument)
				createdValidInstrument, err := testClients.admin.CreateValidInstrument(ctx, exampleValidInstrumentInput)
				require.NoError(t, err)

				checkValidInstrumentEquality(t, exampleValidInstrument, createdValidInstrument)

				createdValidInstrument, err = testClients.main.GetValidInstrument(ctx, createdValidInstrument.ID)
				requireNotNilAndNoProblems(t, createdValidInstrument, err)
				checkValidInstrumentEquality(t, exampleValidInstrument, createdValidInstrument)
				t.Logf("valid instrument %q created", createdValidInstrument.ID)

				exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
				exampleValidPreparationInstrument.ValidInstrumentID = createdValidInstrument.ID
				exampleValidPreparationInstrument.ValidPreparationID = createdValidPreparation.ID
				exampleValidPreparationInstrumentInput := fakes.BuildFakeValidPreparationInstrumentCreationRequestInputFromValidPreparationInstrument(exampleValidPreparationInstrument)
				createdValidPreparationInstrument, createdValidPreparationInstrumentErr := testClients.admin.CreateValidPreparationInstrument(ctx, exampleValidPreparationInstrumentInput)
				require.NoError(t, createdValidPreparationInstrumentErr)

				checkValidPreparationInstrumentEquality(t, exampleValidPreparationInstrument, createdValidPreparationInstrument)

				expected = append(expected, createdValidPreparationInstrument)
			}

			// assert valid preparation instrument list equality
			actual, err := testClients.main.GetValidPreparationInstruments(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.ValidPreparationInstruments),
				"expected %d to be <= %d",
				len(expected),
				len(actual.ValidPreparationInstruments),
			)

			t.Log("cleaning up")
			for _, createdValidPreparationInstrument := range expected {
				assert.NoError(t, testClients.admin.ArchiveValidPreparationInstrument(ctx, createdValidPreparationInstrument.ID))
			}
		}
	})
}
