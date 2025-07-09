// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_AdminLoginForToken(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/users/login/jwt/admin"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		data := &TokenResponse{}
		expected := &APIResponse[*TokenResponse]{
			Data: data,
		}

		exampleInput := &UserLoginInput{}

		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.AdminLoginForToken(ctx, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleInput := &UserLoginInput{}

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.AdminLoginForToken(ctx, exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleInput := &UserLoginInput{}

		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.AdminLoginForToken(ctx, exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
