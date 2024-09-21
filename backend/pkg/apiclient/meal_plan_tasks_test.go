package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/pkg/pointer"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestMealPlanTasks(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(mealPlanTasksTestSuite))
}

type mealPlanTasksBaseSuite struct {
	suite.Suite
	ctx                             context.Context
	exampleMealPlan                 *types.MealPlan
	exampleMealPlanTask             *types.MealPlanTask
	exampleMealPlanTaskStatusUpdate *types.MealPlanTaskStatusChangeRequestInput
	exampleMealPlanTaskResponse     *types.APIResponse[*types.MealPlanTask]
	exampleMealPlanTaskListResponse *types.APIResponse[[]*types.MealPlanTask]
	exampleMealPlanTaskList         []*types.MealPlanTask
}

var _ suite.SetupTestSuite = (*mealPlanTasksBaseSuite)(nil)

func (s *mealPlanTasksBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleMealPlan = fakes.BuildFakeMealPlan()
	s.exampleMealPlanTask = fakes.BuildFakeMealPlanTask()
	s.exampleMealPlanTask.MealPlanOption = *s.exampleMealPlan.Events[0].Options[0]
	s.exampleMealPlanTaskResponse = &types.APIResponse[*types.MealPlanTask]{
		Data: s.exampleMealPlanTask,
	}
	s.exampleMealPlanTaskStatusUpdate = fakes.BuildFakeMealPlanTaskStatusChangeRequestInput()
	exampleMealPlanTaskList := fakes.BuildFakeMealPlanTaskList()
	s.exampleMealPlanTaskList = exampleMealPlanTaskList.Data
	s.exampleMealPlanTaskListResponse = &types.APIResponse[[]*types.MealPlanTask]{
		Data:       exampleMealPlanTaskList.Data,
		Pagination: &exampleMealPlanTaskList.Pagination,
	}
}

type mealPlanTasksTestSuite struct {
	suite.Suite
	mealPlanTasksBaseSuite
}

func (s *mealPlanTasksTestSuite) TestClient_GetMealPlanTask() {
	const expectedPathFormat = "/api/v1/meal_plans/%s/tasks/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleMealPlan.ID, s.exampleMealPlanTask.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleMealPlanTaskResponse)
		actual, err := c.GetMealPlanTask(s.ctx, s.exampleMealPlan.ID, s.exampleMealPlanTask.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)

		assert.Equal(t, s.exampleMealPlanTask, actual)
	})

	s.Run("with invalid meal plan ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetMealPlanTask(s.ctx, "", s.exampleMealPlanTask.ID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid meal plan task ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetMealPlanTask(s.ctx, s.exampleMealPlan.ID, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetMealPlanTask(s.ctx, s.exampleMealPlan.ID, s.exampleMealPlanTask.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleMealPlan.ID, s.exampleMealPlanTask.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetMealPlanTask(s.ctx, s.exampleMealPlan.ID, s.exampleMealPlanTask.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *mealPlanTasksTestSuite) TestClient_CreateMealPlanTask() {
	const expectedPath = "/api/v1/meal_plans/%s/tasks"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeMealPlanTaskCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath, s.exampleMealPlan.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleMealPlanTaskResponse)

		actual, err := c.CreateMealPlanTask(s.ctx, s.exampleMealPlan.ID, exampleInput)
		assert.NoError(t, err)

		assert.Equal(t, s.exampleMealPlanTask, actual)
	})

	s.Run("with missing meal plan ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := fakes.BuildFakeMealPlanTaskCreationRequestInput()

		actual, err := c.CreateMealPlanTask(s.ctx, "", exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateMealPlanTask(s.ctx, s.exampleMealPlan.ID, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.MealPlanTaskCreationRequestInput{}

		actual, err := c.CreateMealPlanTask(s.ctx, s.exampleMealPlan.ID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := converters.ConvertMealPlanTaskToMealPlanTaskCreationRequestInput(s.exampleMealPlanTask)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateMealPlanTask(s.ctx, s.exampleMealPlan.ID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := converters.ConvertMealPlanTaskToMealPlanTaskCreationRequestInput(s.exampleMealPlanTask)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateMealPlanTask(s.ctx, s.exampleMealPlan.ID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *mealPlanTasksTestSuite) TestClient_UpdateMealPlanTaskStatus() {
	const expectedPathFormat = "/api/v1/meal_plans/%s/tasks/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPatch, "", expectedPathFormat, s.exampleMealPlan.ID, s.exampleMealPlanTaskStatusUpdate.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleMealPlanTaskResponse)

		task, err := c.UpdateMealPlanTaskStatus(s.ctx, s.exampleMealPlan.ID, s.exampleMealPlanTaskStatusUpdate)
		assert.NoError(t, err)
		assert.NotNil(t, task)
	})

	s.Run("with invalid ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		_, err := c.UpdateMealPlanTaskStatus(s.ctx, "", s.exampleMealPlanTaskStatusUpdate)
		assert.Error(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		_, err := c.UpdateMealPlanTaskStatus(s.ctx, s.exampleMealPlan.ID, nil)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		s.exampleMealPlanTaskStatusUpdate.Status = pointer.To(t.Name())

		_, err := c.UpdateMealPlanTaskStatus(s.ctx, s.exampleMealPlan.ID, s.exampleMealPlanTaskStatusUpdate)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		_, err := c.UpdateMealPlanTaskStatus(s.ctx, s.exampleMealPlan.ID, s.exampleMealPlanTaskStatusUpdate)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		_, err := c.UpdateMealPlanTaskStatus(s.ctx, s.exampleMealPlan.ID, s.exampleMealPlanTaskStatusUpdate)
		assert.Error(t, err)
	})
}
