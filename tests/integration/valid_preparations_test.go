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

func checkValidPreparationEquality(t *testing.T, expected, actual *types.ValidPreparation) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for valid preparation #%d to be %v, but it was %v ", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.Description, actual.Description, "expected Description for valid preparation #%d to be %v, but it was %v ", expected.ID, expected.Description, actual.Description)
	assert.Equal(t, expected.IconPath, actual.IconPath, "expected IconPath for valid preparation #%d to be %v, but it was %v ", expected.ID, expected.IconPath, actual.IconPath)
	assert.NotZero(t, actual.CreatedOn)
}

func (s *TestSuite) TestValidPreparations_Creating() {
	s.runForEachClientExcept("should be creatable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create valid preparation.
			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakes.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparation, err := testClients.main.CreateValidPreparation(ctx, exampleValidPreparationInput)
			requireNotNilAndNoProblems(t, createdValidPreparation, err)

			// assert valid preparation equality
			checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

			auditLogEntries, err := testClients.admin.GetAuditLogForValidPreparation(ctx, createdValidPreparation.ID)
			require.NoError(t, err)

			expectedAuditLogEntries := []*types.AuditLogEntry{
				{EventType: audit.ValidPreparationCreationEvent},
			}
			validateAuditLogEntries(t, expectedAuditLogEntries, auditLogEntries, createdValidPreparation.ID, audit.ValidPreparationAssignmentKey)

			// Clean up valid preparation.
			assert.NoError(t, testClients.main.ArchiveValidPreparation(ctx, createdValidPreparation.ID))
		}
	})
}

func (s *TestSuite) TestValidPreparations_Listing() {
	s.runForEachClientExcept("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create valid preparations
			var expected []*types.ValidPreparation
			for i := 0; i < 5; i++ {
				exampleValidPreparation := fakes.BuildFakeValidPreparation()
				exampleValidPreparationInput := fakes.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)

				createdValidPreparation, validPreparationCreationErr := testClients.main.CreateValidPreparation(ctx, exampleValidPreparationInput)
				requireNotNilAndNoProblems(t, createdValidPreparation, validPreparationCreationErr)

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

			// clean up
			for _, createdValidPreparation := range actual.ValidPreparations {
				assert.NoError(t, testClients.main.ArchiveValidPreparation(ctx, createdValidPreparation.ID))
			}
		}
	})
}

func (s *TestSuite) TestValidPreparations_Searching() {
	s.runForEachClientExcept("should be able to be search for valid preparations", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create valid preparations
			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			var expected []*types.ValidPreparation
			for i := 0; i < 5; i++ {
				exampleValidPreparationInput := fakes.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)
				exampleValidPreparationInput.Name = fmt.Sprintf("%s %d", exampleValidPreparationInput.Name, i)

				createdValidPreparation, validPreparationCreationErr := testClients.main.CreateValidPreparation(ctx, exampleValidPreparationInput)
				requireNotNilAndNoProblems(t, createdValidPreparation, validPreparationCreationErr)

				expected = append(expected, createdValidPreparation)
			}

			exampleLimit := uint8(20)

			// assert valid preparation list equality
			actual, err := testClients.main.SearchValidPreparations(ctx, exampleValidPreparation.Name, exampleLimit)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual),
				"expected results length %d to be <= %d",
				len(expected),
				len(actual),
			)

			// clean up
			for _, createdValidPreparation := range expected {
				assert.NoError(t, testClients.main.ArchiveValidPreparation(ctx, createdValidPreparation.ID))
			}
		}
	})
}

func (s *TestSuite) TestValidPreparations_ExistenceChecking_ReturnsFalseForNonexistentValidPreparation() {
	s.runForEachClientExcept("should not return an error for nonexistent valid preparation", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			actual, err := testClients.main.ValidPreparationExists(ctx, nonexistentID)
			assert.NoError(t, err)
			assert.False(t, actual)
		}
	})
}

func (s *TestSuite) TestValidPreparations_ExistenceChecking_ReturnsTrueForValidValidPreparation() {
	s.runForEachClientExcept("should not return an error for existent valid preparation", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create valid preparation
			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakes.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparation, err := testClients.main.CreateValidPreparation(ctx, exampleValidPreparationInput)
			requireNotNilAndNoProblems(t, createdValidPreparation, err)

			// retrieve valid preparation
			actual, err := testClients.main.ValidPreparationExists(ctx, createdValidPreparation.ID)
			assert.NoError(t, err)
			assert.True(t, actual)

			// clean up valid preparation
			assert.NoError(t, testClients.main.ArchiveValidPreparation(ctx, createdValidPreparation.ID))
		}
	})
}

