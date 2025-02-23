// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_UpdateUserUsername(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/users/username"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()

		data := &User{}
		expected := &APIResponse[*User]{
			Data: data,
		}

		exampleInput := &UsernameUpdateInput{}

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateUserUsername(ctx, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()

		exampleInput := &UsernameUpdateInput{}

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateUserUsername(ctx, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()

		exampleInput := &UsernameUpdateInput{}

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateUserUsername(ctx, exampleInput)

		assert.Error(t, err)
	})
}
