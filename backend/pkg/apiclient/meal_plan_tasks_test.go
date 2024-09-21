package apiclient

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/suite"
)

func TestMealPlanTasks(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(mealPlanTasksTestSuite))
}

type mealPlanTasksBaseSuite struct {
	suite.Suite
	ctx                             context.Context
	exampleMealPlanTask             *types.MealPlanTask
	exampleMealPlanTaskResponse     *types.APIResponse[*types.MealPlanTask]
	exampleMealPlanTaskListResponse *types.APIResponse[[]*types.MealPlanTask]
	exampleMealPlanTaskList         []*types.MealPlanTask
}

var _ suite.SetupTestSuite = (*mealPlanTasksBaseSuite)(nil)

func (s *mealPlanTasksBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleMealPlanTask = fakes.BuildFakeMealPlanTask()
	s.exampleMealPlanTaskResponse = &types.APIResponse[*types.MealPlanTask]{
		Data: s.exampleMealPlanTask,
	}
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

//func (s *mealPlanTasksTestSuite) TestClient_GetMealPlanTask() {
//	const expectedPathFormat = "/api/v1/meal_plans/%s"
//
//	s.Run("standard", func() {
//		t := s.T()
//
//		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleMealPlanTask.ID)
//		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleMealPlanTaskResponse)
//		actual, err := c.GetMealPlanTask(s.ctx, s.exampleMealPlanTask.ID)
//
//		require.NotNil(t, actual)
//		assert.NoError(t, err)
//
//		assert.Equal(t, s.exampleMealPlanTask, actual)
//	})
//
//	s.Run("with invalid meal plan ID", func() {
//		t := s.T()
//
//		c, _ := buildSimpleTestClient(t)
//		actual, err := c.GetMealPlanTask(s.ctx, "")
//
//		require.Nil(t, actual)
//		assert.Error(t, err)
//	})
//
//	s.Run("with error building request", func() {
//		t := s.T()
//
//		c := buildTestClientWithInvalidURL(t)
//		actual, err := c.GetMealPlanTask(s.ctx, s.exampleMealPlanTask.ID)
//
//		assert.Nil(t, actual)
//		assert.Error(t, err)
//	})
//
//	s.Run("with error executing request", func() {
//		t := s.T()
//
//		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleMealPlanTask.ID)
//		c := buildTestClientWithInvalidResponse(t, spec)
//		actual, err := c.GetMealPlanTask(s.ctx, s.exampleMealPlanTask.ID)
//
//		assert.Nil(t, actual)
//		assert.Error(t, err)
//	})
//}
//
//func (s *mealPlanTasksTestSuite) TestClient_CreateMealPlanTask() {
//	const expectedPath = "/api/v1/meal_plans"
//
//	s.Run("standard", func() {
//		t := s.T()
//
//		exampleInput := fakes.BuildFakeMealPlanTaskCreationRequestInput()
//
//		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
//		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleMealPlanTaskResponse)
//
//		actual, err := c.CreateMealPlanTask(s.ctx, exampleInput)
//		assert.NoError(t, err)
//
//		assert.Equal(t, s.exampleMealPlanTask, actual)
//	})
//
//	s.Run("with nil input", func() {
//		t := s.T()
//
//		c, _ := buildSimpleTestClient(t)
//
//		actual, err := c.CreateMealPlanTask(s.ctx, nil)
//		assert.Nil(t, actual)
//		assert.Error(t, err)
//	})
//
//	s.Run("with invalid input", func() {
//		t := s.T()
//
//		c, _ := buildSimpleTestClient(t)
//		exampleInput := &types.MealPlanTaskCreationRequestInput{}
//
//		actual, err := c.CreateMealPlanTask(s.ctx, exampleInput)
//		assert.Nil(t, actual)
//		assert.Error(t, err)
//	})
//
//	s.Run("with error building request", func() {
//		t := s.T()
//
//		exampleInput := converters.ConvertMealPlanTaskToMealPlanTaskCreationRequestInput(s.exampleMealPlanTask)
//
//		c := buildTestClientWithInvalidURL(t)
//
//		actual, err := c.CreateMealPlanTask(s.ctx, exampleInput)
//		assert.Nil(t, actual)
//		assert.Error(t, err)
//	})
//
//	s.Run("with error executing request", func() {
//		t := s.T()
//
//		exampleInput := converters.ConvertMealPlanTaskToMealPlanTaskCreationRequestInput(s.exampleMealPlanTask)
//		c, _ := buildTestClientThatWaitsTooLong(t)
//
//		actual, err := c.CreateMealPlanTask(s.ctx, exampleInput)
//		assert.Nil(t, actual)
//		assert.Error(t, err)
//	})
//}
//
//func (s *mealPlanTasksTestSuite) TestClient_UpdateMealPlanTaskStatus() {
//	const expectedPathFormat = "/api/v1/meal_plans/%s"
//
//	s.Run("standard", func() {
//		t := s.T()
//
//		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, s.exampleMealPlanTask.ID)
//		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleMealPlanTaskResponse)
//
//		err := c.UpdateMealPlanTaskStatus(s.ctx, s.exampleMealPlanTask)
//		assert.NoError(t, err)
//	})
//
//	s.Run("with nil input", func() {
//		t := s.T()
//
//		c, _ := buildSimpleTestClient(t)
//
//		_, err := c.UpdateMealPlanTaskStatus(s.ctx, nil)
//		assert.Error(t, err)
//	})
//
//	s.Run("with error building request", func() {
//		t := s.T()
//
//		c := buildTestClientWithInvalidURL(t)
//
//		_, err := c.UpdateMealPlanTaskStatus(s.ctx, s.exampleMealPlanTask)
//		assert.Error(t, err)
//	})
//
//	s.Run("with error executing request", func() {
//		t := s.T()
//
//		c, _ := buildTestClientThatWaitsTooLong(t)
//
//		_, err := c.UpdateMealPlanTaskStatus(s.ctx, s.exampleMealPlanTask)
//		assert.Error(t, err)
//	})
//}
