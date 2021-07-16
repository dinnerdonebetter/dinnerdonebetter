package integration

import (
	"fmt"
	"testing"

	audit "gitlab.com/prixfixe/prixfixe/internal/audit"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkValidInstrumentEquality(t *testing.T, expected, actual *types.ValidInstrument) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for valid instrument #%d to be %v, but it was %v ", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.Variant, actual.Variant, "expected Variant for valid instrument #%d to be %v, but it was %v ", expected.ID, expected.Variant, actual.Variant)
	assert.Equal(t, expected.Description, actual.Description, "expected Description for valid instrument #%d to be %v, but it was %v ", expected.ID, expected.Description, actual.Description)
	assert.Equal(t, expected.IconPath, actual.IconPath, "expected IconPath for valid instrument #%d to be %v, but it was %v ", expected.ID, expected.IconPath, actual.IconPath)
	assert.NotZero(t, actual.CreatedOn)
}

func (s *TestSuite) TestValidInstruments_Creating() {
	s.runForEachClientExcept("should be creatable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create valid instrument.
			exampleValidInstrument := fakes.BuildFakeValidInstrument()
			exampleValidInstrumentInput := fakes.BuildFakeValidInstrumentCreationInputFromValidInstrument(exampleValidInstrument)
			createdValidInstrument, err := testClients.main.CreateValidInstrument(ctx, exampleValidInstrumentInput)
			requireNotNilAndNoProblems(t, createdValidInstrument, err)

			// assert valid instrument equality
			checkValidInstrumentEquality(t, exampleValidInstrument, createdValidInstrument)

			auditLogEntries, err := testClients.admin.GetAuditLogForValidInstrument(ctx, createdValidInstrument.ID)
			require.NoError(t, err)

			expectedAuditLogEntries := []*types.AuditLogEntry{
				{EventType: audit.ValidInstrumentCreationEvent},
			}
			validateAuditLogEntries(t, expectedAuditLogEntries, auditLogEntries, createdValidInstrument.ID, audit.ValidInstrumentAssignmentKey)

			// Clean up valid instrument.
			assert.NoError(t, testClients.main.ArchiveValidInstrument(ctx, createdValidInstrument.ID))
		}
	})
}

func (s *TestSuite) TestValidInstruments_Listing() {
	s.runForEachClientExcept("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create valid instruments
			var expected []*types.ValidInstrument
			for i := 0; i < 5; i++ {
				exampleValidInstrument := fakes.BuildFakeValidInstrument()
				exampleValidInstrumentInput := fakes.BuildFakeValidInstrumentCreationInputFromValidInstrument(exampleValidInstrument)

				createdValidInstrument, validInstrumentCreationErr := testClients.main.CreateValidInstrument(ctx, exampleValidInstrumentInput)
				requireNotNilAndNoProblems(t, createdValidInstrument, validInstrumentCreationErr)

				expected = append(expected, createdValidInstrument)
			}

			// assert valid instrument list equality
			actual, err := testClients.main.GetValidInstruments(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.ValidInstruments),
				"expected %d to be <= %d",
				len(expected),
				len(actual.ValidInstruments),
			)

			// clean up
			for _, createdValidInstrument := range actual.ValidInstruments {
				assert.NoError(t, testClients.main.ArchiveValidInstrument(ctx, createdValidInstrument.ID))
			}
		}
	})
}

