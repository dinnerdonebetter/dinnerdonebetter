package mock

import (
	"context"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/mock"
)

var _ models.RecipeStepEventDataManager = (*RecipeStepEventDataManager)(nil)

// RecipeStepEventDataManager is a mocked models.RecipeStepEventDataManager for testing.
type RecipeStepEventDataManager struct {
	mock.Mock
}

// RecipeStepEventExists is a mock function.
func (m *RecipeStepEventDataManager) RecipeStepEventExists(ctx context.Context, recipeID, recipeStepID, recipeStepEventID uint64) (bool, error) {
	args := m.Called(ctx, recipeID, recipeStepID, recipeStepEventID)
	return args.Bool(0), args.Error(1)
}

// GetRecipeStepEvent is a mock function.
func (m *RecipeStepEventDataManager) GetRecipeStepEvent(ctx context.Context, recipeID, recipeStepID, recipeStepEventID uint64) (*models.RecipeStepEvent, error) {
	args := m.Called(ctx, recipeID, recipeStepID, recipeStepEventID)
	return args.Get(0).(*models.RecipeStepEvent), args.Error(1)
}

// GetAllRecipeStepEventsCount is a mock function.
func (m *RecipeStepEventDataManager) GetAllRecipeStepEventsCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetAllRecipeStepEvents is a mock function.
func (m *RecipeStepEventDataManager) GetAllRecipeStepEvents(ctx context.Context, results chan []models.RecipeStepEvent) error {
	args := m.Called(ctx, results)
	return args.Error(0)
}

// GetRecipeStepEvents is a mock function.
func (m *RecipeStepEventDataManager) GetRecipeStepEvents(ctx context.Context, recipeID, recipeStepID uint64, filter *models.QueryFilter) (*models.RecipeStepEventList, error) {
	args := m.Called(ctx, recipeID, recipeStepID, filter)
	return args.Get(0).(*models.RecipeStepEventList), args.Error(1)
}

// GetRecipeStepEventsWithIDs is a mock function.
func (m *RecipeStepEventDataManager) GetRecipeStepEventsWithIDs(ctx context.Context, recipeID, recipeStepID uint64, limit uint8, ids []uint64) ([]models.RecipeStepEvent, error) {
	args := m.Called(ctx, recipeID, recipeStepID, limit, ids)
	return args.Get(0).([]models.RecipeStepEvent), args.Error(1)
}

// CreateRecipeStepEvent is a mock function.
func (m *RecipeStepEventDataManager) CreateRecipeStepEvent(ctx context.Context, input *models.RecipeStepEventCreationInput) (*models.RecipeStepEvent, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*models.RecipeStepEvent), args.Error(1)
}

// UpdateRecipeStepEvent is a mock function.
func (m *RecipeStepEventDataManager) UpdateRecipeStepEvent(ctx context.Context, updated *models.RecipeStepEvent) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveRecipeStepEvent is a mock function.
func (m *RecipeStepEventDataManager) ArchiveRecipeStepEvent(ctx context.Context, recipeStepID, recipeStepEventID uint64) error {
	return m.Called(ctx, recipeStepID, recipeStepEventID).Error(0)
}
