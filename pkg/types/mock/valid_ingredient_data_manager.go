package mock

import (
	"context"

	"gitlab.com/prixfixe/prixfixe/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.ValidIngredientDataManager = (*ValidIngredientDataManager)(nil)

// ValidIngredientDataManager is a mocked types.ValidIngredientDataManager for testing.
type ValidIngredientDataManager struct {
	mock.Mock
}

// ValidIngredientExists is a mock function.
func (m *ValidIngredientDataManager) ValidIngredientExists(ctx context.Context, validIngredientID uint64) (bool, error) {
	args := m.Called(ctx, validIngredientID)
	return args.Bool(0), args.Error(1)
}

// GetValidIngredient is a mock function.
func (m *ValidIngredientDataManager) GetValidIngredient(ctx context.Context, validIngredientID uint64) (*types.ValidIngredient, error) {
	args := m.Called(ctx, validIngredientID)
	return args.Get(0).(*types.ValidIngredient), args.Error(1)
}

// GetAllValidIngredientsCount is a mock function.
func (m *ValidIngredientDataManager) GetAllValidIngredientsCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// SearchForValidIngredients is a mock function.
func (m *ValidIngredientDataManager) SearchForValidIngredients(ctx context.Context, sessionCtxData *types.SessionContextData, validPreparationID uint64, query string, filter *types.QueryFilter) ([]*types.ValidIngredient, error) {
	args := m.Called(ctx, sessionCtxData, validPreparationID, query, filter)
	return args.Get(0).([]*types.ValidIngredient), args.Error(1)
}

// GetAllValidIngredients is a mock function.
func (m *ValidIngredientDataManager) GetAllValidIngredients(ctx context.Context, results chan []*types.ValidIngredient, bucketSize uint16) error {
	args := m.Called(ctx, results, bucketSize)
	return args.Error(0)
}

// GetValidIngredients is a mock function.
func (m *ValidIngredientDataManager) GetValidIngredients(ctx context.Context, filter *types.QueryFilter) (*types.ValidIngredientList, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*types.ValidIngredientList), args.Error(1)
}

// GetValidIngredientsWithIDs is a mock function.
func (m *ValidIngredientDataManager) GetValidIngredientsWithIDs(ctx context.Context, limit uint8, validPreparationID uint64, ids []uint64) ([]*types.ValidIngredient, error) {
	args := m.Called(ctx, limit, validPreparationID, ids)
	return args.Get(0).([]*types.ValidIngredient), args.Error(1)
}

// CreateValidIngredient is a mock function.
func (m *ValidIngredientDataManager) CreateValidIngredient(ctx context.Context, input *types.ValidIngredientCreationInput, createdByUser uint64) (*types.ValidIngredient, error) {
	args := m.Called(ctx, input, createdByUser)
	return args.Get(0).(*types.ValidIngredient), args.Error(1)
}

// UpdateValidIngredient is a mock function.
func (m *ValidIngredientDataManager) UpdateValidIngredient(ctx context.Context, updated *types.ValidIngredient, changedByUser uint64, changes []*types.FieldChangeSummary) error {
	return m.Called(ctx, updated, changedByUser, changes).Error(0)
}

// ArchiveValidIngredient is a mock function.
func (m *ValidIngredientDataManager) ArchiveValidIngredient(ctx context.Context, validIngredientID, archivedBy uint64) error {
	return m.Called(ctx, validIngredientID, archivedBy).Error(0)
}

// GetAuditLogEntriesForValidIngredient is a mock function.
func (m *ValidIngredientDataManager) GetAuditLogEntriesForValidIngredient(ctx context.Context, validIngredientID uint64) ([]*types.AuditLogEntry, error) {
	args := m.Called(ctx, validIngredientID)
	return args.Get(0).([]*types.AuditLogEntry), args.Error(1)
}