func (s *TestSuite) TestValidInstruments_Searching() {
	s.runForEachClientExcept("should be able to be search for valid instruments", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create valid instruments
			exampleValidInstrument := fakes.BuildFakeValidInstrument()
			var expected []*types.ValidInstrument
			for i := 0; i < 5; i++ {
				exampleValidInstrumentInput := fakes.BuildFakeValidInstrumentCreationInputFromValidInstrument(exampleValidInstrument)
				exampleValidInstrumentInput.Name = fmt.Sprintf("%s %d", exampleValidInstrumentInput.Name, i)

				createdValidInstrument, validInstrumentCreationErr := testClients.main.CreateValidInstrument(ctx, exampleValidInstrumentInput)
				requireNotNilAndNoProblems(t, createdValidInstrument, validInstrumentCreationErr)

				expected = append(expected, createdValidInstrument)
			}

			exampleLimit := uint8(20)

			// assert valid instrument list equality
			actual, err := testClients.main.SearchValidInstruments(ctx, exampleValidInstrument.Name, exampleLimit)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual),
				"expected results length %d to be <= %d",
				len(expected),
				len(actual),
			)

			// clean up
			for _, createdValidInstrument := range expected {
				assert.NoError(t, testClients.main.ArchiveValidInstrument(ctx, createdValidInstrument.ID))
			}
		}
	})
}

func (s *TestSuite) TestValidInstruments_ExistenceChecking_ReturnsFalseForNonexistentValidInstrument() {
	s.runForEachClientExcept("should not return an error for nonexistent valid instrument", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			actual, err := testClients.main.ValidInstrumentExists(ctx, nonexistentID)
			assert.NoError(t, err)
			assert.False(t, actual)
		}
	})
}

func (s *TestSuite) TestValidInstruments_ExistenceChecking_ReturnsTrueForValidValidInstrument() {
	s.runForEachClientExcept("should not return an error for existent valid instrument", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create valid instrument
			exampleValidInstrument := fakes.BuildFakeValidInstrument()
			exampleValidInstrumentInput := fakes.BuildFakeValidInstrumentCreationInputFromValidInstrument(exampleValidInstrument)
			createdValidInstrument, err := testClients.main.CreateValidInstrument(ctx, exampleValidInstrumentInput)
			requireNotNilAndNoProblems(t, createdValidInstrument, err)

			// retrieve valid instrument
			actual, err := testClients.main.ValidInstrumentExists(ctx, createdValidInstrument.ID)
			assert.NoError(t, err)
			assert.True(t, actual)

			// clean up valid instrument
			assert.NoError(t, testClients.main.ArchiveValidInstrument(ctx, createdValidInstrument.ID))
		}
	})
}

func (s *TestSuite) TestValidInstruments_Reading_Returns404ForNonexistentValidInstrument() {
	s.runForEachClientExcept("it should return an error when trying to read a valid instrument that does not exist", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			_, err := testClients.main.GetValidInstrument(ctx, nonexistentID)
			assert.Error(t, err)
		}
	})
}

func (s *TestSuite) TestValidInstruments_Reading() {
	s.runForEachClientExcept("it should be readable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create valid instrument
			exampleValidInstrument := fakes.BuildFakeValidInstrument()
			exampleValidInstrumentInput := fakes.BuildFakeValidInstrumentCreationInputFromValidInstrument(exampleValidInstrument)
			createdValidInstrument, err := testClients.main.CreateValidInstrument(ctx, exampleValidInstrumentInput)
			requireNotNilAndNoProblems(t, createdValidInstrument, err)

			// retrieve valid instrument
			actual, err := testClients.main.GetValidInstrument(ctx, createdValidInstrument.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert valid instrument equality
			checkValidInstrumentEquality(t, exampleValidInstrument, actual)

			// clean up valid instrument
			assert.NoError(t, testClients.main.ArchiveValidInstrument(ctx, createdValidInstrument.ID))
		}
	})
}

func (s *TestSuite) TestValidInstruments_Updating_Returns404ForNonexistentValidInstrument() {
	s.runForEachClientExcept("it should return an error when trying to update something that does not exist", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			exampleValidInstrument := fakes.BuildFakeValidInstrument()
			exampleValidInstrument.ID = nonexistentID

			assert.Error(t, testClients.main.UpdateValidInstrument(ctx, exampleValidInstrument))
		}
	})
}

