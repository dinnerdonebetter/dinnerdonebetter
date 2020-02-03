package mock

import (
	"context"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/mock"
)

var _ models.RecipeStepEventDataManager = (*RecipeStepEventDataManager)(nil)

// RecipeStepEventDataManager is a mocked models.RecipeStepEventDataManager for testing
type RecipeStepEventDataManager struct {
	mock.Mock
}

// GetRecipeStepEvent is a mock function
func (m *RecipeStepEventDataManager) GetRecipeStepEvent(ctx context.Context, recipeStepEventID, userID uint64) (*models.RecipeStepEvent, error) {
	args := m.Called(ctx, recipeStepEventID, userID)
	return args.Get(0).(*models.RecipeStepEvent), args.Error(1)
}

// GetRecipeStepEventCount is a mock function
func (m *RecipeStepEventDataManager) GetRecipeStepEventCount(ctx context.Context, filter *models.QueryFilter, userID uint64) (uint64, error) {
	args := m.Called(ctx, filter, userID)
	return args.Get(0).(uint64), args.Error(1)
}

// GetAllRecipeStepEventsCount is a mock function
func (m *RecipeStepEventDataManager) GetAllRecipeStepEventsCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetRecipeStepEvents is a mock function
func (m *RecipeStepEventDataManager) GetRecipeStepEvents(ctx context.Context, filter *models.QueryFilter, userID uint64) (*models.RecipeStepEventList, error) {
	args := m.Called(ctx, filter, userID)
	return args.Get(0).(*models.RecipeStepEventList), args.Error(1)
}

// GetAllRecipeStepEventsForUser is a mock function
func (m *RecipeStepEventDataManager) GetAllRecipeStepEventsForUser(ctx context.Context, userID uint64) ([]models.RecipeStepEvent, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.RecipeStepEvent), args.Error(1)
}

// CreateRecipeStepEvent is a mock function
func (m *RecipeStepEventDataManager) CreateRecipeStepEvent(ctx context.Context, input *models.RecipeStepEventCreationInput) (*models.RecipeStepEvent, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*models.RecipeStepEvent), args.Error(1)
}

// UpdateRecipeStepEvent is a mock function
func (m *RecipeStepEventDataManager) UpdateRecipeStepEvent(ctx context.Context, updated *models.RecipeStepEvent) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveRecipeStepEvent is a mock function
func (m *RecipeStepEventDataManager) ArchiveRecipeStepEvent(ctx context.Context, id, userID uint64) error {
	return m.Called(ctx, id, userID).Error(0)
}
