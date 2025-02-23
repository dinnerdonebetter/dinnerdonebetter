// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
)

func TestClient_UpdateRecipeStepCompletionCondition(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%s/steps/%s/completion_conditions/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fake.BuildFakeID()
		recipeStepID := fake.BuildFakeID()
		recipeStepCompletionConditionID := fake.BuildFakeID()

		data := &RecipeStepCompletionCondition{}
		expected := &APIResponse[*RecipeStepCompletionCondition]{
			Data: data,
		}

		exampleInput := &RecipeStepCompletionConditionUpdateRequestInput{}

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, recipeID, recipeStepID, recipeStepCompletionConditionID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateRecipeStepCompletionCondition(ctx, recipeID, recipeStepID, recipeStepCompletionConditionID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		recipeStepID := fake.BuildFakeID()
		recipeStepCompletionConditionID := fake.BuildFakeID()

		exampleInput := &RecipeStepCompletionConditionUpdateRequestInput{}

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateRecipeStepCompletionCondition(ctx, "", recipeStepID, recipeStepCompletionConditionID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with invalid recipeStep ID", func(t *testing.T) {
		t.Parallel()

		recipeID := fake.BuildFakeID()

		recipeStepCompletionConditionID := fake.BuildFakeID()

		exampleInput := &RecipeStepCompletionConditionUpdateRequestInput{}

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateRecipeStepCompletionCondition(ctx, recipeID, "", recipeStepCompletionConditionID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with invalid recipeStepCompletionCondition ID", func(t *testing.T) {
		t.Parallel()

		recipeID := fake.BuildFakeID()
		recipeStepID := fake.BuildFakeID()

		exampleInput := &RecipeStepCompletionConditionUpdateRequestInput{}

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateRecipeStepCompletionCondition(ctx, recipeID, recipeStepID, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fake.BuildFakeID()
		recipeStepID := fake.BuildFakeID()
		recipeStepCompletionConditionID := fake.BuildFakeID()

		exampleInput := &RecipeStepCompletionConditionUpdateRequestInput{}

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateRecipeStepCompletionCondition(ctx, recipeID, recipeStepID, recipeStepCompletionConditionID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fake.BuildFakeID()
		recipeStepID := fake.BuildFakeID()
		recipeStepCompletionConditionID := fake.BuildFakeID()

		exampleInput := &RecipeStepCompletionConditionUpdateRequestInput{}

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, recipeID, recipeStepID, recipeStepCompletionConditionID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateRecipeStepCompletionCondition(ctx, recipeID, recipeStepID, recipeStepCompletionConditionID, exampleInput)

		assert.Error(t, err)
	})
}
