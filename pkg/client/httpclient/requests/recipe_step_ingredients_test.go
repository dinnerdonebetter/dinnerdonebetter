package requests

import (
	"net/http"
	"testing"

	"gitlab.com/prixfixe/prixfixe/pkg/types"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuilder_BuildRecipeStepIngredientExistsRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%d/recipe_steps/%d/recipe_step_ingredients/%d"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		actual, err := helper.builder.BuildRecipeStepIngredientExistsRequest(helper.ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepIngredient.ID)
		spec := newRequestSpec(true, http.MethodHead, "", expectedPathFormat, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepIngredient.ID)

		assert.NoError(t, err)
		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		actual, err := helper.builder.BuildRecipeStepIngredientExistsRequest(helper.ctx, 0, exampleRecipeStepID, exampleRecipeStepIngredient.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		actual, err := helper.builder.BuildRecipeStepIngredientExistsRequest(helper.ctx, exampleRecipeID, 0, exampleRecipeStepIngredient.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid recipe step ingredient ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildRecipeStepIngredientExistsRequest(helper.ctx, exampleRecipeID, exampleRecipeStepID, 0)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		actual, err := helper.builder.BuildRecipeStepIngredientExistsRequest(helper.ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepIngredient.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetRecipeStepIngredientRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%d/recipe_steps/%d/recipe_step_ingredients/%d"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepIngredient.ID)

		actual, err := helper.builder.BuildGetRecipeStepIngredientRequest(helper.ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepIngredient.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		actual, err := helper.builder.BuildGetRecipeStepIngredientRequest(helper.ctx, 0, exampleRecipeStepID, exampleRecipeStepIngredient.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		actual, err := helper.builder.BuildGetRecipeStepIngredientRequest(helper.ctx, exampleRecipeID, 0, exampleRecipeStepIngredient.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid recipe step ingredient ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildGetRecipeStepIngredientRequest(helper.ctx, exampleRecipeID, exampleRecipeStepID, 0)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		actual, err := helper.builder.BuildGetRecipeStepIngredientRequest(helper.ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepIngredient.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetRecipeStepIngredientsRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%d/recipe_steps/%d/recipe_step_ingredients"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		filter := (*types.QueryFilter)(nil)
		spec := newRequestSpec(true, http.MethodGet, "includeArchived=false&limit=20&page=1&sortBy=asc", expectedPathFormat, exampleRecipeID, exampleRecipeStepID)

		actual, err := helper.builder.BuildGetRecipeStepIngredientsRequest(helper.ctx, exampleRecipeID, exampleRecipeStepID, filter)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeStepID := fakes.BuildFakeID()

		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetRecipeStepIngredientsRequest(helper.ctx, 0, exampleRecipeStepID, filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()

		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetRecipeStepIngredientsRequest(helper.ctx, exampleRecipeID, 0, filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetRecipeStepIngredientsRequest(helper.ctx, exampleRecipeID, exampleRecipeStepID, filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildCreateRecipeStepIngredientRequest(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/recipes/%d/recipe_steps/%d/recipe_step_ingredients"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()
		exampleInput := fakes.BuildFakeRecipeStepIngredientCreationInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath, exampleRecipeID, exampleInput.BelongsToRecipeStep)

		actual, err := helper.builder.BuildCreateRecipeStepIngredientRequest(helper.ctx, exampleRecipeID, exampleInput)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleInput := fakes.BuildFakeRecipeStepIngredientCreationInput()

		actual, err := helper.builder.BuildCreateRecipeStepIngredientRequest(helper.ctx, 0, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildCreateRecipeStepIngredientRequest(helper.ctx, exampleRecipeID, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildCreateRecipeStepIngredientRequest(helper.ctx, exampleRecipeID, &types.RecipeStepIngredientCreationInput{})
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleRecipeID := fakes.BuildFakeID()
		exampleInput := fakes.BuildFakeRecipeStepIngredientCreationInput()

		actual, err := helper.builder.BuildCreateRecipeStepIngredientRequest(helper.ctx, exampleRecipeID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildUpdateRecipeStepIngredientRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%d/recipe_steps/%d/recipe_step_ingredients/%d"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, exampleRecipeID, exampleRecipeStepIngredient.BelongsToRecipeStep, exampleRecipeStepIngredient.ID)

		actual, err := helper.builder.BuildUpdateRecipeStepIngredientRequest(helper.ctx, exampleRecipeID, exampleRecipeStepIngredient)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		actual, err := helper.builder.BuildUpdateRecipeStepIngredientRequest(helper.ctx, 0, exampleRecipeStepIngredient)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildUpdateRecipeStepIngredientRequest(helper.ctx, exampleRecipeID, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		actual, err := helper.builder.BuildUpdateRecipeStepIngredientRequest(helper.ctx, exampleRecipeID, exampleRecipeStepIngredient)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildArchiveRecipeStepIngredientRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%d/recipe_steps/%d/recipe_step_ingredients/%d"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepIngredient.ID)

		actual, err := helper.builder.BuildArchiveRecipeStepIngredientRequest(helper.ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepIngredient.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		actual, err := helper.builder.BuildArchiveRecipeStepIngredientRequest(helper.ctx, 0, exampleRecipeStepID, exampleRecipeStepIngredient.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		actual, err := helper.builder.BuildArchiveRecipeStepIngredientRequest(helper.ctx, exampleRecipeID, 0, exampleRecipeStepIngredient.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid recipe step ingredient ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildArchiveRecipeStepIngredientRequest(helper.ctx, exampleRecipeID, exampleRecipeStepID, 0)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		actual, err := helper.builder.BuildArchiveRecipeStepIngredientRequest(helper.ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepIngredient.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetAuditLogForRecipeStepIngredientRequest(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/recipes/%d/recipe_steps/%d/recipe_step_ingredients/%d/audit"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		actual, err := helper.builder.BuildGetAuditLogForRecipeStepIngredientRequest(helper.ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepIngredient.ID)
		require.NotNil(t, actual)
		assert.NoError(t, err)

		spec := newRequestSpec(true, http.MethodGet, "", expectedPath, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepIngredient.ID)
		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		actual, err := helper.builder.BuildGetAuditLogForRecipeStepIngredientRequest(helper.ctx, 0, exampleRecipeStepID, exampleRecipeStepIngredient.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		actual, err := helper.builder.BuildGetAuditLogForRecipeStepIngredientRequest(helper.ctx, exampleRecipeID, 0, exampleRecipeStepIngredient.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid recipe step ingredient ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildGetAuditLogForRecipeStepIngredientRequest(helper.ctx, exampleRecipeID, exampleRecipeStepID, 0)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		actual, err := helper.builder.BuildGetAuditLogForRecipeStepIngredientRequest(helper.ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepIngredient.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
