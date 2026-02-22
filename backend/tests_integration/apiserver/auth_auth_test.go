package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	"github.com/dinnerdonebetter/backend/internal/domain/identity/fakes"
	authsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/auth"
	identitysvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/identity"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/auditlogentries"
	authrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/auth"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuth_LoginForToken_DesiredAccount(T *testing.T) {
	T.Parallel()

	T.Run("login for non-default account", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		// Create inviter user + client; inviter gets account A (default from registration).
		_, testClient := createUserAndClientForTest(t)
		accountRes, err := testClient.GetActiveAccount(ctx, &authsvc.GetActiveAccountRequest{})
		require.NoError(t, err)
		inviterAccountID := accountRes.Result.Id

		// Create invitee user + client; invitee gets account B (their default from registration).
		inviteeEmailAddress := fmt.Sprintf("invitee%d@testing.com", time.Now().UnixMicro())
		input := &identity.UserRegistrationInput{
			Birthday:              pointer.To(time.Now()),
			EmailAddress:          inviteeEmailAddress,
			FirstName:             fmt.Sprintf("invitee_%d", hashStringToNumber(t.Name()+time.Now().Format(time.RFC3339Nano))),
			AccountName:           fmt.Sprintf("invitee_%d", hashStringToNumber(t.Name()+time.Now().Format(time.RFC3339Nano))),
			LastName:              fmt.Sprintf("invitee_%d", hashStringToNumber(t.Name()+time.Now().Format(time.RFC3339Nano))),
			Password:              fmt.Sprintf("invitee_%d", hashStringToNumber(t.Name()+time.Now().Format(time.RFC3339Nano))),
			Username:              fmt.Sprintf("invitee_%d", hashStringToNumber(t.Name()+time.Now().Format(time.RFC3339Nano))),
			AcceptedPrivacyPolicy: true,
			AcceptedTOS:           true,
		}
		invitee, inviteeClient := createUserAndClientForTestWithRegistrationInput(t, input)

		// Inviter creates invitation for invitee.
		invitation, err := testClient.CreateAccountInvitation(ctx, &identitysvc.CreateAccountInvitationRequest{
			Input: &identitysvc.AccountInvitationCreationRequestInput{
				Note:    t.Name(),
				ToName:  t.Name(),
				ToEmail: inviteeEmailAddress,
			},
		})
		require.NoError(t, err)

		// Invitee accepts invitation. Invitee now has accounts A and B; B remains default.
		_, err = inviteeClient.AcceptAccountInvitation(ctx, &identitysvc.AcceptAccountInvitationRequest{
			AccountInvitationId: invitation.Created.Id,
			Input: &identitysvc.AccountInvitationUpdateRequestInput{
				Token: invitation.Created.Token,
				Note:  t.Name(),
			},
		})
		require.NoError(t, err)

		// Login with DesiredAccountId set to inviter's account (non-default for invitee).
		unauthedClient := buildUnauthenticatedGRPCClientForTest(t)
		tokenRes, err := unauthedClient.LoginForToken(ctx, &authsvc.LoginForTokenRequest{
			Input: &authsvc.UserLoginInput{
				Username:         invitee.Username,
				Password:         invitee.HashedPassword,
				TotpToken:        generateTOTPCodeForUserForTest(t, invitee),
				DesiredAccountId: inviterAccountID,
			},
		})
		require.NoError(t, err)
		require.NotNil(t, tokenRes)
		require.NotEmpty(t, tokenRes.Result.AccessToken)
		assert.Equal(t, inviterAccountID, tokenRes.Result.AccountId, "token response should include the desired (non-default) account")

		// Use JWT directly as Bearer so session context uses the token's account_id claim.
		jwtClient, err := buildAuthedGRPCClientWithBearerToken(tokenRes.Result.AccessToken)
		require.NoError(t, err)
		status, err := jwtClient.GetAuthStatus(ctx, &authsvc.GetAuthStatusRequest{})
		require.NoError(t, err)
		require.NotNil(t, status)
		assert.Equal(t, inviterAccountID, status.ActiveAccount, "session context should use account from token")
	})

	T.Run("login with invalid desired account returns error", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		user, _ := createUserAndClientForTest(t)
		unauthedClient := buildUnauthenticatedGRPCClientForTest(t)

		// Use a made-up account ID the user is not a member of.
		tokenRes, err := unauthedClient.LoginForToken(ctx, &authsvc.LoginForTokenRequest{
			Input: &authsvc.UserLoginInput{
				Username:         user.Username,
				Password:         user.HashedPassword,
				TotpToken:        generateTOTPCodeForUserForTest(t, user),
				DesiredAccountId: nonexistentID,
			},
		})
		assert.Error(t, err)
		assert.Nil(t, tokenRes)
	})
}

