package mocktypes

import (
	"context"

	"github.com/prixfixeco/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.ValidIngredientStateIngredientDataManager = (*ValidIngredientStateIngredientDataManager)(nil)

// ValidIngredientStateIngredientDataManager is a mocked types.ValidIngredientStateIngredientDataManager for testing.
type ValidIngredientStateIngredientDataManager struct {
	mock.Mock
}

// ValidIngredientStateIngredientExists is a mock function.
func (m *ValidIngredientStateIngredientDataManager) ValidIngredientStateIngredientExists(ctx context.Context, validIngredientStateIngredientID string) (bool, error) {
	args := m.Called(ctx, validIngredientStateIngredientID)
	return args.Bool(0), args.Error(1)
}

// GetValidIngredientStateIngredient is a mock function.
func (m *ValidIngredientStateIngredientDataManager) GetValidIngredientStateIngredient(ctx context.Context, validIngredientStateIngredientID string) (*types.ValidIngredientStateIngredient, error) {
	args := m.Called(ctx, validIngredientStateIngredientID)
	return args.Get(0).(*types.ValidIngredientStateIngredient), args.Error(1)
}

// GetValidIngredientStateIngredients is a mock function.
func (m *ValidIngredientStateIngredientDataManager) GetValidIngredientStateIngredients(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidIngredientStateIngredient], error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*types.QueryFilteredResult[types.ValidIngredientStateIngredient]), args.Error(1)
}

// GetValidIngredientStateIngredientsForIngredient is a mock function.
func (m *ValidIngredientStateIngredientDataManager) GetValidIngredientStateIngredientsForIngredient(ctx context.Context, ingredientID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidIngredientStateIngredient], error) {
	args := m.Called(ctx, ingredientID, filter)
	return args.Get(0).(*types.QueryFilteredResult[types.ValidIngredientStateIngredient]), args.Error(1)
}

// GetValidIngredientStateIngredientsForIngredientState is a mock function.
func (m *ValidIngredientStateIngredientDataManager) GetValidIngredientStateIngredientsForIngredientState(ctx context.Context, ingredientStateID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidIngredientStateIngredient], error) {
	args := m.Called(ctx, ingredientStateID, filter)
	return args.Get(0).(*types.QueryFilteredResult[types.ValidIngredientStateIngredient]), args.Error(1)
}

// CreateValidIngredientStateIngredient is a mock function.
func (m *ValidIngredientStateIngredientDataManager) CreateValidIngredientStateIngredient(ctx context.Context, input *types.ValidIngredientStateIngredientDatabaseCreationInput) (*types.ValidIngredientStateIngredient, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*types.ValidIngredientStateIngredient), args.Error(1)
}

// UpdateValidIngredientStateIngredient is a mock function.
func (m *ValidIngredientStateIngredientDataManager) UpdateValidIngredientStateIngredient(ctx context.Context, updated *types.ValidIngredientStateIngredient) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveValidIngredientStateIngredient is a mock function.
func (m *ValidIngredientStateIngredientDataManager) ArchiveValidIngredientStateIngredient(ctx context.Context, validIngredientStateIngredientID string) error {
	return m.Called(ctx, validIngredientStateIngredientID).Error(0)
}
