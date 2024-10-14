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

func TestClient_UpdateRecipeRating(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%s/ratings/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fakes.BuildFakeID()
		recipeRatingID := fakes.BuildFakeID()

		data := fakes.BuildFakeRecipeRating()
		expected := &types.APIResponse[*types.RecipeRating]{
			Data: data,
		}

		exampleInput := fakes.BuildFakeRecipeRatingUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, recipeID, recipeRatingID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateRecipeRating(ctx, recipeID, recipeRatingID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		recipeRatingID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeRecipeRatingUpdateRequestInput()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateRecipeRating(ctx, "", recipeRatingID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with invalid recipeRating ID", func(t *testing.T) {
		t.Parallel()

		recipeID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeRecipeRatingUpdateRequestInput()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateRecipeRating(ctx, recipeID, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fakes.BuildFakeID()
		recipeRatingID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeRecipeRatingUpdateRequestInput()

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateRecipeRating(ctx, recipeID, recipeRatingID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fakes.BuildFakeID()
		recipeRatingID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeRecipeRatingUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, recipeID, recipeRatingID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateRecipeRating(ctx, recipeID, recipeRatingID, exampleInput)

		assert.Error(t, err)
	})
}
