// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_CreateRecipeStepCompletionCondition(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%s/steps/%s/completion_conditions"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		recipeID := fake.BuildFakeID()
		recipeStepID := fake.BuildFakeID()

		data := &RecipeStepCompletionCondition{}
		expected := &APIResponse[*RecipeStepCompletionCondition]{
			Data: data,
		}

		exampleInput := &RecipeStepCompletionConditionForExistingRecipeCreationRequestInput{}

		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat, recipeID, recipeStepID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.CreateRecipeStepCompletionCondition(ctx, recipeID, recipeStepID, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		recipeStepID := fake.BuildFakeID()

		exampleInput := &RecipeStepCompletionConditionForExistingRecipeCreationRequestInput{}

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.CreateRecipeStepCompletionCondition(ctx, "", recipeStepID, exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid recipeStep ID", func(t *testing.T) {
		t.Parallel()

		recipeID := fake.BuildFakeID()

		exampleInput := &RecipeStepCompletionConditionForExistingRecipeCreationRequestInput{}

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.CreateRecipeStepCompletionCondition(ctx, recipeID, "", exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		recipeID := fake.BuildFakeID()
		recipeStepID := fake.BuildFakeID()

		exampleInput := &RecipeStepCompletionConditionForExistingRecipeCreationRequestInput{}

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.CreateRecipeStepCompletionCondition(ctx, recipeID, recipeStepID, exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		recipeID := fake.BuildFakeID()
		recipeStepID := fake.BuildFakeID()

		exampleInput := &RecipeStepCompletionConditionForExistingRecipeCreationRequestInput{}

		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat, recipeID, recipeStepID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.CreateRecipeStepCompletionCondition(ctx, recipeID, recipeStepID, exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
