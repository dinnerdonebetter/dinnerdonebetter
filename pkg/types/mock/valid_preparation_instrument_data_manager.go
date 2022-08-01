package mocktypes

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/pkg/types"
)

var _ types.ValidPreparationInstrumentDataManager = (*ValidPreparationInstrumentDataManager)(nil)

// ValidPreparationInstrumentDataManager is a mocked types.ValidPreparationInstrumentDataManager for testing.
type ValidPreparationInstrumentDataManager struct {
	mock.Mock
}

// ValidPreparationInstrumentExists is a mock function.
func (m *ValidPreparationInstrumentDataManager) ValidPreparationInstrumentExists(ctx context.Context, validPreparationInstrumentID string) (bool, error) {
	args := m.Called(ctx, validPreparationInstrumentID)
	return args.Bool(0), args.Error(1)
}

// GetValidPreparationInstrument is a mock function.
func (m *ValidPreparationInstrumentDataManager) GetValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID string) (*types.ValidPreparationInstrument, error) {
	args := m.Called(ctx, validPreparationInstrumentID)
	return args.Get(0).(*types.ValidPreparationInstrument), args.Error(1)
}

// GetValidPreparationInstruments is a mock function.
func (m *ValidPreparationInstrumentDataManager) GetValidPreparationInstruments(ctx context.Context, filter *types.QueryFilter) (*types.ValidPreparationInstrumentList, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*types.ValidPreparationInstrumentList), args.Error(1)
}

// GetValidPreparationInstrumentsWithIDs is a mock function.
func (m *ValidPreparationInstrumentDataManager) GetValidPreparationInstrumentsWithIDs(ctx context.Context, limit uint8, ids []string) ([]*types.ValidPreparationInstrument, error) {
	args := m.Called(ctx, limit, ids)
	return args.Get(0).([]*types.ValidPreparationInstrument), args.Error(1)
}

// GetValidInstrumentsForPreparations is a mock function.
func (m *ValidPreparationInstrumentDataManager) GetValidInstrumentsForPreparations(ctx context.Context, preparationName string, filter *types.QueryFilter) (*types.ValidPreparationInstrumentList, error) {
	args := m.Called(ctx, preparationName, filter)
	return args.Get(0).(*types.ValidPreparationInstrumentList), args.Error(1)
}

// CreateValidPreparationInstrument is a mock function.
func (m *ValidPreparationInstrumentDataManager) CreateValidPreparationInstrument(ctx context.Context, input *types.ValidPreparationInstrumentDatabaseCreationInput) (*types.ValidPreparationInstrument, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*types.ValidPreparationInstrument), args.Error(1)
}

// UpdateValidPreparationInstrument is a mock function.
func (m *ValidPreparationInstrumentDataManager) UpdateValidPreparationInstrument(ctx context.Context, updated *types.ValidPreparationInstrument) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveValidPreparationInstrument is a mock function.
func (m *ValidPreparationInstrumentDataManager) ArchiveValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID string) error {
	return m.Called(ctx, validPreparationInstrumentID).Error(0)
}
