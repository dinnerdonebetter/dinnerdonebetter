package mock

import (
	"context"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/mock"
)

var _ models.RecipeStepIngredientDataManager = (*RecipeStepIngredientDataManager)(nil)

// RecipeStepIngredientDataManager is a mocked models.RecipeStepIngredientDataManager for testing
type RecipeStepIngredientDataManager struct {
	mock.Mock
}

// GetRecipeStepIngredient is a mock function
func (m *RecipeStepIngredientDataManager) GetRecipeStepIngredient(ctx context.Context, recipeStepIngredientID, userID uint64) (*models.RecipeStepIngredient, error) {
	args := m.Called(ctx, recipeStepIngredientID, userID)
	return args.Get(0).(*models.RecipeStepIngredient), args.Error(1)
}

// GetRecipeStepIngredientCount is a mock function
func (m *RecipeStepIngredientDataManager) GetRecipeStepIngredientCount(ctx context.Context, filter *models.QueryFilter, userID uint64) (uint64, error) {
	args := m.Called(ctx, filter, userID)
	return args.Get(0).(uint64), args.Error(1)
}

// GetAllRecipeStepIngredientsCount is a mock function
func (m *RecipeStepIngredientDataManager) GetAllRecipeStepIngredientsCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetRecipeStepIngredients is a mock function
func (m *RecipeStepIngredientDataManager) GetRecipeStepIngredients(ctx context.Context, filter *models.QueryFilter, userID uint64) (*models.RecipeStepIngredientList, error) {
	args := m.Called(ctx, filter, userID)
	return args.Get(0).(*models.RecipeStepIngredientList), args.Error(1)
}

// GetAllRecipeStepIngredientsForUser is a mock function
func (m *RecipeStepIngredientDataManager) GetAllRecipeStepIngredientsForUser(ctx context.Context, userID uint64) ([]models.RecipeStepIngredient, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.RecipeStepIngredient), args.Error(1)
}

// CreateRecipeStepIngredient is a mock function
func (m *RecipeStepIngredientDataManager) CreateRecipeStepIngredient(ctx context.Context, input *models.RecipeStepIngredientCreationInput) (*models.RecipeStepIngredient, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*models.RecipeStepIngredient), args.Error(1)
}

// UpdateRecipeStepIngredient is a mock function
func (m *RecipeStepIngredientDataManager) UpdateRecipeStepIngredient(ctx context.Context, updated *models.RecipeStepIngredient) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveRecipeStepIngredient is a mock function
func (m *RecipeStepIngredientDataManager) ArchiveRecipeStepIngredient(ctx context.Context, id, userID uint64) error {
	return m.Called(ctx, id, userID).Error(0)
}
