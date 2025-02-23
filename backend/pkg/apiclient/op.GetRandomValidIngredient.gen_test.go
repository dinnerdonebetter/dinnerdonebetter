// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetRandomValidIngredient(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_ingredients/random"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()

		data := &ValidIngredient{}
		expected := &APIResponse[*ValidIngredient]{
			Data: data,
		}

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.GetRandomValidIngredient(ctx)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetRandomValidIngredient(ctx)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetRandomValidIngredient(ctx)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
