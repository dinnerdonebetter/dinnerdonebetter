package mocktypes

import (
	"context"

	"github.com/stretchr/testify/mock"

	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

var _ types.ValidIngredientDataManager = (*ValidIngredientDataManager)(nil)

// ValidIngredientDataManager is a mocked types.ValidIngredientDataManager for testing.
type ValidIngredientDataManager struct {
	mock.Mock
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

// GetTotalValidIngredientCount is a mock function.
func (m *ValidIngredientDataManager) GetTotalValidIngredientCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetValidIngredients is a mock function.
func (m *ValidIngredientDataManager) GetValidIngredients(ctx context.Context, filter *types.QueryFilter) (*types.ValidIngredientList, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*types.ValidIngredientList), args.Error(1)
}

// GetValidIngredientsWithIDs is a mock function.
func (m *ValidIngredientDataManager) GetValidIngredientsWithIDs(ctx context.Context, limit uint8, ids []string) ([]*types.ValidIngredient, error) {
	args := m.Called(ctx, limit, ids)
	return args.Get(0).([]*types.ValidIngredient), args.Error(1)
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
