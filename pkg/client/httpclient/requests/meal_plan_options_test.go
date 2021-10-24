package requests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.com/prixfixe/prixfixe/pkg/types"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"
)

func TestBuilder_BuildGetMealPlanOptionRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plans/%s/meal_plan_options/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanOption := fakes.BuildFakeMealPlanOption()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, exampleMealPlanID, exampleMealPlanOption.ID)

		actual, err := helper.builder.BuildGetMealPlanOptionRequest(helper.ctx, exampleMealPlanID, exampleMealPlanOption.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanOption := fakes.BuildFakeMealPlanOption()

		actual, err := helper.builder.BuildGetMealPlanOptionRequest(helper.ctx, "", exampleMealPlanOption.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid meal plan option ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildGetMealPlanOptionRequest(helper.ctx, exampleMealPlanID, "")
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanOption := fakes.BuildFakeMealPlanOption()

		actual, err := helper.builder.BuildGetMealPlanOptionRequest(helper.ctx, exampleMealPlanID, exampleMealPlanOption.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetMealPlanOptionsRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plans/%s/meal_plan_options"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanID := fakes.BuildFakeID()

		filter := (*types.QueryFilter)(nil)
		spec := newRequestSpec(true, http.MethodGet, "includeArchived=false&limit=20&page=1&sortBy=asc", expectedPathFormat, exampleMealPlanID)

		actual, err := helper.builder.BuildGetMealPlanOptionsRequest(helper.ctx, exampleMealPlanID, filter)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetMealPlanOptionsRequest(helper.ctx, "", filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleMealPlanID := fakes.BuildFakeID()

		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetMealPlanOptionsRequest(helper.ctx, exampleMealPlanID, filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildCreateMealPlanOptionRequest(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/meal_plans/%s/meal_plan_options"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleInput := fakes.BuildFakeMealPlanOptionCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath, exampleInput.BelongsToMealPlan)

		actual, err := helper.builder.BuildCreateMealPlanOptionRequest(helper.ctx, exampleInput)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildCreateMealPlanOptionRequest(helper.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildCreateMealPlanOptionRequest(helper.ctx, &types.MealPlanOptionCreationRequestInput{})
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleInput := fakes.BuildFakeMealPlanOptionCreationRequestInput()

		actual, err := helper.builder.BuildCreateMealPlanOptionRequest(helper.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildUpdateMealPlanOptionRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plans/%s/meal_plan_options/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanOption := fakes.BuildFakeMealPlanOption()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, exampleMealPlanOption.BelongsToMealPlan, exampleMealPlanOption.ID)

		actual, err := helper.builder.BuildUpdateMealPlanOptionRequest(helper.ctx, exampleMealPlanOption)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildUpdateMealPlanOptionRequest(helper.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleMealPlanOption := fakes.BuildFakeMealPlanOption()

		actual, err := helper.builder.BuildUpdateMealPlanOptionRequest(helper.ctx, exampleMealPlanOption)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildArchiveMealPlanOptionRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plans/%s/meal_plan_options/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanOption := fakes.BuildFakeMealPlanOption()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, exampleMealPlanID, exampleMealPlanOption.ID)

		actual, err := helper.builder.BuildArchiveMealPlanOptionRequest(helper.ctx, exampleMealPlanID, exampleMealPlanOption.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanOption := fakes.BuildFakeMealPlanOption()

		actual, err := helper.builder.BuildArchiveMealPlanOptionRequest(helper.ctx, "", exampleMealPlanOption.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid meal plan option ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildArchiveMealPlanOptionRequest(helper.ctx, exampleMealPlanID, "")
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanOption := fakes.BuildFakeMealPlanOption()

		actual, err := helper.builder.BuildArchiveMealPlanOptionRequest(helper.ctx, exampleMealPlanID, exampleMealPlanOption.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
