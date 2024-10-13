// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestClient_GetValidIngredientPreparationsByPreparation(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_ingredient_preparations/by_preparation/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q := fakes.BuildFakeID()
		validPreparationID := fakes.BuildFakeID()

		list := fakes.BuildFakeValidIngredientPreparationsList()

		expected := &types.APIResponse[[]*types.ValidIngredientPreparation]{
			Pagination: &list.Pagination,
			Data:       list.Data,
		}

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPathFormat, validPreparationID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.GetValidIngredientPreparationsByPreparation(ctx, q, validPreparationID, nil)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, list, actual)
	})

	T.Run("with invalid query ", func(t *testing.T) {
		t.Parallel()

		validPreparationID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetValidIngredientPreparationsByPreparation(ctx, "", validPreparationID, nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid validPreparation ID", func(t *testing.T) {
		t.Parallel()

		q := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetValidIngredientPreparationsByPreparation(ctx, q, "", nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q := fakes.BuildFakeID()
		validPreparationID := fakes.BuildFakeID()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidIngredientPreparationsByPreparation(ctx, q, validPreparationID, nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q := fakes.BuildFakeID()
		validPreparationID := fakes.BuildFakeID()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPathFormat, validPreparationID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidIngredientPreparationsByPreparation(ctx, q, validPreparationID, nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
