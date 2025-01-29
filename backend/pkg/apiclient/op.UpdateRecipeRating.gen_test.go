// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
)

func TestClient_UpdateRecipeRating(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%s/ratings/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fake.BuildFakeID()
		recipeRatingID := fake.BuildFakeID()

		data := fake.BuildFakeForTest[*RecipeRating](t)

		expected := &APIResponse[*RecipeRating]{
			Data: data,
		}

		exampleInput := fake.BuildFakeForTest[*RecipeRatingUpdateRequestInput](t)

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, recipeID, recipeRatingID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateRecipeRating(ctx, recipeID, recipeRatingID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		recipeRatingID := fake.BuildFakeID()

		exampleInput := fake.BuildFakeForTest[*RecipeRatingUpdateRequestInput](t)

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateRecipeRating(ctx, "", recipeRatingID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with invalid recipeRating ID", func(t *testing.T) {
		t.Parallel()

		recipeID := fake.BuildFakeID()

		exampleInput := fake.BuildFakeForTest[*RecipeRatingUpdateRequestInput](t)

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateRecipeRating(ctx, recipeID, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fake.BuildFakeID()
		recipeRatingID := fake.BuildFakeID()

		exampleInput := fake.BuildFakeForTest[*RecipeRatingUpdateRequestInput](t)

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateRecipeRating(ctx, recipeID, recipeRatingID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fake.BuildFakeID()
		recipeRatingID := fake.BuildFakeID()

		exampleInput := fake.BuildFakeForTest[*RecipeRatingUpdateRequestInput](t)

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, recipeID, recipeRatingID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateRecipeRating(ctx, recipeID, recipeRatingID, exampleInput)

		assert.Error(t, err)
	})
}
