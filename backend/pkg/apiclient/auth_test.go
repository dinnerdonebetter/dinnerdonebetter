package apiclient

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestAuth(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(authTestSuite))
}

type authTestSuite struct {
	suite.Suite

	ctx           context.Context
	exampleUser   *types.User
	exampleCookie *http.Cookie
}

var _ suite.SetupTestSuite = (*authTestSuite)(nil)

func (s *authTestSuite) SetupTest() {
	s.ctx = context.Background()

	s.exampleCookie = &http.Cookie{Name: s.T().Name()}

	s.exampleUser = fakes.BuildFakeUser()
	// the hashed passwords is never transmitted over the wire.
	s.exampleUser.HashedPassword = ""
	// the two factor secret is transmitted over the wire only on creation.
	s.exampleUser.TwoFactorSecret = ""
	// the two factor secret validation is never transmitted over the wire.
	s.exampleUser.TwoFactorSecretVerifiedAt = nil
}

func (s *authTestSuite) TestClient_UserStatus() {
	const expectedPath = "/auth/status"

	s.Run("standard", func() {
		t := s.T()

		expected := &types.UserStatusResponse{
			AccountStatus:            s.exampleUser.AccountStatus,
			AccountStatusExplanation: s.exampleUser.AccountStatusExplanation,
			UserIsAuthenticated:      true,
		}
		expectedResponse := &types.APIResponse[*types.UserStatusResponse]{
			Data: expected,
		}

		spec := newRequestSpec(false, http.MethodGet, "", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, expectedResponse)

		actual, err := c.UserStatus(s.ctx)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.UserStatus(s.ctx)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.UserStatus(s.ctx)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func (s *authTestSuite) TestClient_BeginSession() {
	const expectedPath = "/users/login"

	spec := newRequestSpec(false, http.MethodPost, "", expectedPath)

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeUserLoginInputFromUser(s.exampleUser)

		ts := httptest.NewTLSServer(http.HandlerFunc(
			func(res http.ResponseWriter, req *http.Request) {
				assertRequestQuality(t, req, spec)

				http.SetCookie(res, &http.Cookie{Name: s.exampleUser.Username})
			},
		))
		c := buildTestClient(t, ts)

		cookie, err := c.BeginSession(s.ctx, exampleInput)
		require.NotNil(t, cookie)
		assert.NoError(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		cookie, err := c.BeginSession(s.ctx, nil)
		assert.Nil(t, cookie)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeUserLoginInputFromUser(s.exampleUser)

		c := buildTestClientWithInvalidURL(t)

		cookie, err := c.BeginSession(s.ctx, exampleInput)
		assert.Nil(t, cookie)
		assert.Error(t, err)
	})

	s.Run("with timeout", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeUserLoginInputFromUser(s.exampleUser)
		c, _ := buildTestClientThatWaitsTooLong(t)

		cookie, err := c.BeginSession(s.ctx, exampleInput)
		require.Nil(t, cookie)
		assert.Error(t, err)
	})

	s.Run("with missing cookie", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeUserLoginInputFromUser(s.exampleUser)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusOK)

		cookie, err := c.BeginSession(s.ctx, exampleInput)
		require.Nil(t, cookie)
		assert.Error(t, err)
	})
}

func (s *authTestSuite) TestClient_EndSession() {
	const expectedPath = "/users/logout"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusAccepted)

		err := c.EndSession(s.ctx)
		assert.NoError(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.EndSession(s.ctx)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.EndSession(s.ctx)
		assert.Error(t, err)
	})
}

func (s *authTestSuite) TestClient_ChangePassword() {
	const expectedPath = "/api/v1/users/password/new"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPath)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusAccepted)
		exampleInput := fakes.BuildFakePasswordUpdateInput()

		err := c.ChangePassword(s.ctx, s.exampleCookie, exampleInput)
		assert.NoError(t, err)
	})

	s.Run("with nil cookie", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := fakes.BuildFakePasswordUpdateInput()

		err := c.ChangePassword(s.ctx, nil, exampleInput)
		assert.Error(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ChangePassword(s.ctx, s.exampleCookie, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		exampleInput := fakes.BuildFakePasswordUpdateInput()

		err := c.ChangePassword(s.ctx, s.exampleCookie, exampleInput)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)
		exampleInput := fakes.BuildFakePasswordUpdateInput()

		err := c.ChangePassword(s.ctx, s.exampleCookie, exampleInput)
		assert.Error(t, err)
	})

	s.Run("with unsatisfactory response code", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPath)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusBadRequest)
		exampleInput := fakes.BuildFakePasswordUpdateInput()

		err := c.ChangePassword(s.ctx, s.exampleCookie, exampleInput)
		assert.Error(t, err)
	})
}

