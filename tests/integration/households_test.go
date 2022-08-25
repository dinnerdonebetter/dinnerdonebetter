package integration

import (
	"testing"

	"github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/internal/authorization"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/client/httpclient"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
	testutils "github.com/prixfixeco/api_server/tests/utils"
)

func checkHouseholdEquality(t *testing.T, expected, actual *types.Household) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for household %s to be %v, but it was %v ", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.TimeZone, actual.TimeZone, "expected TimeZone for household %s to be %v, but it was %v ", expected.ID, expected.TimeZone, actual.TimeZone)
	assert.NotZero(t, actual.CreatedOn)
}

func (s *TestSuite) TestHouseholds_Creating() {
	s.runForEachClient("should be possible to create households", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create household.
			exampleHousehold := fakes.BuildFakeHousehold()
			exampleHouseholdInput := fakes.BuildFakeHouseholdCreationRequestInputFromHousehold(exampleHousehold)
			createdHousehold, err := testClients.user.CreateHousehold(ctx, exampleHouseholdInput)
			requireNotNilAndNoProblems(t, createdHousehold, err)

			// Assert household equality.
			checkHouseholdEquality(t, exampleHousehold, createdHousehold)

			// Clean up.
			assert.NoError(t, testClients.user.ArchiveHousehold(ctx, createdHousehold.ID))
		}
	})
}

func (s *TestSuite) TestHouseholds_Listing() {
	s.runForEachClient("should be possible to list households", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create households.
			var expected []*types.Household
			for i := 0; i < 5; i++ {
				// Create household.
				exampleHousehold := fakes.BuildFakeHousehold()
				exampleHouseholdInput := fakes.BuildFakeHouseholdCreationRequestInputFromHousehold(exampleHousehold)
				createdHousehold, householdCreationErr := testClients.user.CreateHousehold(ctx, exampleHouseholdInput)
				requireNotNilAndNoProblems(t, createdHousehold, householdCreationErr)

				expected = append(expected, createdHousehold)
			}

			// Assert household list equality.
			actual, err := testClients.user.GetHouseholds(ctx, nil)
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
				assert.NoError(t, testClients.user.ArchiveHousehold(ctx, createdHousehold.ID))
			}
		}
	})
}

func (s *TestSuite) TestHouseholds_Reading_Returns404ForNonexistentHousehold() {
	s.runForEachClient("should not be possible to read a non-existent household", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Attempt to fetch nonexistent household.
			_, err := testClients.user.GetHousehold(ctx, nonexistentID)
			assert.Error(t, err)
		}
	})
}

func (s *TestSuite) TestHouseholds_Reading() {
	s.runForEachClient("should be possible to read a household", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create household.
			exampleHousehold := fakes.BuildFakeHousehold()
			exampleHouseholdInput := fakes.BuildFakeHouseholdCreationRequestInputFromHousehold(exampleHousehold)
			createdHousehold, err := testClients.user.CreateHousehold(ctx, exampleHouseholdInput)
			requireNotNilAndNoProblems(t, createdHousehold, err)

			// Fetch household.
			actual, err := testClients.user.GetHousehold(ctx, createdHousehold.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// Assert household equality.
			checkHouseholdEquality(t, exampleHousehold, actual)

			// Clean up household.
			assert.NoError(t, testClients.user.ArchiveHousehold(ctx, createdHousehold.ID))
		}
	})
}

func (s *TestSuite) TestHouseholds_Updating_Returns404ForNonexistentHousehold() {
	s.runForEachClient("should not be possible to update a non-existent household", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			exampleHousehold := fakes.BuildFakeHousehold()
			exampleHousehold.ID = nonexistentID

			assert.Error(t, testClients.user.UpdateHousehold(ctx, exampleHousehold))
		}
	})
}

// convertHouseholdToHouseholdUpdateInput creates a householdUpdateInput struct from a household.
func convertHouseholdToHouseholdUpdateInput(x *types.Household) *types.HouseholdUpdateRequestInput {
	return &types.HouseholdUpdateRequestInput{
		Name:          &x.Name,
		BelongsToUser: x.BelongsToUser,
	}
}

