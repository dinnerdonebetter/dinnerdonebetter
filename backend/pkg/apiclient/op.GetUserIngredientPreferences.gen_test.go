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

func TestClient_GetUserIngredientPreferences(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/user_ingredient_preferences"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()

		list := []*UserIngredientPreference{}
		exampleResponse := &APIResponse[[]*UserIngredientPreference]{
			Pagination: fake.BuildFakeForTest[*Pagination](t),
			Data:       list,
		}
		expected := &QueryFilteredResult[UserIngredientPreference]{
			Pagination: *exampleResponse.Pagination,
			Data:       list,
		}

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPathFormat)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleResponse)
		actual, err := c.GetUserIngredientPreferences(ctx, nil)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetUserIngredientPreferences(ctx, nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPathFormat)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetUserIngredientPreferences(ctx, nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
