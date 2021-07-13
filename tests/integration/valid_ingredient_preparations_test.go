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

func checkValidIngredientPreparationEquality(t *testing.T, expected, actual *types.ValidIngredientPreparation) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Notes, actual.Notes, "expected Notes for valid ingredient preparation #%d to be %v, but it was %v ", expected.ID, expected.Notes, actual.Notes)
	assert.Equal(t, expected.ValidIngredientID, actual.ValidIngredientID, "expected ValidIngredientID for valid ingredient preparation #%d to be %v, but it was %v ", expected.ID, expected.ValidIngredientID, actual.ValidIngredientID)
	assert.Equal(t, expected.ValidPreparationID, actual.ValidPreparationID, "expected ValidPreparationID for valid ingredient preparation #%d to be %v, but it was %v ", expected.ID, expected.ValidPreparationID, actual.ValidPreparationID)
	assert.NotZero(t, actual.CreatedOn)
}

func (s *TestSuite) TestValidIngredientPreparations_Creating() {
	s.runForEachClientExcept("should be creatable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create valid ingredient preparation.
			exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
			exampleValidIngredientPreparationInput := fakes.BuildFakeValidIngredientPreparationCreationInputFromValidIngredientPreparation(exampleValidIngredientPreparation)
			createdValidIngredientPreparation, err := testClients.main.CreateValidIngredientPreparation(ctx, exampleValidIngredientPreparationInput)
			requireNotNilAndNoProblems(t, createdValidIngredientPreparation, err)

			// assert valid ingredient preparation equality
			checkValidIngredientPreparationEquality(t, exampleValidIngredientPreparation, createdValidIngredientPreparation)

			auditLogEntries, err := testClients.admin.GetAuditLogForValidIngredientPreparation(ctx, createdValidIngredientPreparation.ID)
			require.NoError(t, err)

			expectedAuditLogEntries := []*types.AuditLogEntry{
				{EventType: audit.ValidIngredientPreparationCreationEvent},
			}
			validateAuditLogEntries(t, expectedAuditLogEntries, auditLogEntries, createdValidIngredientPreparation.ID, audit.ValidIngredientPreparationAssignmentKey)

			// Clean up valid ingredient preparation.
			assert.NoError(t, testClients.main.ArchiveValidIngredientPreparation(ctx, createdValidIngredientPreparation.ID))
		}
	})
}

func (s *TestSuite) TestValidIngredientPreparations_Listing() {
	s.runForEachClientExcept("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create valid ingredient preparations
			var expected []*types.ValidIngredientPreparation
			for i := 0; i < 5; i++ {
				exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
				exampleValidIngredientPreparationInput := fakes.BuildFakeValidIngredientPreparationCreationInputFromValidIngredientPreparation(exampleValidIngredientPreparation)

				createdValidIngredientPreparation, validIngredientPreparationCreationErr := testClients.main.CreateValidIngredientPreparation(ctx, exampleValidIngredientPreparationInput)
				requireNotNilAndNoProblems(t, createdValidIngredientPreparation, validIngredientPreparationCreationErr)

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

			// clean up
			for _, createdValidIngredientPreparation := range actual.ValidIngredientPreparations {
				assert.NoError(t, testClients.main.ArchiveValidIngredientPreparation(ctx, createdValidIngredientPreparation.ID))
			}
		}
	})
}

func (s *TestSuite) TestValidIngredientPreparations_ExistenceChecking_ReturnsFalseForNonexistentValidIngredientPreparation() {
	s.runForEachClientExcept("should not return an error for nonexistent valid ingredient preparation", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			actual, err := testClients.main.ValidIngredientPreparationExists(ctx, nonexistentID)
			assert.NoError(t, err)
			assert.False(t, actual)
		}
	})
}

func (s *TestSuite) TestValidIngredientPreparations_ExistenceChecking_ReturnsTrueForValidValidIngredientPreparation() {
	s.runForEachClientExcept("should not return an error for existent valid ingredient preparation", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create valid ingredient preparation
			exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
			exampleValidIngredientPreparationInput := fakes.BuildFakeValidIngredientPreparationCreationInputFromValidIngredientPreparation(exampleValidIngredientPreparation)
			createdValidIngredientPreparation, err := testClients.main.CreateValidIngredientPreparation(ctx, exampleValidIngredientPreparationInput)
			requireNotNilAndNoProblems(t, createdValidIngredientPreparation, err)

			// retrieve valid ingredient preparation
			actual, err := testClients.main.ValidIngredientPreparationExists(ctx, createdValidIngredientPreparation.ID)
			assert.NoError(t, err)
			assert.True(t, actual)

			// clean up valid ingredient preparation
			assert.NoError(t, testClients.main.ArchiveValidIngredientPreparation(ctx, createdValidIngredientPreparation.ID))
		}
	})
}

func (s *TestSuite) TestValidIngredientPreparations_Reading_Returns404ForNonexistentValidIngredientPreparation() {
	s.runForEachClientExcept("it should return an error when trying to read a valid ingredient preparation that does not exist", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			_, err := testClients.main.GetValidIngredientPreparation(ctx, nonexistentID)
			assert.Error(t, err)
		}
	})
}

