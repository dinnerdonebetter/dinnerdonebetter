package mocktypes

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/pkg/types"
)

var _ types.MealPlanOptionDataManager = (*MealPlanOptionDataManager)(nil)

// MealPlanOptionDataManager is a mocked types.MealPlanOptionDataManager for testing.
type MealPlanOptionDataManager struct {
	mock.Mock
}

// MealPlanOptionExists is a mock function.
func (m *MealPlanOptionDataManager) MealPlanOptionExists(ctx context.Context, mealPlanID, mealPlanOptionID string) (bool, error) {
	args := m.Called(ctx, mealPlanID, mealPlanOptionID)
	return args.Bool(0), args.Error(1)
}

// GetMealPlanOption is a mock function.
func (m *MealPlanOptionDataManager) GetMealPlanOption(ctx context.Context, mealPlanID, mealPlanOptionID string) (*types.MealPlanOption, error) {
	args := m.Called(ctx, mealPlanID, mealPlanOptionID)
	return args.Get(0).(*types.MealPlanOption), args.Error(1)
}

// GetTotalMealPlanOptionCount is a mock function.
func (m *MealPlanOptionDataManager) GetTotalMealPlanOptionCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetMealPlanOptions is a mock function.
func (m *MealPlanOptionDataManager) GetMealPlanOptions(ctx context.Context, mealPlanID string, filter *types.QueryFilter) (*types.MealPlanOptionList, error) {
	args := m.Called(ctx, mealPlanID, filter)
	return args.Get(0).(*types.MealPlanOptionList), args.Error(1)
}

// GetMealPlanOptionsWithIDs is a mock function.
func (m *MealPlanOptionDataManager) GetMealPlanOptionsWithIDs(ctx context.Context, mealPlanID string, limit uint8, ids []string) ([]*types.MealPlanOption, error) {
	args := m.Called(ctx, mealPlanID, limit, ids)
	return args.Get(0).([]*types.MealPlanOption), args.Error(1)
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
func (m *MealPlanOptionDataManager) ArchiveMealPlanOption(ctx context.Context, mealPlanID, mealPlanOptionID string) error {
	return m.Called(ctx, mealPlanID, mealPlanOptionID).Error(0)
}

// FinalizeMealPlanOption is a mock function.
func (m *MealPlanOptionDataManager) FinalizeMealPlanOption(ctx context.Context, mealPlanID, mealPlanOptionID, householdID string, winnerRequired bool) (bool, error) {
	args := m.Called(ctx, mealPlanID, mealPlanOptionID, householdID, winnerRequired)

	return args.Bool(0), args.Error(1)
}
