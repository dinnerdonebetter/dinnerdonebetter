package integration

import (
	"context"
	"database/sql"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestLoginForJWT() {
	s.Run("logging in and out works", func() {
		t := s.T()

		ctx, span := tracing.StartCustomSpan(context.Background(), t.Name())
		defer span.End()

		testUser, testClient := createUserAndClientForTest(ctx, t, nil)
		token, err := testClient.LoginForJWT(ctx, &types.UserLoginInput{
			Username:  testUser.Username,
			Password:  testUser.HashedPassword,
			TOTPToken: generateTOTPTokenForUser(t, testUser),
		})

		assert.NotEmpty(t, token)
		assert.NoError(t, err)
	})
}

func (s *TestSuite) TestAdminLoginForJWT() {
	s.Run("logging in and out works", func() {
		t := s.T()

		ctx, span := tracing.StartCustomSpan(context.Background(), t.Name())
		defer span.End()

		loginInput := &types.UserLoginInput{
			Username:  premadeAdminUser.Username,
			Password:  premadeAdminUser.HashedPassword,
			TOTPToken: generateTOTPTokenForUser(t, premadeAdminUser),
		}

		adminClient, err := initializeOAuth2PoweredClient(ctx, loginInput)
		require.NoError(t, err)

		token, err := adminClient.AdminLoginForJWT(ctx, &types.UserLoginInput{
			Username:  premadeAdminUser.Username,
			Password:  premadeAdminUser.HashedPassword,
			TOTPToken: generateTOTPTokenForUser(t, premadeAdminUser),
		})

		assert.NotEmpty(t, token)
		assert.NoError(t, err)
	})
}

func (s *TestSuite) TestLogin_WithoutBodyReturnsError() {
	s.Run("login request without body returns error", func() {
		t := s.T()

		ctx, span := tracing.StartSpan(context.Background())
		defer span.End()

		_, testClient := createUserAndClientForTest(ctx, t, nil)

		u, err := url.Parse(testClient.BuildURL(ctx, nil))
		require.NoError(t, err)
		u.Path = "/users/login/jwt"

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), http.NoBody)
		requireNotNilAndNoProblems(t, req, err)

		// execute login request.
		res, err := testClient.PlainClient().Do(req)
		requireNotNilAndNoProblems(t, res, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
}

func (s *TestSuite) TestLogin_ShouldNotBeAbleToLoginWithInvalidPassword() {
	s.Run("should not be able to log in with the wrong password", func() {
		t := s.T()

		ctx, span := tracing.StartCustomSpan(context.Background(), t.Name())
		defer span.End()

		testUser, testClient := createUserAndClientForTest(ctx, t, nil)

		// create login request.
		var badPassword string
		for _, v := range testUser.HashedPassword {
			badPassword = string(v) + badPassword
		}

		r := &types.UserLoginInput{
			Username:  testUser.Username,
			Password:  badPassword,
			TOTPToken: generateTOTPTokenForUser(t, testUser),
		}

		cookie, err := testClient.LoginForJWT(ctx, r)
		assert.Nil(t, cookie)
		assert.Error(t, err)
	})
}

func (s *TestSuite) TestLogin_ShouldNotBeAbleToLoginAsANonexistentUser() {
	s.Run("should not be able to login as someone that does not exist", func() {
		t := s.T()

		ctx, span := tracing.StartCustomSpan(context.Background(), t.Name())
		defer span.End()

		testUser, testClient := createUserAndClientForTest(ctx, t, nil)

		exampleUserCreationInput := fakes.BuildFakeUserCreationInput()
		r := &types.UserLoginInput{
			Username:  exampleUserCreationInput.Username,
			Password:  testUser.HashedPassword,
			TOTPToken: "123456",
		}

		cookie, err := testClient.LoginForJWT(ctx, r)
		assert.Nil(t, cookie)
		assert.Error(t, err)
	})
}

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

		cookie, err := testClient.LoginForJWT(ctx, r)
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
		cookie, err := testClient.LoginForJWT(ctx, &types.UserLoginInput{
			Username:  testUser.Username,
			Password:  testUser.HashedPassword,
			TOTPToken: generateTOTPTokenForUser(t, testUser),
		})

		require.NotNil(t, cookie)
		assert.NoError(t, err)

		actual, err := testClient.GetAuthStatus(ctx)
		assert.NoError(t, err)

		assert.Equal(t, true, actual.UserIsAuthenticated, "expected UserIsAuthenticated to equal %v, but got %v", true, actual.UserIsAuthenticated)
		assert.Equal(t, string(types.UnverifiedHouseholdStatus), actual.AccountStatus, "expected AccountStatus to equal %v, but got %v", types.GoodStandingUserAccountStatus, actual.AccountStatus)
		assert.Equal(t, "", actual.AccountStatusExplanation, "expected AccountStatusExplanation to equal %v, but got %v", "", actual.AccountStatusExplanation)
		assert.NotZero(t, actual.ActiveHousehold)
	})
}

func (s *TestSuite) TestPasswordChanging() {
	s.Run("should be possible to change your password", func() {
		t := s.T()

		ctx, span := tracing.StartSpan(context.Background())
		defer span.End()

		testUser, testClient := createUserAndClientForTest(ctx, t, nil)

		// login.
		cookie, err := testClient.LoginForJWT(ctx, &types.UserLoginInput{
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
		cookie, err = testClient.LoginForJWT(ctx, &types.UserLoginInput{
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

		secondCookie, err := testClient.LoginForJWT(ctx, &types.UserLoginInput{
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

		resetToken, err := dbmanager.GetPasswordResetTokenByToken(ctx, token)
		requireNotNilAndNoProblems(t, resetToken, err)

		fakeInput := fakes.BuildFakeUserCreationInput()

		_, err = testClient.RedeemPasswordResetToken(ctx, &types.PasswordResetTokenRedemptionRequestInput{
			Token:       resetToken.Token,
			NewPassword: fakeInput.Password,
		})
		require.NoError(t, err)

		cookie, err := testClient.LoginForJWT(ctx, &types.UserLoginInput{
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
