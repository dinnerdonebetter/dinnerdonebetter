// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_CreateRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%s/steps/%s/ingredients"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		recipeID := fake.BuildFakeID()
		recipeStepID := fake.BuildFakeID()

		data := &RecipeStepIngredient{}
		expected := &APIResponse[*RecipeStepIngredient]{
			Data: data,
		}

		exampleInput := &RecipeStepIngredientCreationRequestInput{}

		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat, recipeID, recipeStepID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.CreateRecipeStepIngredient(ctx, recipeID, recipeStepID, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		recipeStepID := fake.BuildFakeID()

		exampleInput := &RecipeStepIngredientCreationRequestInput{}

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.CreateRecipeStepIngredient(ctx, "", recipeStepID, exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid recipeStep ID", func(t *testing.T) {
		t.Parallel()

		recipeID := fake.BuildFakeID()

		exampleInput := &RecipeStepIngredientCreationRequestInput{}

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.CreateRecipeStepIngredient(ctx, recipeID, "", exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		recipeID := fake.BuildFakeID()
		recipeStepID := fake.BuildFakeID()

		exampleInput := &RecipeStepIngredientCreationRequestInput{}

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.CreateRecipeStepIngredient(ctx, recipeID, recipeStepID, exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		recipeID := fake.BuildFakeID()
		recipeStepID := fake.BuildFakeID()

		exampleInput := &RecipeStepIngredientCreationRequestInput{}

		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat, recipeID, recipeStepID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.CreateRecipeStepIngredient(ctx, recipeID, recipeStepID, exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
