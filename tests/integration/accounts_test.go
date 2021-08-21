package integration

import (
	"fmt"
	"testing"
	"time"

	"gitlab.com/prixfixe/prixfixe/internal/audit"
	"gitlab.com/prixfixe/prixfixe/internal/authorization"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/pkg/client/httpclient"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkHouseholdEquality(t *testing.T, expected, actual *types.Household) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected BucketName for household #%d to be %v, but it was %v ", expected.ID, expected.Name, actual.Name)
	assert.NotZero(t, actual.CreatedOn)
}

func (s *TestSuite) TestHouseholds_Creating() {
	s.runForEachClientExcept("should be possible to create households", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create household.
			exampleHousehold := fakes.BuildFakeHousehold()
			exampleHouseholdInput := fakes.BuildFakeHouseholdCreationInputFromHousehold(exampleHousehold)
			createdHousehold, err := testClients.main.CreateHousehold(ctx, exampleHouseholdInput)
			requireNotNilAndNoProblems(t, createdHousehold, err)

			// Assert household equality.
			checkHouseholdEquality(t, exampleHousehold, createdHousehold)

			auditLogEntries, err := testClients.admin.GetAuditLogForHousehold(ctx, createdHousehold.ID)
			require.NoError(t, err)

			expectedAuditLogEntries := []*types.AuditLogEntry{
				{EventType: audit.HouseholdCreationEvent},
				{EventType: audit.UserAddedToHouseholdEvent},
			}
			validateAuditLogEntries(t, expectedAuditLogEntries, auditLogEntries, createdHousehold.ID, audit.HouseholdAssignmentKey)

			// Clean up.
			assert.NoError(t, testClients.main.ArchiveHousehold(ctx, createdHousehold.ID))
		}
	})
}

func (s *TestSuite) TestHouseholds_Listing() {
	s.runForEachClientExcept("should be possible to list households", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create households.
			var expected []*types.Household
			for i := 0; i < 5; i++ {
				// Create household.
				exampleHousehold := fakes.BuildFakeHousehold()
				exampleHouseholdInput := fakes.BuildFakeHouseholdCreationInputFromHousehold(exampleHousehold)
				createdHousehold, householdCreationErr := testClients.main.CreateHousehold(ctx, exampleHouseholdInput)
				requireNotNilAndNoProblems(t, createdHousehold, householdCreationErr)

				expected = append(expected, createdHousehold)
			}

			// Assert household list equality.
			actual, err := testClients.main.GetHouseholds(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Households),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Households),
			)

			// Clean up.
			for _, createdHousehold := range actual.Households {
				assert.NoError(t, testClients.main.ArchiveHousehold(ctx, createdHousehold.ID))
			}
		}
	})
}

func (s *TestSuite) TestHouseholds_Reading_Returns404ForNonexistentHousehold() {
	s.runForEachClientExcept("should not be possible to read a non-existent household", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Attempt to fetch nonexistent household.
			_, err := testClients.main.GetHousehold(ctx, nonexistentID)
			assert.Error(t, err)
		}
	})
}

func (s *TestSuite) TestHouseholds_Reading() {
	s.runForEachClientExcept("should be possible to read an household", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create household.
			exampleHousehold := fakes.BuildFakeHousehold()
			exampleHouseholdInput := fakes.BuildFakeHouseholdCreationInputFromHousehold(exampleHousehold)
			createdHousehold, err := testClients.main.CreateHousehold(ctx, exampleHouseholdInput)
			requireNotNilAndNoProblems(t, createdHousehold, err)

			// Fetch household.
			actual, err := testClients.main.GetHousehold(ctx, createdHousehold.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// Assert household equality.
			checkHouseholdEquality(t, exampleHousehold, actual)

			// Clean up household.
			assert.NoError(t, testClients.main.ArchiveHousehold(ctx, createdHousehold.ID))
		}
	})
}