func (s *TestSuite) TestHouseholds_Updating() {
	s.runForEachClient("should be possible to update a household", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create household.
			exampleHousehold := fakes.BuildFakeHousehold()
			exampleHouseholdInput := fakes.BuildFakeHouseholdCreationRequestInputFromHousehold(exampleHousehold)
			createdHousehold, err := testClients.user.CreateHousehold(ctx, exampleHouseholdInput)
			requireNotNilAndNoProblems(t, createdHousehold, err)

			// Change household.
			createdHousehold.Update(convertHouseholdToHouseholdUpdateInput(exampleHousehold))
			assert.NoError(t, testClients.user.UpdateHousehold(ctx, createdHousehold))

			// Fetch household.
			actual, err := testClients.user.GetHousehold(ctx, createdHousehold.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// Assert household equality.
			checkHouseholdEquality(t, exampleHousehold, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			// Clean up household.
			assert.NoError(t, testClients.user.ArchiveHousehold(ctx, createdHousehold.ID))
		}
	})
}

func (s *TestSuite) TestHouseholds_Archiving_Returns404ForNonexistentHousehold() {
	s.runForEachClient("should not be possible to archive a non-existent household", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			assert.Error(t, testClients.user.ArchiveHousehold(ctx, nonexistentID))
		}
	})
}

func (s *TestSuite) TestHouseholds_Archiving() {
	s.runForEachClient("should be possible to archive a household", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create household.
			exampleHousehold := fakes.BuildFakeHousehold()
			exampleHouseholdInput := fakes.BuildFakeHouseholdCreationRequestInputFromHousehold(exampleHousehold)
			createdHousehold, err := testClients.user.CreateHousehold(ctx, exampleHouseholdInput)
			requireNotNilAndNoProblems(t, createdHousehold, err)

			// Clean up household.
			assert.NoError(t, testClients.user.ArchiveHousehold(ctx, createdHousehold.ID))
		}
	})
}

func (s *TestSuite) TestHouseholds_InvitingPreExistentUser() {
	s.runForEachClient("should be possible to invite an already-registered user", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Logf("determining household ID")
			currentStatus, statusErr := testClients.user.UserStatus(s.ctx)
			requireNotNilAndNoProblems(t, currentStatus, statusErr)
			relevantHouseholdID := currentStatus.ActiveHousehold
			t.Logf("initial household is %s; initial user ID is %s", relevantHouseholdID, s.user.ID)

			// Create webhook.
			exampleWebhook := fakes.BuildFakeWebhook()
			exampleWebhookInput := fakes.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)
			createdWebhook, err := testClients.user.CreateWebhook(ctx, exampleWebhookInput)
			require.NoError(t, err)

			checkWebhookEquality(t, exampleWebhook, createdWebhook)

			createdWebhook, err = testClients.user.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, createdWebhook, err)
			require.Equal(t, relevantHouseholdID, createdWebhook.BelongsToHousehold)

			t.Logf("creating user to invite")
			u, _, c, _ := createUserAndClientForTest(ctx, t, nil)

			t.Logf("inviting user")
			invitation, err := testClients.user.InviteUserToHousehold(ctx, &types.HouseholdInvitationCreationRequestInput{
				FromUser:               s.user.ID,
				Note:                   t.Name(),
				ToEmail:                u.EmailAddress,
				DestinationHouseholdID: relevantHouseholdID,
			})
			require.NoError(t, err)

			t.Logf("checking for sent invitation")
			sentInvitations, err := testClients.user.GetPendingHouseholdInvitationsFromUser(ctx, nil)
			requireNotNilAndNoProblems(t, sentInvitations, err)
			assert.NotEmpty(t, sentInvitations.HouseholdInvitations)

			t.Logf("checking for received invitation")
			invitations, err := c.GetPendingHouseholdInvitationsForUser(ctx, nil)
			requireNotNilAndNoProblems(t, invitations, err)
			assert.NotEmpty(t, invitations.HouseholdInvitations)

			t.Logf("accepting invitation")
			err = c.AcceptHouseholdInvitation(ctx, invitation.ID, invitation.Token, t.Name())
			require.NoError(t, err)

			t.Logf("checking for sent invitation")
			sentInvitations, err = testClients.user.GetPendingHouseholdInvitationsFromUser(ctx, nil)
			requireNotNilAndNoProblems(t, sentInvitations, err)
			assert.Empty(t, sentInvitations.HouseholdInvitations)

			t.Logf("fetching households")
			households, err := c.GetHouseholds(ctx, nil)

			var found bool
			for _, household := range households.Households {
				if !found {
					found = household.ID == relevantHouseholdID
				}
			}

			require.True(t, found)
			require.NoError(t, c.SwitchActiveHousehold(ctx, relevantHouseholdID))

			_, err = c.GetWebhook(ctx, createdWebhook.ID)
			require.NoError(t, err)
		}
	})
}

