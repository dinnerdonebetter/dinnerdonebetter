// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
)

func TestClient_UpdateValidIngredientMeasurementUnit(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_ingredient_measurement_units/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validIngredientMeasurementUnitID := fakes.BuildFakeID()

		data := fakes.BuildFakeValidIngredientMeasurementUnit()
		expected := &types.APIResponse[*types.ValidIngredientMeasurementUnit]{
			Data: data,
		}

		exampleInput := fakes.BuildFakeValidIngredientMeasurementUnitUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, validIngredientMeasurementUnitID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateValidIngredientMeasurementUnit(ctx, validIngredientMeasurementUnitID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid validIngredientMeasurementUnit ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := fakes.BuildFakeValidIngredientMeasurementUnitUpdateRequestInput()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateValidIngredientMeasurementUnit(ctx, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validIngredientMeasurementUnitID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeValidIngredientMeasurementUnitUpdateRequestInput()

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateValidIngredientMeasurementUnit(ctx, validIngredientMeasurementUnitID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validIngredientMeasurementUnitID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeValidIngredientMeasurementUnitUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, validIngredientMeasurementUnitID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateValidIngredientMeasurementUnit(ctx, validIngredientMeasurementUnitID, exampleInput)

		assert.Error(t, err)
	})
}
