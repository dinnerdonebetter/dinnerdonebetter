package mocktypes

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/pkg/types"
)

var _ types.MealDataManager = (*MealDataManager)(nil)

// MealDataManager is a mocked types.MealDataManager for testing.
type MealDataManager struct {
	mock.Mock
}

// MealExists is a mock function.
func (m *MealDataManager) MealExists(ctx context.Context, recipeID string) (bool, error) {
	args := m.Called(ctx, recipeID)
	return args.Bool(0), args.Error(1)
}

// GetMeal is a mock function.
func (m *MealDataManager) GetMeal(ctx context.Context, recipeID string) (*types.Meal, error) {
	args := m.Called(ctx, recipeID)
	return args.Get(0).(*types.Meal), args.Error(1)
}

// GetMealByIDAndUser is a mock function.
func (m *MealDataManager) GetMealByIDAndUser(ctx context.Context, recipeID, userID string) (*types.Meal, error) {
	args := m.Called(ctx, recipeID, userID)
	return args.Get(0).(*types.Meal), args.Error(1)
}

// GetTotalMealCount is a mock function.
func (m *MealDataManager) GetTotalMealCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetMeals is a mock function.
func (m *MealDataManager) GetMeals(ctx context.Context, filter *types.QueryFilter) (*types.MealList, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*types.MealList), args.Error(1)
}

// SearchForMeals is a mock function.
func (m *MealDataManager) SearchForMeals(ctx context.Context, query string, filter *types.QueryFilter) (*types.MealList, error) {
	args := m.Called(ctx, query, filter)
	return args.Get(0).(*types.MealList), args.Error(1)
}

// GetMealsWithIDs is a mock function.
func (m *MealDataManager) GetMealsWithIDs(ctx context.Context, householdID string, limit uint8, ids []string) ([]*types.Meal, error) {
	args := m.Called(ctx, householdID, limit, ids)
	return args.Get(0).([]*types.Meal), args.Error(1)
}

// CreateMeal is a mock function.
func (m *MealDataManager) CreateMeal(ctx context.Context, input *types.MealDatabaseCreationInput) (*types.Meal, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*types.Meal), args.Error(1)
}

// UpdateMeal is a mock function.
func (m *MealDataManager) UpdateMeal(ctx context.Context, updated *types.Meal) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveMeal is a mock function.
func (m *MealDataManager) ArchiveMeal(ctx context.Context, recipeID, householdID string) error {
	return m.Called(ctx, recipeID, householdID).Error(0)
}
