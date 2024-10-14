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

func TestClient_UpdateRecipeStepInstrument(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%s/steps/%s/instruments/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fakes.BuildFakeID()
		recipeStepID := fakes.BuildFakeID()
		recipeStepInstrumentID := fakes.BuildFakeID()

		data := fakes.BuildFakeRecipeStepInstrument()
		expected := &types.APIResponse[*types.RecipeStepInstrument]{
			Data: data,
		}

		exampleInput := fakes.BuildFakeRecipeStepInstrumentUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, recipeID, recipeStepID, recipeStepInstrumentID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateRecipeStepInstrument(ctx, recipeID, recipeStepID, recipeStepInstrumentID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		recipeStepID := fakes.BuildFakeID()
		recipeStepInstrumentID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeRecipeStepInstrumentUpdateRequestInput()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateRecipeStepInstrument(ctx, "", recipeStepID, recipeStepInstrumentID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with invalid recipeStep ID", func(t *testing.T) {
		t.Parallel()

		recipeID := fakes.BuildFakeID()

		recipeStepInstrumentID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeRecipeStepInstrumentUpdateRequestInput()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateRecipeStepInstrument(ctx, recipeID, "", recipeStepInstrumentID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with invalid recipeStepInstrument ID", func(t *testing.T) {
		t.Parallel()

		recipeID := fakes.BuildFakeID()
		recipeStepID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeRecipeStepInstrumentUpdateRequestInput()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateRecipeStepInstrument(ctx, recipeID, recipeStepID, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fakes.BuildFakeID()
		recipeStepID := fakes.BuildFakeID()
		recipeStepInstrumentID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeRecipeStepInstrumentUpdateRequestInput()

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateRecipeStepInstrument(ctx, recipeID, recipeStepID, recipeStepInstrumentID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fakes.BuildFakeID()
		recipeStepID := fakes.BuildFakeID()
		recipeStepInstrumentID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeRecipeStepInstrumentUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, recipeID, recipeStepID, recipeStepInstrumentID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateRecipeStepInstrument(ctx, recipeID, recipeStepID, recipeStepInstrumentID, exampleInput)

		assert.Error(t, err)
	})
}
