// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_CreateRecipeRating(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%s/ratings"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		recipeID := fake.BuildFakeID()

		data := &RecipeRating{}
		expected := &APIResponse[*RecipeRating]{
			Data: data,
		}

		exampleInput := &RecipeRatingCreationRequestInput{}

		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat, recipeID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.CreateRecipeRating(ctx, recipeID, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := &RecipeRatingCreationRequestInput{}

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.CreateRecipeRating(ctx, "", exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		recipeID := fake.BuildFakeID()

		exampleInput := &RecipeRatingCreationRequestInput{}

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.CreateRecipeRating(ctx, recipeID, exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		recipeID := fake.BuildFakeID()

		exampleInput := &RecipeRatingCreationRequestInput{}

		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat, recipeID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.CreateRecipeRating(ctx, recipeID, exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
