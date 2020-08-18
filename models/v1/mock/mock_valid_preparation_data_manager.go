package mock

import (
	"context"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/mock"
)

var _ models.ValidPreparationDataManager = (*ValidPreparationDataManager)(nil)

// ValidPreparationDataManager is a mocked models.ValidPreparationDataManager for testing.
type ValidPreparationDataManager struct {
	mock.Mock
}

// ValidPreparationExists is a mock function.
func (m *ValidPreparationDataManager) ValidPreparationExists(ctx context.Context, validPreparationID uint64) (bool, error) {
	args := m.Called(ctx, validPreparationID)
	return args.Bool(0), args.Error(1)
}

// GetValidPreparation is a mock function.
func (m *ValidPreparationDataManager) GetValidPreparation(ctx context.Context, validPreparationID uint64) (*models.ValidPreparation, error) {
	args := m.Called(ctx, validPreparationID)
	return args.Get(0).(*models.ValidPreparation), args.Error(1)
}

// GetAllValidPreparationsCount is a mock function.
func (m *ValidPreparationDataManager) GetAllValidPreparationsCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetAllValidPreparations is a mock function.
func (m *ValidPreparationDataManager) GetAllValidPreparations(ctx context.Context, results chan []models.ValidPreparation) error {
	args := m.Called(ctx, results)
	return args.Error(0)
}

// GetValidPreparations is a mock function.
func (m *ValidPreparationDataManager) GetValidPreparations(ctx context.Context, filter *models.QueryFilter) (*models.ValidPreparationList, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*models.ValidPreparationList), args.Error(1)
}

// GetValidPreparationsWithIDs is a mock function.
func (m *ValidPreparationDataManager) GetValidPreparationsWithIDs(ctx context.Context, limit uint8, ids []uint64) ([]models.ValidPreparation, error) {
	args := m.Called(ctx, limit, ids)
	return args.Get(0).([]models.ValidPreparation), args.Error(1)
}

// CreateValidPreparation is a mock function.
func (m *ValidPreparationDataManager) CreateValidPreparation(ctx context.Context, input *models.ValidPreparationCreationInput) (*models.ValidPreparation, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*models.ValidPreparation), args.Error(1)
}

// UpdateValidPreparation is a mock function.
func (m *ValidPreparationDataManager) UpdateValidPreparation(ctx context.Context, updated *models.ValidPreparation) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveValidPreparation is a mock function.
func (m *ValidPreparationDataManager) ArchiveValidPreparation(ctx context.Context, validPreparationID uint64) error {
	return m.Called(ctx, validPreparationID).Error(0)
}