func (s *TestSuite) TestHouseholds_InvitingUserWhoSignsUpIndependently() {
	s.runForEachClient("should be possible to invite a user before they sign up", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Logf("determining household ID")
			currentStatus, statusErr := testClients.user.UserStatus(s.ctx)
			requireNotNilAndNoProblems(t, currentStatus, statusErr)
			relevantHouseholdID := currentStatus.ActiveHousehold
			t.Logf("initial household is %s; initial user ID is %s", relevantHouseholdID, s.user.ID)

			// Create webhook.
			exampleWebhook := fakes.BuildFakeWebhook()
			exampleWebhookInput := fakes.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)
			createdWebhook, err := testClients.user.CreateWebhook(ctx, exampleWebhookInput)
			require.NoError(t, err)

			checkWebhookEquality(t, exampleWebhook, createdWebhook)

			createdWebhook, err = testClients.user.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, createdWebhook, err)
			require.Equal(t, relevantHouseholdID, createdWebhook.BelongsToHousehold)

			t.Logf("inviting user")
			inviteReq := &types.HouseholdInvitationCreationRequestInput{
				FromUser:               s.user.ID,
				Note:                   t.Name(),
				ToEmail:                gofakeit.Email(),
				DestinationHouseholdID: relevantHouseholdID,
			}
			invitation, err := testClients.user.InviteUserToHousehold(ctx, inviteReq)
			require.NoError(t, err)

			t.Logf("checking for sent invitation")
			sentInvitations, err := testClients.user.GetPendingHouseholdInvitationsFromUser(ctx, nil)
			requireNotNilAndNoProblems(t, sentInvitations, err)
			assert.NotEmpty(t, sentInvitations.HouseholdInvitations)

			t.Logf("creating user to invite")
			_, _, c, _ := createUserAndClientForTest(ctx, t, &types.UserRegistrationInput{
				EmailAddress: inviteReq.ToEmail,
				Username:     fakes.BuildFakeUser().Username,
				Password:     gofakeit.Password(true, true, true, true, false, 64),
			})

			t.Logf("checking for invitation")
			invitations, err := c.GetPendingHouseholdInvitationsForUser(ctx, nil)
			requireNotNilAndNoProblems(t, invitations, err)
			assert.NotEmpty(t, invitations.HouseholdInvitations)

			t.Logf("accepting invitation")
			err = c.AcceptHouseholdInvitation(ctx, invitation.ID, invitation.Token, t.Name())
			require.NoError(t, err)

			t.Logf("fetching households")
			households, err := c.GetHouseholds(ctx, nil)

			var found bool
			for _, household := range households.Households {
				if !found {
					found = household.ID == relevantHouseholdID
				}
			}

			require.True(t, found)
			require.NoError(t, c.SwitchActiveHousehold(ctx, relevantHouseholdID))

			_, err = c.GetWebhook(ctx, createdWebhook.ID)
			require.NoError(t, err)

			t.Logf("checking for sent invitation")
			sentInvitations, err = testClients.user.GetPendingHouseholdInvitationsFromUser(ctx, nil)
			requireNotNilAndNoProblems(t, sentInvitations, err)
			assert.Empty(t, sentInvitations.HouseholdInvitations)
		}
	})
}