func (s *TestSuite) TestHouseholds_Updating_Returns404ForNonexistentHousehold() {
	s.runForEachClientExcept("should not be possible to update a non-existent household", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			exampleHousehold := fakes.BuildFakeHousehold()
			exampleHousehold.ID = nonexistentID

			assert.Error(t, testClients.main.UpdateHousehold(ctx, exampleHousehold))
		}
	})
}

// convertHouseholdToHouseholdUpdateInput creates an HouseholdUpdateInput struct from an household.
func convertHouseholdToHouseholdUpdateInput(x *types.Household) *types.HouseholdUpdateInput {
	return &types.HouseholdUpdateInput{
		Name:          x.Name,
		BelongsToUser: x.BelongsToUser,
	}
}

func (s *TestSuite) TestHouseholds_Updating() {
	s.runForEachClientExcept("should be possible to update an household", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create household.
			exampleHousehold := fakes.BuildFakeHousehold()
			exampleHouseholdInput := fakes.BuildFakeHouseholdCreationInputFromHousehold(exampleHousehold)
			createdHousehold, err := testClients.main.CreateHousehold(ctx, exampleHouseholdInput)
			requireNotNilAndNoProblems(t, createdHousehold, err)

			// Change household.
			createdHousehold.Update(convertHouseholdToHouseholdUpdateInput(exampleHousehold))
			assert.NoError(t, testClients.main.UpdateHousehold(ctx, createdHousehold))

			// Fetch household.
			actual, err := testClients.main.GetHousehold(ctx, createdHousehold.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// Assert household equality.
			checkHouseholdEquality(t, exampleHousehold, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			auditLogEntries, err := testClients.admin.GetAuditLogForHousehold(ctx, createdHousehold.ID)
			require.NoError(t, err)

			expectedAuditLogEntries := []*types.AuditLogEntry{
				{EventType: audit.HouseholdCreationEvent},
				{EventType: audit.UserAddedToHouseholdEvent},
				{EventType: audit.HouseholdUpdateEvent},
			}
			validateAuditLogEntries(t, expectedAuditLogEntries, auditLogEntries, createdHousehold.ID, audit.HouseholdAssignmentKey)

			// Clean up household.
			assert.NoError(t, testClients.main.ArchiveHousehold(ctx, createdHousehold.ID))
		}
	})
}

func (s *TestSuite) TestHouseholds_Archiving_Returns404ForNonexistentHousehold() {
	s.runForEachClientExcept("should not be possible to archiv a non-existent household", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			assert.Error(t, testClients.main.ArchiveHousehold(ctx, nonexistentID))
		}
	})
}

func (s *TestSuite) TestHouseholds_Archiving() {
	s.runForEachClientExcept("should be possible to archive an household", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create household.
			exampleHousehold := fakes.BuildFakeHousehold()
			exampleHouseholdInput := fakes.BuildFakeHouseholdCreationInputFromHousehold(exampleHousehold)
			createdHousehold, err := testClients.main.CreateHousehold(ctx, exampleHouseholdInput)
			requireNotNilAndNoProblems(t, createdHousehold, err)

			// Clean up household.
			assert.NoError(t, testClients.main.ArchiveHousehold(ctx, createdHousehold.ID))

			auditLogEntries, err := testClients.admin.GetAuditLogForHousehold(ctx, createdHousehold.ID)
			require.NoError(t, err)

			expectedAuditLogEntries := []*types.AuditLogEntry{
				{EventType: audit.HouseholdCreationEvent},
				{EventType: audit.UserAddedToHouseholdEvent},
				{EventType: audit.HouseholdArchiveEvent},
			}
			validateAuditLogEntries(t, expectedAuditLogEntries, auditLogEntries, createdHousehold.ID, audit.HouseholdAssignmentKey)
		}
	})
}

