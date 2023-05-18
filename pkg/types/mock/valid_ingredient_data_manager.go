package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.ValidIngredientDataManager = (*ValidIngredientDataManager)(nil)

// ValidIngredientDataManager is a mocked types.ValidIngredientDataManager for testing.
type ValidIngredientDataManager struct {
	mock.Mock
}

func (m *ValidIngredientDataManager) SearchForValidIngredientsForIngredientState(ctx context.Context, ingredientStateID, query string, filter *types.QueryFilter) ([]*types.ValidIngredient, error) {
	args := m.Called(ctx, ingredientStateID, query, filter)
	return args.Get(0).([]*types.ValidIngredient), args.Error(1)
}

// ValidIngredientExists is a mock function.
func (m *ValidIngredientDataManager) ValidIngredientExists(ctx context.Context, validIngredientID string) (bool, error) {
	args := m.Called(ctx, validIngredientID)
	return args.Bool(0), args.Error(1)
}

// GetValidIngredient is a mock function.
func (m *ValidIngredientDataManager) GetValidIngredient(ctx context.Context, validIngredientID string) (*types.ValidIngredient, error) {
	args := m.Called(ctx, validIngredientID)
	return args.Get(0).(*types.ValidIngredient), args.Error(1)
}

// GetRandomValidIngredient is a mock function.
func (m *ValidIngredientDataManager) GetRandomValidIngredient(ctx context.Context) (*types.ValidIngredient, error) {
	args := m.Called(ctx)
	return args.Get(0).(*types.ValidIngredient), args.Error(1)
}

// SearchForValidIngredients is a mock function.
func (m *ValidIngredientDataManager) SearchForValidIngredients(ctx context.Context, query string, filter *types.QueryFilter) ([]*types.ValidIngredient, error) {
	args := m.Called(ctx, query, filter)
	return args.Get(0).([]*types.ValidIngredient), args.Error(1)
}

// SearchForValidIngredientsForPreparation is a mock function.
func (m *ValidIngredientDataManager) SearchForValidIngredientsForPreparation(ctx context.Context, preparationID, query string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidIngredient], error) {
	args := m.Called(ctx, preparationID, query, filter)
	return args.Get(0).(*types.QueryFilteredResult[types.ValidIngredient]), args.Error(1)
}

// GetValidIngredients is a mock function.
func (m *ValidIngredientDataManager) GetValidIngredients(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidIngredient], error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*types.QueryFilteredResult[types.ValidIngredient]), args.Error(1)
}

// CreateValidIngredient is a mock function.
func (m *ValidIngredientDataManager) CreateValidIngredient(ctx context.Context, input *types.ValidIngredientDatabaseCreationInput) (*types.ValidIngredient, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*types.ValidIngredient), args.Error(1)
}

// UpdateValidIngredient is a mock function.
func (m *ValidIngredientDataManager) UpdateValidIngredient(ctx context.Context, updated *types.ValidIngredient) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveValidIngredient is a mock function.
func (m *ValidIngredientDataManager) ArchiveValidIngredient(ctx context.Context, validIngredientID string) error {
	return m.Called(ctx, validIngredientID).Error(0)
}
