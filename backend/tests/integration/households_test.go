package integration

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/apiclient"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkHouseholdEquality(t *testing.T, expected, actual *types.Household) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for household %s to be %v, but it was %v ", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.AddressLine1, actual.AddressLine1, "expected AddressLine1 for household %s to be %v, but it was %v", expected.ID, expected.AddressLine1, actual.AddressLine1)
	assert.Equal(t, expected.AddressLine2, actual.AddressLine2, "expected AddressLine2 for household %s to be %v, but it was %v", expected.ID, expected.AddressLine2, actual.AddressLine2)
	assert.Equal(t, expected.City, actual.City, "expected City for household %s to be %v, but it was %v", expected.ID, expected.City, actual.City)
	assert.Equal(t, expected.State, actual.State, "expected State for household %s to be %v, but it was %v", expected.ID, expected.State, actual.State)
	assert.Equal(t, expected.ZipCode, actual.ZipCode, "expected ZipCode for household %s to be %v, but it was %v", expected.ID, expected.ZipCode, actual.ZipCode)
	assert.Equal(t, expected.Country, actual.Country, "expected Country for household %s to be %v, but it was %v", expected.ID, expected.Country, actual.Country)
	assert.Equal(t, expected.Latitude, actual.Latitude, "expected Latitude for household %s to be %v, but it was %v", expected.ID, *expected.Latitude, *actual.Latitude)
	assert.Equal(t, expected.Longitude, actual.Longitude, "expected Longitude for household %s to be %v, but it was %v", expected.ID, *expected.Longitude, *actual.Longitude)
	assert.NotZero(t, actual.CreatedAt)
}

func (s *TestSuite) TestHouseholds_Creating() {
	s.runTest("should be possible to create households", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create household.
			exampleHousehold := fakes.BuildFakeHousehold()
			exampleHouseholdInput := converters.ConvertHouseholdToHouseholdCreationRequestInput(exampleHousehold)
			createdHousehold, err := testClients.userClient.CreateHousehold(ctx, exampleHouseholdInput)
			requireNotNilAndNoProblems(t, createdHousehold, err)

			// Assert household equality.
			checkHouseholdEquality(t, exampleHousehold, createdHousehold)

			// Clean up.
			assert.NoError(t, testClients.userClient.ArchiveHousehold(ctx, createdHousehold.ID))
		}
	})
}

func (s *TestSuite) TestHouseholds_Listing() {
	s.runTest("should be possible to list households", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create households.
			var expected []*types.Household
			for i := 0; i < 5; i++ {
				// Create household.
				exampleHousehold := fakes.BuildFakeHousehold()
				exampleHouseholdInput := converters.ConvertHouseholdToHouseholdCreationRequestInput(exampleHousehold)
				createdHousehold, householdCreationErr := testClients.userClient.CreateHousehold(ctx, exampleHouseholdInput)
				requireNotNilAndNoProblems(t, createdHousehold, householdCreationErr)

				expected = append(expected, createdHousehold)
			}

			// Assert household list equality.
			actual, err := testClients.userClient.GetHouseholds(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Data),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Data),
			)

			// Clean up.
			for _, createdHousehold := range actual.Data {
				assert.NoError(t, testClients.userClient.ArchiveHousehold(ctx, createdHousehold.ID))
			}
		}
	})
}

func (s *TestSuite) TestHouseholds_Reading_Returns404ForNonexistentHousehold() {
	s.runTest("should not be possible to read a non-existent household", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Attempt to fetch nonexistent household.
			_, err := testClients.userClient.GetHousehold(ctx, nonexistentID)
			assert.Error(t, err)
		}
	})
}

func (s *TestSuite) TestHouseholds_Reading() {
	s.runTest("should be possible to read a household", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create household.
			exampleHousehold := fakes.BuildFakeHousehold()
			exampleHouseholdInput := converters.ConvertHouseholdToHouseholdCreationRequestInput(exampleHousehold)
			createdHousehold, err := testClients.userClient.CreateHousehold(ctx, exampleHouseholdInput)
			requireNotNilAndNoProblems(t, createdHousehold, err)

			// Fetch household.
			actual, err := testClients.userClient.GetHousehold(ctx, createdHousehold.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// Assert household equality.
			checkHouseholdEquality(t, exampleHousehold, actual)

			// Clean up household.
			assert.NoError(t, testClients.userClient.ArchiveHousehold(ctx, createdHousehold.ID))
		}
	})
}

