// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
)

func TestClient_UpdateValidPreparationVessel(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_preparation_vessels/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validPreparationVesselID := fake.BuildFakeID()

		data := fake.BuildFakeForTest[*ValidPreparationVessel](t)

		expected := &APIResponse[*ValidPreparationVessel]{
			Data: data,
		}

		exampleInput := fake.BuildFakeForTest[*ValidPreparationVesselUpdateRequestInput](t)

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, validPreparationVesselID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateValidPreparationVessel(ctx, validPreparationVesselID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid validPreparationVessel ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := fake.BuildFakeForTest[*ValidPreparationVesselUpdateRequestInput](t)

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateValidPreparationVessel(ctx, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validPreparationVesselID := fake.BuildFakeID()

		exampleInput := fake.BuildFakeForTest[*ValidPreparationVesselUpdateRequestInput](t)

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateValidPreparationVessel(ctx, validPreparationVesselID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validPreparationVesselID := fake.BuildFakeID()

		exampleInput := fake.BuildFakeForTest[*ValidPreparationVesselUpdateRequestInput](t)

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, validPreparationVesselID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateValidPreparationVessel(ctx, validPreparationVesselID, exampleInput)

		assert.Error(t, err)
	})
}
