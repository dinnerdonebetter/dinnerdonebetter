// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
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

		ctx := t.Context()
		recipeID := fake.BuildFakeID()
		recipeRatingID := fake.BuildFakeID()

		data := &RecipeRating{}
		expected := &APIResponse[*RecipeRating]{
			Data: data,
		}

		exampleInput := &RecipeRatingUpdateRequestInput{}

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, recipeID, recipeRatingID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateRecipeRating(ctx, recipeID, recipeRatingID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		recipeRatingID := fake.BuildFakeID()

		exampleInput := &RecipeRatingUpdateRequestInput{}

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateRecipeRating(ctx, "", recipeRatingID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with invalid recipeRating ID", func(t *testing.T) {
		t.Parallel()

		recipeID := fake.BuildFakeID()

		exampleInput := &RecipeRatingUpdateRequestInput{}

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateRecipeRating(ctx, recipeID, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		recipeID := fake.BuildFakeID()
		recipeRatingID := fake.BuildFakeID()

		exampleInput := &RecipeRatingUpdateRequestInput{}

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateRecipeRating(ctx, recipeID, recipeRatingID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		recipeID := fake.BuildFakeID()
		recipeRatingID := fake.BuildFakeID()

		exampleInput := &RecipeRatingUpdateRequestInput{}

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, recipeID, recipeRatingID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateRecipeRating(ctx, recipeID, recipeRatingID, exampleInput)

		assert.Error(t, err)
	})
}
