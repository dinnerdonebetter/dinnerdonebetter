// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
)

func TestClient_UpdateValidMeasurementUnit(T *testing.T) {
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

		exampleInput := &ValidMeasurementUnitUpdateRequestInput{}

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, validMeasurementUnitID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateValidMeasurementUnit(ctx, validMeasurementUnitID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid validMeasurementUnit ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := &ValidMeasurementUnitUpdateRequestInput{}

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateValidMeasurementUnit(ctx, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		validMeasurementUnitID := fake.BuildFakeID()

		exampleInput := &ValidMeasurementUnitUpdateRequestInput{}

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateValidMeasurementUnit(ctx, validMeasurementUnitID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		validMeasurementUnitID := fake.BuildFakeID()

		exampleInput := &ValidMeasurementUnitUpdateRequestInput{}

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, validMeasurementUnitID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateValidMeasurementUnit(ctx, validMeasurementUnitID, exampleInput)

		assert.Error(t, err)
	})
}