func (s *TestSuite) TestValidPreparations_Reading_Returns404ForNonexistentValidPreparation() {
	s.runForEachClientExcept("it should return an error when trying to read a valid preparation that does not exist", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			_, err := testClients.main.GetValidPreparation(ctx, nonexistentID)
			assert.Error(t, err)
		}
	})
}

func (s *TestSuite) TestValidPreparations_Reading() {
	s.runForEachClientExcept("it should be readable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create valid preparation
			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakes.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparation, err := testClients.main.CreateValidPreparation(ctx, exampleValidPreparationInput)
			requireNotNilAndNoProblems(t, createdValidPreparation, err)

			// retrieve valid preparation
			actual, err := testClients.main.GetValidPreparation(ctx, createdValidPreparation.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert valid preparation equality
			checkValidPreparationEquality(t, exampleValidPreparation, actual)

			// clean up valid preparation
			assert.NoError(t, testClients.main.ArchiveValidPreparation(ctx, createdValidPreparation.ID))
		}
	})
}

func (s *TestSuite) TestValidPreparations_Updating_Returns404ForNonexistentValidPreparation() {
	s.runForEachClientExcept("it should return an error when trying to update something that does not exist", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparation.ID = nonexistentID

			assert.Error(t, testClients.main.UpdateValidPreparation(ctx, exampleValidPreparation))
		}
	})
}

// convertValidPreparationToValidPreparationUpdateInput creates an ValidPreparationUpdateInput struct from a valid preparation.
func convertValidPreparationToValidPreparationUpdateInput(x *types.ValidPreparation) *types.ValidPreparationUpdateInput {
	return &types.ValidPreparationUpdateInput{
		Name:        x.Name,
		Description: x.Description,
		IconPath:    x.IconPath,
	}
}

func (s *TestSuite) TestValidPreparations_Updating() {
	s.runForEachClientExcept("it should be possible to update a valid preparation", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create valid preparation
			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakes.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparation, err := testClients.main.CreateValidPreparation(ctx, exampleValidPreparationInput)
			requireNotNilAndNoProblems(t, createdValidPreparation, err)

			// change valid preparation
			createdValidPreparation.Update(convertValidPreparationToValidPreparationUpdateInput(exampleValidPreparation))
			assert.NoError(t, testClients.main.UpdateValidPreparation(ctx, createdValidPreparation))

			// retrieve changed valid preparation
			actual, err := testClients.main.GetValidPreparation(ctx, createdValidPreparation.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert valid preparation equality
			checkValidPreparationEquality(t, exampleValidPreparation, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			auditLogEntries, err := testClients.admin.GetAuditLogForValidPreparation(ctx, createdValidPreparation.ID)
			require.NoError(t, err)

			expectedAuditLogEntries := []*types.AuditLogEntry{
				{EventType: audit.ValidPreparationCreationEvent},
				{EventType: audit.ValidPreparationUpdateEvent},
			}
			validateAuditLogEntries(t, expectedAuditLogEntries, auditLogEntries, createdValidPreparation.ID, audit.ValidPreparationAssignmentKey)

			// clean up valid preparation
			assert.NoError(t, testClients.main.ArchiveValidPreparation(ctx, createdValidPreparation.ID))
		}
	})
}

func (s *TestSuite) TestValidPreparations_Archiving_Returns404ForNonexistentValidPreparation() {
	s.runForEachClientExcept("it should return an error when trying to delete something that does not exist", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			assert.Error(t, testClients.main.ArchiveValidPreparation(ctx, nonexistentID))
		}
	})
}

func (s *TestSuite) TestValidPreparations_Archiving() {
	s.runForEachClientExcept("it should be possible to delete a valid preparation", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create valid preparation
			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakes.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparation, err := testClients.main.CreateValidPreparation(ctx, exampleValidPreparationInput)
			requireNotNilAndNoProblems(t, createdValidPreparation, err)

			// clean up valid preparation
			assert.NoError(t, testClients.main.ArchiveValidPreparation(ctx, createdValidPreparation.ID))

			auditLogEntries, err := testClients.admin.GetAuditLogForValidPreparation(ctx, createdValidPreparation.ID)
			require.NoError(t, err)

			expectedAuditLogEntries := []*types.AuditLogEntry{
				{EventType: audit.ValidPreparationCreationEvent},
				{EventType: audit.ValidPreparationArchiveEvent},
			}
			validateAuditLogEntries(t, expectedAuditLogEntries, auditLogEntries, createdValidPreparation.ID, audit.ValidPreparationAssignmentKey)
		}
	})
}

func (s *TestSuite) TestValidPreparations_Auditing_Returns404ForNonexistentValidPreparation() {
	s.runForEachClientExcept("it should return an error when trying to audit something that does not exist", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			x, err := testClients.admin.GetAuditLogForValidPreparation(ctx, nonexistentID)

			assert.NoError(t, err)
			assert.Empty(t, x)
		}
	})
}
