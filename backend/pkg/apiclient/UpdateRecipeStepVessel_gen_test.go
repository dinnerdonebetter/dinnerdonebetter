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

func TestClient_UpdateRecipeStepVessel(T *testing.T) {
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

		exampleInput := fakes.BuildFakeRecipeStepVesselUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, recipeID, recipeStepID, recipeStepVesselID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateRecipeStepVessel(ctx, recipeID, recipeStepID, recipeStepVesselID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		recipeStepID := fakes.BuildFakeID()
		recipeStepVesselID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeRecipeStepVesselUpdateRequestInput()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateRecipeStepVessel(ctx, "", recipeStepID, recipeStepVesselID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with invalid recipeStep ID", func(t *testing.T) {
		t.Parallel()

		recipeID := fakes.BuildFakeID()

		recipeStepVesselID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeRecipeStepVesselUpdateRequestInput()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateRecipeStepVessel(ctx, recipeID, "", recipeStepVesselID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with invalid recipeStepVessel ID", func(t *testing.T) {
		t.Parallel()

		recipeID := fakes.BuildFakeID()
		recipeStepID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeRecipeStepVesselUpdateRequestInput()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateRecipeStepVessel(ctx, recipeID, recipeStepID, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fakes.BuildFakeID()
		recipeStepID := fakes.BuildFakeID()
		recipeStepVesselID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeRecipeStepVesselUpdateRequestInput()

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateRecipeStepVessel(ctx, recipeID, recipeStepID, recipeStepVesselID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fakes.BuildFakeID()
		recipeStepID := fakes.BuildFakeID()
		recipeStepVesselID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeRecipeStepVesselUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, recipeID, recipeStepID, recipeStepVesselID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateRecipeStepVessel(ctx, recipeID, recipeStepID, recipeStepVesselID, exampleInput)

		assert.Error(t, err)
	})
}
