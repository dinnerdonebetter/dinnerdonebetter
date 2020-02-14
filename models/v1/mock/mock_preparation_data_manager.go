package mock

import (
	"context"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/mock"
)

var _ models.PreparationDataManager = (*PreparationDataManager)(nil)

// PreparationDataManager is a mocked models.PreparationDataManager for testing
type PreparationDataManager struct {
	mock.Mock
}

// GetPreparation is a mock function
func (m *PreparationDataManager) GetPreparation(ctx context.Context, preparationID uint64) (*models.Preparation, error) {
	args := m.Called(ctx, preparationID)
	return args.Get(0).(*models.Preparation), args.Error(1)
}

// GetPreparationCount is a mock function
func (m *PreparationDataManager) GetPreparationCount(ctx context.Context, filter *models.QueryFilter) (uint64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(uint64), args.Error(1)
}

// GetAllPreparationsCount is a mock function
func (m *PreparationDataManager) GetAllPreparationsCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetPreparations is a mock function
func (m *PreparationDataManager) GetPreparations(ctx context.Context, filter *models.QueryFilter) (*models.PreparationList, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*models.PreparationList), args.Error(1)
}

// CreatePreparation is a mock function
func (m *PreparationDataManager) CreatePreparation(ctx context.Context, input *models.PreparationCreationInput) (*models.Preparation, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*models.Preparation), args.Error(1)
}

// UpdatePreparation is a mock function
func (m *PreparationDataManager) UpdatePreparation(ctx context.Context, updated *models.Preparation) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchivePreparation is a mock function
func (m *PreparationDataManager) ArchivePreparation(ctx context.Context, id uint64) error {
	return m.Called(ctx, id).Error(0)
}
