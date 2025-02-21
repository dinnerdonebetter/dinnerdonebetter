package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/services/eating/types"

	"github.com/stretchr/testify/mock"
)

var _ types.MealPlanTaskDataManager = (*MealPlanTaskDataManagerMock)(nil)

// MealPlanTaskDataManagerMock is a mocked types.MealPlanTaskDataManager for testing.
type MealPlanTaskDataManagerMock struct {
	mock.Mock
}

// MarkMealPlanAsHavingTasksCreated is a mock function.
func (m *MealPlanTaskDataManagerMock) MarkMealPlanAsHavingTasksCreated(ctx context.Context, mealPlanID string) error {
	return m.Called(ctx, mealPlanID).Error(0)
}

// MealPlanTaskExists is a mock function.
func (m *MealPlanTaskDataManagerMock) MealPlanTaskExists(ctx context.Context, mealPlanID, mealPlanTaskID string) (bool, error) {
	returnValues := m.Called(ctx, mealPlanID, mealPlanTaskID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetMealPlanTask is a mock function.
func (m *MealPlanTaskDataManagerMock) GetMealPlanTask(ctx context.Context, mealPlanTaskID string) (*types.MealPlanTask, error) {
	returnValues := m.Called(ctx, mealPlanTaskID)
	return returnValues.Get(0).(*types.MealPlanTask), returnValues.Error(1)
}

// CreateMealPlanTask is a mock function.
func (m *MealPlanTaskDataManagerMock) CreateMealPlanTask(ctx context.Context, input *types.MealPlanTaskDatabaseCreationInput) (*types.MealPlanTask, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*types.MealPlanTask), returnValues.Error(1)
}

// GetMealPlanTasksForMealPlan is a mock function.
func (m *MealPlanTaskDataManagerMock) GetMealPlanTasksForMealPlan(ctx context.Context, mealPlanID string) ([]*types.MealPlanTask, error) {
	returnValues := m.Called(ctx, mealPlanID)
	return returnValues.Get(0).([]*types.MealPlanTask), returnValues.Error(1)
}

// CreateMealPlanTasksForMealPlanOption is a mock function.
func (m *MealPlanTaskDataManagerMock) CreateMealPlanTasksForMealPlanOption(ctx context.Context, inputs []*types.MealPlanTaskDatabaseCreationInput) ([]*types.MealPlanTask, error) {
	returnValues := m.Called(ctx, inputs)
	return returnValues.Get(0).([]*types.MealPlanTask), returnValues.Error(1)
}

// ChangeMealPlanTaskStatus is a mock function.
func (m *MealPlanTaskDataManagerMock) ChangeMealPlanTaskStatus(ctx context.Context, input *types.MealPlanTaskStatusChangeRequestInput) error {
	return m.Called(ctx, input).Error(0)
}