func (s *TestSuite) TestHouseholds_ChangingMemberships() {
	s.runForEachClientExcept("should be possible to change members of an household", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			const userCount = 1

			currentStatus, statusErr := testClients.main.UserStatus(s.ctx)
			requireNotNilAndNoProblems(t, currentStatus, statusErr)
			t.Logf("initial household is #%d; initial user ID is #%d", currentStatus.ActiveHousehold, s.user.ID)

			// fetch household data
			householdCreationInput := &types.HouseholdCreationInput{
				Name: fakes.BuildFakeHousehold().Name,
			}
			household, householdCreationErr := testClients.main.CreateHousehold(ctx, householdCreationInput)
			require.NoError(t, householdCreationErr)
			require.NotNil(t, household)

			t.Logf("created household #%d", household.ID)

			require.NoError(t, testClients.main.SwitchActiveHousehold(ctx, household.ID))

			t.Logf("switched main test client active household to #%d, creating webhook", household.ID)

			// create a webhook
			exampleWebhook := fakes.BuildFakeWebhook()
			exampleWebhookInput := fakes.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)
			createdWebhook, creationErr := testClients.main.CreateWebhook(ctx, exampleWebhookInput)
			requireNotNilAndNoProblems(t, createdWebhook, creationErr)
			require.Equal(t, household.ID, createdWebhook.BelongsToHousehold)

			t.Logf("created webhook #%d for household #%d", createdWebhook.ID, createdWebhook.BelongsToHousehold)

			expectedAuditLogEntries := []*types.AuditLogEntry{
				{EventType: audit.HouseholdCreationEvent},
				{EventType: audit.UserAddedToHouseholdEvent},
				{EventType: audit.WebhookCreationEvent},
			}

			// create dummy users
			users := []*types.User{}
			clients := []*httpclient.Client{}

			// create users
			for i := 0; i < userCount; i++ {
				u, _, c, _ := createUserAndClientForTest(ctx, t)
				users = append(users, u)
				clients = append(clients, c)

				currentStatus, statusErr = c.UserStatus(s.ctx)
				requireNotNilAndNoProblems(t, currentStatus, statusErr)
				t.Logf("created user user #%d with household #%d", u.ID, currentStatus.ActiveHousehold)
			}

			// check that each user cannot see the unreachable webhook
			for i := 0; i < userCount; i++ {
				t.Logf("checking that user #%d CANNOT see webhook #%d belonging to household #%d", users[i].ID, createdWebhook.ID, createdWebhook.BelongsToHousehold)
				webhook, err := clients[i].GetWebhook(ctx, createdWebhook.ID)
				require.Nil(t, webhook)
				require.Error(t, err)
			}

			// add them to the household
			for i := 0; i < userCount; i++ {
				t.Logf("adding user #%d to household #%d", users[i].ID, household.ID)
				require.NoError(t, testClients.main.AddUserToHousehold(ctx, &types.AddUserToHouseholdInput{
					UserID:         users[i].ID,
					HouseholdID:    household.ID,
					Reason:         t.Name(),
					HouseholdRoles: []string{authorization.HouseholdAdminRole.String()},
				}))
				t.Logf("added user #%d to household #%d", users[i].ID, household.ID)
				expectedAuditLogEntries = append(expectedAuditLogEntries, &types.AuditLogEntry{EventType: audit.UserAddedToHouseholdEvent})

				t.Logf("setting user #%d's client to household #%d", users[i].ID, household.ID)
				require.NoError(t, clients[i].SwitchActiveHousehold(ctx, household.ID))

				currentStatus, statusErr = clients[i].UserStatus(s.ctx)
				requireNotNilAndNoProblems(t, currentStatus, statusErr)
				require.Equal(t, currentStatus.ActiveHousehold, household.ID)
				t.Logf("set user #%d's current active household to #%d", users[i].ID, household.ID)
			}

			// grant all permissions
			for i := 0; i < userCount; i++ {
				input := &types.ModifyUserPermissionsInput{
					Reason:   t.Name(),
					NewRoles: []string{authorization.HouseholdAdminRole.String()},
				}
				require.NoError(t, testClients.main.ModifyMemberPermissions(ctx, household.ID, users[i].ID, input))
				expectedAuditLogEntries = append(expectedAuditLogEntries, &types.AuditLogEntry{EventType: audit.UserHouseholdPermissionsModifiedEvent})
			}

			// check that each user can see the webhook
			for i := 0; i < userCount; i++ {
				t.Logf("checking if user #%d CAN now see webhook #%d belonging to household #%d", users[i].ID, createdWebhook.ID, createdWebhook.BelongsToHousehold)
				webhook, err := clients[i].GetWebhook(ctx, createdWebhook.ID)
				requireNotNilAndNoProblems(t, webhook, err)
			}

			originalWebhookName := createdWebhook.Name
			// check that each user can update the webhook
			for i := 0; i < userCount; i++ {
				createdWebhook.Name = fmt.Sprintf("%s_%d", originalWebhookName, time.Now().UnixNano())
				require.NoError(t, clients[i].UpdateWebhook(ctx, createdWebhook))
				expectedAuditLogEntries = append(expectedAuditLogEntries, &types.AuditLogEntry{EventType: audit.WebhookUpdateEvent})
			}

			// remove users from household
			for i := 0; i < userCount; i++ {
				require.NoError(t, testClients.main.RemoveUserFromHousehold(ctx, household.ID, users[i].ID, t.Name()))
			}

			// check that each user cannot see the webhook
			for i := 0; i < userCount; i++ {
				webhook, err := clients[i].GetWebhook(ctx, createdWebhook.ID)
				require.Nil(t, webhook)
				require.Error(t, err)
			}

			// check audit log entries
			auditLogEntries, err := testClients.admin.GetAuditLogForHousehold(ctx, household.ID)
			require.NoError(t, err)

			validateAuditLogEntries(t, expectedAuditLogEntries, auditLogEntries, household.ID, audit.HouseholdAssignmentKey)

			// Clean up.
			require.NoError(t, testClients.main.ArchiveWebhook(ctx, createdWebhook.ID))

			for i := 0; i < userCount; i++ {
				require.NoError(t, testClients.admin.ArchiveUser(ctx, users[i].ID))
			}
		}
	})
}

