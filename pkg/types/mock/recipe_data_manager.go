package mocktypes

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/pkg/types"
)

var _ types.RecipeDataManager = (*RecipeDataManager)(nil)

// RecipeDataManager is a mocked types.RecipeDataManager for testing.
type RecipeDataManager struct {
	mock.Mock
}

// RecipeExists is a mock function.
func (m *RecipeDataManager) RecipeExists(ctx context.Context, recipeID string) (bool, error) {
	args := m.Called(ctx, recipeID)
	return args.Bool(0), args.Error(1)
}

// GetRecipe is a mock function.
func (m *RecipeDataManager) GetRecipe(ctx context.Context, recipeID string) (*types.Recipe, error) {
	args := m.Called(ctx, recipeID)
	return args.Get(0).(*types.Recipe), args.Error(1)
}

// GetRecipeByIDAndUser is a mock function.
func (m *RecipeDataManager) GetRecipeByIDAndUser(ctx context.Context, recipeID, userID string) (*types.Recipe, error) {
	args := m.Called(ctx, recipeID, userID)
	return args.Get(0).(*types.Recipe), args.Error(1)
}

// SearchForRecipes is a mock function.
func (m *RecipeDataManager) SearchForRecipes(ctx context.Context, query string, filter *types.QueryFilter) (*types.RecipeList, error) {
	args := m.Called(ctx, query, filter)
	return args.Get(0).(*types.RecipeList), args.Error(1)
}

// GetRecipes is a mock function.
func (m *RecipeDataManager) GetRecipes(ctx context.Context, filter *types.QueryFilter) (*types.RecipeList, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*types.RecipeList), args.Error(1)
}

// CreateRecipe is a mock function.
func (m *RecipeDataManager) CreateRecipe(ctx context.Context, input *types.RecipeDatabaseCreationInput) (*types.Recipe, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*types.Recipe), args.Error(1)
}

// UpdateRecipe is a mock function.
func (m *RecipeDataManager) UpdateRecipe(ctx context.Context, updated *types.Recipe) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveRecipe is a mock function.
func (m *RecipeDataManager) ArchiveRecipe(ctx context.Context, recipeID, householdID string) error {
	return m.Called(ctx, recipeID, householdID).Error(0)
}
