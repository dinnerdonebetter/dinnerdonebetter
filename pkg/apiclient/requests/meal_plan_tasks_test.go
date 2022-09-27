package requests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
)

func TestBuilder_BuildGetMealPlanTaskRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plan_tasks/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanTask := fakes.BuildFakeMealPlanTask()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, exampleMealPlanTask.ID)

		actual, err := helper.builder.BuildGetMealPlanTaskRequest(helper.ctx, exampleMealPlanTask.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid advanced prep step ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildGetMealPlanTaskRequest(helper.ctx, "")
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleMealPlanTask := fakes.BuildFakeMealPlanTask()

		actual, err := helper.builder.BuildGetMealPlanTaskRequest(helper.ctx, exampleMealPlanTask.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetMealPlanTasksRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plan_tasks"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		filter := (*types.QueryFilter)(nil)
		spec := newRequestSpec(true, http.MethodGet, "limit=20&page=1&sortBy=asc", expectedPathFormat)

		actual, err := helper.builder.BuildGetMealPlanTasksRequest(helper.ctx, filter)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetMealPlanTasksRequest(helper.ctx, filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildChangeMealPlanTaskStatusRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plan_tasks/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleInput := fakes.BuildFakeMealPlanTaskStatusChangeRequestInput()
		spec := newRequestSpec(false, http.MethodPatch, "", expectedPathFormat, exampleInput.ID)

		actual, err := helper.builder.BuildChangeMealPlanTaskStatusRequest(helper.ctx, exampleInput)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid advanced prep step ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildChangeMealPlanTaskStatusRequest(helper.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleInput := fakes.BuildFakeMealPlanTaskStatusChangeRequestInput()

		actual, err := helper.builder.BuildChangeMealPlanTaskStatusRequest(helper.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
