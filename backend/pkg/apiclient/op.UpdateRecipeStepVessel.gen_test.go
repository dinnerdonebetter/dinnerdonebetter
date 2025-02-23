// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
)

func TestClient_UpdateRecipeStepVessel(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%s/steps/%s/vessels/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		recipeID := fake.BuildFakeID()
		recipeStepID := fake.BuildFakeID()
		recipeStepVesselID := fake.BuildFakeID()

		data := &RecipeStepVessel{}
		expected := &APIResponse[*RecipeStepVessel]{
			Data: data,
		}

		exampleInput := &RecipeStepVesselUpdateRequestInput{}

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, recipeID, recipeStepID, recipeStepVesselID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateRecipeStepVessel(ctx, recipeID, recipeStepID, recipeStepVesselID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		recipeStepID := fake.BuildFakeID()
		recipeStepVesselID := fake.BuildFakeID()

		exampleInput := &RecipeStepVesselUpdateRequestInput{}

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateRecipeStepVessel(ctx, "", recipeStepID, recipeStepVesselID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with invalid recipeStep ID", func(t *testing.T) {
		t.Parallel()

		recipeID := fake.BuildFakeID()

		recipeStepVesselID := fake.BuildFakeID()

		exampleInput := &RecipeStepVesselUpdateRequestInput{}

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateRecipeStepVessel(ctx, recipeID, "", recipeStepVesselID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with invalid recipeStepVessel ID", func(t *testing.T) {
		t.Parallel()

		recipeID := fake.BuildFakeID()
		recipeStepID := fake.BuildFakeID()

		exampleInput := &RecipeStepVesselUpdateRequestInput{}

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateRecipeStepVessel(ctx, recipeID, recipeStepID, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		recipeID := fake.BuildFakeID()
		recipeStepID := fake.BuildFakeID()
		recipeStepVesselID := fake.BuildFakeID()

		exampleInput := &RecipeStepVesselUpdateRequestInput{}

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateRecipeStepVessel(ctx, recipeID, recipeStepID, recipeStepVesselID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		recipeID := fake.BuildFakeID()
		recipeStepID := fake.BuildFakeID()
		recipeStepVesselID := fake.BuildFakeID()

		exampleInput := &RecipeStepVesselUpdateRequestInput{}

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, recipeID, recipeStepID, recipeStepVesselID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateRecipeStepVessel(ctx, recipeID, recipeStepID, recipeStepVesselID, exampleInput)

		assert.Error(t, err)
	})
}