func (s *TestSuite) TestHouseholds_Updating_Returns404ForNonexistentHousehold() {
	s.runTest("should not be possible to update a non-existent household", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			assert.Error(t, testClients.userClient.UpdateHousehold(ctx, nonexistentID, &types.HouseholdUpdateRequestInput{}))
		}
	})
}

func (s *TestSuite) TestHouseholds_Updating() {
	s.runTest("should be possible to update a household", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create household.
			exampleHousehold := fakes.BuildFakeHousehold()
			exampleHouseholdInput := converters.ConvertHouseholdToHouseholdCreationRequestInput(exampleHousehold)
			createdHousehold, err := testClients.userClient.CreateHousehold(ctx, exampleHouseholdInput)
			requireNotNilAndNoProblems(t, createdHousehold, err)

			// Change household.
			updateInput := converters.ConvertHouseholdToHouseholdUpdateRequestInput(exampleHousehold)
			createdHousehold.Update(updateInput)
			assert.NoError(t, testClients.userClient.UpdateHousehold(ctx, createdHousehold.ID, updateInput))

			// Fetch household.
			actual, err := testClients.userClient.GetHousehold(ctx, createdHousehold.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// Assert household equality.
			checkHouseholdEquality(t, exampleHousehold, actual)
			assert.NotNil(t, actual.LastUpdatedAt)

			// Clean up household.
			assert.NoError(t, testClients.userClient.ArchiveHousehold(ctx, createdHousehold.ID))
		}
	})
}

func (s *TestSuite) TestHouseholds_Archiving() {
	s.runTest("should be possible to archive a household", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create household.
			exampleHousehold := fakes.BuildFakeHousehold()
			exampleHouseholdInput := converters.ConvertHouseholdToHouseholdCreationRequestInput(exampleHousehold)
			createdHousehold, err := testClients.userClient.CreateHousehold(ctx, exampleHouseholdInput)
			requireNotNilAndNoProblems(t, createdHousehold, err)

			// Clean up household.
			assert.NoError(t, testClients.userClient.ArchiveHousehold(ctx, createdHousehold.ID))
		}
	})
}

func (s *TestSuite) TestHouseholds_InvitingPreExistentUser() {
	s.runTest("should be possible to invite an already-registered userClient", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			currentStatus, statusErr := testClients.userClient.GetAuthStatus(s.ctx)
			requireNotNilAndNoProblems(t, currentStatus, statusErr)
			relevantHouseholdID := currentStatus.ActiveHousehold

			// Create webhook.
			exampleWebhook := fakes.BuildFakeWebhook()
			exampleWebhookInput := converters.ConvertWebhookToWebhookCreationRequestInput(exampleWebhook)
			createdWebhook, err := testClients.userClient.CreateWebhook(ctx, exampleWebhookInput)
			require.NoError(t, err)

			checkWebhookEquality(t, exampleWebhook, createdWebhook)

			createdWebhook, err = testClients.userClient.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, createdWebhook, err)
			require.Equal(t, relevantHouseholdID, createdWebhook.BelongsToHousehold)

			u, c := createUserAndClientForTest(ctx, t, nil)

			invitation, err := testClients.userClient.CreateHouseholdInvitation(ctx, relevantHouseholdID, &types.HouseholdInvitationCreationRequestInput{
				Note:    t.Name(),
				ToName:  t.Name(),
				ToEmail: u.EmailAddress,
			})
			require.NoError(t, err)

			sentInvitations, err := testClients.userClient.GetSentHouseholdInvitations(ctx, nil)
			requireNotNilAndNoProblems(t, sentInvitations, err)
			assert.NotEmpty(t, sentInvitations.Data)

			invitations, err := c.GetReceivedHouseholdInvitations(ctx, nil)
			requireNotNilAndNoProblems(t, invitations, err)
			assert.NotEmpty(t, invitations.Data)

			err = c.AcceptHouseholdInvitation(ctx, invitation.ID, &types.HouseholdInvitationUpdateRequestInput{
				Token: invitation.Token,
				Note:  t.Name(),
			})
			require.NoError(t, err)

			sentInvitations, err = testClients.userClient.GetSentHouseholdInvitations(ctx, nil)
			requireNotNilAndNoProblems(t, sentInvitations, err)
			assert.Empty(t, sentInvitations.Data)

			households, err := c.GetHouseholds(ctx, nil)

			var found bool
			for _, household := range households.Data {
				if !found {
					found = household.ID == relevantHouseholdID
				}
			}

			require.True(t, found)
			_, err = c.SetDefaultHousehold(ctx, relevantHouseholdID)
			require.NoError(t, err)

			tokenResponse, err := c.LoginForJWT(ctx, &types.UserLoginInput{Username: u.Username, Password: u.HashedPassword, TOTPToken: generateTOTPTokenForUser(t, u)})
			require.NoError(t, err)

			require.NoError(t, c.SetOptions(apiclient.UsingOAuth2(ctx, createdClientID, createdClientSecret, []string{"household_member"}, tokenResponse.Token)))

			webhook, err := c.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, webhook, err)
		}
	})
}

