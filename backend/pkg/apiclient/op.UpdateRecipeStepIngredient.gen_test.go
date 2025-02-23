// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
)

func TestClient_UpdateRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%s/steps/%s/ingredients/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		recipeID := fake.BuildFakeID()
		recipeStepID := fake.BuildFakeID()
		recipeStepIngredientID := fake.BuildFakeID()

		data := &RecipeStepIngredient{}
		expected := &APIResponse[*RecipeStepIngredient]{
			Data: data,
		}

		exampleInput := &RecipeStepIngredientUpdateRequestInput{}

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, recipeID, recipeStepID, recipeStepIngredientID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateRecipeStepIngredient(ctx, recipeID, recipeStepID, recipeStepIngredientID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		recipeStepID := fake.BuildFakeID()
		recipeStepIngredientID := fake.BuildFakeID()

		exampleInput := &RecipeStepIngredientUpdateRequestInput{}

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateRecipeStepIngredient(ctx, "", recipeStepID, recipeStepIngredientID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with invalid recipeStep ID", func(t *testing.T) {
		t.Parallel()

		recipeID := fake.BuildFakeID()

		recipeStepIngredientID := fake.BuildFakeID()

		exampleInput := &RecipeStepIngredientUpdateRequestInput{}

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateRecipeStepIngredient(ctx, recipeID, "", recipeStepIngredientID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with invalid recipeStepIngredient ID", func(t *testing.T) {
		t.Parallel()

		recipeID := fake.BuildFakeID()
		recipeStepID := fake.BuildFakeID()

		exampleInput := &RecipeStepIngredientUpdateRequestInput{}

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateRecipeStepIngredient(ctx, recipeID, recipeStepID, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		recipeID := fake.BuildFakeID()
		recipeStepID := fake.BuildFakeID()
		recipeStepIngredientID := fake.BuildFakeID()

		exampleInput := &RecipeStepIngredientUpdateRequestInput{}

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateRecipeStepIngredient(ctx, recipeID, recipeStepID, recipeStepIngredientID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		recipeID := fake.BuildFakeID()
		recipeStepID := fake.BuildFakeID()
		recipeStepIngredientID := fake.BuildFakeID()

		exampleInput := &RecipeStepIngredientUpdateRequestInput{}

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, recipeID, recipeStepID, recipeStepIngredientID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateRecipeStepIngredient(ctx, recipeID, recipeStepID, recipeStepIngredientID, exampleInput)

		assert.Error(t, err)
	})
}