func (s *TestSuite) TestHouseholds_InvitingUserWhoSignsUpIndependentlyAndThenCancelling() {
	s.runForEachClient("should be possible to invite a user before they sign up and cancel before they can accept", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Logf("determining household ID")
			currentStatus, statusErr := testClients.user.UserStatus(s.ctx)
			requireNotNilAndNoProblems(t, currentStatus, statusErr)
			relevantHouseholdID := currentStatus.ActiveHousehold
			t.Logf("initial household is %s; initial user ID is %s", relevantHouseholdID, s.user.ID)

			// Create webhook.
			exampleWebhook := fakes.BuildFakeWebhook()
			exampleWebhookInput := fakes.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)
			createdWebhook, err := testClients.user.CreateWebhook(ctx, exampleWebhookInput)
			require.NoError(t, err)

			checkWebhookEquality(t, exampleWebhook, createdWebhook)

			createdWebhook, err = testClients.user.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, createdWebhook, err)
			require.Equal(t, relevantHouseholdID, createdWebhook.BelongsToHousehold)

			t.Logf("inviting user")
			inviteReq := &types.HouseholdInvitationCreationRequestInput{
				FromUser:               s.user.ID,
				Note:                   t.Name(),
				ToEmail:                gofakeit.Email(),
				DestinationHouseholdID: relevantHouseholdID,
			}
			invitation, err := testClients.user.InviteUserToHousehold(ctx, inviteReq)
			require.NoError(t, err)

			t.Logf("checking for sent invitation")
			sentInvitations, err := testClients.user.GetPendingHouseholdInvitationsFromUser(ctx, nil)
			requireNotNilAndNoProblems(t, sentInvitations, err)
			assert.NotEmpty(t, sentInvitations.HouseholdInvitations)

			t.Logf("creating user to invite")
			_, _, c, _ := createUserAndClientForTest(ctx, t, &types.UserRegistrationInput{
				EmailAddress: inviteReq.ToEmail,
				Username:     fakes.BuildFakeUser().Username,
				Password:     gofakeit.Password(true, true, true, true, false, 64),
			})

			t.Logf("checking for invitation")
			invitations, err := c.GetPendingHouseholdInvitationsForUser(ctx, nil)
			requireNotNilAndNoProblems(t, invitations, err)
			assert.NotEmpty(t, invitations.HouseholdInvitations)

			t.Logf("cancelling invitation")
			err = testClients.user.CancelHouseholdInvitation(ctx, invitation.ID, invitation.Token, t.Name())
			require.NoError(t, err)

			t.Logf("checking for sent invitation")
			sentInvitations, err = testClients.user.GetPendingHouseholdInvitationsFromUser(ctx, nil)
			requireNotNilAndNoProblems(t, sentInvitations, err)
			assert.Empty(t, sentInvitations.HouseholdInvitations)
		}
	})
}

func (s *TestSuite) TestHouseholds_InvitingNewUserWithInviteLink() {
	s.runForEachClient("should be possible to invite a user via referral link", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Logf("determining household ID")
			currentStatus, statusErr := testClients.user.UserStatus(s.ctx)
			requireNotNilAndNoProblems(t, currentStatus, statusErr)
			relevantHouseholdID := currentStatus.ActiveHousehold
			t.Logf("initial household is %s; initial user ID is %s", relevantHouseholdID, s.user.ID)

			// Create webhook.
			exampleWebhook := fakes.BuildFakeWebhook()
			exampleWebhookInput := fakes.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)
			createdWebhook, err := testClients.user.CreateWebhook(ctx, exampleWebhookInput)
			require.NoError(t, err)

			checkWebhookEquality(t, exampleWebhook, createdWebhook)

			createdWebhook, err = testClients.user.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, createdWebhook, err)
			require.Equal(t, relevantHouseholdID, createdWebhook.BelongsToHousehold)

			t.Logf("inviting user")
			inviteReq := &types.HouseholdInvitationCreationRequestInput{
				FromUser:               s.user.ID,
				Note:                   t.Name(),
				ToEmail:                gofakeit.Email(),
				DestinationHouseholdID: relevantHouseholdID,
			}
			createdInvitation, err := testClients.user.InviteUserToHousehold(ctx, inviteReq)
			require.NoError(t, err)

			createdInvitation, err = testClients.user.GetHouseholdInvitation(ctx, relevantHouseholdID, createdInvitation.ID)
			requireNotNilAndNoProblems(t, createdInvitation, err)

			t.Logf("checking for sent invitation")
			sentInvitations, err := testClients.user.GetPendingHouseholdInvitationsFromUser(ctx, nil)
			requireNotNilAndNoProblems(t, sentInvitations, err)
			assert.NotEmpty(t, sentInvitations.HouseholdInvitations)

			t.Logf("creating user to invite")
			_, _, c, _ := createUserAndClientForTest(ctx, t, &types.UserRegistrationInput{
				EmailAddress:    inviteReq.ToEmail,
				Username:        fakes.BuildFakeUser().Username,
				Password:        gofakeit.Password(true, true, true, true, false, 64),
				InvitationID:    createdInvitation.ID,
				InvitationToken: createdInvitation.Token,
			})

			t.Logf("fetching households")
			households, err := c.GetHouseholds(ctx, nil)
			require.NoError(t, err)

			var found bool
			for _, household := range households.Households {
				if !found {
					found = household.ID == relevantHouseholdID
				}
			}

			require.True(t, found)

			_, err = c.GetWebhook(ctx, createdWebhook.ID)
			require.NoError(t, err)
		}
	})
}

