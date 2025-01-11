package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.RecipeStepIngredientDataManager = (*RecipeStepIngredientDataManagerMock)(nil)

// RecipeStepIngredientDataManagerMock is a mocked types.RecipeStepIngredientDataManager for testing.
type RecipeStepIngredientDataManagerMock struct {
	mock.Mock
}

// RecipeStepIngredientExists is a mock function.
func (m *RecipeStepIngredientDataManagerMock) RecipeStepIngredientExists(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) (bool, error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID, recipeStepIngredientID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetRecipeStepIngredient is a mock function.
func (m *RecipeStepIngredientDataManagerMock) GetRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) (*types.RecipeStepIngredient, error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID, recipeStepIngredientID)
	return returnValues.Get(0).(*types.RecipeStepIngredient), returnValues.Error(1)
}

// GetRecipeStepIngredients is a mock function.
func (m *RecipeStepIngredientDataManagerMock) GetRecipeStepIngredients(ctx context.Context, recipeID, recipeStepID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.RecipeStepIngredient], error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID, filter)
	return returnValues.Get(0).(*types.QueryFilteredResult[types.RecipeStepIngredient]), returnValues.Error(1)
}

// CreateRecipeStepIngredient is a mock function.
func (m *RecipeStepIngredientDataManagerMock) CreateRecipeStepIngredient(ctx context.Context, input *types.RecipeStepIngredientDatabaseCreationInput) (*types.RecipeStepIngredient, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*types.RecipeStepIngredient), returnValues.Error(1)
}

// UpdateRecipeStepIngredient is a mock function.
func (m *RecipeStepIngredientDataManagerMock) UpdateRecipeStepIngredient(ctx context.Context, updated *types.RecipeStepIngredient) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveRecipeStepIngredient is a mock function.
func (m *RecipeStepIngredientDataManagerMock) ArchiveRecipeStepIngredient(ctx context.Context, recipeStepID, recipeStepIngredientID string) error {
	return m.Called(ctx, recipeStepID, recipeStepIngredientID).Error(0)
}