func (s *TestSuite) TestValidIngredientPreparations_Reading() {
	s.runForEachClientExcept("it should be readable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create valid ingredient preparation
			exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
			exampleValidIngredientPreparationInput := fakes.BuildFakeValidIngredientPreparationCreationInputFromValidIngredientPreparation(exampleValidIngredientPreparation)
			createdValidIngredientPreparation, err := testClients.main.CreateValidIngredientPreparation(ctx, exampleValidIngredientPreparationInput)
			requireNotNilAndNoProblems(t, createdValidIngredientPreparation, err)

			// retrieve valid ingredient preparation
			actual, err := testClients.main.GetValidIngredientPreparation(ctx, createdValidIngredientPreparation.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert valid ingredient preparation equality
			checkValidIngredientPreparationEquality(t, exampleValidIngredientPreparation, actual)

			// clean up valid ingredient preparation
			assert.NoError(t, testClients.main.ArchiveValidIngredientPreparation(ctx, createdValidIngredientPreparation.ID))
		}
	})
}

func (s *TestSuite) TestValidIngredientPreparations_Updating_Returns404ForNonexistentValidIngredientPreparation() {
	s.runForEachClientExcept("it should return an error when trying to update something that does not exist", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
			exampleValidIngredientPreparation.ID = nonexistentID

			assert.Error(t, testClients.main.UpdateValidIngredientPreparation(ctx, exampleValidIngredientPreparation))
		}
	})
}

// convertValidIngredientPreparationToValidIngredientPreparationUpdateInput creates an ValidIngredientPreparationUpdateInput struct from a valid ingredient preparation.
func convertValidIngredientPreparationToValidIngredientPreparationUpdateInput(x *types.ValidIngredientPreparation) *types.ValidIngredientPreparationUpdateInput {
	return &types.ValidIngredientPreparationUpdateInput{
		Notes:              x.Notes,
		ValidIngredientID:  x.ValidIngredientID,
		ValidPreparationID: x.ValidPreparationID,
	}
}

func (s *TestSuite) TestValidIngredientPreparations_Updating() {
	s.runForEachClientExcept("it should be possible to update a valid ingredient preparation", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create valid ingredient preparation
			exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
			exampleValidIngredientPreparationInput := fakes.BuildFakeValidIngredientPreparationCreationInputFromValidIngredientPreparation(exampleValidIngredientPreparation)
			createdValidIngredientPreparation, err := testClients.main.CreateValidIngredientPreparation(ctx, exampleValidIngredientPreparationInput)
			requireNotNilAndNoProblems(t, createdValidIngredientPreparation, err)

			// change valid ingredient preparation
			createdValidIngredientPreparation.Update(convertValidIngredientPreparationToValidIngredientPreparationUpdateInput(exampleValidIngredientPreparation))
			assert.NoError(t, testClients.main.UpdateValidIngredientPreparation(ctx, createdValidIngredientPreparation))

			// retrieve changed valid ingredient preparation
			actual, err := testClients.main.GetValidIngredientPreparation(ctx, createdValidIngredientPreparation.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert valid ingredient preparation equality
			checkValidIngredientPreparationEquality(t, exampleValidIngredientPreparation, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			auditLogEntries, err := testClients.admin.GetAuditLogForValidIngredientPreparation(ctx, createdValidIngredientPreparation.ID)
			require.NoError(t, err)

			expectedAuditLogEntries := []*types.AuditLogEntry{
				{EventType: audit.ValidIngredientPreparationCreationEvent},
				{EventType: audit.ValidIngredientPreparationUpdateEvent},
			}
			validateAuditLogEntries(t, expectedAuditLogEntries, auditLogEntries, createdValidIngredientPreparation.ID, audit.ValidIngredientPreparationAssignmentKey)

			// clean up valid ingredient preparation
			assert.NoError(t, testClients.main.ArchiveValidIngredientPreparation(ctx, createdValidIngredientPreparation.ID))
		}
	})
}

func (s *TestSuite) TestValidIngredientPreparations_Archiving_Returns404ForNonexistentValidIngredientPreparation() {
	s.runForEachClientExcept("it should return an error when trying to delete something that does not exist", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			assert.Error(t, testClients.main.ArchiveValidIngredientPreparation(ctx, nonexistentID))
		}
	})
}

func (s *TestSuite) TestValidIngredientPreparations_Archiving() {
	s.runForEachClientExcept("it should be possible to delete a valid ingredient preparation", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create valid ingredient preparation
			exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
			exampleValidIngredientPreparationInput := fakes.BuildFakeValidIngredientPreparationCreationInputFromValidIngredientPreparation(exampleValidIngredientPreparation)
			createdValidIngredientPreparation, err := testClients.main.CreateValidIngredientPreparation(ctx, exampleValidIngredientPreparationInput)
			requireNotNilAndNoProblems(t, createdValidIngredientPreparation, err)

			// clean up valid ingredient preparation
			assert.NoError(t, testClients.main.ArchiveValidIngredientPreparation(ctx, createdValidIngredientPreparation.ID))

			auditLogEntries, err := testClients.admin.GetAuditLogForValidIngredientPreparation(ctx, createdValidIngredientPreparation.ID)
			require.NoError(t, err)

			expectedAuditLogEntries := []*types.AuditLogEntry{
				{EventType: audit.ValidIngredientPreparationCreationEvent},
				{EventType: audit.ValidIngredientPreparationArchiveEvent},
			}
			validateAuditLogEntries(t, expectedAuditLogEntries, auditLogEntries, createdValidIngredientPreparation.ID, audit.ValidIngredientPreparationAssignmentKey)
		}
	})
}

func (s *TestSuite) TestValidIngredientPreparations_Auditing_Returns404ForNonexistentValidIngredientPreparation() {
	s.runForEachClientExcept("it should return an error when trying to audit something that does not exist", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			x, err := testClients.admin.GetAuditLogForValidIngredientPreparation(ctx, nonexistentID)

			assert.NoError(t, err)
			assert.Empty(t, x)
		}
	})
}