func (s *TestSuite) TestHouseholds_InviteCanBeCancelled() {
	s.runForEachClient("should be possible to invite an already-registered user", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Logf("determining household ID")
			currentStatus, statusErr := testClients.user.UserStatus(s.ctx)
			requireNotNilAndNoProblems(t, currentStatus, statusErr)
			relevantHouseholdID := currentStatus.ActiveHousehold
			t.Logf("initial household is %s; initial user ID is %s", relevantHouseholdID, s.user.ID)

			// Create webhook.
			exampleWebhook := fakes.BuildFakeWebhook()
			exampleWebhookInput := fakes.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)
			createdWebhook, err := testClients.user.CreateWebhook(ctx, exampleWebhookInput)
			require.NoError(t, err)

			checkWebhookEquality(t, exampleWebhook, createdWebhook)

			createdWebhook, err = testClients.user.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, createdWebhook, err)
			require.Equal(t, relevantHouseholdID, createdWebhook.BelongsToHousehold)

			t.Logf("inviting user")
			inviteReq := &types.HouseholdInvitationCreationRequestInput{
				FromUser:               s.user.ID,
				Note:                   t.Name(),
				ToEmail:                gofakeit.Email(),
				DestinationHouseholdID: relevantHouseholdID,
			}
			invitation, err := testClients.user.InviteUserToHousehold(ctx, inviteReq)
			require.NoError(t, err)

			require.NoError(t, testClients.user.CancelHouseholdInvitation(ctx, invitation.ID, invitation.Token, t.Name()))

			t.Logf("checking for sent invitation")
			sentInvitations, err := testClients.user.GetPendingHouseholdInvitationsFromUser(ctx, nil)
			requireNotNilAndNoProblems(t, sentInvitations, err)
			assert.Empty(t, sentInvitations.HouseholdInvitations)

			t.Logf("creating user to invite")
			_, _, c, _ := createUserAndClientForTest(ctx, t, &types.UserRegistrationInput{
				EmailAddress: inviteReq.ToEmail,
				Username:     fakes.BuildFakeUser().Username,
				Password:     gofakeit.Password(true, true, true, true, false, 64),
			})

			t.Logf("checking for invitation")
			invitations, err := c.GetPendingHouseholdInvitationsForUser(ctx, nil)
			requireNotNilAndNoProblems(t, invitations, err)
			assert.Empty(t, invitations.HouseholdInvitations)
		}
	})
}

