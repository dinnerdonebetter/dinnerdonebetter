// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
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

		ctx := context.Background()
		q := fake.BuildFakeID()
		validPreparationID := fake.BuildFakeID()

		list := fake.BuildFakeForTest[[]*ValidIngredient](t)

		expected := &APIResponse[[]*ValidIngredient]{
			Pagination: fake.BuildFakeForTest[*Pagination](t),
			Data:       list,
		}

		spec := newRequestSpec(true, http.MethodGet, fmt.Sprintf("limit=50&page=1&q=%s&sortBy=asc", q), expectedPathFormat, validPreparationID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.SearchValidIngredientsByPreparation(ctx, q, validPreparationID, nil)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, list, actual)
	})

	T.Run("with empty validPreparation ID", func(t *testing.T) {
		t.Parallel()

		q := fake.BuildFakeID()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.SearchValidIngredientsByPreparation(ctx, q, "", nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q := fake.BuildFakeID()
		validPreparationID := fake.BuildFakeID()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.SearchValidIngredientsByPreparation(ctx, q, validPreparationID, nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q := fake.BuildFakeID()
		validPreparationID := fake.BuildFakeID()

		spec := newRequestSpec(true, http.MethodGet, fmt.Sprintf("limit=50&page=1&q=%s&sortBy=asc", q), expectedPathFormat, validPreparationID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.SearchValidIngredientsByPreparation(ctx, q, validPreparationID, nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
