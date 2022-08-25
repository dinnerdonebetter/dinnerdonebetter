package requests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
)

func TestBuilder_BuildGetRecipeStepRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%s/recipe_steps/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStep := fakes.BuildFakeRecipeStep()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, exampleRecipeID, exampleRecipeStep.ID)

		actual, err := helper.builder.BuildGetRecipeStepRequest(helper.ctx, exampleRecipeID, exampleRecipeStep.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeStep := fakes.BuildFakeRecipeStep()

		actual, err := helper.builder.BuildGetRecipeStepRequest(helper.ctx, "", exampleRecipeStep.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildGetRecipeStepRequest(helper.ctx, exampleRecipeID, "")
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStep := fakes.BuildFakeRecipeStep()

		actual, err := helper.builder.BuildGetRecipeStepRequest(helper.ctx, exampleRecipeID, exampleRecipeStep.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetRecipeStepsRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%s/recipe_steps"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()

		filter := (*types.QueryFilter)(nil)
		spec := newRequestSpec(true, http.MethodGet, "limit=20&page=1&sortBy=asc", expectedPathFormat, exampleRecipeID)

		actual, err := helper.builder.BuildGetRecipeStepsRequest(helper.ctx, exampleRecipeID, filter)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetRecipeStepsRequest(helper.ctx, "", filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleRecipeID := fakes.BuildFakeID()

		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetRecipeStepsRequest(helper.ctx, exampleRecipeID, filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildCreateRecipeStepRequest(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/recipes/%s/recipe_steps"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleInput := fakes.BuildFakeRecipeStepCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath, exampleInput.BelongsToRecipe)

		actual, err := helper.builder.BuildCreateRecipeStepRequest(helper.ctx, exampleInput)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildCreateRecipeStepRequest(helper.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildCreateRecipeStepRequest(helper.ctx, &types.RecipeStepCreationRequestInput{})
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleInput := fakes.BuildFakeRecipeStepCreationRequestInput()

		actual, err := helper.builder.BuildCreateRecipeStepRequest(helper.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildUpdateRecipeStepRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%s/recipe_steps/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeStep := fakes.BuildFakeRecipeStep()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, exampleRecipeStep.BelongsToRecipe, exampleRecipeStep.ID)

		actual, err := helper.builder.BuildUpdateRecipeStepRequest(helper.ctx, exampleRecipeStep)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildUpdateRecipeStepRequest(helper.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleRecipeStep := fakes.BuildFakeRecipeStep()

		actual, err := helper.builder.BuildUpdateRecipeStepRequest(helper.ctx, exampleRecipeStep)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildArchiveRecipeStepRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%s/recipe_steps/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStep := fakes.BuildFakeRecipeStep()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, exampleRecipeID, exampleRecipeStep.ID)

		actual, err := helper.builder.BuildArchiveRecipeStepRequest(helper.ctx, exampleRecipeID, exampleRecipeStep.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeStep := fakes.BuildFakeRecipeStep()

		actual, err := helper.builder.BuildArchiveRecipeStepRequest(helper.ctx, "", exampleRecipeStep.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildArchiveRecipeStepRequest(helper.ctx, exampleRecipeID, "")
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStep := fakes.BuildFakeRecipeStep()

		actual, err := helper.builder.BuildArchiveRecipeStepRequest(helper.ctx, exampleRecipeID, exampleRecipeStep.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
