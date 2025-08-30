package integration

import (
	"github.com/dinnerdonebetter/backend/internal/domain/identity/fakes"
	authsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAuth_LoginForToken(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()

		user := createServiceUserForTest(t, grpcTestServerAddress, true, fakes.BuildFakeUserRegistrationInput())
		actual := fetchLoginTokenForUserForTest(t, grpcTestServerAddress, user)

		assert.NotEmpty(t, actual)
	})

	T.Run("2FA is not required for non-admin users who haven't verified their secrets", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		user := createServiceUserForTest(t, grpcTestServerAddress, false, fakes.BuildFakeUserRegistrationInput())

		loginInput := &authsvc.UserLoginInput{
			Username: user.Username,
			Password: user.HashedPassword,
		}

		unauthedClient := buildUnauthenticatedGRPCClientForTest(t, grpcTestServerAddress)

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
			TOTPToken: "otp scode",
		}

		unauthedClient := buildUnauthenticatedGRPCClientForTest(t, grpcTestServerAddress)

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
			TOTPToken: generateTOTPCodeForUserForTest(t, premadeAdminUser),
		}

		unauthedClient := buildUnauthenticatedGRPCClientForTest(t, grpcTestServerAddress)

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

		user := createServiceUserForTest(t, grpcTestServerAddress, true, fakes.BuildFakeUserRegistrationInput())

		loginInput := &authsvc.UserLoginInput{
			Username:  user.Username,
			Password:  user.HashedPassword,
			TOTPToken: generateTOTPCodeForUserForTest(t, user),
		}

		unauthedClient := buildUnauthenticatedGRPCClientForTest(t, grpcTestServerAddress)

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
			TOTPToken: generateTOTPCodeForUserForTest(t, premadeAdminUser),
		}

		unauthedClient := buildUnauthenticatedGRPCClientForTest(t, grpcTestServerAddress)

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
			TOTPToken: "000000",
		}

		unauthedClient := buildUnauthenticatedGRPCClientForTest(t, grpcTestServerAddress)

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
			TOTPToken: "otp scode",
		}

		unauthedClient := buildUnauthenticatedGRPCClientForTest(t, grpcTestServerAddress)

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

		unauthedClient := buildUnauthenticatedGRPCClientForTest(t, grpcTestServerAddress)

		res, err := unauthedClient.GetAuthStatus(ctx, &authsvc.GetAuthStatusRequest{})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	T.Run("for logged in user", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		user, testClient := createUserAndClientForTest(T, httpTestServerAddress, grpcTestServerAddress)

		res, err := testClient.GetAuthStatus(ctx, &authsvc.GetAuthStatusRequest{})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, user.ID, res.UserID)
	})
}

func TestAuth_ChangingPassword(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		user, testClient := createUserAndClientForTest(T, httpTestServerAddress, grpcTestServerAddress)

		_, err := testClient.UpdatePassword(ctx, &authsvc.UpdatePasswordRequest{
			NewPassword:     user.HashedPassword + "blah",
			CurrentPassword: user.HashedPassword,
			TOTPToken:       generateTOTPCodeForUserForTest(t, user),
		})
		require.NoError(t, err)

		unauthedClient := buildUnauthenticatedGRPCClientForTest(t, grpcTestServerAddress)

		tokenRes, err := unauthedClient.LoginForToken(ctx, &authsvc.LoginForTokenRequest{
			Input: &authsvc.UserLoginInput{
				Username:  user.Username,
				Password:  user.HashedPassword + "blah",
				TOTPToken: generateTOTPCodeForUserForTest(t, user),
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, tokenRes)
		assert.NotEmpty(t, tokenRes.Result.AccessToken)
	})

	T.Run("with inadequate new password", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		user, testClient := createUserAndClientForTest(T, httpTestServerAddress, grpcTestServerAddress)

		_, err := testClient.UpdatePassword(ctx, &authsvc.UpdatePasswordRequest{
			NewPassword:     "b",
			CurrentPassword: user.HashedPassword,
			TOTPToken:       generateTOTPCodeForUserForTest(t, user),
		})
		assert.Error(t, err)
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
		t.Skipf("TODO")
	})
}

