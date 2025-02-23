// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetValidMeasurementUnitConversion(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_measurement_conversions/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		validMeasurementUnitConversionID := fake.BuildFakeID()

		data := &ValidMeasurementUnitConversion{}
		expected := &APIResponse[*ValidMeasurementUnitConversion]{
			Data: data,
		}

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, validMeasurementUnitConversionID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.GetValidMeasurementUnitConversion(ctx, validMeasurementUnitConversionID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with invalid validMeasurementUnitConversion ID", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetValidMeasurementUnitConversion(ctx, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		validMeasurementUnitConversionID := fake.BuildFakeID()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidMeasurementUnitConversion(ctx, validMeasurementUnitConversionID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		validMeasurementUnitConversionID := fake.BuildFakeID()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, validMeasurementUnitConversionID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidMeasurementUnitConversion(ctx, validMeasurementUnitConversionID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