func (s *TestSuite) TestHouseholds_InvitingUserWhoSignsUpIndependently() {
	s.runTest("should be possible to invite a userClient before they sign up", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			currentStatus, statusErr := testClients.userClient.GetAuthStatus(s.ctx)
			requireNotNilAndNoProblems(t, currentStatus, statusErr)
			relevantHouseholdID := currentStatus.ActiveHousehold

			// Create webhook.
			exampleWebhook := fakes.BuildFakeWebhook()
			exampleWebhookInput := converters.ConvertWebhookToWebhookCreationRequestInput(exampleWebhook)
			createdWebhook, err := testClients.userClient.CreateWebhook(ctx, exampleWebhookInput)
			require.NoError(t, err)

			checkWebhookEquality(t, exampleWebhook, createdWebhook)

			createdWebhook, err = testClients.userClient.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, createdWebhook, err)
			require.Equal(t, relevantHouseholdID, createdWebhook.BelongsToHousehold)

			inviteReq := &types.HouseholdInvitationCreationRequestInput{
				Note:    t.Name(),
				ToEmail: gofakeit.Email(),
			}
			invitation, err := testClients.userClient.CreateHouseholdInvitation(ctx, relevantHouseholdID, inviteReq)
			require.NoError(t, err)

			sentInvitations, err := testClients.userClient.GetSentHouseholdInvitations(ctx, nil)
			requireNotNilAndNoProblems(t, sentInvitations, err)
			assert.NotEmpty(t, sentInvitations.Data)

			u, c := createUserAndClientForTest(ctx, t, &types.UserRegistrationInput{
				EmailAddress: inviteReq.ToEmail,
				Username:     fakes.BuildFakeUser().Username,
				Password:     gofakeit.Password(true, true, true, true, false, 64),
			})

			invitations, err := c.GetReceivedHouseholdInvitations(ctx, nil)
			requireNotNilAndNoProblems(t, invitations, err)
			assert.NotEmpty(t, invitations.Data)

			err = c.AcceptHouseholdInvitation(ctx, invitation.ID, &types.HouseholdInvitationUpdateRequestInput{
				Token: invitation.Token,
				Note:  t.Name(),
			})
			require.NoError(t, err)

			households, err := c.GetHouseholds(ctx, nil)

			var found bool
			for _, household := range households.Data {
				if !found {
					found = household.ID == relevantHouseholdID
				}
			}

			require.True(t, found)
			_, err = c.SetDefaultHousehold(ctx, relevantHouseholdID)
			require.NoError(t, err)

			tokenResponse, err := c.LoginForJWT(ctx, &types.UserLoginInput{Username: u.Username, Password: u.HashedPassword, TOTPToken: generateTOTPTokenForUser(t, u)})
			require.NoError(t, err)

			require.NoError(t, c.SetOptions(apiclient.UsingOAuth2(ctx, createdClientID, createdClientSecret, []string{"household_member"}, tokenResponse.Token)))

			_, err = c.GetWebhook(ctx, createdWebhook.ID)
			require.NoError(t, err)

			sentInvitations, err = testClients.userClient.GetSentHouseholdInvitations(ctx, nil)
			requireNotNilAndNoProblems(t, sentInvitations, err)
			assert.Empty(t, sentInvitations.Data)
		}
	})
}