/*

func (s *TestSuite) TestTOTPSecretChanging() {
	s.Run("should be possible to change your TOTP secret", func() {
		t := s.T()

		ctx, span := tracing.StartCustomSpan(context.Background(), t.Name())
		defer span.End()

		testUser, testClient := createUserAndClientForTest(ctx, t, nil)

		r, err := testClient.RefreshTOTPSecret(ctx, &types.TOTPSecretRefreshInput{
			CurrentPassword: testUser.HashedPassword,
			TOTPToken:       generateTOTPTokenForUser(t, testUser),
		})
		require.NoError(t, err)

		testUser.TwoFactorSecret = r.TwoFactorSecret
		_, err = testClient.VerifyTOTPSecret(ctx, &types.TOTPSecretVerificationInput{
			TOTPToken: generateTOTPTokenForUser(t, testUser),
			UserID:    testUser.ID,
		})
		require.NoError(t, err)

		// create login request.
		newToken, err := totp.GenerateCode(r.TwoFactorSecret, time.Now().UTC())
		requireNotNilAndNoProblems(t, newToken, err)

		secondCookie, err := testClient.LoginForToken(ctx, &types.UserLoginInput{
			Username:  testUser.Username,
			Password:  testUser.HashedPassword,
			TOTPToken: newToken,
		})
		assert.NoError(t, err)
		assert.NotNil(t, secondCookie)
	})
}

func (s *TestSuite) TestTOTPTokenValidation() {
	s.Run("should be possible to validate TOTP", func() {
		t := s.T()

		ctx, span := tracing.StartSpan(context.Background())
		defer span.End()

		testClient := buildSimpleClient(t)

		userCR, err := testClient.CreateUser(ctx, fakes.BuildFakeUserRegistrationInput())
		requireNotNilAndNoProblems(t, userCR, err)

		user, err := premadeAdminClient.GetUser(ctx, userCR.CreatedUserID)
		requireNotNilAndNoProblems(t, user, err)
		user.TwoFactorSecret = userCR.TwoFactorSecret

		_, err = testClient.VerifyTOTPSecret(ctx, &types.TOTPSecretVerificationInput{
			TOTPToken: generateTOTPTokenForUser(t, user),
			UserID:    user.ID,
		})
		assert.NoError(t, err)
	})

	s.Run("should not be possible to validate an invalid TOTP", func() {
		t := s.T()

		ctx, span := tracing.StartCustomSpan(context.Background(), t.Name())
		defer span.End()

		testClient := buildSimpleClient(t)

		userCR, err := testClient.CreateUser(ctx, fakes.BuildFakeUserRegistrationInput())
		requireNotNilAndNoProblems(t, userCR, err)

		user, err := premadeAdminClient.GetUser(ctx, userCR.CreatedUserID)
		requireNotNilAndNoProblems(t, user, err)

		_, err = testClient.VerifyTOTPSecret(ctx, &types.TOTPSecretVerificationInput{
			TOTPToken: "NOTREAL",
			UserID:    user.ID,
		})
		assert.Error(t, err)
	})
}

func (s *TestSuite) TestLogin_RequestingPasswordReset() {
	s.Run("able to reset one's password and then redeem it", func() {
		t := s.T()

		ctx, span := tracing.StartSpan(context.Background())
		defer span.End()

		u, testClient := createUserAndClientForTest(ctx, t, nil)

		_, err := testClient.RequestPasswordResetToken(ctx, &types.PasswordResetTokenCreationRequestInput{EmailAddress: u.EmailAddress})
		require.NoError(t, err)

		dbAddr := os.Getenv("TARGET_DATABASE")
		if dbAddr == "" {
			panic("empty database address provided")
		}

		db, err := sql.Open("pgx", dbAddr)
		if err != nil {
			panic(err)
		}

		var token string
		queryErr := db.QueryRow(`SELECT token FROM password_reset_tokens WHERE belongs_to_user = $1`, u.ID).Scan(&token)
		require.NoError(t, queryErr)

		resetToken, err := dbManager.GetPasswordResetTokenByToken(ctx, token)
		requireNotNilAndNoProblems(t, resetToken, err)

		fakeInput := fakes.BuildFakeUserCreationInput()

		_, err = testClient.RedeemPasswordResetToken(ctx, &types.PasswordResetTokenRedemptionRequestInput{
			Token:       resetToken.Token,
			NewPassword: fakeInput.Password,
		})
		require.NoError(t, err)

		cookie, err := testClient.LoginForToken(ctx, &types.UserLoginInput{
			Username:  u.Username,
			Password:  fakeInput.Password,
			TOTPToken: generateTOTPTokenForUser(t, u),
		})
		requireNotNilAndNoProblems(t, cookie, err)

		_, err = testClient.RedeemPasswordResetToken(ctx, &types.PasswordResetTokenRedemptionRequestInput{
			Token:       resetToken.Token,
			NewPassword: fakeInput.Password,
		})
		require.Error(t, err)
	})
}

*/