// convertValidInstrumentToValidInstrumentUpdateInput creates an ValidInstrumentUpdateInput struct from a valid instrument.
func convertValidInstrumentToValidInstrumentUpdateInput(x *types.ValidInstrument) *types.ValidInstrumentUpdateInput {
	return &types.ValidInstrumentUpdateInput{
		Name:        x.Name,
		Variant:     x.Variant,
		Description: x.Description,
		IconPath:    x.IconPath,
	}
}

func (s *TestSuite) TestValidInstruments_Updating() {
	s.runForEachClientExcept("it should be possible to update a valid instrument", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create valid instrument
			exampleValidInstrument := fakes.BuildFakeValidInstrument()
			exampleValidInstrumentInput := fakes.BuildFakeValidInstrumentCreationInputFromValidInstrument(exampleValidInstrument)
			createdValidInstrument, err := testClients.main.CreateValidInstrument(ctx, exampleValidInstrumentInput)
			requireNotNilAndNoProblems(t, createdValidInstrument, err)

			// change valid instrument
			createdValidInstrument.Update(convertValidInstrumentToValidInstrumentUpdateInput(exampleValidInstrument))
			assert.NoError(t, testClients.main.UpdateValidInstrument(ctx, createdValidInstrument))

			// retrieve changed valid instrument
			actual, err := testClients.main.GetValidInstrument(ctx, createdValidInstrument.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert valid instrument equality
			checkValidInstrumentEquality(t, exampleValidInstrument, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			auditLogEntries, err := testClients.admin.GetAuditLogForValidInstrument(ctx, createdValidInstrument.ID)
			require.NoError(t, err)

			expectedAuditLogEntries := []*types.AuditLogEntry{
				{EventType: audit.ValidInstrumentCreationEvent},
				{EventType: audit.ValidInstrumentUpdateEvent},
			}
			validateAuditLogEntries(t, expectedAuditLogEntries, auditLogEntries, createdValidInstrument.ID, audit.ValidInstrumentAssignmentKey)

			// clean up valid instrument
			assert.NoError(t, testClients.main.ArchiveValidInstrument(ctx, createdValidInstrument.ID))
		}
	})
}

func (s *TestSuite) TestValidInstruments_Archiving_Returns404ForNonexistentValidInstrument() {
	s.runForEachClientExcept("it should return an error when trying to delete something that does not exist", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			assert.Error(t, testClients.main.ArchiveValidInstrument(ctx, nonexistentID))
		}
	})
}

func (s *TestSuite) TestValidInstruments_Archiving() {
	s.runForEachClientExcept("it should be possible to delete a valid instrument", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create valid instrument
			exampleValidInstrument := fakes.BuildFakeValidInstrument()
			exampleValidInstrumentInput := fakes.BuildFakeValidInstrumentCreationInputFromValidInstrument(exampleValidInstrument)
			createdValidInstrument, err := testClients.main.CreateValidInstrument(ctx, exampleValidInstrumentInput)
			requireNotNilAndNoProblems(t, createdValidInstrument, err)

			// clean up valid instrument
			assert.NoError(t, testClients.main.ArchiveValidInstrument(ctx, createdValidInstrument.ID))

			auditLogEntries, err := testClients.admin.GetAuditLogForValidInstrument(ctx, createdValidInstrument.ID)
			require.NoError(t, err)

			expectedAuditLogEntries := []*types.AuditLogEntry{
				{EventType: audit.ValidInstrumentCreationEvent},
				{EventType: audit.ValidInstrumentArchiveEvent},
			}
			validateAuditLogEntries(t, expectedAuditLogEntries, auditLogEntries, createdValidInstrument.ID, audit.ValidInstrumentAssignmentKey)
		}
	})
}

func (s *TestSuite) TestValidInstruments_Auditing_Returns404ForNonexistentValidInstrument() {
	s.runForEachClientExcept("it should return an error when trying to audit something that does not exist", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			x, err := testClients.admin.GetAuditLogForValidInstrument(ctx, nonexistentID)

			assert.NoError(t, err)
			assert.Empty(t, x)
		}
	})
}