func (s *TestSuite) TestHouseholds_InvitingUserWhoSignsUpIndependentlyAndThenCancelling() {
	s.runTest("should be possible to invite a userClient before they sign up and cancel before they can accept", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			currentStatus, statusErr := testClients.userClient.GetAuthStatus(s.ctx)
			requireNotNilAndNoProblems(t, currentStatus, statusErr)
			relevantHouseholdID := currentStatus.ActiveHousehold

			// Create webhook.
			exampleWebhook := fakes.BuildFakeWebhook()
			exampleWebhookInput := converters.ConvertWebhookToWebhookCreationRequestInput(exampleWebhook)
			createdWebhook, err := testClients.userClient.CreateWebhook(ctx, exampleWebhookInput)
			require.NoError(t, err)

			checkWebhookEquality(t, exampleWebhook, createdWebhook)

			createdWebhook, err = testClients.userClient.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, createdWebhook, err)
			require.Equal(t, relevantHouseholdID, createdWebhook.BelongsToHousehold)

			inviteReq := &types.HouseholdInvitationCreationRequestInput{
				Note:    t.Name(),
				ToEmail: gofakeit.Email(),
			}
			invitation, err := testClients.userClient.CreateHouseholdInvitation(ctx, relevantHouseholdID, inviteReq)
			require.NoError(t, err)

			sentInvitations, err := testClients.userClient.GetSentHouseholdInvitations(ctx, nil)
			requireNotNilAndNoProblems(t, sentInvitations, err)
			assert.NotEmpty(t, sentInvitations.Data)

			_, c := createUserAndClientForTest(ctx, t, &types.UserRegistrationInput{
				EmailAddress: inviteReq.ToEmail,
				Username:     fakes.BuildFakeUser().Username,
				Password:     gofakeit.Password(true, true, true, true, false, 64),
			})

			invitations, err := c.GetReceivedHouseholdInvitations(ctx, nil)
			requireNotNilAndNoProblems(t, invitations, err)
			assert.NotEmpty(t, invitations.Data)

			err = testClients.userClient.CancelHouseholdInvitation(ctx, invitation.ID, &types.HouseholdInvitationUpdateRequestInput{
				Token: invitation.Token,
				Note:  t.Name(),
			})
			require.NoError(t, err)

			sentInvitations, err = testClients.userClient.GetSentHouseholdInvitations(ctx, nil)
			requireNotNilAndNoProblems(t, sentInvitations, err)
			assert.Empty(t, sentInvitations.Data)
		}
	})
}

