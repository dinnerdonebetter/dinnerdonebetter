package mock

import (
	"context"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/mock"
)

var _ models.RecipeTagDataManager = (*RecipeTagDataManager)(nil)

// RecipeTagDataManager is a mocked models.RecipeTagDataManager for testing.
type RecipeTagDataManager struct {
	mock.Mock
}

// RecipeTagExists is a mock function.
func (m *RecipeTagDataManager) RecipeTagExists(ctx context.Context, recipeID, recipeTagID uint64) (bool, error) {
	args := m.Called(ctx, recipeID, recipeTagID)
	return args.Bool(0), args.Error(1)
}

// GetRecipeTag is a mock function.
func (m *RecipeTagDataManager) GetRecipeTag(ctx context.Context, recipeID, recipeTagID uint64) (*models.RecipeTag, error) {
	args := m.Called(ctx, recipeID, recipeTagID)
	return args.Get(0).(*models.RecipeTag), args.Error(1)
}

// GetAllRecipeTagsCount is a mock function.
func (m *RecipeTagDataManager) GetAllRecipeTagsCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetRecipeTags is a mock function.
func (m *RecipeTagDataManager) GetRecipeTags(ctx context.Context, recipeID uint64, filter *models.QueryFilter) (*models.RecipeTagList, error) {
	args := m.Called(ctx, recipeID, filter)
	return args.Get(0).(*models.RecipeTagList), args.Error(1)
}

// CreateRecipeTag is a mock function.
func (m *RecipeTagDataManager) CreateRecipeTag(ctx context.Context, input *models.RecipeTagCreationInput) (*models.RecipeTag, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*models.RecipeTag), args.Error(1)
}

// UpdateRecipeTag is a mock function.
func (m *RecipeTagDataManager) UpdateRecipeTag(ctx context.Context, updated *models.RecipeTag) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveRecipeTag is a mock function.
func (m *RecipeTagDataManager) ArchiveRecipeTag(ctx context.Context, recipeID, recipeTagID uint64) error {
	return m.Called(ctx, recipeID, recipeTagID).Error(0)
}
