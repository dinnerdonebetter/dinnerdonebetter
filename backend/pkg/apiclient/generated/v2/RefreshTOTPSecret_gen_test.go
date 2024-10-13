// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestClient_RefreshTOTPSecret(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/users/totp_secret/new"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		data := fakes.BuildFakeTOTPSecretRefreshResponse()
		expected := &types.APIResponse[*types.TOTPSecretRefreshResponse]{
			Data: data,
		}

		exampleInput := fakes.BuildFakeTOTPSecretRefreshInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.RefreshTOTPSecret(ctx, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleInput := fakes.BuildFakeTOTPSecretRefreshInput()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.RefreshTOTPSecret(ctx, exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleInput := fakes.BuildFakeTOTPSecretRefreshInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.RefreshTOTPSecret(ctx, exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
