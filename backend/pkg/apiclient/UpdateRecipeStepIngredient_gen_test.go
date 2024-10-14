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

func TestClient_UpdateRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%s/steps/%s/ingredients/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fakes.BuildFakeID()
		recipeStepID := fakes.BuildFakeID()
		recipeStepIngredientID := fakes.BuildFakeID()

		data := fakes.BuildFakeRecipeStepIngredient()
		expected := &types.APIResponse[*types.RecipeStepIngredient]{
			Data: data,
		}

		exampleInput := fakes.BuildFakeRecipeStepIngredientUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, recipeID, recipeStepID, recipeStepIngredientID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateRecipeStepIngredient(ctx, recipeID, recipeStepID, recipeStepIngredientID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		recipeStepID := fakes.BuildFakeID()
		recipeStepIngredientID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeRecipeStepIngredientUpdateRequestInput()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateRecipeStepIngredient(ctx, "", recipeStepID, recipeStepIngredientID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with invalid recipeStep ID", func(t *testing.T) {
		t.Parallel()

		recipeID := fakes.BuildFakeID()

		recipeStepIngredientID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeRecipeStepIngredientUpdateRequestInput()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateRecipeStepIngredient(ctx, recipeID, "", recipeStepIngredientID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with invalid recipeStepIngredient ID", func(t *testing.T) {
		t.Parallel()

		recipeID := fakes.BuildFakeID()
		recipeStepID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeRecipeStepIngredientUpdateRequestInput()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateRecipeStepIngredient(ctx, recipeID, recipeStepID, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fakes.BuildFakeID()
		recipeStepID := fakes.BuildFakeID()
		recipeStepIngredientID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeRecipeStepIngredientUpdateRequestInput()

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateRecipeStepIngredient(ctx, recipeID, recipeStepID, recipeStepIngredientID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fakes.BuildFakeID()
		recipeStepID := fakes.BuildFakeID()
		recipeStepIngredientID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeRecipeStepIngredientUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, recipeID, recipeStepID, recipeStepIngredientID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateRecipeStepIngredient(ctx, recipeID, recipeStepID, recipeStepIngredientID, exampleInput)

		assert.Error(t, err)
	})
}
