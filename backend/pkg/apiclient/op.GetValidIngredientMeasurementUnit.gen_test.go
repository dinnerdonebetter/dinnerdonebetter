// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetValidIngredientMeasurementUnit(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_ingredient_measurement_units/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		validIngredientMeasurementUnitID := fake.BuildFakeID()

		data := &ValidIngredientMeasurementUnit{}
		expected := &APIResponse[*ValidIngredientMeasurementUnit]{
			Data: data,
		}

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, validIngredientMeasurementUnitID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.GetValidIngredientMeasurementUnit(ctx, validIngredientMeasurementUnitID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with invalid validIngredientMeasurementUnit ID", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetValidIngredientMeasurementUnit(ctx, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		validIngredientMeasurementUnitID := fake.BuildFakeID()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidIngredientMeasurementUnit(ctx, validIngredientMeasurementUnitID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		validIngredientMeasurementUnitID := fake.BuildFakeID()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, validIngredientMeasurementUnitID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidIngredientMeasurementUnit(ctx, validIngredientMeasurementUnitID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
