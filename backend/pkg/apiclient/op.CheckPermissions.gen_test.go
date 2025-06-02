// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_CheckPermissions(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/users/permissions/check"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		data := &UserPermissionsResponse{}
		expected := &APIResponse[*UserPermissionsResponse]{
			Data: data,
		}

		exampleInput := &UserPermissionsRequestInput{}

		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.CheckPermissions(ctx, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleInput := &UserPermissionsRequestInput{}

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.CheckPermissions(ctx, exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleInput := &UserPermissionsRequestInput{}

		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.CheckPermissions(ctx, exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
