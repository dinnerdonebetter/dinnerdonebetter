// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_CreateValidPreparationVessel(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_preparation_vessels"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()

		data := &ValidPreparationVessel{}
		expected := &APIResponse[*ValidPreparationVessel]{
			Data: data,
		}

		exampleInput := &ValidPreparationVesselCreationRequestInput{}

		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.CreateValidPreparationVessel(ctx, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()

		exampleInput := &ValidPreparationVesselCreationRequestInput{}

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.CreateValidPreparationVessel(ctx, exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()

		exampleInput := &ValidPreparationVesselCreationRequestInput{}

		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.CreateValidPreparationVessel(ctx, exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
