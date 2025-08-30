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

		unauthedClient, err := buildUnauthenticatedGRPCClient(grpcTestServerAddress)
		require.NoError(t, err)

		tokenRes, err := unauthedClient.LoginForToken(ctx, &authsvc.LoginForTokenRequest{
			Input: loginInput,
		})
		assert.NoError(t, err)
		assert.NotNil(t, tokenRes)
	})

	T.Run("with bogus input", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		loginInput := &authsvc.UserLoginInput{
			Username:  " ",
			Password:  "1",
			TOTPToken: "otp scode",
		}

		unauthedClient, err := buildUnauthenticatedGRPCClient(grpcTestServerAddress)
		require.NoError(t, err)

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

		unauthedClient, err := buildUnauthenticatedGRPCClient(grpcTestServerAddress)
		require.NoError(t, err)

		tokenRes, err := unauthedClient.AdminLoginForToken(ctx, &authsvc.AdminLoginForTokenRequest{
			Input: loginInput,
		})
		assert.NoError(t, err)
		assert.NotNil(t, tokenRes)
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

		unauthedClient, err := buildUnauthenticatedGRPCClient(grpcTestServerAddress)
		require.NoError(t, err)

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

		unauthedClient, err := buildUnauthenticatedGRPCClient(grpcTestServerAddress)
		require.NoError(t, err)

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

		unauthedClient, err := buildUnauthenticatedGRPCClient(grpcTestServerAddress)
		require.NoError(t, err)

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

		unauthedClient, err := buildUnauthenticatedGRPCClient(grpcTestServerAddress)
		require.NoError(t, err)

		tokenRes, err := unauthedClient.AdminLoginForToken(ctx, &authsvc.AdminLoginForTokenRequest{
			Input: loginInput,
		})
		assert.Error(t, err)
		assert.Nil(t, tokenRes)
	})
}

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

func (s *TestSuite) TestLogin_ShouldNotBeAbleToLoginWithoutValidating2FASecret() {
	s.Run("should be able to login without validating 2FA secret", func() {
		t := s.T()

		ctx, span := tracing.StartSpan(context.Background())
		defer span.End()

		testClient := buildSimpleClient(t)

		// create a userClient.
		exampleUser := fakes.BuildFakeUser()
		exampleUserCreationInput := fakes.BuildFakeUserRegistrationInputFromUser(exampleUser)
		ucr, err := testClient.CreateUser(ctx, exampleUserCreationInput)
		requireNotNilAndNoProblems(t, ucr, err)

		// create login request.
		r := &types.UserLoginInput{
			Username: exampleUserCreationInput.Username,
			Password: exampleUserCreationInput.Password,
		}

		cookie, err := testClient.LoginForToken(ctx, r)
		assert.NotNil(t, cookie)
		assert.NoError(t, err)
	})
}

func (s *TestSuite) TestCheckingAuthStatus() {
	s.Run("checking auth status", func() {
		t := s.T()

		ctx, span := tracing.StartCustomSpan(context.Background(), t.Name())
		defer span.End()

		testUser, testClient := createUserAndClientForTest(ctx, t, nil)
		cookie, err := testClient.LoginForToken(ctx, &types.UserLoginInput{
			Username:  testUser.Username,
			Password:  testUser.HashedPassword,
			TOTPToken: generateTOTPTokenForUser(t, testUser),
		})

		require.NotNil(t, cookie)
		assert.NoError(t, err)

		actual, err := testClient.GetAuthStatus(ctx)
		assert.NoError(t, err)

		assert.Equal(t, true, actual.UserIsAuthenticated, "expected UserIsAuthenticated to equal %v, but got %v", true, actual.UserIsAuthenticated)
		assert.Equal(t, string(types.UnverifiedAccountStatus), actual.AccountStatus, "expected AccountStatus to equal %v, but got %v", types.GoodStandingUserAccountStatus, actual.AccountStatus)
		assert.Equal(t, "", actual.AccountStatusExplanation, "expected AccountStatusExplanation to equal %v, but got %v", "", actual.AccountStatusExplanation)
		assert.NotZero(t, actual.ActiveAccount)
	})
}

func (s *TestSuite) TestPasswordChanging() {
	s.Run("should be possible to change your password", func() {
		t := s.T()

		ctx, span := tracing.StartSpan(context.Background())
		defer span.End()

		testUser, testClient := createUserAndClientForTest(ctx, t, nil)

		// login.
		cookie, err := testClient.LoginForToken(ctx, &types.UserLoginInput{
			Username:  testUser.Username,
			Password:  testUser.HashedPassword,
			TOTPToken: generateTOTPTokenForUser(t, testUser),
		})
		require.NotNil(t, cookie)
		assert.NoError(t, err)

		// create new passwords.
		var backwardsPass string
		for _, v := range testUser.HashedPassword {
			backwardsPass = string(v) + backwardsPass
		}

		// update passwords.
		assert.NoError(t, testClient.UpdatePassword(ctx, &types.PasswordUpdateInput{
			CurrentPassword: testUser.HashedPassword,
			TOTPToken:       generateTOTPTokenForUser(t, testUser),
			NewPassword:     backwardsPass,
		}))

		// login again with new passwords.
		cookie, err = testClient.LoginForToken(ctx, &types.UserLoginInput{
			Username:  testUser.Username,
			Password:  backwardsPass,
			TOTPToken: generateTOTPTokenForUser(t, testUser),
		})

		assert.NotNil(t, cookie)
		assert.NoError(t, err)
	})
}

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
