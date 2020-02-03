package mock

import (
	"context"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/mock"
)

var _ models.InstrumentDataManager = (*InstrumentDataManager)(nil)

// InstrumentDataManager is a mocked models.InstrumentDataManager for testing
type InstrumentDataManager struct {
	mock.Mock
}

// GetInstrument is a mock function
func (m *InstrumentDataManager) GetInstrument(ctx context.Context, instrumentID, userID uint64) (*models.Instrument, error) {
	args := m.Called(ctx, instrumentID, userID)
	return args.Get(0).(*models.Instrument), args.Error(1)
}

// GetInstrumentCount is a mock function
func (m *InstrumentDataManager) GetInstrumentCount(ctx context.Context, filter *models.QueryFilter, userID uint64) (uint64, error) {
	args := m.Called(ctx, filter, userID)
	return args.Get(0).(uint64), args.Error(1)
}

// GetAllInstrumentsCount is a mock function
func (m *InstrumentDataManager) GetAllInstrumentsCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetInstruments is a mock function
func (m *InstrumentDataManager) GetInstruments(ctx context.Context, filter *models.QueryFilter, userID uint64) (*models.InstrumentList, error) {
	args := m.Called(ctx, filter, userID)
	return args.Get(0).(*models.InstrumentList), args.Error(1)
}

// GetAllInstrumentsForUser is a mock function
func (m *InstrumentDataManager) GetAllInstrumentsForUser(ctx context.Context, userID uint64) ([]models.Instrument, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Instrument), args.Error(1)
}

// CreateInstrument is a mock function
func (m *InstrumentDataManager) CreateInstrument(ctx context.Context, input *models.InstrumentCreationInput) (*models.Instrument, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*models.Instrument), args.Error(1)
}

// UpdateInstrument is a mock function
func (m *InstrumentDataManager) UpdateInstrument(ctx context.Context, updated *models.Instrument) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveInstrument is a mock function
func (m *InstrumentDataManager) ArchiveInstrument(ctx context.Context, id, userID uint64) error {
	return m.Called(ctx, id, userID).Error(0)
}
