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

func checkInvitationEquality(t *testing.T, expected, actual *types.Invitation) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Code, actual.Code, "expected Code for invitation #%d to be %v, but it was %v ", expected.ID, expected.Code, actual.Code)
	assert.Equal(t, expected.Consumed, actual.Consumed, "expected Consumed for invitation #%d to be %v, but it was %v ", expected.ID, expected.Consumed, actual.Consumed)
	assert.NotZero(t, actual.CreatedOn)
}

func (s *TestSuite) TestInvitations_Creating() {
	s.runForEachClientExcept("should be creatable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create invitation.
			exampleInvitation := fakes.BuildFakeInvitation()
			exampleInvitationInput := fakes.BuildFakeInvitationCreationInputFromInvitation(exampleInvitation)
			createdInvitation, err := testClients.main.CreateInvitation(ctx, exampleInvitationInput)
			requireNotNilAndNoProblems(t, createdInvitation, err)

			// assert invitation equality
			checkInvitationEquality(t, exampleInvitation, createdInvitation)

			auditLogEntries, err := testClients.admin.GetAuditLogForInvitation(ctx, createdInvitation.ID)
			require.NoError(t, err)

			expectedAuditLogEntries := []*types.AuditLogEntry{
				{EventType: audit.InvitationCreationEvent},
			}
			validateAuditLogEntries(t, expectedAuditLogEntries, auditLogEntries, createdInvitation.ID, audit.InvitationAssignmentKey)

			// Clean up invitation.
			assert.NoError(t, testClients.main.ArchiveInvitation(ctx, createdInvitation.ID))
		}
	})
}

func (s *TestSuite) TestInvitations_Listing() {
	s.runForEachClientExcept("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create invitations
			var expected []*types.Invitation
			for i := 0; i < 5; i++ {
				exampleInvitation := fakes.BuildFakeInvitation()
				exampleInvitationInput := fakes.BuildFakeInvitationCreationInputFromInvitation(exampleInvitation)

				createdInvitation, invitationCreationErr := testClients.main.CreateInvitation(ctx, exampleInvitationInput)
				requireNotNilAndNoProblems(t, createdInvitation, invitationCreationErr)

				expected = append(expected, createdInvitation)
			}

			// assert invitation list equality
			actual, err := testClients.main.GetInvitations(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Invitations),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Invitations),
			)

			// clean up
			for _, createdInvitation := range actual.Invitations {
				assert.NoError(t, testClients.main.ArchiveInvitation(ctx, createdInvitation.ID))
			}
		}
	})
}

func (s *TestSuite) TestInvitations_ExistenceChecking_ReturnsFalseForNonexistentInvitation() {
	s.runForEachClientExcept("should not return an error for nonexistent invitation", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			actual, err := testClients.main.InvitationExists(ctx, nonexistentID)
			assert.NoError(t, err)
			assert.False(t, actual)
		}
	})
}

func (s *TestSuite) TestInvitations_ExistenceChecking_ReturnsTrueForValidInvitation() {
	s.runForEachClientExcept("should not return an error for existent invitation", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create invitation
			exampleInvitation := fakes.BuildFakeInvitation()
			exampleInvitationInput := fakes.BuildFakeInvitationCreationInputFromInvitation(exampleInvitation)
			createdInvitation, err := testClients.main.CreateInvitation(ctx, exampleInvitationInput)
			requireNotNilAndNoProblems(t, createdInvitation, err)

			// retrieve invitation
			actual, err := testClients.main.InvitationExists(ctx, createdInvitation.ID)
			assert.NoError(t, err)
			assert.True(t, actual)

			// clean up invitation
			assert.NoError(t, testClients.main.ArchiveInvitation(ctx, createdInvitation.ID))
		}
	})
}

func (s *TestSuite) TestInvitations_Reading_Returns404ForNonexistentInvitation() {
	s.runForEachClientExcept("it should return an error when trying to read an invitation that does not exist", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			_, err := testClients.main.GetInvitation(ctx, nonexistentID)
			assert.Error(t, err)
		}
	})
}

