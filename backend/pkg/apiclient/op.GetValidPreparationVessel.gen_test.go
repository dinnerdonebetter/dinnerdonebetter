// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetValidPreparationVessel(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_preparation_vessels/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		validPreparationVesselID := fake.BuildFakeID()

		data := &ValidPreparationVessel{}
		expected := &APIResponse[*ValidPreparationVessel]{
			Data: data,
		}

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, validPreparationVesselID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.GetValidPreparationVessel(ctx, validPreparationVesselID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with invalid validPreparationVessel ID", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetValidPreparationVessel(ctx, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		validPreparationVesselID := fake.BuildFakeID()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidPreparationVessel(ctx, validPreparationVesselID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		validPreparationVesselID := fake.BuildFakeID()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, validPreparationVesselID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidPreparationVessel(ctx, validPreparationVesselID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