func TestAuth_LoginForToken(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()

		user := createServiceUserForTest(t, true, fakes.BuildFakeUserRegistrationInput())
		actual := fetchLoginTokenForUserForTest(t, user)

		assert.NotEmpty(t, actual)
	})

	T.Run("2FA is not required for non-admin users who haven't verified their secrets", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		user := createServiceUserForTest(t, false, fakes.BuildFakeUserRegistrationInput())

		loginInput := &authsvc.UserLoginInput{
			Username: user.Username,
			Password: user.HashedPassword,
		}

		unauthedClient := buildUnauthenticatedGRPCClientForTest(t)

		tokenRes, err := unauthedClient.LoginForToken(ctx, &authsvc.LoginForTokenRequest{
			Input: loginInput,
		})
		assert.NoError(t, err)
		assert.NotNil(t, tokenRes)
		assert.NotEmpty(t, tokenRes.Result.AccessToken)
	})

	T.Run("with bogus input", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		loginInput := &authsvc.UserLoginInput{
			Username:  " ",
			Password:  "1",
			TotpToken: "otp scode",
		}

		unauthedClient := buildUnauthenticatedGRPCClientForTest(t)

		tokenRes, err := unauthedClient.LoginForToken(ctx, &authsvc.LoginForTokenRequest{
			Input: loginInput,
		})
		assert.Error(t, err)
		assert.Nil(t, tokenRes)
	})
}

func TestAuth_AdminLoginForToken(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		loginInput := &authsvc.UserLoginInput{
			Username:  premadeAdminUser.Username,
			Password:  adminUserPassword,
			TotpToken: generateTOTPCodeForUserForTest(t, premadeAdminUser),
		}

		unauthedClient := buildUnauthenticatedGRPCClientForTest(t)

		tokenRes, err := unauthedClient.AdminLoginForToken(ctx, &authsvc.AdminLoginForTokenRequest{
			Input: loginInput,
		})
		assert.NoError(t, err)
		assert.NotNil(t, tokenRes)
		assert.NotEmpty(t, tokenRes.Result.AccessToken)
	})

	T.Run("non-admin users cannot login via this route", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		user := createServiceUserForTest(t, true, fakes.BuildFakeUserRegistrationInput())

		loginInput := &authsvc.UserLoginInput{
			Username:  user.Username,
			Password:  user.HashedPassword,
			TotpToken: generateTOTPCodeForUserForTest(t, user),
		}

		unauthedClient := buildUnauthenticatedGRPCClientForTest(t)

		tokenRes, err := unauthedClient.AdminLoginForToken(ctx, &authsvc.AdminLoginForTokenRequest{
			Input: loginInput,
		})
		assert.Error(t, err)
		assert.Nil(t, tokenRes)
	})

	T.Run("with incorrect password", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		loginInput := &authsvc.UserLoginInput{
			Username:  premadeAdminUser.Username,
			Password:  adminUserPassword + "blah",
			TotpToken: generateTOTPCodeForUserForTest(t, premadeAdminUser),
		}

		unauthedClient := buildUnauthenticatedGRPCClientForTest(t)

		tokenRes, err := unauthedClient.AdminLoginForToken(ctx, &authsvc.AdminLoginForTokenRequest{
			Input: loginInput,
		})
		assert.Error(t, err)
		assert.Nil(t, tokenRes)
	})

	T.Run("with incorrect TOTP Code", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		loginInput := &authsvc.UserLoginInput{
			Username:  premadeAdminUser.Username,
			Password:  adminUserPassword,
			TotpToken: "000000",
		}

		unauthedClient := buildUnauthenticatedGRPCClientForTest(t)

		tokenRes, err := unauthedClient.AdminLoginForToken(ctx, &authsvc.AdminLoginForTokenRequest{
			Input: loginInput,
		})
		assert.Error(t, err)
		assert.Nil(t, tokenRes)
	})

	T.Run("with bogus input", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		loginInput := &authsvc.UserLoginInput{
			Username:  " ",
			Password:  "1",
			TotpToken: "otp scode",
		}

		unauthedClient := buildUnauthenticatedGRPCClientForTest(t)

		tokenRes, err := unauthedClient.AdminLoginForToken(ctx, &authsvc.AdminLoginForTokenRequest{
			Input: loginInput,
		})
		assert.Error(t, err)
		assert.Nil(t, tokenRes)
	})
}