func (s *authTestSuite) TestClient_CycleTwoFactorSecret() {
	const expectedPath = "/api/v1/users/totp_secret/new"

	s.Run("standard", func() {
		t := s.T()

		expected := &types.TOTPSecretRefreshResponse{
			TwoFactorQRCode: t.Name(),
			TwoFactorSecret: t.Name(),
		}
		expectedResponse := &types.APIResponse[*types.TOTPSecretRefreshResponse]{
			Data: expected,
		}

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, expectedResponse)
		exampleInput := fakes.BuildFakeTOTPSecretRefreshInput()

		actual, err := c.CycleTwoFactorSecret(s.ctx, s.exampleCookie, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	s.Run("with nil cookie", func() {
		t := s.T()

		expected := &types.TOTPSecretRefreshResponse{
			TwoFactorQRCode: t.Name(),
			TwoFactorSecret: t.Name(),
		}
		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		exampleInput := fakes.BuildFakeTOTPSecretRefreshInput()

		actual, err := c.CycleTwoFactorSecret(s.ctx, nil, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CycleTwoFactorSecret(s.ctx, s.exampleCookie, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.TOTPSecretRefreshInput{}

		actual, err := c.CycleTwoFactorSecret(s.ctx, s.exampleCookie, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		exampleInput := fakes.BuildFakeTOTPSecretRefreshInput()

		actual, err := c.CycleTwoFactorSecret(s.ctx, s.exampleCookie, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)
		exampleInput := fakes.BuildFakeTOTPSecretRefreshInput()

		actual, err := c.CycleTwoFactorSecret(s.ctx, s.exampleCookie, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func (s *authTestSuite) TestClient_VerifyTOTPSecret() {
	const expectedPath = "/users/totp_secret/verify"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		exampleInput := fakes.BuildFakeTOTPSecretVerificationInputForUser(s.exampleUser)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusAccepted)

		err := c.VerifyTOTPSecret(s.ctx, s.exampleUser.ID, exampleInput.TOTPToken)
		assert.NoError(t, err)
	})

	s.Run("with invalid user ID", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeTOTPSecretVerificationInputForUser(s.exampleUser)
		c, _ := buildSimpleTestClient(t)

		err := c.VerifyTOTPSecret(s.ctx, "", exampleInput.TOTPToken)
		assert.Error(t, err)
	})

	s.Run("with invalid token", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.VerifyTOTPSecret(s.ctx, s.exampleUser.ID, " doesn't parse lol ")
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeTOTPSecretVerificationInputForUser(s.exampleUser)

		c := buildTestClientWithInvalidURL(t)

		err := c.VerifyTOTPSecret(s.ctx, s.exampleUser.ID, exampleInput.TOTPToken)
		assert.Error(t, err)
	})

	s.Run("with bad request response", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		exampleInput := fakes.BuildFakeTOTPSecretVerificationInputForUser(s.exampleUser)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusBadRequest)

		err := c.VerifyTOTPSecret(s.ctx, s.exampleUser.ID, exampleInput.TOTPToken)
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidTOTPToken, err)
	})

	s.Run("with otherwise invalid status code response", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		exampleInput := fakes.BuildFakeTOTPSecretVerificationInputForUser(s.exampleUser)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusInternalServerError)

		err := c.VerifyTOTPSecret(s.ctx, s.exampleUser.ID, exampleInput.TOTPToken)
		assert.Error(t, err)
	})

	s.Run("with timeout", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)
		c.unauthenticatedClient.Timeout = time.Millisecond
		exampleInput := fakes.BuildFakeTOTPSecretVerificationInputForUser(s.exampleUser)

		err := c.VerifyTOTPSecret(s.ctx, s.exampleUser.ID, exampleInput.TOTPToken)
		assert.Error(t, err)
	})
}

func (s *authTestSuite) TestClient_RequestPasswordResetToken() {
	const expectedPath = "/users/password/reset"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusAccepted)

		err := c.RequestPasswordResetToken(s.ctx, s.exampleUser.EmailAddress)
		assert.NoError(t, err)
	})

	s.Run("with invalid user ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.RequestPasswordResetToken(s.ctx, "")
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.RequestPasswordResetToken(s.ctx, s.exampleUser.EmailAddress)
		assert.Error(t, err)
	})

	s.Run("with bad request response", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusBadRequest)

		err := c.RequestPasswordResetToken(s.ctx, s.exampleUser.EmailAddress)
		assert.Error(t, err)
	})

	s.Run("with otherwise invalid status code response", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusInternalServerError)

		err := c.RequestPasswordResetToken(s.ctx, s.exampleUser.EmailAddress)
		assert.Error(t, err)
	})

	s.Run("with timeout", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)
		c.unauthenticatedClient.Timeout = time.Millisecond

		err := c.RequestPasswordResetToken(s.ctx, s.exampleUser.EmailAddress)
		assert.Error(t, err)
	})
}

func (s *authTestSuite) TestClient_RedeemPasswordResetToken() {
	const expectedPath = "/users/password/reset/redeem"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusAccepted)

		exampleInput := fakes.BuildFakePasswordResetTokenRedemptionRequestInput()

		err := c.RedeemPasswordResetToken(s.ctx, exampleInput)
		assert.NoError(t, err)
	})

	s.Run("with invalid user ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.RedeemPasswordResetToken(s.ctx, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		exampleInput := fakes.BuildFakePasswordResetTokenRedemptionRequestInput()

		err := c.RedeemPasswordResetToken(s.ctx, exampleInput)
		assert.Error(t, err)
	})

	s.Run("with bad request response", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusBadRequest)

		exampleInput := fakes.BuildFakePasswordResetTokenRedemptionRequestInput()

		err := c.RedeemPasswordResetToken(s.ctx, exampleInput)
		assert.Error(t, err)
	})

	s.Run("with otherwise invalid status code response", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusInternalServerError)

		exampleInput := fakes.BuildFakePasswordResetTokenRedemptionRequestInput()

		err := c.RedeemPasswordResetToken(s.ctx, exampleInput)
		assert.Error(t, err)
	})

	s.Run("with timeout", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)
		c.unauthenticatedClient.Timeout = time.Millisecond

		exampleInput := fakes.BuildFakePasswordResetTokenRedemptionRequestInput()

		err := c.RedeemPasswordResetToken(s.ctx, exampleInput)
		assert.Error(t, err)
	})
}
