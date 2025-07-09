// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_RequestPasswordResetToken(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/users/password/reset"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		data := &PasswordResetToken{}
		expected := &APIResponse[*PasswordResetToken]{
			Data: data,
		}

		exampleInput := &PasswordResetTokenCreationRequestInput{}

		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.RequestPasswordResetToken(ctx, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleInput := &PasswordResetTokenCreationRequestInput{}

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.RequestPasswordResetToken(ctx, exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleInput := &PasswordResetTokenCreationRequestInput{}

		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.RequestPasswordResetToken(ctx, exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