func (s *TestSuite) TestHouseholds_OwnershipTransfer() {
	s.runForEachClientExcept("should be possible to transfer ownership of an household", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create users
			futureOwner, _, _, futureOwnerClient := createUserAndClientForTest(ctx, t)

			// fetch household data
			householdCreationInput := &types.HouseholdCreationInput{
				Name: fakes.BuildFakeHousehold().Name,
			}
			household, householdCreationErr := testClients.main.CreateHousehold(ctx, householdCreationInput)
			require.NoError(t, householdCreationErr)
			require.NotNil(t, household)

			t.Logf("created household #%d", household.ID)

			require.NoError(t, testClients.main.SwitchActiveHousehold(ctx, household.ID))

			t.Logf("switched to active household: %d", household.ID)

			// create a webhook
			exampleWebhook := fakes.BuildFakeWebhook()
			exampleWebhookInput := fakes.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)
			createdWebhook, creationErr := testClients.main.CreateWebhook(ctx, exampleWebhookInput)
			requireNotNilAndNoProblems(t, createdWebhook, creationErr)

			t.Logf("created webhook #%d belonging to household #%d", createdWebhook.ID, createdWebhook.BelongsToHousehold)
			require.Equal(t, household.ID, createdWebhook.BelongsToHousehold)

			// check that user cannot see the webhook
			webhook, err := futureOwnerClient.GetWebhook(ctx, createdWebhook.ID)
			require.Nil(t, webhook)
			require.Error(t, err)

			// add them to the household
			require.NoError(t, testClients.main.TransferHouseholdOwnership(ctx, household.ID, &types.HouseholdOwnershipTransferInput{
				Reason:       t.Name(),
				CurrentOwner: household.BelongsToUser,
				NewOwner:     futureOwner.ID,
			}))

			t.Logf("transferred household %d from user %d to user %d", household.ID, household.BelongsToUser, futureOwner.ID)

			require.NoError(t, futureOwnerClient.SwitchActiveHousehold(ctx, household.ID))

			// check that user can see the webhook
			webhook, err = futureOwnerClient.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, webhook, err)

			// check that old user cannot see the webhook
			webhook, err = testClients.main.GetWebhook(ctx, createdWebhook.ID)
			require.Nil(t, webhook)
			require.Error(t, err)

			// check that new owner can update the webhook
			require.NoError(t, futureOwnerClient.UpdateWebhook(ctx, createdWebhook))

			// check audit log entries
			auditLogEntries, err := testClients.admin.GetAuditLogForHousehold(ctx, household.ID)
			require.NoError(t, err)

			expectedAuditLogEntries := []*types.AuditLogEntry{
				{EventType: audit.HouseholdCreationEvent},
				{EventType: audit.UserAddedToHouseholdEvent},
				{EventType: audit.WebhookCreationEvent},
				{EventType: audit.HouseholdTransferredEvent},
				{EventType: audit.WebhookUpdateEvent},
			}
			validateAuditLogEntries(t, expectedAuditLogEntries, auditLogEntries, household.ID, audit.HouseholdAssignmentKey)

			// Clean up.
			require.Error(t, testClients.main.ArchiveWebhook(ctx, createdWebhook.ID))
			require.NoError(t, futureOwnerClient.ArchiveWebhook(ctx, createdWebhook.ID))

			require.NoError(t, testClients.admin.ArchiveUser(ctx, futureOwner.ID))
		}
	})
}