func (s *TestSuite) TestHouseholds_InvitingNewUserWithInviteLink() {
	s.runTest("should be possible to invite a userClient via referral link", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			currentStatus, statusErr := testClients.userClient.GetAuthStatus(s.ctx)
			requireNotNilAndNoProblems(t, currentStatus, statusErr)
			relevantHouseholdID := currentStatus.ActiveHousehold

			// Create webhook.
			exampleWebhook := fakes.BuildFakeWebhook()
			exampleWebhookInput := converters.ConvertWebhookToWebhookCreationRequestInput(exampleWebhook)
			createdWebhook, err := testClients.userClient.CreateWebhook(ctx, exampleWebhookInput)
			require.NoError(t, err)

			checkWebhookEquality(t, exampleWebhook, createdWebhook)

			createdWebhook, err = testClients.userClient.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, createdWebhook, err)
			require.Equal(t, relevantHouseholdID, createdWebhook.BelongsToHousehold)

			inviteReq := &types.HouseholdInvitationCreationRequestInput{
				Note:    t.Name(),
				ToEmail: gofakeit.Email(),
			}
			createdInvitation, err := testClients.userClient.CreateHouseholdInvitation(ctx, relevantHouseholdID, inviteReq)
			require.NoError(t, err)

			createdInvitation, err = testClients.userClient.GetHouseholdInvitation(ctx, createdInvitation.ID)
			requireNotNilAndNoProblems(t, createdInvitation, err)

			sentInvitations, err := testClients.userClient.GetSentHouseholdInvitations(ctx, nil)
			requireNotNilAndNoProblems(t, sentInvitations, err)
			assert.NotEmpty(t, sentInvitations.Data)

			_, c := createUserAndClientForTest(ctx, t, &types.UserRegistrationInput{
				EmailAddress:    inviteReq.ToEmail,
				Username:        fakes.BuildFakeUser().Username,
				Password:        gofakeit.Password(true, true, true, true, false, 64),
				InvitationID:    createdInvitation.ID,
				InvitationToken: createdInvitation.Token,
			})

			households, err := c.GetHouseholds(ctx, nil)
			require.NoError(t, err)

			var found bool
			for _, household := range households.Data {
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
	s.runTest("should be possible to invite an already-registered userClient", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			currentStatus, statusErr := testClients.userClient.GetAuthStatus(s.ctx)
			requireNotNilAndNoProblems(t, currentStatus, statusErr)
			relevantHouseholdID := currentStatus.ActiveHousehold

			// Create webhook.
			exampleWebhook := fakes.BuildFakeWebhook()
			exampleWebhookInput := converters.ConvertWebhookToWebhookCreationRequestInput(exampleWebhook)
			createdWebhook, err := testClients.userClient.CreateWebhook(ctx, exampleWebhookInput)
			require.NoError(t, err)

			checkWebhookEquality(t, exampleWebhook, createdWebhook)

			createdWebhook, err = testClients.userClient.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, createdWebhook, err)
			require.Equal(t, relevantHouseholdID, createdWebhook.BelongsToHousehold)

			inviteReq := &types.HouseholdInvitationCreationRequestInput{
				Note:    t.Name(),
				ToEmail: gofakeit.Email(),
			}
			invitation, err := testClients.userClient.CreateHouseholdInvitation(ctx, relevantHouseholdID, inviteReq)
			require.NoError(t, err)

			require.NoError(t, testClients.userClient.CancelHouseholdInvitation(ctx, invitation.ID, &types.HouseholdInvitationUpdateRequestInput{
				Token: invitation.Token,
				Note:  t.Name(),
			}))

			sentInvitations, err := testClients.userClient.GetSentHouseholdInvitations(ctx, nil)
			requireNotNilAndNoProblems(t, sentInvitations, err)
			assert.Empty(t, sentInvitations.Data)

			_, c := createUserAndClientForTest(ctx, t, &types.UserRegistrationInput{
				EmailAddress: inviteReq.ToEmail,
				Username:     fakes.BuildFakeUser().Username,
				Password:     gofakeit.Password(true, true, true, true, false, 64),
			})

			invitations, err := c.GetReceivedHouseholdInvitations(ctx, nil)
			requireNotNilAndNoProblems(t, invitations, err)
			assert.Empty(t, invitations.Data)
		}
	})
}

func (s *TestSuite) TestHouseholds_InviteCanBeRejected() {
	s.runTest("should be possible to invite an already-registered userClient", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			currentStatus, statusErr := testClients.userClient.GetAuthStatus(s.ctx)
			requireNotNilAndNoProblems(t, currentStatus, statusErr)
			relevantHouseholdID := currentStatus.ActiveHousehold

			// Create webhook.
			exampleWebhook := fakes.BuildFakeWebhook()
			exampleWebhookInput := converters.ConvertWebhookToWebhookCreationRequestInput(exampleWebhook)
			createdWebhook, err := testClients.userClient.CreateWebhook(ctx, exampleWebhookInput)
			require.NoError(t, err)

			checkWebhookEquality(t, exampleWebhook, createdWebhook)

			createdWebhook, err = testClients.userClient.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, createdWebhook, err)
			require.Equal(t, relevantHouseholdID, createdWebhook.BelongsToHousehold)

			u, c := createUserAndClientForTest(ctx, t, nil)

			invitation, err := testClients.userClient.CreateHouseholdInvitation(ctx, relevantHouseholdID, &types.HouseholdInvitationCreationRequestInput{
				Note:    t.Name(),
				ToEmail: u.EmailAddress,
			})
			require.NoError(t, err)

			sentInvitations, err := testClients.userClient.GetSentHouseholdInvitations(ctx, nil)
			requireNotNilAndNoProblems(t, sentInvitations, err)
			assert.NotEmpty(t, sentInvitations.Data)

			invitations, err := c.GetReceivedHouseholdInvitations(ctx, nil)
			requireNotNilAndNoProblems(t, invitations, err)
			assert.NotEmpty(t, invitations.Data)

			err = c.RejectHouseholdInvitation(ctx, invitation.ID, &types.HouseholdInvitationUpdateRequestInput{
				Token: invitation.Token,
				Note:  t.Name(),
			})
			require.NoError(t, err)

			sentInvitations, err = testClients.userClient.GetSentHouseholdInvitations(ctx, nil)
			requireNotNilAndNoProblems(t, sentInvitations, err)
			assert.Empty(t, sentInvitations.Data)
		}
	})
}

