package mock

import (
	"context"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/mock"
)

var _ models.RequiredPreparationInstrumentDataManager = (*RequiredPreparationInstrumentDataManager)(nil)

// RequiredPreparationInstrumentDataManager is a mocked models.RequiredPreparationInstrumentDataManager for testing.
type RequiredPreparationInstrumentDataManager struct {
	mock.Mock
}

// RequiredPreparationInstrumentExists is a mock function.
func (m *RequiredPreparationInstrumentDataManager) RequiredPreparationInstrumentExists(ctx context.Context, requiredPreparationInstrumentID uint64) (bool, error) {
	args := m.Called(ctx, requiredPreparationInstrumentID)
	return args.Bool(0), args.Error(1)
}

// GetRequiredPreparationInstrument is a mock function.
func (m *RequiredPreparationInstrumentDataManager) GetRequiredPreparationInstrument(ctx context.Context, requiredPreparationInstrumentID uint64) (*models.RequiredPreparationInstrument, error) {
	args := m.Called(ctx, requiredPreparationInstrumentID)
	return args.Get(0).(*models.RequiredPreparationInstrument), args.Error(1)
}

// GetAllRequiredPreparationInstrumentsCount is a mock function.
func (m *RequiredPreparationInstrumentDataManager) GetAllRequiredPreparationInstrumentsCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetAllRequiredPreparationInstruments is a mock function.
func (m *RequiredPreparationInstrumentDataManager) GetAllRequiredPreparationInstruments(ctx context.Context, results chan []models.RequiredPreparationInstrument) error {
	args := m.Called(ctx, results)
	return args.Error(0)
}

// GetRequiredPreparationInstruments is a mock function.
func (m *RequiredPreparationInstrumentDataManager) GetRequiredPreparationInstruments(ctx context.Context, filter *models.QueryFilter) (*models.RequiredPreparationInstrumentList, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*models.RequiredPreparationInstrumentList), args.Error(1)
}

// GetRequiredPreparationInstrumentsWithIDs is a mock function.
func (m *RequiredPreparationInstrumentDataManager) GetRequiredPreparationInstrumentsWithIDs(ctx context.Context, limit uint8, ids []uint64) ([]models.RequiredPreparationInstrument, error) {
	args := m.Called(ctx, limit, ids)
	return args.Get(0).([]models.RequiredPreparationInstrument), args.Error(1)
}

// CreateRequiredPreparationInstrument is a mock function.
func (m *RequiredPreparationInstrumentDataManager) CreateRequiredPreparationInstrument(ctx context.Context, input *models.RequiredPreparationInstrumentCreationInput) (*models.RequiredPreparationInstrument, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*models.RequiredPreparationInstrument), args.Error(1)
}

// UpdateRequiredPreparationInstrument is a mock function.
func (m *RequiredPreparationInstrumentDataManager) UpdateRequiredPreparationInstrument(ctx context.Context, updated *models.RequiredPreparationInstrument) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveRequiredPreparationInstrument is a mock function.
func (m *RequiredPreparationInstrumentDataManager) ArchiveRequiredPreparationInstrument(ctx context.Context, requiredPreparationInstrumentID uint64) error {
	return m.Called(ctx, requiredPreparationInstrumentID).Error(0)
}
