package requests

import (
	"fmt"
	"net/http"
	"testing"

	"gitlab.com/prixfixe/prixfixe/pkg/types"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuilder_BuildValidIngredientExistsRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_ingredients/%d"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleValidIngredient := fakes.BuildFakeValidIngredient()

		actual, err := helper.builder.BuildValidIngredientExistsRequest(helper.ctx, exampleValidIngredient.ID)
		spec := newRequestSpec(true, http.MethodHead, "", expectedPathFormat, exampleValidIngredient.ID)

		assert.NoError(t, err)
		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid valid ingredient ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildValidIngredientExistsRequest(helper.ctx, 0)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleValidIngredient := fakes.BuildFakeValidIngredient()

		actual, err := helper.builder.BuildValidIngredientExistsRequest(helper.ctx, exampleValidIngredient.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetValidIngredientRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_ingredients/%d"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleValidIngredient := fakes.BuildFakeValidIngredient()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, exampleValidIngredient.ID)

		actual, err := helper.builder.BuildGetValidIngredientRequest(helper.ctx, exampleValidIngredient.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid valid ingredient ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildGetValidIngredientRequest(helper.ctx, 0)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleValidIngredient := fakes.BuildFakeValidIngredient()

		actual, err := helper.builder.BuildGetValidIngredientRequest(helper.ctx, exampleValidIngredient.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetValidIngredientsRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_ingredients"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		filter := (*types.QueryFilter)(nil)
		spec := newRequestSpec(true, http.MethodGet, "includeArchived=false&limit=20&page=1&sortBy=asc", expectedPathFormat)

		actual, err := helper.builder.BuildGetValidIngredientsRequest(helper.ctx, filter)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetValidIngredientsRequest(helper.ctx, filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildSearchValidIngredientsRequest(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/valid_ingredients/search"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleValidPreparation := fakes.BuildFakeValidPreparation()
		limit := types.DefaultQueryFilter().Limit
		exampleQuery := "whatever"
		spec := newRequestSpec(true, http.MethodGet, fmt.Sprintf("limit=20&pid=%d&q=whatever", exampleValidPreparation.ID), expectedPath)

		actual, err := helper.builder.BuildSearchValidIngredientsRequest(helper.ctx, exampleValidPreparation.ID, exampleQuery, limit)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleValidPreparation := fakes.BuildFakeValidPreparation()
		limit := types.DefaultQueryFilter().Limit
		exampleQuery := "whatever"

		actual, err := helper.builder.BuildSearchValidIngredientsRequest(helper.ctx, exampleValidPreparation.ID, exampleQuery, limit)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildCreateValidIngredientRequest(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/valid_ingredients"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleInput := fakes.BuildFakeValidIngredientCreationInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)

		actual, err := helper.builder.BuildCreateValidIngredientRequest(helper.ctx, exampleInput)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildCreateValidIngredientRequest(helper.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildCreateValidIngredientRequest(helper.ctx, &types.ValidIngredientCreationInput{})
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleInput := fakes.BuildFakeValidIngredientCreationInput()

		actual, err := helper.builder.BuildCreateValidIngredientRequest(helper.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildUpdateValidIngredientRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_ingredients/%d"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleValidIngredient := fakes.BuildFakeValidIngredient()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, exampleValidIngredient.ID)

		actual, err := helper.builder.BuildUpdateValidIngredientRequest(helper.ctx, exampleValidIngredient)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildUpdateValidIngredientRequest(helper.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleValidIngredient := fakes.BuildFakeValidIngredient()

		actual, err := helper.builder.BuildUpdateValidIngredientRequest(helper.ctx, exampleValidIngredient)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildArchiveValidIngredientRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_ingredients/%d"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleValidIngredient := fakes.BuildFakeValidIngredient()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, exampleValidIngredient.ID)

		actual, err := helper.builder.BuildArchiveValidIngredientRequest(helper.ctx, exampleValidIngredient.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid valid ingredient ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildArchiveValidIngredientRequest(helper.ctx, 0)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleValidIngredient := fakes.BuildFakeValidIngredient()

		actual, err := helper.builder.BuildArchiveValidIngredientRequest(helper.ctx, exampleValidIngredient.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetAuditLogForValidIngredientRequest(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/valid_ingredients/%d/audit"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleValidIngredient := fakes.BuildFakeValidIngredient()

		actual, err := helper.builder.BuildGetAuditLogForValidIngredientRequest(helper.ctx, exampleValidIngredient.ID)
		require.NotNil(t, actual)
		assert.NoError(t, err)

		spec := newRequestSpec(true, http.MethodGet, "", expectedPath, exampleValidIngredient.ID)
		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid valid ingredient ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildGetAuditLogForValidIngredientRequest(helper.ctx, 0)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleValidIngredient := fakes.BuildFakeValidIngredient()

		actual, err := helper.builder.BuildGetAuditLogForValidIngredientRequest(helper.ctx, exampleValidIngredient.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
