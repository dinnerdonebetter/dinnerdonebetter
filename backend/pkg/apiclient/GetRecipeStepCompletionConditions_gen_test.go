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

func TestClient_GetRecipeStepCompletionConditions(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%s/steps/%s/completion_conditions"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fakes.BuildFakeID()
		recipeStepID := fakes.BuildFakeID()

		list := fakes.BuildFakeRecipeStepCompletionConditionsList()

		expected := &types.APIResponse[[]*types.RecipeStepCompletionCondition]{
			Pagination: &list.Pagination,
			Data:       list.Data,
		}

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPathFormat, recipeID, recipeStepID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.GetRecipeStepCompletionConditions(ctx, recipeID, recipeStepID, nil)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, list, actual)
	})

	T.Run("with empty recipe ID", func(t *testing.T) {
		t.Parallel()

		recipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepCompletionConditions(ctx, "", recipeStepID, nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with empty recipeStep ID", func(t *testing.T) {
		t.Parallel()

		recipeID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepCompletionConditions(ctx, recipeID, "", nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fakes.BuildFakeID()
		recipeStepID := fakes.BuildFakeID()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetRecipeStepCompletionConditions(ctx, recipeID, recipeStepID, nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fakes.BuildFakeID()
		recipeStepID := fakes.BuildFakeID()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPathFormat, recipeID, recipeStepID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetRecipeStepCompletionConditions(ctx, recipeID, recipeStepID, nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
