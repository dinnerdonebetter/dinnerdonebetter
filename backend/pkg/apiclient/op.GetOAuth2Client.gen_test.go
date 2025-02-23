// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetOAuth2Client(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/oauth2_clients/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		oauth2ClientID := fake.BuildFakeID()

		data := &OAuth2Client{}
		expected := &APIResponse[*OAuth2Client]{
			Data: data,
		}

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, oauth2ClientID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.GetOAuth2Client(ctx, oauth2ClientID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with invalid oauth2Client ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetOAuth2Client(ctx, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		oauth2ClientID := fake.BuildFakeID()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetOAuth2Client(ctx, oauth2ClientID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		oauth2ClientID := fake.BuildFakeID()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, oauth2ClientID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetOAuth2Client(ctx, oauth2ClientID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
