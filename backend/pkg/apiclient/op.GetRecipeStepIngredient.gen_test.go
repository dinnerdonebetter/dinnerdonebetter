// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetRecipeStepIngredient(T *testing.T) {
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

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, recipeID, recipeStepID, recipeStepIngredientID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.GetRecipeStepIngredient(ctx, recipeID, recipeStepID, recipeStepIngredientID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		recipeStepID := fake.BuildFakeID()
		recipeStepIngredientID := fake.BuildFakeID()

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepIngredient(ctx, "", recipeStepID, recipeStepIngredientID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid recipeStep ID", func(t *testing.T) {
		t.Parallel()

		recipeID := fake.BuildFakeID()

		recipeStepIngredientID := fake.BuildFakeID()

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepIngredient(ctx, recipeID, "", recipeStepIngredientID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid recipeStepIngredient ID", func(t *testing.T) {
		t.Parallel()

		recipeID := fake.BuildFakeID()
		recipeStepID := fake.BuildFakeID()

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepIngredient(ctx, recipeID, recipeStepID, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		recipeID := fake.BuildFakeID()
		recipeStepID := fake.BuildFakeID()
		recipeStepIngredientID := fake.BuildFakeID()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetRecipeStepIngredient(ctx, recipeID, recipeStepID, recipeStepIngredientID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		recipeID := fake.BuildFakeID()
		recipeStepID := fake.BuildFakeID()
		recipeStepIngredientID := fake.BuildFakeID()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, recipeID, recipeStepID, recipeStepIngredientID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetRecipeStepIngredient(ctx, recipeID, recipeStepID, recipeStepIngredientID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
