package requests

import (
	"net/http"
	"testing"

	"gitlab.com/prixfixe/prixfixe/pkg/types"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuilder_BuildRecipeStepProductExistsRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%d/recipe_steps/%d/recipe_step_products/%d"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()

		actual, err := helper.builder.BuildRecipeStepProductExistsRequest(helper.ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepProduct.ID)
		spec := newRequestSpec(true, http.MethodHead, "", expectedPathFormat, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepProduct.ID)

		assert.NoError(t, err)
		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()

		actual, err := helper.builder.BuildRecipeStepProductExistsRequest(helper.ctx, 0, exampleRecipeStepID, exampleRecipeStepProduct.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()

		actual, err := helper.builder.BuildRecipeStepProductExistsRequest(helper.ctx, exampleRecipeID, 0, exampleRecipeStepProduct.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid recipe step product ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildRecipeStepProductExistsRequest(helper.ctx, exampleRecipeID, exampleRecipeStepID, 0)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()

		actual, err := helper.builder.BuildRecipeStepProductExistsRequest(helper.ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepProduct.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetRecipeStepProductRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%d/recipe_steps/%d/recipe_step_products/%d"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepProduct.ID)

		actual, err := helper.builder.BuildGetRecipeStepProductRequest(helper.ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepProduct.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()

		actual, err := helper.builder.BuildGetRecipeStepProductRequest(helper.ctx, 0, exampleRecipeStepID, exampleRecipeStepProduct.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()

		actual, err := helper.builder.BuildGetRecipeStepProductRequest(helper.ctx, exampleRecipeID, 0, exampleRecipeStepProduct.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid recipe step product ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildGetRecipeStepProductRequest(helper.ctx, exampleRecipeID, exampleRecipeStepID, 0)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()

		actual, err := helper.builder.BuildGetRecipeStepProductRequest(helper.ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepProduct.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetRecipeStepProductsRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%d/recipe_steps/%d/recipe_step_products"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		filter := (*types.QueryFilter)(nil)
		spec := newRequestSpec(true, http.MethodGet, "includeArchived=false&limit=20&page=1&sortBy=asc", expectedPathFormat, exampleRecipeID, exampleRecipeStepID)

		actual, err := helper.builder.BuildGetRecipeStepProductsRequest(helper.ctx, exampleRecipeID, exampleRecipeStepID, filter)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeStepID := fakes.BuildFakeID()

		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetRecipeStepProductsRequest(helper.ctx, 0, exampleRecipeStepID, filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()

		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetRecipeStepProductsRequest(helper.ctx, exampleRecipeID, 0, filter)
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

		actual, err := helper.builder.BuildGetRecipeStepProductsRequest(helper.ctx, exampleRecipeID, exampleRecipeStepID, filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildCreateRecipeStepProductRequest(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/recipes/%d/recipe_steps/%d/recipe_step_products"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()
		exampleInput := fakes.BuildFakeRecipeStepProductCreationInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath, exampleRecipeID, exampleInput.BelongsToRecipeStep)

		actual, err := helper.builder.BuildCreateRecipeStepProductRequest(helper.ctx, exampleRecipeID, exampleInput)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleInput := fakes.BuildFakeRecipeStepProductCreationInput()

		actual, err := helper.builder.BuildCreateRecipeStepProductRequest(helper.ctx, 0, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildCreateRecipeStepProductRequest(helper.ctx, exampleRecipeID, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildCreateRecipeStepProductRequest(helper.ctx, exampleRecipeID, &types.RecipeStepProductCreationInput{})
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleRecipeID := fakes.BuildFakeID()
		exampleInput := fakes.BuildFakeRecipeStepProductCreationInput()

		actual, err := helper.builder.BuildCreateRecipeStepProductRequest(helper.ctx, exampleRecipeID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildUpdateRecipeStepProductRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%d/recipe_steps/%d/recipe_step_products/%d"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, exampleRecipeID, exampleRecipeStepProduct.BelongsToRecipeStep, exampleRecipeStepProduct.ID)

		actual, err := helper.builder.BuildUpdateRecipeStepProductRequest(helper.ctx, exampleRecipeID, exampleRecipeStepProduct)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()

		actual, err := helper.builder.BuildUpdateRecipeStepProductRequest(helper.ctx, 0, exampleRecipeStepProduct)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildUpdateRecipeStepProductRequest(helper.ctx, exampleRecipeID, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()

		actual, err := helper.builder.BuildUpdateRecipeStepProductRequest(helper.ctx, exampleRecipeID, exampleRecipeStepProduct)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildArchiveRecipeStepProductRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%d/recipe_steps/%d/recipe_step_products/%d"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepProduct.ID)

		actual, err := helper.builder.BuildArchiveRecipeStepProductRequest(helper.ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepProduct.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()

		actual, err := helper.builder.BuildArchiveRecipeStepProductRequest(helper.ctx, 0, exampleRecipeStepID, exampleRecipeStepProduct.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()

		actual, err := helper.builder.BuildArchiveRecipeStepProductRequest(helper.ctx, exampleRecipeID, 0, exampleRecipeStepProduct.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid recipe step product ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildArchiveRecipeStepProductRequest(helper.ctx, exampleRecipeID, exampleRecipeStepID, 0)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()

		actual, err := helper.builder.BuildArchiveRecipeStepProductRequest(helper.ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepProduct.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetAuditLogForRecipeStepProductRequest(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/recipes/%d/recipe_steps/%d/recipe_step_products/%d/audit"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()

		actual, err := helper.builder.BuildGetAuditLogForRecipeStepProductRequest(helper.ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepProduct.ID)
		require.NotNil(t, actual)
		assert.NoError(t, err)

		spec := newRequestSpec(true, http.MethodGet, "", expectedPath, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepProduct.ID)
		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()

		actual, err := helper.builder.BuildGetAuditLogForRecipeStepProductRequest(helper.ctx, 0, exampleRecipeStepID, exampleRecipeStepProduct.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()

		actual, err := helper.builder.BuildGetAuditLogForRecipeStepProductRequest(helper.ctx, exampleRecipeID, 0, exampleRecipeStepProduct.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid recipe step product ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildGetAuditLogForRecipeStepProductRequest(helper.ctx, exampleRecipeID, exampleRecipeStepID, 0)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()

		actual, err := helper.builder.BuildGetAuditLogForRecipeStepProductRequest(helper.ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepProduct.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
