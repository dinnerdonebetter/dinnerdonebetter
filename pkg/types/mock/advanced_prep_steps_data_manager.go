package mocktypes

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/pkg/types"
)

var _ types.MealPlanTaskDataManager = (*MealPlanTaskDataManager)(nil)

// MealPlanTaskDataManager is a mocked types.MealPlanTaskDataManager for testing.
type MealPlanTaskDataManager struct {
	mock.Mock
}

// MealPlanTaskExists is a mock function.
func (m *MealPlanTaskDataManager) MealPlanTaskExists(ctx context.Context, mealPlanID, mealPlanTaskID string) (bool, error) {
	args := m.Called(ctx, mealPlanID, mealPlanTaskID)
	return args.Bool(0), args.Error(1)
}

// GetMealPlanTask is a mock function.
func (m *MealPlanTaskDataManager) GetMealPlanTask(ctx context.Context, mealPlanTaskID string) (*types.MealPlanTask, error) {
	args := m.Called(ctx, mealPlanTaskID)
	return args.Get(0).(*types.MealPlanTask), args.Error(1)
}

// CreateMealPlanTask is a mock function.
func (m *MealPlanTaskDataManager) CreateMealPlanTask(ctx context.Context, mealPlanID string, input *types.MealPlanTaskDatabaseCreationInput) (*types.MealPlanTask, error) {
	args := m.Called(ctx, mealPlanID, input)
	return args.Get(0).(*types.MealPlanTask), args.Error(1)
}

// GetMealPlanTasksForMealPlan is a mock function.
func (m *MealPlanTaskDataManager) GetMealPlanTasksForMealPlan(ctx context.Context, mealPlanID string) ([]*types.MealPlanTask, error) {
	args := m.Called(ctx, mealPlanID)
	return args.Get(0).([]*types.MealPlanTask), args.Error(1)
}

// CreateMealPlanTasksForMealPlanOption is a mock function.
func (m *MealPlanTaskDataManager) CreateMealPlanTasksForMealPlanOption(ctx context.Context, mealPlanOptionID string, inputs []*types.MealPlanTaskDatabaseCreationInput) ([]*types.MealPlanTask, error) {
	args := m.Called(ctx, mealPlanOptionID, inputs)
	return args.Get(0).([]*types.MealPlanTask), args.Error(1)
}

// ChangeMealPlanTaskStatus is a mock function.
func (m *MealPlanTaskDataManager) ChangeMealPlanTaskStatus(ctx context.Context, input *types.MealPlanTaskStatusChangeRequestInput) error {
	return m.Called(ctx, input).Error(0)
}
