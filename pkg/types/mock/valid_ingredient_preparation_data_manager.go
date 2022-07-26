package mocktypes

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/pkg/types"
)

var _ types.ValidIngredientPreparationDataManager = (*ValidIngredientPreparationDataManager)(nil)

// ValidIngredientPreparationDataManager is a mocked types.ValidIngredientPreparationDataManager for testing.
type ValidIngredientPreparationDataManager struct {
	mock.Mock
}

// ValidIngredientPreparationExists is a mock function.
func (m *ValidIngredientPreparationDataManager) ValidIngredientPreparationExists(ctx context.Context, validIngredientPreparationID string) (bool, error) {
	args := m.Called(ctx, validIngredientPreparationID)
	return args.Bool(0), args.Error(1)
}

// GetValidIngredientPreparation is a mock function.
func (m *ValidIngredientPreparationDataManager) GetValidIngredientPreparation(ctx context.Context, validIngredientPreparationID string) (*types.ValidIngredientPreparation, error) {
	args := m.Called(ctx, validIngredientPreparationID)
	return args.Get(0).(*types.ValidIngredientPreparation), args.Error(1)
}

// GetValidIngredientPreparations is a mock function.
func (m *ValidIngredientPreparationDataManager) GetValidIngredientPreparations(ctx context.Context, filter *types.QueryFilter) (*types.ValidIngredientPreparationList, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*types.ValidIngredientPreparationList), args.Error(1)
}

// GetValidIngredientPreparationsWithIDs is a mock function.
func (m *ValidIngredientPreparationDataManager) GetValidIngredientPreparationsWithIDs(ctx context.Context, limit uint8, ids []string) ([]*types.ValidIngredientPreparation, error) {
	args := m.Called(ctx, limit, ids)
	return args.Get(0).([]*types.ValidIngredientPreparation), args.Error(1)
}

// CreateValidIngredientPreparation is a mock function.
func (m *ValidIngredientPreparationDataManager) CreateValidIngredientPreparation(ctx context.Context, input *types.ValidIngredientPreparationDatabaseCreationInput) (*types.ValidIngredientPreparation, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*types.ValidIngredientPreparation), args.Error(1)
}

// UpdateValidIngredientPreparation is a mock function.
func (m *ValidIngredientPreparationDataManager) UpdateValidIngredientPreparation(ctx context.Context, updated *types.ValidIngredientPreparation) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveValidIngredientPreparation is a mock function.
func (m *ValidIngredientPreparationDataManager) ArchiveValidIngredientPreparation(ctx context.Context, validIngredientPreparationID string) error {
	return m.Called(ctx, validIngredientPreparationID).Error(0)
}
