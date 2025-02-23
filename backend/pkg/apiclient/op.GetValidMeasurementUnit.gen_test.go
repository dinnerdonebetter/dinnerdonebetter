// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetValidMeasurementUnit(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_measurement_units/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		validMeasurementUnitID := fake.BuildFakeID()

		data := &ValidMeasurementUnit{}
		expected := &APIResponse[*ValidMeasurementUnit]{
			Data: data,
		}

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, validMeasurementUnitID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.GetValidMeasurementUnit(ctx, validMeasurementUnitID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with invalid validMeasurementUnit ID", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetValidMeasurementUnit(ctx, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		validMeasurementUnitID := fake.BuildFakeID()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidMeasurementUnit(ctx, validMeasurementUnitID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		validMeasurementUnitID := fake.BuildFakeID()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, validMeasurementUnitID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidMeasurementUnit(ctx, validMeasurementUnitID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
