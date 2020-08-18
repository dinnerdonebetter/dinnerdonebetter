package mock

import (
	"context"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/mock"
)

var _ models.ValidInstrumentDataManager = (*ValidInstrumentDataManager)(nil)

// ValidInstrumentDataManager is a mocked models.ValidInstrumentDataManager for testing.
type ValidInstrumentDataManager struct {
	mock.Mock
}

// ValidInstrumentExists is a mock function.
func (m *ValidInstrumentDataManager) ValidInstrumentExists(ctx context.Context, validInstrumentID uint64) (bool, error) {
	args := m.Called(ctx, validInstrumentID)
	return args.Bool(0), args.Error(1)
}

// GetValidInstrument is a mock function.
func (m *ValidInstrumentDataManager) GetValidInstrument(ctx context.Context, validInstrumentID uint64) (*models.ValidInstrument, error) {
	args := m.Called(ctx, validInstrumentID)
	return args.Get(0).(*models.ValidInstrument), args.Error(1)
}

// GetAllValidInstrumentsCount is a mock function.
func (m *ValidInstrumentDataManager) GetAllValidInstrumentsCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetAllValidInstruments is a mock function.
func (m *ValidInstrumentDataManager) GetAllValidInstruments(ctx context.Context, results chan []models.ValidInstrument) error {
	args := m.Called(ctx, results)
	return args.Error(0)
}

// GetValidInstruments is a mock function.
func (m *ValidInstrumentDataManager) GetValidInstruments(ctx context.Context, filter *models.QueryFilter) (*models.ValidInstrumentList, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*models.ValidInstrumentList), args.Error(1)
}

// GetValidInstrumentsWithIDs is a mock function.
func (m *ValidInstrumentDataManager) GetValidInstrumentsWithIDs(ctx context.Context, limit uint8, ids []uint64) ([]models.ValidInstrument, error) {
	args := m.Called(ctx, limit, ids)
	return args.Get(0).([]models.ValidInstrument), args.Error(1)
}

// CreateValidInstrument is a mock function.
func (m *ValidInstrumentDataManager) CreateValidInstrument(ctx context.Context, input *models.ValidInstrumentCreationInput) (*models.ValidInstrument, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*models.ValidInstrument), args.Error(1)
}

// UpdateValidInstrument is a mock function.
func (m *ValidInstrumentDataManager) UpdateValidInstrument(ctx context.Context, updated *models.ValidInstrument) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveValidInstrument is a mock function.
func (m *ValidInstrumentDataManager) ArchiveValidInstrument(ctx context.Context, validInstrumentID uint64) error {
	return m.Called(ctx, validInstrumentID).Error(0)
}
