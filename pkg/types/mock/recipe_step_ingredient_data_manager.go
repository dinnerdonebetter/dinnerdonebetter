package mocktypes

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/backend/pkg/types"
)

var _ types.RecipeStepIngredientDataManager = (*RecipeStepIngredientDataManager)(nil)

// RecipeStepIngredientDataManager is a mocked types.RecipeStepIngredientDataManager for testing.
type RecipeStepIngredientDataManager struct {
	mock.Mock
}

// RecipeStepIngredientExists is a mock function.
func (m *RecipeStepIngredientDataManager) RecipeStepIngredientExists(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) (bool, error) {
	args := m.Called(ctx, recipeID, recipeStepID, recipeStepIngredientID)
	return args.Bool(0), args.Error(1)
}

// GetRecipeStepIngredient is a mock function.
func (m *RecipeStepIngredientDataManager) GetRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) (*types.RecipeStepIngredient, error) {
	args := m.Called(ctx, recipeID, recipeStepID, recipeStepIngredientID)
	return args.Get(0).(*types.RecipeStepIngredient), args.Error(1)
}

// GetRecipeStepIngredients is a mock function.
func (m *RecipeStepIngredientDataManager) GetRecipeStepIngredients(ctx context.Context, recipeID, recipeStepID string, filter *types.QueryFilter) (*types.RecipeStepIngredientList, error) {
	args := m.Called(ctx, recipeID, recipeStepID, filter)
	return args.Get(0).(*types.RecipeStepIngredientList), args.Error(1)
}

// CreateRecipeStepIngredient is a mock function.
func (m *RecipeStepIngredientDataManager) CreateRecipeStepIngredient(ctx context.Context, input *types.RecipeStepIngredientDatabaseCreationInput) (*types.RecipeStepIngredient, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*types.RecipeStepIngredient), args.Error(1)
}

// UpdateRecipeStepIngredient is a mock function.
func (m *RecipeStepIngredientDataManager) UpdateRecipeStepIngredient(ctx context.Context, updated *types.RecipeStepIngredient) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveRecipeStepIngredient is a mock function.
func (m *RecipeStepIngredientDataManager) ArchiveRecipeStepIngredient(ctx context.Context, recipeStepID, recipeStepIngredientID string) error {
	return m.Called(ctx, recipeStepID, recipeStepIngredientID).Error(0)
}