func (s *TestSuite) TestHouseholds_Auditing_Returns404ForNonexistentHousehold() {
	s.runForEachClientExcept("should not be possible to audit a non-existent household", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			x, err := testClients.admin.GetAuditLogForHousehold(ctx, nonexistentID)

			assert.NoError(t, err)
			assert.Empty(t, x)
		}
	})
}

func (s *TestSuite) TestHouseholds_Auditing_InaccessibleToNonAdmins() {
	s.runForEachClientExcept("should not be possible to audit an household as non-admin", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create household.
			exampleHousehold := fakes.BuildFakeHousehold()
			exampleHouseholdInput := fakes.BuildFakeHouseholdCreationInputFromHousehold(exampleHousehold)
			createdHousehold, err := testClients.main.CreateHousehold(ctx, exampleHouseholdInput)
			requireNotNilAndNoProblems(t, createdHousehold, err)

			// fetch audit log entries
			actual, err := testClients.main.GetAuditLogForHousehold(ctx, createdHousehold.ID)
			assert.Error(t, err)
			assert.Nil(t, actual)

			// Clean up household.
			assert.NoError(t, testClients.main.ArchiveHousehold(ctx, createdHousehold.ID))
		}
	})
}

func (s *TestSuite) TestHouseholds_Auditing() {
	s.runForEachClientExcept("should be possible to audit an household", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create household.
			exampleHousehold := fakes.BuildFakeHousehold()
			exampleHouseholdInput := fakes.BuildFakeHouseholdCreationInputFromHousehold(exampleHousehold)
			createdHousehold, err := testClients.main.CreateHousehold(ctx, exampleHouseholdInput)
			requireNotNilAndNoProblems(t, createdHousehold, err)

			// fetch audit log entries
			actual, err := testClients.admin.GetAuditLogForHousehold(ctx, createdHousehold.ID)
			assert.NoError(t, err)
			assert.NotNil(t, actual)

			// Clean up household.
			assert.NoError(t, testClients.main.ArchiveHousehold(ctx, createdHousehold.ID))
		}
	})
}