func (s *TestSuite) TestInvitations_Reading() {
	s.runForEachClientExcept("it should be readable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create invitation
			exampleInvitation := fakes.BuildFakeInvitation()
			exampleInvitationInput := fakes.BuildFakeInvitationCreationInputFromInvitation(exampleInvitation)
			createdInvitation, err := testClients.main.CreateInvitation(ctx, exampleInvitationInput)
			requireNotNilAndNoProblems(t, createdInvitation, err)

			// retrieve invitation
			actual, err := testClients.main.GetInvitation(ctx, createdInvitation.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert invitation equality
			checkInvitationEquality(t, exampleInvitation, actual)

			// clean up invitation
			assert.NoError(t, testClients.main.ArchiveInvitation(ctx, createdInvitation.ID))
		}
	})
}

func (s *TestSuite) TestInvitations_Updating_Returns404ForNonexistentInvitation() {
	s.runForEachClientExcept("it should return an error when trying to update something that does not exist", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			exampleInvitation := fakes.BuildFakeInvitation()
			exampleInvitation.ID = nonexistentID

			assert.Error(t, testClients.main.UpdateInvitation(ctx, exampleInvitation))
		}
	})
}

// convertInvitationToInvitationUpdateInput creates an InvitationUpdateInput struct from an invitation.
func convertInvitationToInvitationUpdateInput(x *types.Invitation) *types.InvitationUpdateInput {
	return &types.InvitationUpdateInput{
		Code:     x.Code,
		Consumed: x.Consumed,
	}
}

func (s *TestSuite) TestInvitations_Updating() {
	s.runForEachClientExcept("it should be possible to update an invitation", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create invitation
			exampleInvitation := fakes.BuildFakeInvitation()
			exampleInvitationInput := fakes.BuildFakeInvitationCreationInputFromInvitation(exampleInvitation)
			createdInvitation, err := testClients.main.CreateInvitation(ctx, exampleInvitationInput)
			requireNotNilAndNoProblems(t, createdInvitation, err)

			// change invitation
			createdInvitation.Update(convertInvitationToInvitationUpdateInput(exampleInvitation))
			assert.NoError(t, testClients.main.UpdateInvitation(ctx, createdInvitation))

			// retrieve changed invitation
			actual, err := testClients.main.GetInvitation(ctx, createdInvitation.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert invitation equality
			checkInvitationEquality(t, exampleInvitation, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			auditLogEntries, err := testClients.admin.GetAuditLogForInvitation(ctx, createdInvitation.ID)
			require.NoError(t, err)

			expectedAuditLogEntries := []*types.AuditLogEntry{
				{EventType: audit.InvitationCreationEvent},
				{EventType: audit.InvitationUpdateEvent},
			}
			validateAuditLogEntries(t, expectedAuditLogEntries, auditLogEntries, createdInvitation.ID, audit.InvitationAssignmentKey)

			// clean up invitation
			assert.NoError(t, testClients.main.ArchiveInvitation(ctx, createdInvitation.ID))
		}
	})
}

func (s *TestSuite) TestInvitations_Archiving_Returns404ForNonexistentInvitation() {
	s.runForEachClientExcept("it should return an error when trying to delete something that does not exist", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			assert.Error(t, testClients.main.ArchiveInvitation(ctx, nonexistentID))
		}
	})
}

func (s *TestSuite) TestInvitations_Archiving() {
	s.runForEachClientExcept("it should be possible to delete an invitation", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create invitation
			exampleInvitation := fakes.BuildFakeInvitation()
			exampleInvitationInput := fakes.BuildFakeInvitationCreationInputFromInvitation(exampleInvitation)
			createdInvitation, err := testClients.main.CreateInvitation(ctx, exampleInvitationInput)
			requireNotNilAndNoProblems(t, createdInvitation, err)

			// clean up invitation
			assert.NoError(t, testClients.main.ArchiveInvitation(ctx, createdInvitation.ID))

			auditLogEntries, err := testClients.admin.GetAuditLogForInvitation(ctx, createdInvitation.ID)
			require.NoError(t, err)

			expectedAuditLogEntries := []*types.AuditLogEntry{
				{EventType: audit.InvitationCreationEvent},
				{EventType: audit.InvitationArchiveEvent},
			}
			validateAuditLogEntries(t, expectedAuditLogEntries, auditLogEntries, createdInvitation.ID, audit.InvitationAssignmentKey)
		}
	})
}

func (s *TestSuite) TestInvitations_Auditing_Returns404ForNonexistentInvitation() {
	s.runForEachClientExcept("it should return an error when trying to audit something that does not exist", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			x, err := testClients.admin.GetAuditLogForInvitation(ctx, nonexistentID)

			assert.NoError(t, err)
			assert.Empty(t, x)
		}
	})
}
