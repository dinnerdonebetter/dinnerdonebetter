package requests

import (
	"net/http"
	"testing"

	"gitlab.com/prixfixe/prixfixe/pkg/types"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuilder_BuildValidIngredientPreparationExistsRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_ingredient_preparations/%d"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()

		actual, err := helper.builder.BuildValidIngredientPreparationExistsRequest(helper.ctx, exampleValidIngredientPreparation.ID)
		spec := newRequestSpec(true, http.MethodHead, "", expectedPathFormat, exampleValidIngredientPreparation.ID)

		assert.NoError(t, err)
		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid valid ingredient preparation ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildValidIngredientPreparationExistsRequest(helper.ctx, 0)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()

		actual, err := helper.builder.BuildValidIngredientPreparationExistsRequest(helper.ctx, exampleValidIngredientPreparation.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetValidIngredientPreparationRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_ingredient_preparations/%d"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, exampleValidIngredientPreparation.ID)

		actual, err := helper.builder.BuildGetValidIngredientPreparationRequest(helper.ctx, exampleValidIngredientPreparation.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid valid ingredient preparation ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildGetValidIngredientPreparationRequest(helper.ctx, 0)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()

		actual, err := helper.builder.BuildGetValidIngredientPreparationRequest(helper.ctx, exampleValidIngredientPreparation.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetValidIngredientPreparationsRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_ingredient_preparations"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		filter := (*types.QueryFilter)(nil)
		spec := newRequestSpec(true, http.MethodGet, "includeArchived=false&limit=20&page=1&sortBy=asc", expectedPathFormat)

		actual, err := helper.builder.BuildGetValidIngredientPreparationsRequest(helper.ctx, filter)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetValidIngredientPreparationsRequest(helper.ctx, filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildCreateValidIngredientPreparationRequest(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/valid_ingredient_preparations"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleInput := fakes.BuildFakeValidIngredientPreparationCreationInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)

		actual, err := helper.builder.BuildCreateValidIngredientPreparationRequest(helper.ctx, exampleInput)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildCreateValidIngredientPreparationRequest(helper.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildCreateValidIngredientPreparationRequest(helper.ctx, &types.ValidIngredientPreparationCreationInput{})
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleInput := fakes.BuildFakeValidIngredientPreparationCreationInput()

		actual, err := helper.builder.BuildCreateValidIngredientPreparationRequest(helper.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildUpdateValidIngredientPreparationRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_ingredient_preparations/%d"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, exampleValidIngredientPreparation.ID)

		actual, err := helper.builder.BuildUpdateValidIngredientPreparationRequest(helper.ctx, exampleValidIngredientPreparation)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildUpdateValidIngredientPreparationRequest(helper.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()

		actual, err := helper.builder.BuildUpdateValidIngredientPreparationRequest(helper.ctx, exampleValidIngredientPreparation)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildArchiveValidIngredientPreparationRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_ingredient_preparations/%d"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, exampleValidIngredientPreparation.ID)

		actual, err := helper.builder.BuildArchiveValidIngredientPreparationRequest(helper.ctx, exampleValidIngredientPreparation.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid valid ingredient preparation ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildArchiveValidIngredientPreparationRequest(helper.ctx, 0)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()

		actual, err := helper.builder.BuildArchiveValidIngredientPreparationRequest(helper.ctx, exampleValidIngredientPreparation.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetAuditLogForValidIngredientPreparationRequest(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/valid_ingredient_preparations/%d/audit"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()

		actual, err := helper.builder.BuildGetAuditLogForValidIngredientPreparationRequest(helper.ctx, exampleValidIngredientPreparation.ID)
		require.NotNil(t, actual)
		assert.NoError(t, err)

		spec := newRequestSpec(true, http.MethodGet, "", expectedPath, exampleValidIngredientPreparation.ID)
		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid valid ingredient preparation ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildGetAuditLogForValidIngredientPreparationRequest(helper.ctx, 0)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()

		actual, err := helper.builder.BuildGetAuditLogForValidIngredientPreparationRequest(helper.ctx, exampleValidIngredientPreparation.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
