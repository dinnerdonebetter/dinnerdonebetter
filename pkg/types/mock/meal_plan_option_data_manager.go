package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.MealPlanOptionDataManager = (*MealPlanOptionDataManager)(nil)

// MealPlanOptionDataManager is a mocked types.MealPlanOptionDataManager for testing.
type MealPlanOptionDataManager struct {
	mock.Mock
}

// MealPlanOptionExists is a mock function.
func (m *MealPlanOptionDataManager) MealPlanOptionExists(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string) (bool, error) {
	args := m.Called(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID)
	return args.Bool(0), args.Error(1)
}

// GetMealPlanOption is a mock function.
func (m *MealPlanOptionDataManager) GetMealPlanOption(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string) (*types.MealPlanOption, error) {
	args := m.Called(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID)
	return args.Get(0).(*types.MealPlanOption), args.Error(1)
}

// GetMealPlanOptions is a mock function.
func (m *MealPlanOptionDataManager) GetMealPlanOptions(ctx context.Context, mealPlanID, mealPlanEventID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.MealPlanOption], error) {
	args := m.Called(ctx, mealPlanID, mealPlanEventID, filter)
	return args.Get(0).(*types.QueryFilteredResult[types.MealPlanOption]), args.Error(1)
}

// CreateMealPlanOption is a mock function.
func (m *MealPlanOptionDataManager) CreateMealPlanOption(ctx context.Context, input *types.MealPlanOptionDatabaseCreationInput) (*types.MealPlanOption, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*types.MealPlanOption), args.Error(1)
}

// UpdateMealPlanOption is a mock function.
func (m *MealPlanOptionDataManager) UpdateMealPlanOption(ctx context.Context, updated *types.MealPlanOption) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveMealPlanOption is a mock function.
func (m *MealPlanOptionDataManager) ArchiveMealPlanOption(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string) error {
	return m.Called(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID).Error(0)
}

// FinalizeMealPlanOption is a mock function.
func (m *MealPlanOptionDataManager) FinalizeMealPlanOption(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID, householdID string) (bool, error) {
	args := m.Called(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, householdID)

	return args.Bool(0), args.Error(1)
}
