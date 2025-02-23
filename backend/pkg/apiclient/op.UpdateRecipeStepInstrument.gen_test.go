// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
)

func TestClient_UpdateRecipeStepInstrument(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%s/steps/%s/instruments/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		recipeID := fake.BuildFakeID()
		recipeStepID := fake.BuildFakeID()
		recipeStepInstrumentID := fake.BuildFakeID()

		data := &RecipeStepInstrument{}
		expected := &APIResponse[*RecipeStepInstrument]{
			Data: data,
		}

		exampleInput := &RecipeStepInstrumentUpdateRequestInput{}

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, recipeID, recipeStepID, recipeStepInstrumentID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateRecipeStepInstrument(ctx, recipeID, recipeStepID, recipeStepInstrumentID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		recipeStepID := fake.BuildFakeID()
		recipeStepInstrumentID := fake.BuildFakeID()

		exampleInput := &RecipeStepInstrumentUpdateRequestInput{}

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateRecipeStepInstrument(ctx, "", recipeStepID, recipeStepInstrumentID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with invalid recipeStep ID", func(t *testing.T) {
		t.Parallel()

		recipeID := fake.BuildFakeID()

		recipeStepInstrumentID := fake.BuildFakeID()

		exampleInput := &RecipeStepInstrumentUpdateRequestInput{}

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateRecipeStepInstrument(ctx, recipeID, "", recipeStepInstrumentID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with invalid recipeStepInstrument ID", func(t *testing.T) {
		t.Parallel()

		recipeID := fake.BuildFakeID()
		recipeStepID := fake.BuildFakeID()

		exampleInput := &RecipeStepInstrumentUpdateRequestInput{}

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateRecipeStepInstrument(ctx, recipeID, recipeStepID, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		recipeID := fake.BuildFakeID()
		recipeStepID := fake.BuildFakeID()
		recipeStepInstrumentID := fake.BuildFakeID()

		exampleInput := &RecipeStepInstrumentUpdateRequestInput{}

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateRecipeStepInstrument(ctx, recipeID, recipeStepID, recipeStepInstrumentID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		recipeID := fake.BuildFakeID()
		recipeStepID := fake.BuildFakeID()
		recipeStepInstrumentID := fake.BuildFakeID()

		exampleInput := &RecipeStepInstrumentUpdateRequestInput{}

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, recipeID, recipeStepID, recipeStepInstrumentID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateRecipeStepInstrument(ctx, recipeID, recipeStepID, recipeStepInstrumentID, exampleInput)

		assert.Error(t, err)
	})
}
