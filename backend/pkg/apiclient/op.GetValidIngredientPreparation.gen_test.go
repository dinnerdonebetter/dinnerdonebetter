// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_ingredient_preparations/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		validIngredientPreparationID := fake.BuildFakeID()

		data := &ValidIngredientPreparation{}
		expected := &APIResponse[*ValidIngredientPreparation]{
			Data: data,
		}

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, validIngredientPreparationID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.GetValidIngredientPreparation(ctx, validIngredientPreparationID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with invalid validIngredientPreparation ID", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetValidIngredientPreparation(ctx, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		validIngredientPreparationID := fake.BuildFakeID()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidIngredientPreparation(ctx, validIngredientPreparationID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		validIngredientPreparationID := fake.BuildFakeID()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, validIngredientPreparationID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidIngredientPreparation(ctx, validIngredientPreparationID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