func (s *TestSuite) TestHouseholds_ChangingMemberships() {
	s.runTest("should be possible to change members of a household", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			const userCount = 1

			currentStatus, statusErr := testClients.userClient.GetAuthStatus(s.ctx)
			requireNotNilAndNoProblems(t, currentStatus, statusErr)

			// fetch household data
			householdCreationInput := &types.HouseholdCreationRequestInput{
				Name: fakes.BuildFakeHousehold().Name,
			}
			household, householdCreationErr := testClients.userClient.CreateHousehold(ctx, householdCreationInput)
			require.NoError(t, householdCreationErr)
			require.NotNil(t, household)

			_, err := testClients.userClient.SetDefaultHousehold(ctx, household.ID)
			require.NoError(t, err)

			tokenResponse, err := testClients.userClient.LoginForJWT(ctx, &types.UserLoginInput{Username: testClients.user.Username, Password: testClients.user.HashedPassword, TOTPToken: generateTOTPTokenForUser(t, testClients.user)})
			require.NoError(t, err)

			require.NoError(t, testClients.userClient.SetOptions(apiclient.UsingOAuth2(ctx, createdClientID, createdClientSecret, []string{"household_member"}, tokenResponse.Token)))

			// Create webhook.
			exampleWebhook := fakes.BuildFakeWebhook()
			exampleWebhookInput := converters.ConvertWebhookToWebhookCreationRequestInput(exampleWebhook)
			createdWebhook, err := testClients.userClient.CreateWebhook(ctx, exampleWebhookInput)
			require.NoError(t, err)

			checkWebhookEquality(t, exampleWebhook, createdWebhook)

			createdWebhook, err = testClients.userClient.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, createdWebhook, err)
			require.Equal(t, household.ID, createdWebhook.BelongsToHousehold)

			// create dummy users
			users := []*types.User{}
			clients := []*apiclient.Client{}

			// create users
			for i := 0; i < userCount; i++ {
				u, c := createUserAndClientForTest(ctx, t, nil)
				users = append(users, u)
				clients = append(clients, c)

				currentStatus, statusErr = c.GetAuthStatus(s.ctx)
				requireNotNilAndNoProblems(t, currentStatus, statusErr)
			}

			// check that each userClient cannot see the unreachable webhook
			for i := 0; i < userCount; i++ {
				webhook, err := clients[i].GetWebhook(ctx, createdWebhook.ID)
				require.Nil(t, webhook)
				require.Error(t, err)
			}

			// add them to the household
			for i := 0; i < userCount; i++ {
				invitation, invitationErr := testClients.userClient.CreateHouseholdInvitation(ctx, household.ID, &types.HouseholdInvitationCreationRequestInput{
					ToEmail: users[i].EmailAddress,
					Note:    t.Name(),
				})
				require.NoError(t, invitationErr)
				require.NotEmpty(t, invitation.ID)

				invitations, fetchInvitationsErr := clients[i].GetReceivedHouseholdInvitations(ctx, nil)
				requireNotNilAndNoProblems(t, invitations, fetchInvitationsErr)
				assert.NotEmpty(t, invitations.Data)

				err = clients[i].AcceptHouseholdInvitation(ctx, invitation.ID, &types.HouseholdInvitationUpdateRequestInput{
					Token: invitation.Token,
					Note:  t.Name(),
				})
				require.NoError(t, err)

				_, err = clients[i].SetDefaultHousehold(ctx, household.ID)
				require.NoError(t, err)

				tokenResponse, err = clients[i].LoginForJWT(ctx, &types.UserLoginInput{Username: users[i].Username, Password: users[i].HashedPassword, TOTPToken: generateTOTPTokenForUser(t, users[i])})
				require.NoError(t, err)

				require.NoError(t, clients[i].SetOptions(apiclient.UsingOAuth2(ctx, createdClientID, createdClientSecret, []string{"household_member"}, tokenResponse.Token)))

				currentStatus, statusErr = clients[i].GetAuthStatus(s.ctx)
				requireNotNilAndNoProblems(t, currentStatus, statusErr)
				require.Equal(t, currentStatus.ActiveHousehold, household.ID)
			}

			// grant all permissions
			for i := 0; i < userCount; i++ {
				input := &types.ModifyUserPermissionsInput{
					Reason:  t.Name(),
					NewRole: authorization.HouseholdAdminRole.String(),
				}
				require.NoError(t, testClients.userClient.UpdateHouseholdMemberPermissions(ctx, household.ID, users[i].ID, input))
			}

			// check that each userClient can see the webhook
			for i := 0; i < userCount; i++ {
				webhook, webhookRetrievalError := clients[i].GetWebhook(ctx, createdWebhook.ID)
				requireNotNilAndNoProblems(t, webhook, webhookRetrievalError)
			}

			// remove users from household
			for i := 0; i < userCount; i++ {
				require.NoError(t, testClients.userClient.ArchiveUserMembership(ctx, household.ID, users[i].ID))
			}

			// check that each userClient cannot see the webhook
			for i := 0; i < userCount; i++ {
				webhook, webhookRetrievalError := clients[i].GetWebhook(ctx, createdWebhook.ID)
				require.Nil(t, webhook)
				require.Error(t, webhookRetrievalError)
			}

			// Clean up.
			require.NoError(t, testClients.userClient.ArchiveWebhook(ctx, createdWebhook.ID))

			for i := 0; i < userCount; i++ {
				require.NoError(t, testClients.adminClient.ArchiveUser(ctx, users[i].ID))
			}
		}
	})
}