func (s *TestSuite) TestHouseholds_InviteCanBeRejected() {
	s.runForEachClient("should be possible to invite an already-registered user", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Logf("determining household ID")
			currentStatus, statusErr := testClients.user.UserStatus(s.ctx)
			requireNotNilAndNoProblems(t, currentStatus, statusErr)
			relevantHouseholdID := currentStatus.ActiveHousehold
			t.Logf("initial household is %s; initial user ID is %s", relevantHouseholdID, s.user.ID)

			// Create webhook.
			exampleWebhook := fakes.BuildFakeWebhook()
			exampleWebhookInput := fakes.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)
			createdWebhook, err := testClients.user.CreateWebhook(ctx, exampleWebhookInput)
			require.NoError(t, err)

			checkWebhookEquality(t, exampleWebhook, createdWebhook)

			createdWebhook, err = testClients.user.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, createdWebhook, err)
			require.Equal(t, relevantHouseholdID, createdWebhook.BelongsToHousehold)

			t.Logf("creating user to invite")
			u, _, c, _ := createUserAndClientForTest(ctx, t, nil)

			t.Logf("inviting user")
			invitation, err := testClients.user.InviteUserToHousehold(ctx, &types.HouseholdInvitationCreationRequestInput{
				FromUser:               s.user.ID,
				Note:                   t.Name(),
				ToEmail:                u.EmailAddress,
				DestinationHouseholdID: relevantHouseholdID,
			})
			require.NoError(t, err)

			t.Logf("checking for sent invitation")
			sentInvitations, err := testClients.user.GetPendingHouseholdInvitationsFromUser(ctx, nil)
			requireNotNilAndNoProblems(t, sentInvitations, err)
			assert.NotEmpty(t, sentInvitations.HouseholdInvitations)

			t.Logf("checking for received invitation")
			invitations, err := c.GetPendingHouseholdInvitationsForUser(ctx, nil)
			requireNotNilAndNoProblems(t, invitations, err)
			assert.NotEmpty(t, invitations.HouseholdInvitations)

			t.Logf("accepting invitation")
			err = c.RejectHouseholdInvitation(ctx, invitation.ID, invitation.Token, t.Name())
			require.NoError(t, err)

			t.Logf("checking for sent invitation")
			sentInvitations, err = testClients.user.GetPendingHouseholdInvitationsFromUser(ctx, nil)
			requireNotNilAndNoProblems(t, sentInvitations, err)
			assert.Empty(t, sentInvitations.HouseholdInvitations)
		}
	})
}

func (s *TestSuite) TestHouseholds_ChangingMemberships() {
	s.runForEachClient("should be possible to change members of a household", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			const userCount = 1

			currentStatus, statusErr := testClients.user.UserStatus(s.ctx)
			requireNotNilAndNoProblems(t, currentStatus, statusErr)
			t.Logf("initial household is %s; initial user ID is %s", currentStatus.ActiveHousehold, s.user.ID)

			// fetch household data
			householdCreationInput := &types.HouseholdCreationRequestInput{
				Name: fakes.BuildFakeHousehold().Name,
			}
			household, householdCreationErr := testClients.user.CreateHousehold(ctx, householdCreationInput)
			require.NoError(t, householdCreationErr)
			require.NotNil(t, household)

			t.Logf("created household %s", household.ID)

			require.NoError(t, testClients.user.SwitchActiveHousehold(ctx, household.ID))

			t.Logf("switched user test client active household to %s, creating webhook", household.ID)

			// Create webhook.
			exampleWebhook := fakes.BuildFakeWebhook()
			exampleWebhookInput := fakes.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)
			createdWebhook, err := testClients.user.CreateWebhook(ctx, exampleWebhookInput)
			require.NoError(t, err)

			checkWebhookEquality(t, exampleWebhook, createdWebhook)

			createdWebhook, err = testClients.user.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, createdWebhook, err)
			require.Equal(t, household.ID, createdWebhook.BelongsToHousehold)

			t.Logf("created webhook %s for household %s", createdWebhook.ID, createdWebhook.BelongsToHousehold)

			// create dummy users
			users := []*types.User{}
			clients := []*httpclient.Client{}

			// create users
			for i := 0; i < userCount; i++ {
				u, _, c, _ := createUserAndClientForTest(ctx, t, nil)
				users = append(users, u)
				clients = append(clients, c)

				currentStatus, statusErr = c.UserStatus(s.ctx)
				requireNotNilAndNoProblems(t, currentStatus, statusErr)
				t.Logf("created user user %q with household %s", u.ID, currentStatus.ActiveHousehold)
			}

			// check that each user cannot see the unreachable webhook
			for i := 0; i < userCount; i++ {
				t.Logf("checking that user %q CANNOT see webhook %s belonging to household %s", users[i].ID, createdWebhook.ID, createdWebhook.BelongsToHousehold)
				webhook, err := clients[i].GetWebhook(ctx, createdWebhook.ID)
				require.Nil(t, webhook)
				require.Error(t, err)
			}

			// add them to the household
			for i := 0; i < userCount; i++ {
				t.Logf("adding user %q to household %s", users[i].ID, household.ID)
				invitation, invitationErr := testClients.user.InviteUserToHousehold(ctx, &types.HouseholdInvitationCreationRequestInput{
					ToEmail:                users[i].EmailAddress,
					DestinationHouseholdID: household.ID,
					Note:                   t.Name(),
				})
				require.NoError(t, invitationErr)
				require.NotEmpty(t, invitation.ID)
				t.Logf("invited user %q to household %s", users[i].ID, household.ID)

				t.Logf("checking for received invitation")
				invitations, fetchInvitationsErr := clients[i].GetPendingHouseholdInvitationsForUser(ctx, nil)
				requireNotNilAndNoProblems(t, invitations, fetchInvitationsErr)
				assert.NotEmpty(t, invitations.HouseholdInvitations)

				t.Logf("accepting invitation")
				err = clients[i].AcceptHouseholdInvitation(ctx, invitation.ID, invitation.Token, t.Name())
				require.NoError(t, err)

				t.Logf("setting user %q's client to household %s", users[i].ID, household.ID)
				require.NoError(t, clients[i].SwitchActiveHousehold(ctx, household.ID))

				currentStatus, statusErr = clients[i].UserStatus(s.ctx)
				requireNotNilAndNoProblems(t, currentStatus, statusErr)
				require.Equal(t, currentStatus.ActiveHousehold, household.ID)
				t.Logf("set user %q's current active household to %s", users[i].ID, household.ID)
			}

			// grant all permissions
			for i := 0; i < userCount; i++ {
				input := &types.ModifyUserPermissionsInput{
					Reason:   t.Name(),
					NewRoles: []string{authorization.HouseholdAdminRole.String()},
				}
				require.NoError(t, testClients.user.ModifyMemberPermissions(ctx, household.ID, users[i].ID, input))
			}

			// check that each user can see the webhook
			for i := 0; i < userCount; i++ {
				t.Logf("checking if user %q CAN now see webhook %s belonging to household %s", users[i].ID, createdWebhook.ID, createdWebhook.BelongsToHousehold)
				webhook, webhookRetrievalError := clients[i].GetWebhook(ctx, createdWebhook.ID)
				requireNotNilAndNoProblems(t, webhook, webhookRetrievalError)
			}

			// remove users from household
			for i := 0; i < userCount; i++ {
				require.NoError(t, testClients.user.RemoveUserFromHousehold(ctx, household.ID, users[i].ID))
			}

			// check that each user cannot see the webhook
			for i := 0; i < userCount; i++ {
				webhook, webhookRetrievalError := clients[i].GetWebhook(ctx, createdWebhook.ID)
				require.Nil(t, webhook)
				require.Error(t, webhookRetrievalError)
			}

			// Clean up.
			require.NoError(t, testClients.user.ArchiveWebhook(ctx, createdWebhook.ID))

			for i := 0; i < userCount; i++ {
				require.NoError(t, testClients.admin.ArchiveUser(ctx, users[i].ID))
			}
		}
	})
}

