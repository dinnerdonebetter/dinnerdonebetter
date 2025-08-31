package integration

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/domain/identity/fakes"
	identitysvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/identity"
	"github.com/dinnerdonebetter/backend/internal/services/identity/grpc/converters"

	"github.com/stretchr/testify/assert"
)

func TestAccounts_Creating(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, testClient := createUserAndClientForTest(t)

		exampleAccount := fakes.BuildFakeAccountCreationRequestInput()
		exampleAccountInput := converters.ConvertAccountCreationRequestInputToGRPCAccountCreationRequestInput(exampleAccount)

		createdAccount, err := testClient.CreateAccount(ctx, &identitysvc.CreateAccountRequest{Input: exampleAccountInput})
		assert.NoError(t, err)
		assert.NotNil(t, createdAccount)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()
	})
}

func TestAccounts_Listing(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
	})
}

func TestAccounts_Reading(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
	})

	T.Run("for nonexistent account", func(t *testing.T) {
		t.Parallel()
	})
}

func TestAccounts_Updating(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
	})

	T.Run("for nonexistent account", func(t *testing.T) {
		t.Parallel()
	})
}

func TestAccounts_Archiving(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
	})

	T.Run("for nonexistent account", func(t *testing.T) {
		t.Parallel()
	})
}

func TestAccounts_InvitingPreExistentUser(T *testing.T) {
	T.Parallel()

	T.Run("pre existing user", func(t *testing.T) {
		t.Parallel()
	})

	T.Run("user who signed up independently", func(t *testing.T) {
		t.Parallel()
	})

	T.Run("user who signs up independently and then cancels", func(t *testing.T) {
		t.Parallel()
	})

	T.Run("user with invite link", func(t *testing.T) {
		t.Parallel()
	})

	T.Run("invites can be cancelled", func(t *testing.T) {
		t.Parallel()
	})

	T.Run("invites can be rejected", func(t *testing.T) {
		t.Parallel()
	})
}

func TestAccounts_ChangingMemberships(T *testing.T) {
	T.Parallel()

	T.Run("", func(t *testing.T) {
		t.Parallel()
	})
}

func TestAccounts_OwnershipTransfer(T *testing.T) {
	T.Parallel()

	T.Run("", func(t *testing.T) {
		t.Parallel()
	})
}

func TestAccounts_UsersHaveBackupAccountCreatedForThemWhenRemovedFromLastAccount(T *testing.T) {
	T.Parallel()

	T.Run("", func(t *testing.T) {
		t.Parallel()
	})
}

