// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_UpdatePassword(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/users/password/new"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()

		data := &PasswordResetResponse{}
		expected := &APIResponse[*PasswordResetResponse]{
			Data: data,
		}

		exampleInput := &PasswordUpdateInput{}

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdatePassword(ctx, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()

		exampleInput := &PasswordUpdateInput{}

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdatePassword(ctx, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()

		exampleInput := &PasswordUpdateInput{}

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdatePassword(ctx, exampleInput)

		assert.Error(t, err)
	})
}
