// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestClient_UpdatePassword(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/users/password/new"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		data := fakes.BuildFakePasswordResetResponse()
		expected := &types.APIResponse[*types.PasswordResetResponse]{
			Data: data,
		}

		exampleInput := fakes.BuildFakePasswordUpdateInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdatePassword(ctx, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleInput := fakes.BuildFakePasswordUpdateInput()

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdatePassword(ctx, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleInput := fakes.BuildFakePasswordUpdateInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdatePassword(ctx, exampleInput)

		assert.Error(t, err)
	})
}