func TestAuth_GetAuthStatus(T *testing.T) {
	T.Parallel()

	T.Run("for unauthenticated user", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		unauthedClient := buildUnauthenticatedGRPCClientForTest(t)

		res, err := unauthedClient.GetAuthStatus(ctx, &authsvc.GetAuthStatusRequest{})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	T.Run("for logged in user", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		user, testClient := createUserAndClientForTest(T)

		res, err := testClient.GetAuthStatus(ctx, &authsvc.GetAuthStatusRequest{})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, user.ID, res.UserId)
	})
}

func TestAuth_ChangingPassword(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		user, testClient := createUserAndClientForTest(T)

		_, err := testClient.UpdatePassword(ctx, &authsvc.UpdatePasswordRequest{
			NewPassword:     user.HashedPassword + "blah",
			CurrentPassword: user.HashedPassword,
			TotpToken:       generateTOTPCodeForUserForTest(t, user),
		})
		require.NoError(t, err)

		unauthedClient := buildUnauthenticatedGRPCClientForTest(t)

		tokenRes, err := unauthedClient.LoginForToken(ctx, &authsvc.LoginForTokenRequest{
			Input: &authsvc.UserLoginInput{
				Username:  user.Username,
				Password:  user.HashedPassword + "blah",
				TotpToken: generateTOTPCodeForUserForTest(t, user),
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, tokenRes)
		assert.NotEmpty(t, tokenRes.Result.AccessToken)

		AssertAuditLogContainsFuzzyForUser(t, ctx, testClient, user.ID, 10, []*ExpectedAuditEntry{
			{EventType: "updated", ResourceType: "users", RelevantID: user.ID},
		})
	})

	T.Run("with inadequate new password", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		user, testClient := createUserAndClientForTest(T)

		_, err := testClient.UpdatePassword(ctx, &authsvc.UpdatePasswordRequest{
			NewPassword:     "b",
			CurrentPassword: user.HashedPassword,
			TotpToken:       generateTOTPCodeForUserForTest(t, user),
		})
		assert.Error(t, err)
	})
}

