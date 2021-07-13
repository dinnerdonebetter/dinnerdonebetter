package integration

import (
	"testing"

	audit "gitlab.com/prixfixe/prixfixe/internal/audit"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkValidPreparationInstrumentEquality(t *testing.T, expected, actual *types.ValidPreparationInstrument) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.InstrumentID, actual.InstrumentID, "expected InstrumentID for valid preparation instrument #%d to be %v, but it was %v ", expected.ID, expected.InstrumentID, actual.InstrumentID)
	assert.Equal(t, expected.PreparationID, actual.PreparationID, "expected PreparationID for valid preparation instrument #%d to be %v, but it was %v ", expected.ID, expected.PreparationID, actual.PreparationID)
	assert.Equal(t, expected.Notes, actual.Notes, "expected Notes for valid preparation instrument #%d to be %v, but it was %v ", expected.ID, expected.Notes, actual.Notes)
	assert.NotZero(t, actual.CreatedOn)
}

func (s *TestSuite) TestValidPreparationInstruments_Creating() {
	s.runForEachClientExcept("should be creatable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create valid preparation instrument.
			exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
			exampleValidPreparationInstrumentInput := fakes.BuildFakeValidPreparationInstrumentCreationInputFromValidPreparationInstrument(exampleValidPreparationInstrument)
			createdValidPreparationInstrument, err := testClients.main.CreateValidPreparationInstrument(ctx, exampleValidPreparationInstrumentInput)
			requireNotNilAndNoProblems(t, createdValidPreparationInstrument, err)

			// assert valid preparation instrument equality
			checkValidPreparationInstrumentEquality(t, exampleValidPreparationInstrument, createdValidPreparationInstrument)

			auditLogEntries, err := testClients.admin.GetAuditLogForValidPreparationInstrument(ctx, createdValidPreparationInstrument.ID)
			require.NoError(t, err)

			expectedAuditLogEntries := []*types.AuditLogEntry{
				{EventType: audit.ValidPreparationInstrumentCreationEvent},
			}
			validateAuditLogEntries(t, expectedAuditLogEntries, auditLogEntries, createdValidPreparationInstrument.ID, audit.ValidPreparationInstrumentAssignmentKey)

			// Clean up valid preparation instrument.
			assert.NoError(t, testClients.main.ArchiveValidPreparationInstrument(ctx, createdValidPreparationInstrument.ID))
		}
	})
}

func (s *TestSuite) TestValidPreparationInstruments_Listing() {
	s.runForEachClientExcept("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create valid preparation instruments
			var expected []*types.ValidPreparationInstrument
			for i := 0; i < 5; i++ {
				exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
				exampleValidPreparationInstrumentInput := fakes.BuildFakeValidPreparationInstrumentCreationInputFromValidPreparationInstrument(exampleValidPreparationInstrument)

				createdValidPreparationInstrument, validPreparationInstrumentCreationErr := testClients.main.CreateValidPreparationInstrument(ctx, exampleValidPreparationInstrumentInput)
				requireNotNilAndNoProblems(t, createdValidPreparationInstrument, validPreparationInstrumentCreationErr)

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

			// clean up
			for _, createdValidPreparationInstrument := range actual.ValidPreparationInstruments {
				assert.NoError(t, testClients.main.ArchiveValidPreparationInstrument(ctx, createdValidPreparationInstrument.ID))
			}
		}
	})
}

func (s *TestSuite) TestValidPreparationInstruments_ExistenceChecking_ReturnsFalseForNonexistentValidPreparationInstrument() {
	s.runForEachClientExcept("should not return an error for nonexistent valid preparation instrument", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			actual, err := testClients.main.ValidPreparationInstrumentExists(ctx, nonexistentID)
			assert.NoError(t, err)
			assert.False(t, actual)
		}
	})
}

func (s *TestSuite) TestValidPreparationInstruments_ExistenceChecking_ReturnsTrueForValidValidPreparationInstrument() {
	s.runForEachClientExcept("should not return an error for existent valid preparation instrument", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create valid preparation instrument
			exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
			exampleValidPreparationInstrumentInput := fakes.BuildFakeValidPreparationInstrumentCreationInputFromValidPreparationInstrument(exampleValidPreparationInstrument)
			createdValidPreparationInstrument, err := testClients.main.CreateValidPreparationInstrument(ctx, exampleValidPreparationInstrumentInput)
			requireNotNilAndNoProblems(t, createdValidPreparationInstrument, err)

			// retrieve valid preparation instrument
			actual, err := testClients.main.ValidPreparationInstrumentExists(ctx, createdValidPreparationInstrument.ID)
			assert.NoError(t, err)
			assert.True(t, actual)

			// clean up valid preparation instrument
			assert.NoError(t, testClients.main.ArchiveValidPreparationInstrument(ctx, createdValidPreparationInstrument.ID))
		}
	})
}

func (s *TestSuite) TestValidPreparationInstruments_Reading_Returns404ForNonexistentValidPreparationInstrument() {
	s.runForEachClientExcept("it should return an error when trying to read a valid preparation instrument that does not exist", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			_, err := testClients.main.GetValidPreparationInstrument(ctx, nonexistentID)
			assert.Error(t, err)
		}
	})
}

