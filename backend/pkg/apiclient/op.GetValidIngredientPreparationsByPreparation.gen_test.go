// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetValidIngredientPreparationsByPreparation(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_ingredient_preparations/by_preparation/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		validPreparationID := fake.BuildFakeID()

		list := []*ValidIngredientPreparation{}
		exampleResponse := &APIResponse[[]*ValidIngredientPreparation]{
			Pagination: fake.BuildFakeForTest[*Pagination](t),
			Data:       list,
		}
		expected := &QueryFilteredResult[ValidIngredientPreparation]{
			Pagination: *exampleResponse.Pagination,
			Data:       list,
		}

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPathFormat, validPreparationID)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleResponse)
		actual, err := c.GetValidIngredientPreparationsByPreparation(ctx, validPreparationID, nil)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	T.Run("with empty validPreparation ID", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetValidIngredientPreparationsByPreparation(ctx, "", nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		validPreparationID := fake.BuildFakeID()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidIngredientPreparationsByPreparation(ctx, validPreparationID, nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		validPreparationID := fake.BuildFakeID()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPathFormat, validPreparationID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidIngredientPreparationsByPreparation(ctx, validPreparationID, nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