func (s *TestSuite) TestHouseholds_OwnershipTransfer() {
	s.runTest("should be possible to transfer ownership of a household", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create users
			futureOwner, futureOwnerClient := createUserAndClientForTest(ctx, t, nil)

			// fetch household data
			householdCreationInput := &types.HouseholdCreationRequestInput{
				Name: fakes.BuildFakeHousehold().Name,
			}
			household, householdCreationErr := testClients.userClient.CreateHousehold(ctx, householdCreationInput)
			require.NoError(t, householdCreationErr)
			require.NotNil(t, household)

			_, err := testClients.userClient.SetDefaultHousehold(ctx, household.ID)
			require.NoError(t, err)

			tokenResponse, err := testClients.userClient.LoginForJWT(ctx, &types.UserLoginInput{Username: testClients.user.Username, Password: testClients.user.HashedPassword, TOTPToken: generateTOTPTokenForUser(t, testClients.user)})
			require.NoError(t, err)

			require.NoError(t, testClients.userClient.SetOptions(apiclient.UsingOAuth2(ctx, createdClientID, createdClientSecret, []string{"household_member"}, tokenResponse.Token)))

			// create a webhook

			// Create webhook.
			exampleWebhook := fakes.BuildFakeWebhook()
			exampleWebhookInput := converters.ConvertWebhookToWebhookCreationRequestInput(exampleWebhook)
			createdWebhook, err := testClients.userClient.CreateWebhook(ctx, exampleWebhookInput)
			require.NoError(t, err)

			checkWebhookEquality(t, exampleWebhook, createdWebhook)

			createdWebhook, err = testClients.userClient.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, createdWebhook, err)

			require.Equal(t, household.ID, createdWebhook.BelongsToHousehold)

			// check that userClient cannot see the webhook
			webhook, err := futureOwnerClient.GetWebhook(ctx, createdWebhook.ID)
			require.Nil(t, webhook)
			require.Error(t, err)

			// add them to the household
			_, err = testClients.userClient.TransferHouseholdOwnership(ctx, household.ID, &types.HouseholdOwnershipTransferInput{
				Reason:       t.Name(),
				CurrentOwner: household.BelongsToUser,
				NewOwner:     futureOwner.ID,
			})
			require.NoError(t, err)

			_, err = futureOwnerClient.SetDefaultHousehold(ctx, household.ID)
			require.NoError(t, err)

			tokenResponse, err = futureOwnerClient.LoginForJWT(ctx, &types.UserLoginInput{Username: futureOwner.Username, Password: futureOwner.HashedPassword, TOTPToken: generateTOTPTokenForUser(t, futureOwner)})
			require.NoError(t, err)

			require.NoError(t, futureOwnerClient.SetOptions(apiclient.UsingOAuth2(ctx, createdClientID, createdClientSecret, []string{"household_member"}, tokenResponse.Token)))

			// check that userClient can see the webhook
			webhook, err = futureOwnerClient.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, webhook, err)

			// check that old userClient cannot see the webhook
			webhook, err = testClients.userClient.GetWebhook(ctx, createdWebhook.ID)
			require.Nil(t, webhook)
			require.Error(t, err)

			// check that new owner can delete the webhook
			require.NoError(t, futureOwnerClient.ArchiveWebhook(ctx, createdWebhook.ID))

			// Clean up.
			require.Error(t, testClients.userClient.ArchiveWebhook(ctx, createdWebhook.ID))
			require.NoError(t, testClients.adminClient.ArchiveUser(ctx, futureOwner.ID))
		}
	})
}

