package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.RecipeDataManager = (*RecipeDataManagerMock)(nil)

// RecipeDataManagerMock is a mocked types.RecipeDataManager for testing.
type RecipeDataManagerMock struct {
	mock.Mock
}

// RecipeExists is a mock function.
func (m *RecipeDataManagerMock) RecipeExists(ctx context.Context, recipeID string) (bool, error) {
	returnValues := m.Called(ctx, recipeID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetRecipe is a mock function.
func (m *RecipeDataManagerMock) GetRecipe(ctx context.Context, recipeID string) (*types.Recipe, error) {
	returnValues := m.Called(ctx, recipeID)
	return returnValues.Get(0).(*types.Recipe), returnValues.Error(1)
}

// SearchForRecipes is a mock function.
func (m *RecipeDataManagerMock) SearchForRecipes(ctx context.Context, query string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.Recipe], error) {
	returnValues := m.Called(ctx, query, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[types.Recipe]), returnValues.Error(1)
}

// GetRecipes is a mock function.
func (m *RecipeDataManagerMock) GetRecipes(ctx context.Context, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.Recipe], error) {
	returnValues := m.Called(ctx, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[types.Recipe]), returnValues.Error(1)
}

// GetRecipesCreatedByUser is a mock function.
func (m *RecipeDataManagerMock) GetRecipesCreatedByUser(ctx context.Context, userID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.Recipe], error) {
	returnValues := m.Called(ctx, userID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[types.Recipe]), returnValues.Error(1)
}

// CreateRecipe is a mock function.
func (m *RecipeDataManagerMock) CreateRecipe(ctx context.Context, input *types.RecipeDatabaseCreationInput) (*types.Recipe, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*types.Recipe), returnValues.Error(1)
}

// UpdateRecipe is a mock function.
func (m *RecipeDataManagerMock) UpdateRecipe(ctx context.Context, updated *types.Recipe) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveRecipe is a mock function.
func (m *RecipeDataManagerMock) ArchiveRecipe(ctx context.Context, recipeID, householdID string) error {
	return m.Called(ctx, recipeID, householdID).Error(0)
}

// MarkRecipeAsIndexed is a mock function.
func (m *RecipeDataManagerMock) MarkRecipeAsIndexed(ctx context.Context, recipeID string) error {
	return m.Called(ctx, recipeID).Error(0)
}

// GetRecipeIDsThatNeedSearchIndexing is a mock function.
func (m *RecipeDataManagerMock) GetRecipeIDsThatNeedSearchIndexing(ctx context.Context) ([]string, error) {
	returnValues := m.Called(ctx)
	return returnValues.Get(0).([]string), returnValues.Error(1)
}

// GetRecipesWithIDs is a mock function.
func (m *RecipeDataManagerMock) GetRecipesWithIDs(ctx context.Context, ids []string) ([]*types.Recipe, error) {
	returnValues := m.Called(ctx, ids)
	return returnValues.Get(0).([]*types.Recipe), returnValues.Error(1)
}