func TestAuth_ChangingTOTPSecret(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		user, testClient := createUserAndClientForTest(T)

		res, err := testClient.RefreshTOTPSecret(ctx, &authsvc.RefreshTOTPSecretRequest{
			CurrentPassword: user.HashedPassword,
			TotpToken:       generateTOTPCodeForUserForTest(t, user),
		})
		require.NoError(t, err)
		res.Result.TwoFactorSecret = user.TwoFactorSecret

		_, err = testClient.VerifyTOTPSecret(ctx, &authsvc.VerifyTOTPSecretRequest{
			TotpToken: generateTOTPCodeForUserForTest(t, user),
			UserId:    user.ID,
		})
		assert.NoError(t, err)

		AssertAuditLogContainsFuzzyForUser(t, ctx, testClient, user.ID, 10, []*ExpectedAuditEntry{
			{EventType: "updated", ResourceType: "users", RelevantID: user.ID},
		})
	})

	T.Run("fails with invalid token", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		user, testClient := createUserAndClientForTest(T)

		_, err := testClient.RefreshTOTPSecret(ctx, &authsvc.RefreshTOTPSecretRequest{
			CurrentPassword: user.HashedPassword,
			TotpToken:       "000000",
		})
		assert.Error(t, err)
	})

	T.Run("fails with invalid password", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		user, testClient := createUserAndClientForTest(T)

		_, err := testClient.RefreshTOTPSecret(ctx, &authsvc.RefreshTOTPSecretRequest{
			CurrentPassword: user.HashedPassword + "blah",
			TotpToken:       generateTOTPCodeForUserForTest(t, user),
		})
		assert.Error(t, err)
	})
}

func TestAuth_RequestingPasswordReset(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		user, testClient := createUserAndClientForTest(T)

		res, err := testClient.RequestPasswordResetToken(ctx, &authsvc.RequestPasswordResetTokenRequest{
			EmailAddress: user.EmailAddress,
		})
		require.NoError(t, err)
		assert.NotNil(t, res)

		// boo, hiss, we're talking directly to the database in an _integration test?_ for shame, for shame.
		var token string
		queryErr := databaseClient.DB().QueryRow(`SELECT token FROM password_reset_tokens WHERE belongs_to_user = $1`, user.ID).Scan(&token)
		require.NoError(t, queryErr)

		auditLogRepo := auditlogentries.ProvideAuditLogRepository(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), databaseClient)
		authRepo := authrepo.ProvideAuthRepository(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), auditLogRepo, databaseClient)

		resetToken, err := authRepo.GetPasswordResetTokenByToken(ctx, token)
		require.NoError(t, err)
		require.NotNil(t, resetToken)

		_, err = testClient.RedeemPasswordResetToken(ctx, &authsvc.RedeemPasswordResetTokenRequest{
			Token:       resetToken.Token,
			NewPassword: user.HashedPassword + "blah",
		})
		require.NoError(t, err)

		tokenRes, err := testClient.LoginForToken(ctx, &authsvc.LoginForTokenRequest{
			Input: &authsvc.UserLoginInput{
				Username:  user.Username,
				Password:  user.HashedPassword + "blah",
				TotpToken: generateTOTPCodeForUserForTest(t, user),
			},
		})
		require.NoError(t, err)
		require.NotNil(t, tokenRes)
		assert.NotEmpty(t, tokenRes.Result.AccessToken)

		// verify that we can't do it twice
		_, err = testClient.RedeemPasswordResetToken(ctx, &authsvc.RedeemPasswordResetTokenRequest{
			Token:       resetToken.Token,
			NewPassword: user.HashedPassword + "blah",
		})
		require.Error(t, err)

		AssertAuditLogContainsFuzzyForUser(t, ctx, testClient, user.ID, 15, []*ExpectedAuditEntry{
			{EventType: "created", ResourceType: "password_reset_tokens"},
			{EventType: "updated", ResourceType: "password_reset_tokens"},
		})
	})
}

// TODO section below this line

/*
somewhere here we should validate that a user can't just pretend to be an admin via OAuth somehow?
I feel like that's now how it presently works anyway, but the whole thing is haphazard and fucked up, so
just shore all of it iup.
*/

func TestAuth_InvalidateToken(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		/*
			there's currently no logout mechanism surfaced in the client, but there should be, and it should be tested
		*/
		t.Skipf("TODO")
	})
}
