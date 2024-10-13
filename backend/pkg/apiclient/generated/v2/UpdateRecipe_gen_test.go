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

func TestClient_UpdateRecipe(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fakes.BuildFakeID()

		data := fakes.BuildFakeRecipe()
		expected := &types.APIResponse[*types.Recipe]{
			Data: data,
		}

		exampleInput := fakes.BuildFakeRecipeUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, recipeID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateRecipe(ctx, recipeID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := fakes.BuildFakeRecipeUpdateRequestInput()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateRecipe(ctx, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeRecipeUpdateRequestInput()

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateRecipe(ctx, recipeID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeRecipeUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, recipeID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateRecipe(ctx, recipeID, exampleInput)

		assert.Error(t, err)
	})
}
