package mocktypes

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/backend/pkg/types"
)

var _ types.MealPlanGroceryListItemDataManager = (*MealPlanGroceryListItemDataManager)(nil)

// MealPlanGroceryListItemDataManager is a mocked types.MealPlanGroceryListItemDataManager for testing.
type MealPlanGroceryListItemDataManager struct {
	mock.Mock
}

// MealPlanGroceryListItemExists is a mock function.
func (m *MealPlanGroceryListItemDataManager) MealPlanGroceryListItemExists(ctx context.Context, mealPlanID, mealPlanGroceryListItemID string) (bool, error) {
	returnValues := m.Called(ctx, mealPlanID, mealPlanGroceryListItemID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetMealPlanGroceryListItem is a mock function.
func (m *MealPlanGroceryListItemDataManager) GetMealPlanGroceryListItem(ctx context.Context, mealPlanID, mealPlanGroceryListItemID string) (*types.MealPlanGroceryListItem, error) {
	returnValues := m.Called(ctx, mealPlanID, mealPlanGroceryListItemID)
	return returnValues.Get(0).(*types.MealPlanGroceryListItem), returnValues.Error(1)
}

// GetMealPlanGroceryListItemsForMealPlan is a mock function.
func (m *MealPlanGroceryListItemDataManager) GetMealPlanGroceryListItemsForMealPlan(ctx context.Context, mealPlanID string) ([]*types.MealPlanGroceryListItem, error) {
	returnValues := m.Called(ctx, mealPlanID)
	return returnValues.Get(0).([]*types.MealPlanGroceryListItem), returnValues.Error(1)
}

// CreateMealPlanGroceryListItem is a mock function.
func (m *MealPlanGroceryListItemDataManager) CreateMealPlanGroceryListItem(ctx context.Context, input *types.MealPlanGroceryListItemDatabaseCreationInput) (*types.MealPlanGroceryListItem, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*types.MealPlanGroceryListItem), returnValues.Error(1)
}

// UpdateMealPlanGroceryListItem is a mock function.
func (m *MealPlanGroceryListItemDataManager) UpdateMealPlanGroceryListItem(ctx context.Context, updated *types.MealPlanGroceryListItem) error {
	returnValues := m.Called(ctx, updated)
	return returnValues.Error(0)
}

// ArchiveMealPlanGroceryListItem is a mock function.
func (m *MealPlanGroceryListItemDataManager) ArchiveMealPlanGroceryListItem(ctx context.Context, mealPlanGroceryListItemID string) error {
	returnValues := m.Called(ctx, mealPlanGroceryListItemID)
	return returnValues.Error(0)
}

// CreateMealPlanGroceryListItemsForMealPlan is a mock function.
func (m *MealPlanGroceryListItemDataManager) CreateMealPlanGroceryListItemsForMealPlan(ctx context.Context, mealPlanID string, inputs []*types.MealPlanGroceryListItemDatabaseCreationInput) error {
	returnValues := m.Called(ctx, mealPlanID, inputs)
	return returnValues.Error(0)
}
