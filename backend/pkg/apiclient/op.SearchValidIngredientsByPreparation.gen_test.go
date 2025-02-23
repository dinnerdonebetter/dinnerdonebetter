// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_SearchValidIngredientsByPreparation(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_ingredients/by_preparation/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		q := fake.BuildFakeID()
		validPreparationID := fake.BuildFakeID()

		list := []*ValidIngredient{}
		exampleResponse := &APIResponse[[]*ValidIngredient]{
			Pagination: fake.BuildFakeForTest[*Pagination](t),
			Data:       list,
		}
		expected := &QueryFilteredResult[ValidIngredient]{
			Pagination: *exampleResponse.Pagination,
			Data:       list,
		}

		spec := newRequestSpec(true, http.MethodGet, fmt.Sprintf("limit=50&page=1&q=%s&sortBy=asc", q), expectedPathFormat, validPreparationID)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleResponse)
		actual, err := c.SearchValidIngredientsByPreparation(ctx, q, validPreparationID, nil)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	T.Run("with empty validPreparation ID", func(t *testing.T) {
		t.Parallel()

		q := fake.BuildFakeID()

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.SearchValidIngredientsByPreparation(ctx, q, "", nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		q := fake.BuildFakeID()
		validPreparationID := fake.BuildFakeID()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.SearchValidIngredientsByPreparation(ctx, q, validPreparationID, nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		q := fake.BuildFakeID()
		validPreparationID := fake.BuildFakeID()

		spec := newRequestSpec(true, http.MethodGet, fmt.Sprintf("limit=50&page=1&q=%s&sortBy=asc", q), expectedPathFormat, validPreparationID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.SearchValidIngredientsByPreparation(ctx, q, validPreparationID, nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
