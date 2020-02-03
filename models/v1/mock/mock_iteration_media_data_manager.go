package mock

import (
	"context"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/mock"
)

var _ models.IterationMediaDataManager = (*IterationMediaDataManager)(nil)

// IterationMediaDataManager is a mocked models.IterationMediaDataManager for testing
type IterationMediaDataManager struct {
	mock.Mock
}

// GetIterationMedia is a mock function
func (m *IterationMediaDataManager) GetIterationMedia(ctx context.Context, iterationMediaID, userID uint64) (*models.IterationMedia, error) {
	args := m.Called(ctx, iterationMediaID, userID)
	return args.Get(0).(*models.IterationMedia), args.Error(1)
}

// GetIterationMediaCount is a mock function
func (m *IterationMediaDataManager) GetIterationMediaCount(ctx context.Context, filter *models.QueryFilter, userID uint64) (uint64, error) {
	args := m.Called(ctx, filter, userID)
	return args.Get(0).(uint64), args.Error(1)
}

// GetAllIterationMediasCount is a mock function
func (m *IterationMediaDataManager) GetAllIterationMediasCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetIterationMedias is a mock function
func (m *IterationMediaDataManager) GetIterationMedias(ctx context.Context, filter *models.QueryFilter, userID uint64) (*models.IterationMediaList, error) {
	args := m.Called(ctx, filter, userID)
	return args.Get(0).(*models.IterationMediaList), args.Error(1)
}

// GetAllIterationMediasForUser is a mock function
func (m *IterationMediaDataManager) GetAllIterationMediasForUser(ctx context.Context, userID uint64) ([]models.IterationMedia, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.IterationMedia), args.Error(1)
}

// CreateIterationMedia is a mock function
func (m *IterationMediaDataManager) CreateIterationMedia(ctx context.Context, input *models.IterationMediaCreationInput) (*models.IterationMedia, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*models.IterationMedia), args.Error(1)
}

// UpdateIterationMedia is a mock function
func (m *IterationMediaDataManager) UpdateIterationMedia(ctx context.Context, updated *models.IterationMedia) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveIterationMedia is a mock function
func (m *IterationMediaDataManager) ArchiveIterationMedia(ctx context.Context, id, userID uint64) error {
	return m.Called(ctx, id, userID).Error(0)
}
