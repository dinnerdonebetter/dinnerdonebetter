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

func TestClient_UpdateValidPreparationVessel(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_preparation_vessels/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validPreparationVesselID := fakes.BuildFakeID()

		data := fakes.BuildFakeValidPreparationVessel()
		expected := &types.APIResponse[*types.ValidPreparationVessel]{
			Data: data,
		}

		exampleInput := fakes.BuildFakeValidPreparationVesselUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, validPreparationVesselID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateValidPreparationVessel(ctx, validPreparationVesselID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid validPreparationVessel ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := fakes.BuildFakeValidPreparationVesselUpdateRequestInput()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateValidPreparationVessel(ctx, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validPreparationVesselID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeValidPreparationVesselUpdateRequestInput()

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateValidPreparationVessel(ctx, validPreparationVesselID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validPreparationVesselID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeValidPreparationVesselUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, validPreparationVesselID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateValidPreparationVessel(ctx, validPreparationVesselID, exampleInput)

		assert.Error(t, err)
	})
}
