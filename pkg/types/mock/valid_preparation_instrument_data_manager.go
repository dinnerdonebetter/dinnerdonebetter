package mock

import (
	"context"

	"gitlab.com/prixfixe/prixfixe/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.ValidPreparationInstrumentDataManager = (*ValidPreparationInstrumentDataManager)(nil)

// ValidPreparationInstrumentDataManager is a mocked types.ValidPreparationInstrumentDataManager for testing.
type ValidPreparationInstrumentDataManager struct {
	mock.Mock
}

// ValidPreparationInstrumentExists is a mock function.
func (m *ValidPreparationInstrumentDataManager) ValidPreparationInstrumentExists(ctx context.Context, validPreparationInstrumentID uint64) (bool, error) {
	args := m.Called(ctx, validPreparationInstrumentID)
	return args.Bool(0), args.Error(1)
}

// GetValidPreparationInstrument is a mock function.
func (m *ValidPreparationInstrumentDataManager) GetValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID uint64) (*types.ValidPreparationInstrument, error) {
	args := m.Called(ctx, validPreparationInstrumentID)
	return args.Get(0).(*types.ValidPreparationInstrument), args.Error(1)
}

// GetAllValidPreparationInstrumentsCount is a mock function.
func (m *ValidPreparationInstrumentDataManager) GetAllValidPreparationInstrumentsCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetAllValidPreparationInstruments is a mock function.
func (m *ValidPreparationInstrumentDataManager) GetAllValidPreparationInstruments(ctx context.Context, results chan []*types.ValidPreparationInstrument, bucketSize uint16) error {
	args := m.Called(ctx, results, bucketSize)
	return args.Error(0)
}

// GetValidPreparationInstruments is a mock function.
func (m *ValidPreparationInstrumentDataManager) GetValidPreparationInstruments(ctx context.Context, filter *types.QueryFilter) (*types.ValidPreparationInstrumentList, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*types.ValidPreparationInstrumentList), args.Error(1)
}

// GetValidPreparationInstrumentsWithIDs is a mock function.
func (m *ValidPreparationInstrumentDataManager) GetValidPreparationInstrumentsWithIDs(ctx context.Context, limit uint8, ids []uint64) ([]*types.ValidPreparationInstrument, error) {
	args := m.Called(ctx, limit, ids)
	return args.Get(0).([]*types.ValidPreparationInstrument), args.Error(1)
}

// CreateValidPreparationInstrument is a mock function.
func (m *ValidPreparationInstrumentDataManager) CreateValidPreparationInstrument(ctx context.Context, input *types.ValidPreparationInstrumentCreationInput, createdByUser uint64) (*types.ValidPreparationInstrument, error) {
	args := m.Called(ctx, input, createdByUser)
	return args.Get(0).(*types.ValidPreparationInstrument), args.Error(1)
}

// UpdateValidPreparationInstrument is a mock function.
func (m *ValidPreparationInstrumentDataManager) UpdateValidPreparationInstrument(ctx context.Context, updated *types.ValidPreparationInstrument, changedByUser uint64, changes []*types.FieldChangeSummary) error {
	return m.Called(ctx, updated, changedByUser, changes).Error(0)
}

// ArchiveValidPreparationInstrument is a mock function.
func (m *ValidPreparationInstrumentDataManager) ArchiveValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID, archivedBy uint64) error {
	return m.Called(ctx, validPreparationInstrumentID, archivedBy).Error(0)
}

// GetAuditLogEntriesForValidPreparationInstrument is a mock function.
func (m *ValidPreparationInstrumentDataManager) GetAuditLogEntriesForValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID uint64) ([]*types.AuditLogEntry, error) {
	args := m.Called(ctx, validPreparationInstrumentID)
	return args.Get(0).([]*types.AuditLogEntry), args.Error(1)
}
