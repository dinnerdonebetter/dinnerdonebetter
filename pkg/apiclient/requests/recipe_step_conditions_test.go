package requests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/fakes"
)

func TestBuilder_BuildGetRecipeStepConditionRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%s/steps/%s/conditions/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepCondition := fakes.BuildFakeRecipeStepCondition()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepCondition.ID)

		actual, err := helper.builder.BuildGetRecipeStepConditionRequest(helper.ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepCondition.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepCondition := fakes.BuildFakeRecipeStepCondition()

		actual, err := helper.builder.BuildGetRecipeStepConditionRequest(helper.ctx, "", exampleRecipeStepID, exampleRecipeStepCondition.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepCondition := fakes.BuildFakeRecipeStepCondition()

		actual, err := helper.builder.BuildGetRecipeStepConditionRequest(helper.ctx, exampleRecipeID, "", exampleRecipeStepCondition.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid recipe step ingredient ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildGetRecipeStepConditionRequest(helper.ctx, exampleRecipeID, exampleRecipeStepID, "")
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepCondition := fakes.BuildFakeRecipeStepCondition()

		actual, err := helper.builder.BuildGetRecipeStepConditionRequest(helper.ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepCondition.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetRecipeStepConditionsRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%s/steps/%s/conditions"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		filter := (*types.QueryFilter)(nil)
		spec := newRequestSpec(true, http.MethodGet, "limit=20&page=1&sortBy=asc", expectedPathFormat, exampleRecipeID, exampleRecipeStepID)

		actual, err := helper.builder.BuildGetRecipeStepConditionsRequest(helper.ctx, exampleRecipeID, exampleRecipeStepID, filter)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeStepID := fakes.BuildFakeID()

		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetRecipeStepConditionsRequest(helper.ctx, "", exampleRecipeStepID, filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()

		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetRecipeStepConditionsRequest(helper.ctx, exampleRecipeID, "", filter)
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

		actual, err := helper.builder.BuildGetRecipeStepConditionsRequest(helper.ctx, exampleRecipeID, exampleRecipeStepID, filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildCreateRecipeStepConditionRequest(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/recipes/%s/steps/%s/conditions"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleInput := fakes.BuildFakeRecipeStepConditionCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath, exampleRecipeID, exampleRecipeStepID)

		actual, err := helper.builder.BuildCreateRecipeStepConditionRequest(helper.ctx, exampleRecipeID, exampleRecipeStepID, exampleInput)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleInput := fakes.BuildFakeRecipeStepConditionCreationRequestInput()

		actual, err := helper.builder.BuildCreateRecipeStepConditionRequest(helper.ctx, "", exampleRecipeStepID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildCreateRecipeStepConditionRequest(helper.ctx, exampleRecipeID, exampleRecipeStepID, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildCreateRecipeStepConditionRequest(helper.ctx, exampleRecipeID, exampleRecipeStepID, &types.RecipeStepConditionCreationRequestInput{})
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleInput := fakes.BuildFakeRecipeStepConditionCreationRequestInput()

		actual, err := helper.builder.BuildCreateRecipeStepConditionRequest(helper.ctx, exampleRecipeID, exampleRecipeStepID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildUpdateRecipeStepConditionRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%s/steps/%s/conditions/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepCondition := fakes.BuildFakeRecipeStepCondition()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, exampleRecipeID, exampleRecipeStepCondition.BelongsToRecipeStep, exampleRecipeStepCondition.ID)

		actual, err := helper.builder.BuildUpdateRecipeStepConditionRequest(helper.ctx, exampleRecipeID, exampleRecipeStepCondition)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeStepCondition := fakes.BuildFakeRecipeStepCondition()

		actual, err := helper.builder.BuildUpdateRecipeStepConditionRequest(helper.ctx, "", exampleRecipeStepCondition)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildUpdateRecipeStepConditionRequest(helper.ctx, exampleRecipeID, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepCondition := fakes.BuildFakeRecipeStepCondition()

		actual, err := helper.builder.BuildUpdateRecipeStepConditionRequest(helper.ctx, exampleRecipeID, exampleRecipeStepCondition)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildArchiveRecipeStepConditionRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%s/steps/%s/conditions/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepCondition := fakes.BuildFakeRecipeStepCondition()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepCondition.ID)

		actual, err := helper.builder.BuildArchiveRecipeStepConditionRequest(helper.ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepCondition.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepCondition := fakes.BuildFakeRecipeStepCondition()

		actual, err := helper.builder.BuildArchiveRecipeStepConditionRequest(helper.ctx, "", exampleRecipeStepID, exampleRecipeStepCondition.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepCondition := fakes.BuildFakeRecipeStepCondition()

		actual, err := helper.builder.BuildArchiveRecipeStepConditionRequest(helper.ctx, exampleRecipeID, "", exampleRecipeStepCondition.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid recipe step ingredient ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildArchiveRecipeStepConditionRequest(helper.ctx, exampleRecipeID, exampleRecipeStepID, "")
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepCondition := fakes.BuildFakeRecipeStepCondition()

		actual, err := helper.builder.BuildArchiveRecipeStepConditionRequest(helper.ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepCondition.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