func (s *TestSuite) TestHouseholds_UsersHaveBackupHouseholdCreatedForThemWhenRemovedFromLastHousehold() {
	s.runTest("should be possible to invite a userClient via referral link", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			currentStatus, statusErr := testClients.userClient.GetAuthStatus(s.ctx)
			requireNotNilAndNoProblems(t, currentStatus, statusErr)
			relevantHouseholdID := currentStatus.ActiveHousehold

			inviteReq := &types.HouseholdInvitationCreationRequestInput{
				Note:    t.Name(),
				ToEmail: gofakeit.Email(),
			}
			createdInvitation, err := testClients.userClient.CreateHouseholdInvitation(ctx, relevantHouseholdID, inviteReq)
			require.NoError(t, err)

			createdInvitation, err = testClients.userClient.GetHouseholdInvitation(ctx, createdInvitation.ID)
			requireNotNilAndNoProblems(t, createdInvitation, err)

			sentInvitations, err := testClients.userClient.GetSentHouseholdInvitations(ctx, nil)
			requireNotNilAndNoProblems(t, sentInvitations, err)
			assert.NotEmpty(t, sentInvitations.Data)

			regInput := &types.UserRegistrationInput{
				EmailAddress:    inviteReq.ToEmail,
				Username:        fakes.BuildFakeUser().Username,
				Password:        gofakeit.Password(true, true, true, true, false, 64),
				InvitationID:    createdInvitation.ID,
				InvitationToken: createdInvitation.Token,
			}
			u, c := createUserAndClientForTest(ctx, t, regInput)

			households, err := c.GetHouseholds(ctx, nil)
			require.NoError(t, err)

			assert.Len(t, households.Data, 2)

			var (
				found            bool
				otherHouseholdID string
			)

			for _, household := range households.Data {
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

			require.NoError(t, testClients.userClient.ArchiveUserMembership(ctx, relevantHouseholdID, u.ID))

			u.HashedPassword = regInput.Password

			tokenResponse, err := c.LoginForJWT(ctx, &types.UserLoginInput{Username: u.Username, Password: u.HashedPassword, TOTPToken: generateTOTPTokenForUser(t, u)})
			require.NoError(t, err)

			require.NoError(t, c.SetOptions(apiclient.UsingOAuth2(ctx, createdClientID, createdClientSecret, []string{"household_member"}, tokenResponse.Token)))

			household, err := c.GetActiveHousehold(ctx)
			requireNotNilAndNoProblems(t, household, err)
			assert.NotEqual(t, relevantHouseholdID, household.ID)

			require.True(t, found)
		}
	})
}
