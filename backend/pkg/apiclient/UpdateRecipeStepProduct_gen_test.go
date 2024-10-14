// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
)

func TestClient_UpdateRecipeStepProduct(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%s/steps/%s/products/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fakes.BuildFakeID()
		recipeStepID := fakes.BuildFakeID()
		recipeStepProductID := fakes.BuildFakeID()

		data := fakes.BuildFakeRecipeStepProduct()
		expected := &types.APIResponse[*types.RecipeStepProduct]{
			Data: data,
		}

		exampleInput := fakes.BuildFakeRecipeStepProductUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, recipeID, recipeStepID, recipeStepProductID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateRecipeStepProduct(ctx, recipeID, recipeStepID, recipeStepProductID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		recipeStepID := fakes.BuildFakeID()
		recipeStepProductID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeRecipeStepProductUpdateRequestInput()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateRecipeStepProduct(ctx, "", recipeStepID, recipeStepProductID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with invalid recipeStep ID", func(t *testing.T) {
		t.Parallel()

		recipeID := fakes.BuildFakeID()

		recipeStepProductID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeRecipeStepProductUpdateRequestInput()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateRecipeStepProduct(ctx, recipeID, "", recipeStepProductID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with invalid recipeStepProduct ID", func(t *testing.T) {
		t.Parallel()

		recipeID := fakes.BuildFakeID()
		recipeStepID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeRecipeStepProductUpdateRequestInput()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateRecipeStepProduct(ctx, recipeID, recipeStepID, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fakes.BuildFakeID()
		recipeStepID := fakes.BuildFakeID()
		recipeStepProductID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeRecipeStepProductUpdateRequestInput()

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateRecipeStepProduct(ctx, recipeID, recipeStepID, recipeStepProductID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fakes.BuildFakeID()
		recipeStepID := fakes.BuildFakeID()
		recipeStepProductID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeRecipeStepProductUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, recipeID, recipeStepID, recipeStepProductID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateRecipeStepProduct(ctx, recipeID, recipeStepID, recipeStepProductID, exampleInput)

		assert.Error(t, err)
	})
}
