package mock

import (
	"context"

	"gitlab.com/prixfixe/prixfixe/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.ValidInstrumentDataManager = (*ValidInstrumentDataManager)(nil)

// ValidInstrumentDataManager is a mocked types.ValidInstrumentDataManager for testing.
type ValidInstrumentDataManager struct {
	mock.Mock
}

// ValidInstrumentExists is a mock function.
func (m *ValidInstrumentDataManager) ValidInstrumentExists(ctx context.Context, validInstrumentID uint64) (bool, error) {
	args := m.Called(ctx, validInstrumentID)
	return args.Bool(0), args.Error(1)
}

// GetValidInstrument is a mock function.
func (m *ValidInstrumentDataManager) GetValidInstrument(ctx context.Context, validInstrumentID uint64) (*types.ValidInstrument, error) {
	args := m.Called(ctx, validInstrumentID)
	return args.Get(0).(*types.ValidInstrument), args.Error(1)
}

// GetAllValidInstrumentsCount is a mock function.
func (m *ValidInstrumentDataManager) GetAllValidInstrumentsCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// SearchForValidInstruments is a mock function.
func (m *ValidInstrumentDataManager) SearchForValidInstruments(ctx context.Context, sessionCtxData *types.SessionContextData, query string, filter *types.QueryFilter) ([]*types.ValidInstrument, error) {
	args := m.Called(ctx, sessionCtxData, query, filter)
	return args.Get(0).([]*types.ValidInstrument), args.Error(1)
}

// GetAllValidInstruments is a mock function.
func (m *ValidInstrumentDataManager) GetAllValidInstruments(ctx context.Context, results chan []*types.ValidInstrument, bucketSize uint16) error {
	args := m.Called(ctx, results, bucketSize)
	return args.Error(0)
}

// GetValidInstruments is a mock function.
func (m *ValidInstrumentDataManager) GetValidInstruments(ctx context.Context, filter *types.QueryFilter) (*types.ValidInstrumentList, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*types.ValidInstrumentList), args.Error(1)
}

// GetValidInstrumentsWithIDs is a mock function.
func (m *ValidInstrumentDataManager) GetValidInstrumentsWithIDs(ctx context.Context, limit uint8, ids []uint64) ([]*types.ValidInstrument, error) {
	args := m.Called(ctx, limit, ids)
	return args.Get(0).([]*types.ValidInstrument), args.Error(1)
}

// CreateValidInstrument is a mock function.
func (m *ValidInstrumentDataManager) CreateValidInstrument(ctx context.Context, input *types.ValidInstrumentCreationInput, createdByUser uint64) (*types.ValidInstrument, error) {
	args := m.Called(ctx, input, createdByUser)
	return args.Get(0).(*types.ValidInstrument), args.Error(1)
}

// UpdateValidInstrument is a mock function.
func (m *ValidInstrumentDataManager) UpdateValidInstrument(ctx context.Context, updated *types.ValidInstrument, changedByUser uint64, changes []*types.FieldChangeSummary) error {
	return m.Called(ctx, updated, changedByUser, changes).Error(0)
}

// ArchiveValidInstrument is a mock function.
func (m *ValidInstrumentDataManager) ArchiveValidInstrument(ctx context.Context, validInstrumentID, archivedBy uint64) error {
	return m.Called(ctx, validInstrumentID, archivedBy).Error(0)
}

// GetAuditLogEntriesForValidInstrument is a mock function.
func (m *ValidInstrumentDataManager) GetAuditLogEntriesForValidInstrument(ctx context.Context, validInstrumentID uint64) ([]*types.AuditLogEntry, error) {
	args := m.Called(ctx, validInstrumentID)
	return args.Get(0).([]*types.AuditLogEntry), args.Error(1)
}
