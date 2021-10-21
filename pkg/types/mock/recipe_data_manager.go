package mocktypes

import (
	"context"

	"github.com/stretchr/testify/mock"

	"gitlab.com/prixfixe/prixfixe/pkg/types"
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

// GetTotalRecipeCount is a mock function.
func (m *RecipeDataManager) GetTotalRecipeCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetRecipes is a mock function.
func (m *RecipeDataManager) GetRecipes(ctx context.Context, filter *types.QueryFilter) (*types.RecipeList, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*types.RecipeList), args.Error(1)
}

// GetRecipesWithIDs is a mock function.
func (m *RecipeDataManager) GetRecipesWithIDs(ctx context.Context, accountID string, limit uint8, ids []string) ([]*types.Recipe, error) {
	args := m.Called(ctx, accountID, limit, ids)
	return args.Get(0).([]*types.Recipe), args.Error(1)
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
func (m *RecipeDataManager) ArchiveRecipe(ctx context.Context, recipeID, accountID string) error {
	return m.Called(ctx, recipeID, accountID).Error(0)
}
