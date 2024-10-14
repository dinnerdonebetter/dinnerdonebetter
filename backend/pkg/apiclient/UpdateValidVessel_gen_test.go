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

func TestClient_UpdateValidVessel(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_vessels/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validVesselID := fakes.BuildFakeID()

		data := fakes.BuildFakeValidVessel()
		expected := &types.APIResponse[*types.ValidVessel]{
			Data: data,
		}

		exampleInput := fakes.BuildFakeValidVesselUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, validVesselID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateValidVessel(ctx, validVesselID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid validVessel ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := fakes.BuildFakeValidVesselUpdateRequestInput()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateValidVessel(ctx, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validVesselID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeValidVesselUpdateRequestInput()

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateValidVessel(ctx, validVesselID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validVesselID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeValidVesselUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, validVesselID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateValidVessel(ctx, validVesselID, exampleInput)

		assert.Error(t, err)
	})
}
