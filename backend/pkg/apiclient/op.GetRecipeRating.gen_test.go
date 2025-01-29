// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/lib/internal/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetRecipeRating(T *testing.T) {
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

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, recipeID, recipeRatingID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.GetRecipeRating(ctx, recipeID, recipeRatingID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		recipeRatingID := fake.BuildFakeID()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeRating(ctx, "", recipeRatingID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid recipeRating ID", func(t *testing.T) {
		t.Parallel()

		recipeID := fake.BuildFakeID()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeRating(ctx, recipeID, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fake.BuildFakeID()
		recipeRatingID := fake.BuildFakeID()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetRecipeRating(ctx, recipeID, recipeRatingID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fake.BuildFakeID()
		recipeRatingID := fake.BuildFakeID()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, recipeID, recipeRatingID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetRecipeRating(ctx, recipeID, recipeRatingID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
