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

func TestClient_GetValidIngredientMeasurementUnits(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_ingredient_measurement_units"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		list := fake.BuildFakeForTest[[]*ValidIngredientMeasurementUnit](t)

		expected := &APIResponse[[]*ValidIngredientMeasurementUnit]{
			Pagination: fake.BuildFakeForTest[*Pagination](t),
			Data:       list,
		}

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPathFormat)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.GetValidIngredientMeasurementUnits(ctx, nil)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, list, actual)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidIngredientMeasurementUnits(ctx, nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPathFormat)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidIngredientMeasurementUnits(ctx, nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
