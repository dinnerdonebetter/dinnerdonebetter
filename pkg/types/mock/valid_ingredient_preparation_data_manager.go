package mock

import (
	"context"

	"gitlab.com/prixfixe/prixfixe/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.ValidIngredientPreparationDataManager = (*ValidIngredientPreparationDataManager)(nil)

// ValidIngredientPreparationDataManager is a mocked types.ValidIngredientPreparationDataManager for testing.
type ValidIngredientPreparationDataManager struct {
	mock.Mock
}

// ValidIngredientPreparationExists is a mock function.
func (m *ValidIngredientPreparationDataManager) ValidIngredientPreparationExists(ctx context.Context, validIngredientPreparationID uint64) (bool, error) {
	args := m.Called(ctx, validIngredientPreparationID)
	return args.Bool(0), args.Error(1)
}

// GetValidIngredientPreparation is a mock function.
func (m *ValidIngredientPreparationDataManager) GetValidIngredientPreparation(ctx context.Context, validIngredientPreparationID uint64) (*types.ValidIngredientPreparation, error) {
	args := m.Called(ctx, validIngredientPreparationID)
	return args.Get(0).(*types.ValidIngredientPreparation), args.Error(1)
}

// GetAllValidIngredientPreparationsCount is a mock function.
func (m *ValidIngredientPreparationDataManager) GetAllValidIngredientPreparationsCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetAllValidIngredientPreparations is a mock function.
func (m *ValidIngredientPreparationDataManager) GetAllValidIngredientPreparations(ctx context.Context, results chan []*types.ValidIngredientPreparation, bucketSize uint16) error {
	args := m.Called(ctx, results, bucketSize)
	return args.Error(0)
}

// GetValidIngredientPreparations is a mock function.
func (m *ValidIngredientPreparationDataManager) GetValidIngredientPreparations(ctx context.Context, filter *types.QueryFilter) (*types.ValidIngredientPreparationList, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*types.ValidIngredientPreparationList), args.Error(1)
}

// GetValidIngredientPreparationsWithIDs is a mock function.
func (m *ValidIngredientPreparationDataManager) GetValidIngredientPreparationsWithIDs(ctx context.Context, limit uint8, ids []uint64) ([]*types.ValidIngredientPreparation, error) {
	args := m.Called(ctx, limit, ids)
	return args.Get(0).([]*types.ValidIngredientPreparation), args.Error(1)
}

// CreateValidIngredientPreparation is a mock function.
func (m *ValidIngredientPreparationDataManager) CreateValidIngredientPreparation(ctx context.Context, input *types.ValidIngredientPreparationCreationInput, createdByUser uint64) (*types.ValidIngredientPreparation, error) {
	args := m.Called(ctx, input, createdByUser)
	return args.Get(0).(*types.ValidIngredientPreparation), args.Error(1)
}

// UpdateValidIngredientPreparation is a mock function.
func (m *ValidIngredientPreparationDataManager) UpdateValidIngredientPreparation(ctx context.Context, updated *types.ValidIngredientPreparation, changedByUser uint64, changes []*types.FieldChangeSummary) error {
	return m.Called(ctx, updated, changedByUser, changes).Error(0)
}

// ArchiveValidIngredientPreparation is a mock function.
func (m *ValidIngredientPreparationDataManager) ArchiveValidIngredientPreparation(ctx context.Context, validIngredientPreparationID, archivedBy uint64) error {
	return m.Called(ctx, validIngredientPreparationID, archivedBy).Error(0)
}

// GetAuditLogEntriesForValidIngredientPreparation is a mock function.
func (m *ValidIngredientPreparationDataManager) GetAuditLogEntriesForValidIngredientPreparation(ctx context.Context, validIngredientPreparationID uint64) ([]*types.AuditLogEntry, error) {
	args := m.Called(ctx, validIngredientPreparationID)
	return args.Get(0).([]*types.AuditLogEntry), args.Error(1)
}