func (s *TestSuite) TestHouseholds_OwnershipTransfer() {
	s.runForEachClient("should be possible to transfer ownership of a household", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create users
			futureOwner, _, _, futureOwnerClient := createUserAndClientForTest(ctx, t, nil)

			// fetch household data
			householdCreationInput := &types.HouseholdCreationRequestInput{
				Name: fakes.BuildFakeHousehold().Name,
			}
			household, householdCreationErr := testClients.user.CreateHousehold(ctx, householdCreationInput)
			require.NoError(t, householdCreationErr)
			require.NotNil(t, household)

			t.Logf("created household %s", household.ID)

			require.NoError(t, testClients.user.SwitchActiveHousehold(ctx, household.ID))

			t.Logf("switched to active household: %s", household.ID)

			// create a webhook

			// Create webhook.
			exampleWebhook := fakes.BuildFakeWebhook()
			exampleWebhookInput := fakes.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)
			createdWebhook, err := testClients.user.CreateWebhook(ctx, exampleWebhookInput)
			require.NoError(t, err)

			checkWebhookEquality(t, exampleWebhook, createdWebhook)

			createdWebhook, err = testClients.user.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, createdWebhook, err)

			t.Logf("created webhook %s belonging to household %s", createdWebhook.ID, createdWebhook.BelongsToHousehold)
			require.Equal(t, household.ID, createdWebhook.BelongsToHousehold)

			// check that user cannot see the webhook
			webhook, err := futureOwnerClient.GetWebhook(ctx, createdWebhook.ID)
			require.Nil(t, webhook)
			require.Error(t, err)

			// add them to the household
			require.NoError(t, testClients.user.TransferHouseholdOwnership(ctx, household.ID, &types.HouseholdOwnershipTransferInput{
				Reason:       t.Name(),
				CurrentOwner: household.BelongsToUser,
				NewOwner:     futureOwner.ID,
			}))

			t.Logf("transferred household %s from user %s to user %s", household.ID, household.BelongsToUser, futureOwner.ID)

			require.NoError(t, futureOwnerClient.SwitchActiveHousehold(ctx, household.ID))

			// check that user can see the webhook
			webhook, err = futureOwnerClient.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, webhook, err)

			// check that old user cannot see the webhook
			webhook, err = testClients.user.GetWebhook(ctx, createdWebhook.ID)
			require.Nil(t, webhook)
			require.Error(t, err)

			// check that new owner can delete the webhook
			require.NoError(t, futureOwnerClient.ArchiveWebhook(ctx, createdWebhook.ID))

			// Clean up.
			require.Error(t, testClients.user.ArchiveWebhook(ctx, createdWebhook.ID))
			require.NoError(t, testClients.admin.ArchiveUser(ctx, futureOwner.ID))
		}
	})
}

