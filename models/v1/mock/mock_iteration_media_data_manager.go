package mock

import (
	"context"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/mock"
)

var _ models.IterationMediaDataManager = (*IterationMediaDataManager)(nil)

// IterationMediaDataManager is a mocked models.IterationMediaDataManager for testing.
type IterationMediaDataManager struct {
	mock.Mock
}

// IterationMediaExists is a mock function.
func (m *IterationMediaDataManager) IterationMediaExists(ctx context.Context, recipeID, recipeIterationID, iterationMediaID uint64) (bool, error) {
	args := m.Called(ctx, recipeID, recipeIterationID, iterationMediaID)
	return args.Bool(0), args.Error(1)
}

// GetIterationMedia is a mock function.
func (m *IterationMediaDataManager) GetIterationMedia(ctx context.Context, recipeID, recipeIterationID, iterationMediaID uint64) (*models.IterationMedia, error) {
	args := m.Called(ctx, recipeID, recipeIterationID, iterationMediaID)
	return args.Get(0).(*models.IterationMedia), args.Error(1)
}

// GetAllIterationMediasCount is a mock function.
func (m *IterationMediaDataManager) GetAllIterationMediasCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetAllIterationMedias is a mock function.
func (m *IterationMediaDataManager) GetAllIterationMedias(ctx context.Context, results chan []models.IterationMedia) error {
	args := m.Called(ctx, results)
	return args.Error(0)
}

// GetIterationMedias is a mock function.
func (m *IterationMediaDataManager) GetIterationMedias(ctx context.Context, recipeID, recipeIterationID uint64, filter *models.QueryFilter) (*models.IterationMediaList, error) {
	args := m.Called(ctx, recipeID, recipeIterationID, filter)
	return args.Get(0).(*models.IterationMediaList), args.Error(1)
}

// GetIterationMediasWithIDs is a mock function.
func (m *IterationMediaDataManager) GetIterationMediasWithIDs(ctx context.Context, recipeID, recipeIterationID uint64, limit uint8, ids []uint64) ([]models.IterationMedia, error) {
	args := m.Called(ctx, recipeID, recipeIterationID, limit, ids)
	return args.Get(0).([]models.IterationMedia), args.Error(1)
}

// CreateIterationMedia is a mock function.
func (m *IterationMediaDataManager) CreateIterationMedia(ctx context.Context, input *models.IterationMediaCreationInput) (*models.IterationMedia, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*models.IterationMedia), args.Error(1)
}

// UpdateIterationMedia is a mock function.
func (m *IterationMediaDataManager) UpdateIterationMedia(ctx context.Context, updated *models.IterationMedia) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveIterationMedia is a mock function.
func (m *IterationMediaDataManager) ArchiveIterationMedia(ctx context.Context, recipeIterationID, iterationMediaID uint64) error {
	return m.Called(ctx, recipeIterationID, iterationMediaID).Error(0)
}
