// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_CreateRecipeStep(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%s/steps"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fake.BuildFakeID()

		data := fake.BuildFakeForTest[*RecipeStep](t)

		expected := &APIResponse[*RecipeStep]{
			Data: data,
		}

		exampleInput := fake.BuildFakeForTest[*RecipeStepCreationRequestInput](t)

		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat, recipeID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.CreateRecipeStep(ctx, recipeID, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := fake.BuildFakeForTest[*RecipeStepCreationRequestInput](t)

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.CreateRecipeStep(ctx, "", exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fake.BuildFakeID()

		exampleInput := fake.BuildFakeForTest[*RecipeStepCreationRequestInput](t)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.CreateRecipeStep(ctx, recipeID, exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fake.BuildFakeID()

		exampleInput := fake.BuildFakeForTest[*RecipeStepCreationRequestInput](t)

		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat, recipeID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.CreateRecipeStep(ctx, recipeID, exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
