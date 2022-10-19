package grocerylistpreparation

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/pkg/types"
)

var _ GroceryListCreator = (*MockGroceryListCreator)(nil)

// MockGroceryListCreator is a mock GroceryListCreator.
type MockGroceryListCreator struct {
	mock.Mock
}

// GenerateGroceryListInputs is a mock function.
func (m *MockGroceryListCreator) GenerateGroceryListInputs(ctx context.Context, mealPlan *types.MealPlan) ([]*types.MealPlanGroceryListItemDatabaseCreationInput, error) {
	returnValues := m.Called(ctx, mealPlan)

	return returnValues.Get(0).([]*types.MealPlanGroceryListItemDatabaseCreationInput), returnValues.Error(1)
}
