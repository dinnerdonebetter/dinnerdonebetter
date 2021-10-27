package mocktypes

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/pkg/types"
)

var _ types.ValidInstrumentDataManager = (*ValidInstrumentDataManager)(nil)

// ValidInstrumentDataManager is a mocked types.ValidInstrumentDataManager for testing.
type ValidInstrumentDataManager struct {
	mock.Mock
}

// ValidInstrumentExists is a mock function.
func (m *ValidInstrumentDataManager) ValidInstrumentExists(ctx context.Context, validInstrumentID string) (bool, error) {
	args := m.Called(ctx, validInstrumentID)
	return args.Bool(0), args.Error(1)
}

// GetValidInstrument is a mock function.
func (m *ValidInstrumentDataManager) GetValidInstrument(ctx context.Context, validInstrumentID string) (*types.ValidInstrument, error) {
	args := m.Called(ctx, validInstrumentID)
	return args.Get(0).(*types.ValidInstrument), args.Error(1)
}

// GetTotalValidInstrumentCount is a mock function.
func (m *ValidInstrumentDataManager) GetTotalValidInstrumentCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetValidInstruments is a mock function.
func (m *ValidInstrumentDataManager) GetValidInstruments(ctx context.Context, filter *types.QueryFilter) (*types.ValidInstrumentList, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*types.ValidInstrumentList), args.Error(1)
}

// GetValidInstrumentsWithIDs is a mock function.
func (m *ValidInstrumentDataManager) GetValidInstrumentsWithIDs(ctx context.Context, limit uint8, ids []string) ([]*types.ValidInstrument, error) {
	args := m.Called(ctx, limit, ids)
	return args.Get(0).([]*types.ValidInstrument), args.Error(1)
}

// CreateValidInstrument is a mock function.
func (m *ValidInstrumentDataManager) CreateValidInstrument(ctx context.Context, input *types.ValidInstrumentDatabaseCreationInput) (*types.ValidInstrument, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*types.ValidInstrument), args.Error(1)
}

// UpdateValidInstrument is a mock function.
func (m *ValidInstrumentDataManager) UpdateValidInstrument(ctx context.Context, updated *types.ValidInstrument) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveValidInstrument is a mock function.
func (m *ValidInstrumentDataManager) ArchiveValidInstrument(ctx context.Context, validInstrumentID string) error {
	return m.Called(ctx, validInstrumentID).Error(0)
}
