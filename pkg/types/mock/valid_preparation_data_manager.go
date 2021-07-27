package mock

import (
	"context"

	"gitlab.com/prixfixe/prixfixe/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.ValidPreparationDataManager = (*ValidPreparationDataManager)(nil)

// ValidPreparationDataManager is a mocked types.ValidPreparationDataManager for testing.
type ValidPreparationDataManager struct {
	mock.Mock
}

// ValidPreparationExists is a mock function.
func (m *ValidPreparationDataManager) ValidPreparationExists(ctx context.Context, validPreparationID uint64) (bool, error) {
	args := m.Called(ctx, validPreparationID)
	return args.Bool(0), args.Error(1)
}

// GetValidPreparation is a mock function.
func (m *ValidPreparationDataManager) GetValidPreparation(ctx context.Context, validPreparationID uint64) (*types.ValidPreparation, error) {
	args := m.Called(ctx, validPreparationID)
	return args.Get(0).(*types.ValidPreparation), args.Error(1)
}

// GetAllValidPreparationsCount is a mock function.
func (m *ValidPreparationDataManager) GetAllValidPreparationsCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetAllValidPreparations is a mock function.
func (m *ValidPreparationDataManager) GetAllValidPreparations(ctx context.Context, results chan []*types.ValidPreparation, bucketSize uint16) error {
	args := m.Called(ctx, results, bucketSize)
	return args.Error(0)
}

// GetValidPreparationsFromQuery is a mock function.
func (m *ValidPreparationDataManager) GetValidPreparationsFromQuery(ctx context.Context, query string) (x *types.ValidPreparationList, err error) {
	args := m.Called(ctx, query)
	return args.Get(0).(*types.ValidPreparationList), args.Error(1)
}

// SearchForValidPreparations is a mock function.
func (m *ValidPreparationDataManager) SearchForValidPreparations(ctx context.Context, sessionCtxData *types.SessionContextData, query string, filter *types.QueryFilter) ([]*types.ValidPreparation, error) {
	args := m.Called(ctx, sessionCtxData, query, filter)
	return args.Get(0).([]*types.ValidPreparation), args.Error(1)
}

// GetValidPreparations is a mock function.
func (m *ValidPreparationDataManager) GetValidPreparations(ctx context.Context, filter *types.QueryFilter) (*types.ValidPreparationList, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*types.ValidPreparationList), args.Error(1)
}

// GetValidPreparationsWithIDs is a mock function.
func (m *ValidPreparationDataManager) GetValidPreparationsWithIDs(ctx context.Context, limit uint8, ids []uint64) ([]*types.ValidPreparation, error) {
	args := m.Called(ctx, limit, ids)
	return args.Get(0).([]*types.ValidPreparation), args.Error(1)
}

// CreateValidPreparation is a mock function.
func (m *ValidPreparationDataManager) CreateValidPreparation(ctx context.Context, input *types.ValidPreparationCreationInput, createdByUser uint64) (*types.ValidPreparation, error) {
	args := m.Called(ctx, input, createdByUser)
	return args.Get(0).(*types.ValidPreparation), args.Error(1)
}

// UpdateValidPreparation is a mock function.
func (m *ValidPreparationDataManager) UpdateValidPreparation(ctx context.Context, updated *types.ValidPreparation, changedByUser uint64, changes []*types.FieldChangeSummary) error {
	return m.Called(ctx, updated, changedByUser, changes).Error(0)
}

// ArchiveValidPreparation is a mock function.
func (m *ValidPreparationDataManager) ArchiveValidPreparation(ctx context.Context, validPreparationID, archivedBy uint64) error {
	return m.Called(ctx, validPreparationID, archivedBy).Error(0)
}

// GetAuditLogEntriesForValidPreparation is a mock function.
func (m *ValidPreparationDataManager) GetAuditLogEntriesForValidPreparation(ctx context.Context, validPreparationID uint64) ([]*types.AuditLogEntry, error) {
	args := m.Called(ctx, validPreparationID)
	return args.Get(0).([]*types.AuditLogEntry), args.Error(1)
}
