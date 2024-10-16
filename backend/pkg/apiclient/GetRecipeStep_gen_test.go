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

func TestClient_GetRecipeStep(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%s/steps/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fakes.BuildFakeID()
		recipeStepID := fakes.BuildFakeID()

		data := fakes.BuildFakeRecipeStep()
		expected := &types.APIResponse[*types.RecipeStep]{
			Data: data,
		}

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, recipeID, recipeStepID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.GetRecipeStep(ctx, recipeID, recipeStepID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		recipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStep(ctx, "", recipeStepID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid recipeStep ID", func(t *testing.T) {
		t.Parallel()

		recipeID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStep(ctx, recipeID, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fakes.BuildFakeID()
		recipeStepID := fakes.BuildFakeID()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetRecipeStep(ctx, recipeID, recipeStepID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fakes.BuildFakeID()
		recipeStepID := fakes.BuildFakeID()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, recipeID, recipeStepID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetRecipeStep(ctx, recipeID, recipeStepID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
