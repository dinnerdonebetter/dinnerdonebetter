package authentication

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/internal/authentication"
	mockauthn "github.com/dinnerdonebetter/backend/internal/authentication/mock"
	testutils "github.com/dinnerdonebetter/backend/tests/utils"

	"github.com/gorilla/securecookie"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAuthenticationService_getUserIDFromCookie(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		helper.ctx, helper.req, _ = attachCookieToRequestForTest(t, helper.service, helper.req, helper.exampleUser)

		_, userID, err := helper.service.getUserIDFromCookie(helper.ctx, helper.req)
		assert.NoError(t, err)
		assert.Equal(t, helper.exampleUser.ID, userID)
	})

	T.Run("with invalid cookie", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		// begin building bad cookie.
		// NOTE: any code here is duplicated from service.buildAuthCookie
		// any changes made there might need to be reflected here.
		c := &http.Cookie{
			Name:     helper.service.config.Cookies.Name,
			Value:    "blah blah blah this is not a real cookie",
			Path:     "/",
			HttpOnly: true,
		}
		// end building bad cookie.
		helper.req.AddCookie(c)

		_, userID, err := helper.service.getUserIDFromCookie(helper.req.Context(), helper.req)
		assert.Error(t, err)
		assert.Zero(t, userID)
	})

	T.Run("without cookie", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		_, userID, err := helper.service.getUserIDFromCookie(helper.req.Context(), helper.req)
		assert.Error(t, err)
		assert.Equal(t, err, http.ErrNoCookie)
		assert.Zero(t, userID)
	})

	T.Run("with error loading from session", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		expectedToken := "blahblah"

		sm := &mockSessionManager{}
		sm.On("Load", testutils.ContextMatcher, expectedToken).Return(helper.ctx, errors.New("blah"))
		helper.service.sessionManager = sm

		c, err := helper.service.buildCookie(helper.ctx, helper.service.config.Cookies.Domain, expectedToken, time.Now().Add(helper.service.config.Cookies.Lifetime))
		require.NoError(t, err)
		helper.req.AddCookie(c)

		_, _, err = helper.service.getUserIDFromCookie(helper.ctx, helper.req)
		assert.Error(t, err)
	})

	T.Run("with no user ID attached to context", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		ctx, sessionErr := helper.service.sessionManager.Load(helper.req.Context(), "")
		require.NoError(t, sessionErr)
		require.NoError(t, helper.service.sessionManager.RenewToken(ctx))

		token, _, err := helper.service.sessionManager.Commit(ctx)
		assert.NotEmpty(t, token)
		assert.NoError(t, err)

		c, err := helper.service.buildCookie(helper.ctx, helper.service.config.Cookies.Domain, token, time.Now().Add(helper.service.config.Cookies.Lifetime))
		require.NoError(t, err)
		helper.req.AddCookie(c)

		_, _, err = helper.service.getUserIDFromCookie(helper.ctx, helper.req)
		assert.Error(t, err)
	})
}

func TestAuthenticationService_validateLogin(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		authenticator := &mockauthn.Authenticator{}
		authenticator.On(
			"CredentialsAreValid",
			testutils.ContextMatcher,
			helper.exampleUser.HashedPassword,
			helper.exampleLoginInput.Password,
			helper.exampleUser.TwoFactorSecret,
			helper.exampleLoginInput.TOTPToken,
		).Return(true, nil)
		helper.service.authenticator = authenticator

		actual, err := helper.service.validateLogin(helper.ctx, helper.exampleUser, helper.exampleLoginInput)
		assert.True(t, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, authenticator)
	})

	T.Run("with invalid two factor code", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		authenticator := &mockauthn.Authenticator{}
		authenticator.On(
			"CredentialsAreValid",
			testutils.ContextMatcher,
			helper.exampleUser.HashedPassword,
			helper.exampleLoginInput.Password,
			helper.exampleUser.TwoFactorSecret,
			helper.exampleLoginInput.TOTPToken,
		).Return(true, authentication.ErrInvalidTOTPToken)
		helper.service.authenticator = authenticator

		actual, err := helper.service.validateLogin(helper.ctx, helper.exampleUser, helper.exampleLoginInput)
		assert.False(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, authenticator)
	})

	T.Run("with error returned from validator", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		expectedErr := errors.New("arbitrary")

		authenticator := &mockauthn.Authenticator{}
		authenticator.On(
			"CredentialsAreValid",
			testutils.ContextMatcher,
			helper.exampleUser.HashedPassword,
			helper.exampleLoginInput.Password,
			helper.exampleUser.TwoFactorSecret,
			helper.exampleLoginInput.TOTPToken,
		).Return(false, expectedErr)
		helper.service.authenticator = authenticator

		actual, err := helper.service.validateLogin(helper.ctx, helper.exampleUser, helper.exampleLoginInput)
		assert.False(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, authenticator)
	})

	T.Run("with invalid login", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		authenticator := &mockauthn.Authenticator{}
		authenticator.On(
			"CredentialsAreValid",
			testutils.ContextMatcher,
			helper.exampleUser.HashedPassword,
			helper.exampleLoginInput.Password,
			helper.exampleUser.TwoFactorSecret,
			helper.exampleLoginInput.TOTPToken,
		).Return(false, nil)
		helper.service.authenticator = authenticator

		actual, err := helper.service.validateLogin(helper.ctx, helper.exampleUser, helper.exampleLoginInput)
		assert.False(t, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, authenticator)
	})
}

func TestAuthenticationService_buildCookie(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		cookie, err := helper.service.buildCookie(helper.ctx, helper.service.config.Cookies.Domain, "example", time.Now().Add(helper.service.config.Cookies.Lifetime))
		assert.NotNil(t, cookie)
		assert.NoError(t, err)
	})

	T.Run("with invalid cookie builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		helper.service.cookieManager = securecookie.New(
			securecookie.GenerateRandomKey(0),
			[]byte(""),
		)

		cookie, err := helper.service.buildCookie(helper.ctx, helper.service.config.Cookies.Domain, "example", time.Now().Add(helper.service.config.Cookies.Lifetime))
		assert.Nil(t, cookie)
		assert.Error(t, err)
	})
}
