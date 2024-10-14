// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestClient_UpdateValidMeasurementUnit(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_measurement_units/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validMeasurementUnitID := fakes.BuildFakeID()

		data := fakes.BuildFakeValidMeasurementUnit()
		expected := &types.APIResponse[*types.ValidMeasurementUnit]{
			Data: data,
		}

		exampleInput := fakes.BuildFakeValidMeasurementUnitUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, validMeasurementUnitID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateValidMeasurementUnit(ctx, validMeasurementUnitID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid validMeasurementUnit ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := fakes.BuildFakeValidMeasurementUnitUpdateRequestInput()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateValidMeasurementUnit(ctx, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validMeasurementUnitID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeValidMeasurementUnitUpdateRequestInput()

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateValidMeasurementUnit(ctx, validMeasurementUnitID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validMeasurementUnitID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeValidMeasurementUnitUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, validMeasurementUnitID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateValidMeasurementUnit(ctx, validMeasurementUnitID, exampleInput)

		assert.Error(t, err)
	})
}
