// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
)

func TestClient_UpdateValidMeasurementUnitConversion(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_measurement_conversions/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validMeasurementUnitConversionID := fake.BuildFakeID()

		data := fake.BuildFakeForTest[*ValidMeasurementUnitConversion](t)

		expected := &APIResponse[*ValidMeasurementUnitConversion]{
			Data: data,
		}

		exampleInput := fake.BuildFakeForTest[*ValidMeasurementUnitConversionUpdateRequestInput](t)

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, validMeasurementUnitConversionID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateValidMeasurementUnitConversion(ctx, validMeasurementUnitConversionID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid validMeasurementUnitConversion ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := fake.BuildFakeForTest[*ValidMeasurementUnitConversionUpdateRequestInput](t)

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateValidMeasurementUnitConversion(ctx, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validMeasurementUnitConversionID := fake.BuildFakeID()

		exampleInput := fake.BuildFakeForTest[*ValidMeasurementUnitConversionUpdateRequestInput](t)

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateValidMeasurementUnitConversion(ctx, validMeasurementUnitConversionID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validMeasurementUnitConversionID := fake.BuildFakeID()

		exampleInput := fake.BuildFakeForTest[*ValidMeasurementUnitConversionUpdateRequestInput](t)

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, validMeasurementUnitConversionID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateValidMeasurementUnitConversion(ctx, validMeasurementUnitConversionID, exampleInput)

		assert.Error(t, err)
	})
}
