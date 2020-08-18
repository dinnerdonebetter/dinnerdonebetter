package mock

import (
	"context"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/mock"
)

var _ models.RecipeStepIngredientDataManager = (*RecipeStepIngredientDataManager)(nil)

// RecipeStepIngredientDataManager is a mocked models.RecipeStepIngredientDataManager for testing.
type RecipeStepIngredientDataManager struct {
	mock.Mock
}

// RecipeStepIngredientExists is a mock function.
func (m *RecipeStepIngredientDataManager) RecipeStepIngredientExists(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID uint64) (bool, error) {
	args := m.Called(ctx, recipeID, recipeStepID, recipeStepIngredientID)
	return args.Bool(0), args.Error(1)
}

// GetRecipeStepIngredient is a mock function.
func (m *RecipeStepIngredientDataManager) GetRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID uint64) (*models.RecipeStepIngredient, error) {
	args := m.Called(ctx, recipeID, recipeStepID, recipeStepIngredientID)
	return args.Get(0).(*models.RecipeStepIngredient), args.Error(1)
}

// GetAllRecipeStepIngredientsCount is a mock function.
func (m *RecipeStepIngredientDataManager) GetAllRecipeStepIngredientsCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetAllRecipeStepIngredients is a mock function.
func (m *RecipeStepIngredientDataManager) GetAllRecipeStepIngredients(ctx context.Context, results chan []models.RecipeStepIngredient) error {
	args := m.Called(ctx, results)
	return args.Error(0)
}

// GetRecipeStepIngredients is a mock function.
func (m *RecipeStepIngredientDataManager) GetRecipeStepIngredients(ctx context.Context, recipeID, recipeStepID uint64, filter *models.QueryFilter) (*models.RecipeStepIngredientList, error) {
	args := m.Called(ctx, recipeID, recipeStepID, filter)
	return args.Get(0).(*models.RecipeStepIngredientList), args.Error(1)
}

// GetRecipeStepIngredientsWithIDs is a mock function.
func (m *RecipeStepIngredientDataManager) GetRecipeStepIngredientsWithIDs(ctx context.Context, recipeID, recipeStepID uint64, limit uint8, ids []uint64) ([]models.RecipeStepIngredient, error) {
	args := m.Called(ctx, recipeID, recipeStepID, limit, ids)
	return args.Get(0).([]models.RecipeStepIngredient), args.Error(1)
}

// CreateRecipeStepIngredient is a mock function.
func (m *RecipeStepIngredientDataManager) CreateRecipeStepIngredient(ctx context.Context, input *models.RecipeStepIngredientCreationInput) (*models.RecipeStepIngredient, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*models.RecipeStepIngredient), args.Error(1)
}

// UpdateRecipeStepIngredient is a mock function.
func (m *RecipeStepIngredientDataManager) UpdateRecipeStepIngredient(ctx context.Context, updated *models.RecipeStepIngredient) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveRecipeStepIngredient is a mock function.
func (m *RecipeStepIngredientDataManager) ArchiveRecipeStepIngredient(ctx context.Context, recipeStepID, recipeStepIngredientID uint64) error {
	return m.Called(ctx, recipeStepID, recipeStepIngredientID).Error(0)
}
