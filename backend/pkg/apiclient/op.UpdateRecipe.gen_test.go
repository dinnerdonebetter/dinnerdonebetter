// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
)

func TestClient_UpdateRecipe(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		recipeID := fake.BuildFakeID()

		data := &Recipe{}
		expected := &APIResponse[*Recipe]{
			Data: data,
		}

		exampleInput := &RecipeUpdateRequestInput{}

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, recipeID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateRecipe(ctx, recipeID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := &RecipeUpdateRequestInput{}

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateRecipe(ctx, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		recipeID := fake.BuildFakeID()

		exampleInput := &RecipeUpdateRequestInput{}

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateRecipe(ctx, recipeID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		recipeID := fake.BuildFakeID()

		exampleInput := &RecipeUpdateRequestInput{}

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, recipeID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateRecipe(ctx, recipeID, exampleInput)

		assert.Error(t, err)
	})
}
