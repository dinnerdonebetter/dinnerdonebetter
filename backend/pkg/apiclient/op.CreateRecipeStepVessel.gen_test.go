// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_CreateRecipeStepVessel(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%s/steps/%s/vessels"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		recipeID := fake.BuildFakeID()
		recipeStepID := fake.BuildFakeID()

		data := &RecipeStepVessel{}
		expected := &APIResponse[*RecipeStepVessel]{
			Data: data,
		}

		exampleInput := &RecipeStepVesselCreationRequestInput{}

		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat, recipeID, recipeStepID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.CreateRecipeStepVessel(ctx, recipeID, recipeStepID, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		recipeStepID := fake.BuildFakeID()

		exampleInput := &RecipeStepVesselCreationRequestInput{}

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.CreateRecipeStepVessel(ctx, "", recipeStepID, exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid recipeStep ID", func(t *testing.T) {
		t.Parallel()

		recipeID := fake.BuildFakeID()

		exampleInput := &RecipeStepVesselCreationRequestInput{}

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.CreateRecipeStepVessel(ctx, recipeID, "", exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		recipeID := fake.BuildFakeID()
		recipeStepID := fake.BuildFakeID()

		exampleInput := &RecipeStepVesselCreationRequestInput{}

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.CreateRecipeStepVessel(ctx, recipeID, recipeStepID, exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		recipeID := fake.BuildFakeID()
		recipeStepID := fake.BuildFakeID()

		exampleInput := &RecipeStepVesselCreationRequestInput{}

		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat, recipeID, recipeStepID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.CreateRecipeStepVessel(ctx, recipeID, recipeStepID, exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
