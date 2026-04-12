package authentication

import (
	"context"
	"errors"
	"net/http"
	"testing"

	mockpublishers "github.com/primandproper/platform/messagequeue/mock"

	"github.com/stretchr/testify/assert"
)

//nolint:paralleltest // pending race condition fix on Goth's part.
func Test_service_SSOProviderHandler(T *testing.T) {
	// T.Parallel()

	T.Run("standard", func(t *testing.T) {
		// t.Parallel()

		helper := buildTestHelper(t)
		helper.service.authProviderFetcher = func(*http.Request) string {
			return "google"
		}

		helper.service.SSOLoginHandler(helper.res, helper.req)

		assert.NotEmpty(t, helper.res.Header().Get("Location"))
		assert.Equal(t, http.StatusTemporaryRedirect, helper.res.Code)
	})

	T.Run("with invalid provider", func(t *testing.T) {
		// t.Parallel()

		helper := buildTestHelper(t)
		helper.service.authProviderFetcher = func(*http.Request) string {
			return "NOT REAL LOL"
		}

		helper.service.SSOLoginHandler(helper.res, helper.req)

		assert.Empty(t, helper.res.Header().Get("Location"))
		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
	})

	T.Run("with missing provider", func(t *testing.T) {
		// t.Parallel()

		helper := buildTestHelper(t)
		helper.service.authProviderFetcher = func(*http.Request) string {
			return ""
		}

		helper.service.SSOLoginHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
	})
}

func TestAuthenticationService_postLogin(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dataChangesPublisher := &mockpublishers.PublisherMock{
			PublishFunc: func(_ context.Context, _ any) error { return nil },
		}
		helper.service.dataChangesPublisher = dataChangesPublisher

		statusCode, err := helper.service.postLogin(helper.ctx, helper.exampleUser, helper.exampleAccount.ID)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusAccepted, statusCode)
	})

	T.Run("with publisher error", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dataChangesPublisher := &mockpublishers.PublisherMock{
			PublishFunc: func(_ context.Context, _ any) error { return errors.New("publisher error") },
		}
		helper.service.dataChangesPublisher = dataChangesPublisher

		statusCode, err := helper.service.postLogin(helper.ctx, helper.exampleUser, helper.exampleAccount.ID)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusAccepted, statusCode)
	})
}

func TestAuthenticationService_AuthorizeHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		req := helper.req
		res := helper.res

		// The AuthorizeHandler delegates to the oauth2Server, so we just test that it calls it
		// Since we can't easily mock the oauth2Server.HandleAuthorizeRequest, we'll test error handling
		helper.service.AuthorizeHandler(res, req)

		// The response code will depend on the oauth2 library's behavior
		// We're mainly testing that the method doesn't panic
		assert.True(t, res.Code >= 400) // Expect some error since we don't have a real OAuth2 setup
	})
}

func TestAuthenticationService_TokenHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		req := helper.req
		res := helper.res

		// The TokenHandler delegates to the oauth2Server, so we just test that it calls it
		// Since we can't easily mock the oauth2Server.HandleTokenRequest, we'll test error handling
		helper.service.TokenHandler(res, req)

		// The response code will depend on the oauth2 library's behavior
		// We're mainly testing that the method doesn't panic
		assert.True(t, res.Code >= 400) // Expect some error since we don't have a real OAuth2 setup
	})
}
