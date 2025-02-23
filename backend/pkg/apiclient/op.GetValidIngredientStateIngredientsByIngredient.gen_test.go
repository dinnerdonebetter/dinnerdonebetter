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

func TestClient_GetValidIngredientStateIngredientsByIngredient(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_ingredient_state_ingredients/by_ingredient/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		validIngredientID := fake.BuildFakeID()

		list := []*ValidIngredientStateIngredient{}
		exampleResponse := &APIResponse[[]*ValidIngredientStateIngredient]{
			Pagination: fake.BuildFakeForTest[*Pagination](t),
			Data:       list,
		}
		expected := &QueryFilteredResult[ValidIngredientStateIngredient]{
			Pagination: *exampleResponse.Pagination,
			Data:       list,
		}

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPathFormat, validIngredientID)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleResponse)
		actual, err := c.GetValidIngredientStateIngredientsByIngredient(ctx, validIngredientID, nil)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	T.Run("with empty validIngredient ID", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetValidIngredientStateIngredientsByIngredient(ctx, "", nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validIngredientID := fake.BuildFakeID()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidIngredientStateIngredientsByIngredient(ctx, validIngredientID, nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validIngredientID := fake.BuildFakeID()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPathFormat, validIngredientID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidIngredientStateIngredientsByIngredient(ctx, validIngredientID, nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
