// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
)

func TestClient_UpdateRecipeStepProduct(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%s/steps/%s/products/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fake.BuildFakeID()
		recipeStepID := fake.BuildFakeID()
		recipeStepProductID := fake.BuildFakeID()

		data := fake.BuildFakeForTest[*RecipeStepProduct](t)

		expected := &APIResponse[*RecipeStepProduct]{
			Data: data,
		}

		exampleInput := fake.BuildFakeForTest[*RecipeStepProductUpdateRequestInput](t)

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, recipeID, recipeStepID, recipeStepProductID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateRecipeStepProduct(ctx, recipeID, recipeStepID, recipeStepProductID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		recipeStepID := fake.BuildFakeID()
		recipeStepProductID := fake.BuildFakeID()

		exampleInput := fake.BuildFakeForTest[*RecipeStepProductUpdateRequestInput](t)

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateRecipeStepProduct(ctx, "", recipeStepID, recipeStepProductID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with invalid recipeStep ID", func(t *testing.T) {
		t.Parallel()

		recipeID := fake.BuildFakeID()

		recipeStepProductID := fake.BuildFakeID()

		exampleInput := fake.BuildFakeForTest[*RecipeStepProductUpdateRequestInput](t)

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateRecipeStepProduct(ctx, recipeID, "", recipeStepProductID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with invalid recipeStepProduct ID", func(t *testing.T) {
		t.Parallel()

		recipeID := fake.BuildFakeID()
		recipeStepID := fake.BuildFakeID()

		exampleInput := fake.BuildFakeForTest[*RecipeStepProductUpdateRequestInput](t)

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateRecipeStepProduct(ctx, recipeID, recipeStepID, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fake.BuildFakeID()
		recipeStepID := fake.BuildFakeID()
		recipeStepProductID := fake.BuildFakeID()

		exampleInput := fake.BuildFakeForTest[*RecipeStepProductUpdateRequestInput](t)

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateRecipeStepProduct(ctx, recipeID, recipeStepID, recipeStepProductID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fake.BuildFakeID()
		recipeStepID := fake.BuildFakeID()
		recipeStepProductID := fake.BuildFakeID()

		exampleInput := fake.BuildFakeForTest[*RecipeStepProductUpdateRequestInput](t)

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, recipeID, recipeStepID, recipeStepProductID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateRecipeStepProduct(ctx, recipeID, recipeStepID, recipeStepProductID, exampleInput)

		assert.Error(t, err)
	})
}
