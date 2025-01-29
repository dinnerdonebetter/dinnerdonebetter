// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
)

func TestClient_UpdateValidVessel(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_vessels/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validVesselID := fake.BuildFakeID()

		data := fake.BuildFakeForTest[*ValidVessel](t)

		expected := &APIResponse[*ValidVessel]{
			Data: data,
		}

		exampleInput := fake.BuildFakeForTest[*ValidVesselUpdateRequestInput](t)

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, validVesselID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateValidVessel(ctx, validVesselID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid validVessel ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := fake.BuildFakeForTest[*ValidVesselUpdateRequestInput](t)

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateValidVessel(ctx, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validVesselID := fake.BuildFakeID()

		exampleInput := fake.BuildFakeForTest[*ValidVesselUpdateRequestInput](t)

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateValidVessel(ctx, validVesselID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validVesselID := fake.BuildFakeID()

		exampleInput := fake.BuildFakeForTest[*ValidVesselUpdateRequestInput](t)

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, validVesselID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateValidVessel(ctx, validVesselID, exampleInput)

		assert.Error(t, err)
	})
}
