package mock

import (
	"context"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/mock"
)

var _ models.RequiredPreparationInstrumentDataManager = (*RequiredPreparationInstrumentDataManager)(nil)

// RequiredPreparationInstrumentDataManager is a mocked models.RequiredPreparationInstrumentDataManager for testing
type RequiredPreparationInstrumentDataManager struct {
	mock.Mock
}

// GetRequiredPreparationInstrument is a mock function
func (m *RequiredPreparationInstrumentDataManager) GetRequiredPreparationInstrument(ctx context.Context, requiredPreparationInstrumentID uint64) (*models.RequiredPreparationInstrument, error) {
	args := m.Called(ctx, requiredPreparationInstrumentID)
	return args.Get(0).(*models.RequiredPreparationInstrument), args.Error(1)
}

// GetRequiredPreparationInstrumentCount is a mock function
func (m *RequiredPreparationInstrumentDataManager) GetRequiredPreparationInstrumentCount(ctx context.Context, filter *models.QueryFilter) (uint64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(uint64), args.Error(1)
}

// GetAllRequiredPreparationInstrumentsCount is a mock function
func (m *RequiredPreparationInstrumentDataManager) GetAllRequiredPreparationInstrumentsCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetRequiredPreparationInstruments is a mock function
func (m *RequiredPreparationInstrumentDataManager) GetRequiredPreparationInstruments(ctx context.Context, filter *models.QueryFilter) (*models.RequiredPreparationInstrumentList, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*models.RequiredPreparationInstrumentList), args.Error(1)
}

// CreateRequiredPreparationInstrument is a mock function
func (m *RequiredPreparationInstrumentDataManager) CreateRequiredPreparationInstrument(ctx context.Context, input *models.RequiredPreparationInstrumentCreationInput) (*models.RequiredPreparationInstrument, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*models.RequiredPreparationInstrument), args.Error(1)
}

// UpdateRequiredPreparationInstrument is a mock function
func (m *RequiredPreparationInstrumentDataManager) UpdateRequiredPreparationInstrument(ctx context.Context, updated *models.RequiredPreparationInstrument) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveRequiredPreparationInstrument is a mock function
func (m *RequiredPreparationInstrumentDataManager) ArchiveRequiredPreparationInstrument(ctx context.Context, id uint64) error {
	return m.Called(ctx, id).Error(0)
}