func (s *TestSuite) TestHouseholds_UsersHaveBackupHouseholdCreatedForThemWhenRemovedFromLastHousehold() {
	s.runForEachClient("should be possible to invite a user via referral link", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Logf("determining household ID")
			currentStatus, statusErr := testClients.user.UserStatus(s.ctx)
			requireNotNilAndNoProblems(t, currentStatus, statusErr)
			relevantHouseholdID := currentStatus.ActiveHousehold
			t.Logf("initial household is %s; initial user ID is %s", relevantHouseholdID, s.user.ID)

			t.Logf("inviting user")
			inviteReq := &types.HouseholdInvitationCreationRequestInput{
				FromUser:               s.user.ID,
				Note:                   t.Name(),
				ToEmail:                gofakeit.Email(),
				DestinationHouseholdID: relevantHouseholdID,
			}
			createdInvitation, err := testClients.user.InviteUserToHousehold(ctx, inviteReq)
			require.NoError(t, err)

			createdInvitation, err = testClients.user.GetHouseholdInvitation(ctx, relevantHouseholdID, createdInvitation.ID)
			requireNotNilAndNoProblems(t, createdInvitation, err)

			t.Logf("checking for sent invitation")
			sentInvitations, err := testClients.user.GetPendingHouseholdInvitationsFromUser(ctx, nil)
			requireNotNilAndNoProblems(t, sentInvitations, err)
			assert.NotEmpty(t, sentInvitations.HouseholdInvitations)

			t.Logf("creating user to invite")

			regInput := &types.UserRegistrationInput{
				EmailAddress:    inviteReq.ToEmail,
				Username:        fakes.BuildFakeUser().Username,
				Password:        gofakeit.Password(true, true, true, true, false, 64),
				InvitationID:    createdInvitation.ID,
				InvitationToken: createdInvitation.Token,
			}
			u, _, c, _ := createUserAndClientForTest(ctx, t, regInput)

			t.Logf("fetching households")
			households, err := c.GetHouseholds(ctx, nil)
			require.NoError(t, err)

			assert.Len(t, households.Households, 2)

			var (
				found            bool
				otherHouseholdID string
			)

			for _, household := range households.Households {
				if household.ID == relevantHouseholdID {
					if !found {
						found = true
					}
				} else {
					otherHouseholdID = household.ID
				}
			}

			require.NotEmpty(t, otherHouseholdID)
			require.True(t, found)

			require.NoError(t, testClients.user.RemoveUserFromHousehold(ctx, relevantHouseholdID, u.ID))

			u.HashedPassword = regInput.Password

			newCookie, err := testutils.GetLoginCookie(ctx, urlToUse, u)
			require.NoError(t, err)

			require.NoError(t, c.SetOptions(httpclient.UsingCookie(newCookie)))

			household, err := c.GetCurrentHousehold(ctx)
			requireNotNilAndNoProblems(t, household, err)
			assert.NotEqual(t, relevantHouseholdID, household.ID)

			require.True(t, found)
		}
	})
}