func (s *TestSuite) TestValidPreparationInstruments_Reading() {
	s.runForEachClientExcept("it should be readable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create valid preparation instrument
			exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
			exampleValidPreparationInstrumentInput := fakes.BuildFakeValidPreparationInstrumentCreationInputFromValidPreparationInstrument(exampleValidPreparationInstrument)
			createdValidPreparationInstrument, err := testClients.main.CreateValidPreparationInstrument(ctx, exampleValidPreparationInstrumentInput)
			requireNotNilAndNoProblems(t, createdValidPreparationInstrument, err)

			// retrieve valid preparation instrument
			actual, err := testClients.main.GetValidPreparationInstrument(ctx, createdValidPreparationInstrument.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert valid preparation instrument equality
			checkValidPreparationInstrumentEquality(t, exampleValidPreparationInstrument, actual)

			// clean up valid preparation instrument
			assert.NoError(t, testClients.main.ArchiveValidPreparationInstrument(ctx, createdValidPreparationInstrument.ID))
		}
	})
}

func (s *TestSuite) TestValidPreparationInstruments_Updating_Returns404ForNonexistentValidPreparationInstrument() {
	s.runForEachClientExcept("it should return an error when trying to update something that does not exist", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
			exampleValidPreparationInstrument.ID = nonexistentID

			assert.Error(t, testClients.main.UpdateValidPreparationInstrument(ctx, exampleValidPreparationInstrument))
		}
	})
}

// convertValidPreparationInstrumentToValidPreparationInstrumentUpdateInput creates an ValidPreparationInstrumentUpdateInput struct from a valid preparation instrument.
func convertValidPreparationInstrumentToValidPreparationInstrumentUpdateInput(x *types.ValidPreparationInstrument) *types.ValidPreparationInstrumentUpdateInput {
	return &types.ValidPreparationInstrumentUpdateInput{
		InstrumentID:  x.InstrumentID,
		PreparationID: x.PreparationID,
		Notes:         x.Notes,
	}
}

func (s *TestSuite) TestValidPreparationInstruments_Updating() {
	s.runForEachClientExcept("it should be possible to update a valid preparation instrument", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create valid preparation instrument
			exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
			exampleValidPreparationInstrumentInput := fakes.BuildFakeValidPreparationInstrumentCreationInputFromValidPreparationInstrument(exampleValidPreparationInstrument)
			createdValidPreparationInstrument, err := testClients.main.CreateValidPreparationInstrument(ctx, exampleValidPreparationInstrumentInput)
			requireNotNilAndNoProblems(t, createdValidPreparationInstrument, err)

			// change valid preparation instrument
			createdValidPreparationInstrument.Update(convertValidPreparationInstrumentToValidPreparationInstrumentUpdateInput(exampleValidPreparationInstrument))
			assert.NoError(t, testClients.main.UpdateValidPreparationInstrument(ctx, createdValidPreparationInstrument))

			// retrieve changed valid preparation instrument
			actual, err := testClients.main.GetValidPreparationInstrument(ctx, createdValidPreparationInstrument.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert valid preparation instrument equality
			checkValidPreparationInstrumentEquality(t, exampleValidPreparationInstrument, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			auditLogEntries, err := testClients.admin.GetAuditLogForValidPreparationInstrument(ctx, createdValidPreparationInstrument.ID)
			require.NoError(t, err)

			expectedAuditLogEntries := []*types.AuditLogEntry{
				{EventType: audit.ValidPreparationInstrumentCreationEvent},
				{EventType: audit.ValidPreparationInstrumentUpdateEvent},
			}
			validateAuditLogEntries(t, expectedAuditLogEntries, auditLogEntries, createdValidPreparationInstrument.ID, audit.ValidPreparationInstrumentAssignmentKey)

			// clean up valid preparation instrument
			assert.NoError(t, testClients.main.ArchiveValidPreparationInstrument(ctx, createdValidPreparationInstrument.ID))
		}
	})
}

func (s *TestSuite) TestValidPreparationInstruments_Archiving_Returns404ForNonexistentValidPreparationInstrument() {
	s.runForEachClientExcept("it should return an error when trying to delete something that does not exist", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			assert.Error(t, testClients.main.ArchiveValidPreparationInstrument(ctx, nonexistentID))
		}
	})
}

func (s *TestSuite) TestValidPreparationInstruments_Archiving() {
	s.runForEachClientExcept("it should be possible to delete a valid preparation instrument", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create valid preparation instrument
			exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
			exampleValidPreparationInstrumentInput := fakes.BuildFakeValidPreparationInstrumentCreationInputFromValidPreparationInstrument(exampleValidPreparationInstrument)
			createdValidPreparationInstrument, err := testClients.main.CreateValidPreparationInstrument(ctx, exampleValidPreparationInstrumentInput)
			requireNotNilAndNoProblems(t, createdValidPreparationInstrument, err)

			// clean up valid preparation instrument
			assert.NoError(t, testClients.main.ArchiveValidPreparationInstrument(ctx, createdValidPreparationInstrument.ID))

			auditLogEntries, err := testClients.admin.GetAuditLogForValidPreparationInstrument(ctx, createdValidPreparationInstrument.ID)
			require.NoError(t, err)

			expectedAuditLogEntries := []*types.AuditLogEntry{
				{EventType: audit.ValidPreparationInstrumentCreationEvent},
				{EventType: audit.ValidPreparationInstrumentArchiveEvent},
			}
			validateAuditLogEntries(t, expectedAuditLogEntries, auditLogEntries, createdValidPreparationInstrument.ID, audit.ValidPreparationInstrumentAssignmentKey)
		}
	})
}

func (s *TestSuite) TestValidPreparationInstruments_Auditing_Returns404ForNonexistentValidPreparationInstrument() {
	s.runForEachClientExcept("it should return an error when trying to audit something that does not exist", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			x, err := testClients.admin.GetAuditLogForValidPreparationInstrument(ctx, nonexistentID)

			assert.NoError(t, err)
			assert.Empty(t, x)
		}
	})
}