/*

func (s *TestSuite) TestAccounts_InvitingPreExistentUser() {
	s.runTest("should be possible to invite an already-registered userClient", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			currentStatus, statusErr := testClients.userClient.GetAuthStatus(s.ctx)
			requireNotNilAndNoProblems(t, currentStatus, statusErr)
			relevantAccountID := currentStatus.ActiveAccount

			// Create webhook.
			exampleWebhook := fakes.BuildFakeWebhook()
			exampleWebhookInput := converters.ConvertWebhookToWebhookCreationRequestInput(exampleWebhook)
			createdWebhook, err := testClients.userClient.CreateWebhook(ctx, exampleWebhookInput)
			require.NoError(t, err)

			checkWebhookEquality(t, exampleWebhook, createdWebhook)

			createdWebhook, err = testClients.userClient.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, createdWebhook, err)
			require.Equal(t, relevantAccountID, createdWebhook.BelongsToAccount)

			u, c := createUserAndClientForTest(ctx, t, nil)

			invitation, err := testClients.userClient.CreateAccountInvitation(ctx, relevantAccountID, &types.AccountInvitationCreationRequestInput{
				Note:    t.Name(),
				ToName:  t.Name(),
				ToEmail: u.EmailAddress,
			})
			require.NoError(t, err)

			sentInvitations, err := testClients.userClient.GetSentAccountInvitations(ctx, nil)
			requireNotNilAndNoProblems(t, sentInvitations, err)
			assert.NotEmpty(t, sentInvitations.Data)

			invitations, err := c.GetReceivedAccountInvitations(ctx, nil)
			requireNotNilAndNoProblems(t, invitations, err)
			assert.NotEmpty(t, invitations.Data)

			err = c.AcceptAccountInvitation(ctx, invitation.ID, &types.AccountInvitationUpdateRequestInput{
				Token: invitation.Token,
				Note:  t.Name(),
			})
			require.NoError(t, err)

			sentInvitations, err = testClients.userClient.GetSentAccountInvitations(ctx, nil)
			requireNotNilAndNoProblems(t, sentInvitations, err)
			assert.Empty(t, sentInvitations.Data)

			accounts, err := c.GetAccounts(ctx, nil)

			var found bool
			for _, account := range accounts.Data {
				if !found {
					found = account.ID == relevantAccountID
				}
			}

			require.True(t, found)
			_, err = c.SetDefaultAccount(ctx, relevantAccountID)
			require.NoError(t, err)

			tokenResponse, err := c.LoginForToken(ctx, &types.UserLoginInput{Username: u.Username, Password: u.HashedPassword, TOTPToken: generateTOTPTokenForUser(t, u)})
			require.NoError(t, err)

			require.NoError(t, c.SetOptions(apiclient.UsingOAuth2(ctx, createdClientID, createdClientSecret, []string{"account_member"}, tokenResponse.Token)))

			webhook, err := c.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, webhook, err)
		}
	})
}

func (s *TestSuite) TestAccounts_InvitingUserWhoSignsUpIndependently() {
	s.runTest("should be possible to invite a user before they sign up", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			currentStatus, statusErr := testClients.userClient.GetAuthStatus(s.ctx)
			requireNotNilAndNoProblems(t, currentStatus, statusErr)
			relevantAccountID := currentStatus.ActiveAccount

			// Create webhook.
			exampleWebhook := fakes.BuildFakeWebhook()
			exampleWebhookInput := converters.ConvertWebhookToWebhookCreationRequestInput(exampleWebhook)
			createdWebhook, err := testClients.userClient.CreateWebhook(ctx, exampleWebhookInput)
			require.NoError(t, err)

			checkWebhookEquality(t, exampleWebhook, createdWebhook)

			createdWebhook, err = testClients.userClient.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, createdWebhook, err)
			require.Equal(t, relevantAccountID, createdWebhook.BelongsToAccount)

			inviteReq := &types.AccountInvitationCreationRequestInput{
				Note:    t.Name(),
				ToEmail: gofakeit.Email(),
			}
			invitation, err := testClients.userClient.CreateAccountInvitation(ctx, relevantAccountID, inviteReq)
			require.NoError(t, err)

			sentInvitations, err := testClients.userClient.GetSentAccountInvitations(ctx, nil)
			requireNotNilAndNoProblems(t, sentInvitations, err)
			assert.NotEmpty(t, sentInvitations.Data)

			u, c := createUserAndClientForTest(ctx, t, &types.UserRegistrationInput{
				EmailAddress: inviteReq.ToEmail,
				Username:     fakes.BuildFakeUser().Username,
				Password:     gofakeit.Password(true, true, true, true, false, 64),
			})

			invitations, err := c.GetReceivedAccountInvitations(ctx, nil)
			requireNotNilAndNoProblems(t, invitations, err)
			assert.NotEmpty(t, invitations.Data)

			err = c.AcceptAccountInvitation(ctx, invitation.ID, &types.AccountInvitationUpdateRequestInput{
				Token: invitation.Token,
				Note:  t.Name(),
			})
			require.NoError(t, err)

			accounts, err := c.GetAccounts(ctx, nil)

			var found bool
			for _, account := range accounts.Data {
				if !found {
					found = account.ID == relevantAccountID
				}
			}

			require.True(t, found)
			_, err = c.SetDefaultAccount(ctx, relevantAccountID)
			require.NoError(t, err)

			tokenResponse, err := c.LoginForToken(ctx, &types.UserLoginInput{Username: u.Username, Password: u.HashedPassword, TOTPToken: generateTOTPTokenForUser(t, u)})
			require.NoError(t, err)

			require.NoError(t, c.SetOptions(apiclient.UsingOAuth2(ctx, createdClientID, createdClientSecret, []string{"account_member"}, tokenResponse.Token)))

			_, err = c.GetWebhook(ctx, createdWebhook.ID)
			require.NoError(t, err)

			sentInvitations, err = testClients.userClient.GetSentAccountInvitations(ctx, nil)
			requireNotNilAndNoProblems(t, sentInvitations, err)
			assert.Empty(t, sentInvitations.Data)
		}
	})
}

func (s *TestSuite) TestAccounts_InvitingUserWhoSignsUpIndependentlyAndThenCancelling() {
	s.runTest("should be possible to invite a user before they sign up and cancel before they can accept", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			currentStatus, statusErr := testClients.userClient.GetAuthStatus(s.ctx)
			requireNotNilAndNoProblems(t, currentStatus, statusErr)
			relevantAccountID := currentStatus.ActiveAccount

			// Create webhook.
			exampleWebhook := fakes.BuildFakeWebhook()
			exampleWebhookInput := converters.ConvertWebhookToWebhookCreationRequestInput(exampleWebhook)
			createdWebhook, err := testClients.userClient.CreateWebhook(ctx, exampleWebhookInput)
			require.NoError(t, err)

			checkWebhookEquality(t, exampleWebhook, createdWebhook)

			createdWebhook, err = testClients.userClient.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, createdWebhook, err)
			require.Equal(t, relevantAccountID, createdWebhook.BelongsToAccount)

			inviteReq := &types.AccountInvitationCreationRequestInput{
				Note:    t.Name(),
				ToEmail: gofakeit.Email(),
			}
			invitation, err := testClients.userClient.CreateAccountInvitation(ctx, relevantAccountID, inviteReq)
			require.NoError(t, err)

			sentInvitations, err := testClients.userClient.GetSentAccountInvitations(ctx, nil)
			requireNotNilAndNoProblems(t, sentInvitations, err)
			assert.NotEmpty(t, sentInvitations.Data)

			_, c := createUserAndClientForTest(ctx, t, &types.UserRegistrationInput{
				EmailAddress: inviteReq.ToEmail,
				Username:     fakes.BuildFakeUser().Username,
				Password:     gofakeit.Password(true, true, true, true, false, 64),
			})

			invitations, err := c.GetReceivedAccountInvitations(ctx, nil)
			requireNotNilAndNoProblems(t, invitations, err)
			assert.NotEmpty(t, invitations.Data)

			err = testClients.userClient.CancelAccountInvitation(ctx, invitation.ID, &types.AccountInvitationUpdateRequestInput{
				Token: invitation.Token,
				Note:  t.Name(),
			})
			require.NoError(t, err)

			sentInvitations, err = testClients.userClient.GetSentAccountInvitations(ctx, nil)
			requireNotNilAndNoProblems(t, sentInvitations, err)
			assert.Empty(t, sentInvitations.Data)
		}
	})
}

func (s *TestSuite) TestAccounts_InvitingNewUserWithInviteLink() {
	s.runTest("should be possible to invite a user via referral link", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			currentStatus, statusErr := testClients.userClient.GetAuthStatus(s.ctx)
			requireNotNilAndNoProblems(t, currentStatus, statusErr)
			relevantAccountID := currentStatus.ActiveAccount

			// Create webhook.
			exampleWebhook := fakes.BuildFakeWebhook()
			exampleWebhookInput := converters.ConvertWebhookToWebhookCreationRequestInput(exampleWebhook)
			createdWebhook, err := testClients.userClient.CreateWebhook(ctx, exampleWebhookInput)
			require.NoError(t, err)

			checkWebhookEquality(t, exampleWebhook, createdWebhook)

			createdWebhook, err = testClients.userClient.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, createdWebhook, err)
			require.Equal(t, relevantAccountID, createdWebhook.BelongsToAccount)

			inviteReq := &types.AccountInvitationCreationRequestInput{
				Note:    t.Name(),
				ToEmail: gofakeit.Email(),
			}
			createdInvitation, err := testClients.userClient.CreateAccountInvitation(ctx, relevantAccountID, inviteReq)
			require.NoError(t, err)

			createdInvitation, err = testClients.userClient.GetAccountInvitation(ctx, createdInvitation.ID)
			requireNotNilAndNoProblems(t, createdInvitation, err)

			sentInvitations, err := testClients.userClient.GetSentAccountInvitations(ctx, nil)
			requireNotNilAndNoProblems(t, sentInvitations, err)
			assert.NotEmpty(t, sentInvitations.Data)

			_, c := createUserAndClientForTest(ctx, t, &types.UserRegistrationInput{
				EmailAddress:    inviteReq.ToEmail,
				Username:        fakes.BuildFakeUser().Username,
				Password:        gofakeit.Password(true, true, true, true, false, 64),
				InvitationID:    createdInvitation.ID,
				InvitationToken: createdInvitation.Token,
			})

			accounts, err := c.GetAccounts(ctx, nil)
			require.NoError(t, err)

			var found bool
			for _, account := range accounts.Data {
				if account.ID == relevantAccountID {
					found = true
					break
				}
			}

			require.True(t, found)

			_, err = c.GetWebhook(ctx, createdWebhook.ID)
			require.NoError(t, err)
		}
	})
}

func (s *TestSuite) TestAccounts_InviteCanBeCancelled() {
	s.runTest("should be possible to invite an already-registered userClient", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			currentStatus, statusErr := testClients.userClient.GetAuthStatus(s.ctx)
			requireNotNilAndNoProblems(t, currentStatus, statusErr)
			relevantAccountID := currentStatus.ActiveAccount

			// Create webhook.
			exampleWebhook := fakes.BuildFakeWebhook()
			exampleWebhookInput := converters.ConvertWebhookToWebhookCreationRequestInput(exampleWebhook)
			createdWebhook, err := testClients.userClient.CreateWebhook(ctx, exampleWebhookInput)
			require.NoError(t, err)

			checkWebhookEquality(t, exampleWebhook, createdWebhook)

			createdWebhook, err = testClients.userClient.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, createdWebhook, err)
			require.Equal(t, relevantAccountID, createdWebhook.BelongsToAccount)

			inviteReq := &types.AccountInvitationCreationRequestInput{
				Note:    t.Name(),
				ToEmail: gofakeit.Email(),
			}
			invitation, err := testClients.userClient.CreateAccountInvitation(ctx, relevantAccountID, inviteReq)
			require.NoError(t, err)

			require.NoError(t, testClients.userClient.CancelAccountInvitation(ctx, invitation.ID, &types.AccountInvitationUpdateRequestInput{
				Token: invitation.Token,
				Note:  t.Name(),
			}))

			sentInvitations, err := testClients.userClient.GetSentAccountInvitations(ctx, nil)
			requireNotNilAndNoProblems(t, sentInvitations, err)
			assert.Empty(t, sentInvitations.Data)

			_, c := createUserAndClientForTest(ctx, t, &types.UserRegistrationInput{
				EmailAddress: inviteReq.ToEmail,
				Username:     fakes.BuildFakeUser().Username,
				Password:     gofakeit.Password(true, true, true, true, false, 64),
			})

			invitations, err := c.GetReceivedAccountInvitations(ctx, nil)
			requireNotNilAndNoProblems(t, invitations, err)
			assert.Empty(t, invitations.Data)
		}
	})
}

func (s *TestSuite) TestAccounts_InviteCanBeRejected() {
	s.runTest("should be possible to invite an already-registered userClient", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			currentStatus, statusErr := testClients.userClient.GetAuthStatus(s.ctx)
			requireNotNilAndNoProblems(t, currentStatus, statusErr)
			relevantAccountID := currentStatus.ActiveAccount

			// Create webhook.
			exampleWebhook := fakes.BuildFakeWebhook()
			exampleWebhookInput := converters.ConvertWebhookToWebhookCreationRequestInput(exampleWebhook)
			createdWebhook, err := testClients.userClient.CreateWebhook(ctx, exampleWebhookInput)
			require.NoError(t, err)

			checkWebhookEquality(t, exampleWebhook, createdWebhook)

			createdWebhook, err = testClients.userClient.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, createdWebhook, err)
			require.Equal(t, relevantAccountID, createdWebhook.BelongsToAccount)

			u, c := createUserAndClientForTest(ctx, t, nil)

			invitation, err := testClients.userClient.CreateAccountInvitation(ctx, relevantAccountID, &types.AccountInvitationCreationRequestInput{
				Note:    t.Name(),
				ToEmail: u.EmailAddress,
			})
			require.NoError(t, err)

			sentInvitations, err := testClients.userClient.GetSentAccountInvitations(ctx, nil)
			requireNotNilAndNoProblems(t, sentInvitations, err)
			assert.NotEmpty(t, sentInvitations.Data)

			invitations, err := c.GetReceivedAccountInvitations(ctx, nil)
			requireNotNilAndNoProblems(t, invitations, err)
			assert.NotEmpty(t, invitations.Data)

			err = c.RejectAccountInvitation(ctx, invitation.ID, &types.AccountInvitationUpdateRequestInput{
				Token: invitation.Token,
				Note:  t.Name(),
			})
			require.NoError(t, err)

			sentInvitations, err = testClients.userClient.GetSentAccountInvitations(ctx, nil)
			requireNotNilAndNoProblems(t, sentInvitations, err)
			assert.Empty(t, sentInvitations.Data)
		}
	})
}

func (s *TestSuite) TestAccounts_ChangingMemberships() {
	s.runTest("should be possible to change members of an account", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			const userCount = 1

			currentStatus, statusErr := testClients.userClient.GetAuthStatus(s.ctx)
			requireNotNilAndNoProblems(t, currentStatus, statusErr)

			// fetch account data
			accountCreationInput := &types.AccountCreationRequestInput{
				Name: fakes.BuildFakeAccount().Name,
			}
			account, accountCreationErr := testClients.userClient.CreateAccount(ctx, accountCreationInput)
			require.NoError(t, accountCreationErr)
			require.NotNil(t, account)

			_, err := testClients.userClient.SetDefaultAccount(ctx, account.ID)
			require.NoError(t, err)

			tokenResponse, err := testClients.userClient.LoginForToken(ctx, &types.UserLoginInput{Username: testClients.user.Username, Password: testClients.user.HashedPassword, TOTPToken: generateTOTPTokenForUser(t, testClients.user)})
			require.NoError(t, err)

			require.NoError(t, testClients.userClient.SetOptions(apiclient.UsingOAuth2(ctx, createdClientID, createdClientSecret, []string{"account_member"}, tokenResponse.Token)))

			// Create webhook.
			exampleWebhook := fakes.BuildFakeWebhook()
			exampleWebhookInput := converters.ConvertWebhookToWebhookCreationRequestInput(exampleWebhook)
			createdWebhook, err := testClients.userClient.CreateWebhook(ctx, exampleWebhookInput)
			require.NoError(t, err)

			checkWebhookEquality(t, exampleWebhook, createdWebhook)

			createdWebhook, err = testClients.userClient.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, createdWebhook, err)
			require.Equal(t, account.ID, createdWebhook.BelongsToAccount)

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

			// add them to the account
			for i := 0; i < userCount; i++ {
				invitation, invitationErr := testClients.userClient.CreateAccountInvitation(ctx, account.ID, &types.AccountInvitationCreationRequestInput{
					ToEmail: users[i].EmailAddress,
					Note:    t.Name(),
				})
				require.NoError(t, invitationErr)
				require.NotEmpty(t, invitation.ID)

				invitations, fetchInvitationsErr := clients[i].GetReceivedAccountInvitations(ctx, nil)
				requireNotNilAndNoProblems(t, invitations, fetchInvitationsErr)
				assert.NotEmpty(t, invitations.Data)

				err = clients[i].AcceptAccountInvitation(ctx, invitation.ID, &types.AccountInvitationUpdateRequestInput{
					Token: invitation.Token,
					Note:  t.Name(),
				})
				require.NoError(t, err)

				_, err = clients[i].SetDefaultAccount(ctx, account.ID)
				require.NoError(t, err)

				tokenResponse, err = clients[i].LoginForToken(ctx, &types.UserLoginInput{Username: users[i].Username, Password: users[i].HashedPassword, TOTPToken: generateTOTPTokenForUser(t, users[i])})
				require.NoError(t, err)

				require.NoError(t, clients[i].SetOptions(apiclient.UsingOAuth2(ctx, createdClientID, createdClientSecret, []string{"account_member"}, tokenResponse.Token)))

				currentStatus, statusErr = clients[i].GetAuthStatus(s.ctx)
				requireNotNilAndNoProblems(t, currentStatus, statusErr)
				require.Equal(t, currentStatus.ActiveAccount, account.ID)
			}

			// grant all permissions
			for i := 0; i < userCount; i++ {
				input := &types.ModifyUserPermissionsInput{
					Reason:  t.Name(),
					NewRole: authorization.AccountAdminRole.String(),
				}
				require.NoError(t, testClients.userClient.UpdateAccountMemberPermissions(ctx, account.ID, users[i].ID, input))
			}

			// check that each userClient can see the webhook
			for i := 0; i < userCount; i++ {
				webhook, webhookRetrievalError := clients[i].GetWebhook(ctx, createdWebhook.ID)
				requireNotNilAndNoProblems(t, webhook, webhookRetrievalError)
			}

			// remove users from account
			for i := 0; i < userCount; i++ {
				require.NoError(t, testClients.userClient.ArchiveUserMembership(ctx, account.ID, users[i].ID))
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

func (s *TestSuite) TestAccounts_OwnershipTransfer() {
	s.runTest("should be possible to transfer ownership of an account", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create users
			futureOwner, futureOwnerClient := createUserAndClientForTest(ctx, t, nil)

			// fetch account data
			accountCreationInput := &types.AccountCreationRequestInput{
				Name: fakes.BuildFakeAccount().Name,
			}
			account, accountCreationErr := testClients.userClient.CreateAccount(ctx, accountCreationInput)
			require.NoError(t, accountCreationErr)
			require.NotNil(t, account)

			_, err := testClients.userClient.SetDefaultAccount(ctx, account.ID)
			require.NoError(t, err)

			tokenResponse, err := testClients.userClient.LoginForToken(ctx, &types.UserLoginInput{Username: testClients.user.Username, Password: testClients.user.HashedPassword, TOTPToken: generateTOTPTokenForUser(t, testClients.user)})
			require.NoError(t, err)

			require.NoError(t, testClients.userClient.SetOptions(apiclient.UsingOAuth2(ctx, createdClientID, createdClientSecret, []string{"account_member"}, tokenResponse.Token)))

			// create a webhook

			// Create webhook.
			exampleWebhook := fakes.BuildFakeWebhook()
			exampleWebhookInput := converters.ConvertWebhookToWebhookCreationRequestInput(exampleWebhook)
			createdWebhook, err := testClients.userClient.CreateWebhook(ctx, exampleWebhookInput)
			require.NoError(t, err)

			checkWebhookEquality(t, exampleWebhook, createdWebhook)

			createdWebhook, err = testClients.userClient.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, createdWebhook, err)

			require.Equal(t, account.ID, createdWebhook.BelongsToAccount)

			// check that userClient cannot see the webhook
			webhook, err := futureOwnerClient.GetWebhook(ctx, createdWebhook.ID)
			require.Nil(t, webhook)
			require.Error(t, err)

			// add them to the account
			_, err = testClients.userClient.TransferAccountOwnership(ctx, account.ID, &types.AccountOwnershipTransferInput{
				Reason:       t.Name(),
				CurrentOwner: account.BelongsToUser,
				NewOwner:     futureOwner.ID,
			})
			require.NoError(t, err)

			_, err = futureOwnerClient.SetDefaultAccount(ctx, account.ID)
			require.NoError(t, err)

			tokenResponse, err = futureOwnerClient.LoginForToken(ctx, &types.UserLoginInput{Username: futureOwner.Username, Password: futureOwner.HashedPassword, TOTPToken: generateTOTPTokenForUser(t, futureOwner)})
			require.NoError(t, err)

			require.NoError(t, futureOwnerClient.SetOptions(apiclient.UsingOAuth2(ctx, createdClientID, createdClientSecret, []string{"account_member"}, tokenResponse.Token)))

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

func (s *TestSuite) TestAccounts_UsersHaveBackupAccountCreatedForThemWhenRemovedFromLastAccount() {
	s.runTest("should be possible to invite a user via referral link", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			currentStatus, statusErr := testClients.userClient.GetAuthStatus(s.ctx)
			requireNotNilAndNoProblems(t, currentStatus, statusErr)
			relevantAccountID := currentStatus.ActiveAccount

			inviteReq := &types.AccountInvitationCreationRequestInput{
				Note:    t.Name(),
				ToEmail: gofakeit.Email(),
			}
			createdInvitation, err := testClients.userClient.CreateAccountInvitation(ctx, relevantAccountID, inviteReq)
			require.NoError(t, err)

			createdInvitation, err = testClients.userClient.GetAccountInvitation(ctx, createdInvitation.ID)
			requireNotNilAndNoProblems(t, createdInvitation, err)

			sentInvitations, err := testClients.userClient.GetSentAccountInvitations(ctx, nil)
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

			accounts, err := c.GetAccounts(ctx, nil)
			require.NoError(t, err)

			assert.Len(t, accounts.Data, 2)

			var (
				found            bool
				otherAccountID string
			)

			for _, account := range accounts.Data {
				if account.ID == relevantAccountID {
					if !found {
						found = true
					}
				} else {
					otherAccountID = account.ID
				}
			}

			require.NotEmpty(t, otherAccountID)
			require.True(t, found)

			require.NoError(t, testClients.userClient.ArchiveUserMembership(ctx, relevantAccountID, u.ID))

			u.HashedPassword = regInput.Password

			tokenResponse, err := c.LoginForToken(ctx, &types.UserLoginInput{Username: u.Username, Password: u.HashedPassword, TOTPToken: generateTOTPTokenForUser(t, u)})
			require.NoError(t, err)

			require.NoError(t, c.SetOptions(apiclient.UsingOAuth2(ctx, createdClientID, createdClientSecret, []string{"account_member"}, tokenResponse.Token)))

			account, err := c.GetActiveAccount(ctx)
			requireNotNilAndNoProblems(t, account, err)
			assert.NotEqual(t, relevantAccountID, account.ID)

			require.True(t, found)
		}
	})
}

*/
