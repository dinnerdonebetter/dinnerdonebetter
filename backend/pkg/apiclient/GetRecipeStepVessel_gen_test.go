// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetRecipeStepVessel(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%s/steps/%s/vessels/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fakes.BuildFakeID()
		recipeStepID := fakes.BuildFakeID()
		recipeStepVesselID := fakes.BuildFakeID()

		data := fakes.BuildFakeRecipeStepVessel()
		expected := &types.APIResponse[*types.RecipeStepVessel]{
			Data: data,
		}

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, recipeID, recipeStepID, recipeStepVesselID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.GetRecipeStepVessel(ctx, recipeID, recipeStepID, recipeStepVesselID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		recipeStepID := fakes.BuildFakeID()
		recipeStepVesselID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepVessel(ctx, "", recipeStepID, recipeStepVesselID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid recipeStep ID", func(t *testing.T) {
		t.Parallel()

		recipeID := fakes.BuildFakeID()

		recipeStepVesselID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepVessel(ctx, recipeID, "", recipeStepVesselID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid recipeStepVessel ID", func(t *testing.T) {
		t.Parallel()

		recipeID := fakes.BuildFakeID()
		recipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepVessel(ctx, recipeID, recipeStepID, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fakes.BuildFakeID()
		recipeStepID := fakes.BuildFakeID()
		recipeStepVesselID := fakes.BuildFakeID()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetRecipeStepVessel(ctx, recipeID, recipeStepID, recipeStepVesselID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fakes.BuildFakeID()
		recipeStepID := fakes.BuildFakeID()
		recipeStepVesselID := fakes.BuildFakeID()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, recipeID, recipeStepID, recipeStepVesselID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetRecipeStepVessel(ctx, recipeID, recipeStepID, recipeStepVesselID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
