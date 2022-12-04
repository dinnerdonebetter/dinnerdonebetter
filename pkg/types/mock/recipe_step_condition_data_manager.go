package mocktypes

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/backend/pkg/types"
)

var _ types.RecipeStepConditionDataManager = (*RecipeStepConditionDataManager)(nil)

// RecipeStepConditionDataManager is a mocked types.RecipeStepConditionDataManager for testing.
type RecipeStepConditionDataManager struct {
	mock.Mock
}

// RecipeStepConditionExists is a mock function.
func (m *RecipeStepConditionDataManager) RecipeStepConditionExists(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) (bool, error) {
	args := m.Called(ctx, recipeID, recipeStepID, recipeStepIngredientID)
	return args.Bool(0), args.Error(1)
}

// GetRecipeStepCondition is a mock function.
func (m *RecipeStepConditionDataManager) GetRecipeStepCondition(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) (*types.RecipeStepCondition, error) {
	args := m.Called(ctx, recipeID, recipeStepID, recipeStepIngredientID)
	return args.Get(0).(*types.RecipeStepCondition), args.Error(1)
}

// GetRecipeStepConditions is a mock function.
func (m *RecipeStepConditionDataManager) GetRecipeStepConditions(ctx context.Context, recipeID, recipeStepID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.RecipeStepCondition], error) {
	args := m.Called(ctx, recipeID, recipeStepID, filter)
	return args.Get(0).(*types.QueryFilteredResult[types.RecipeStepCondition]), args.Error(1)
}

// CreateRecipeStepCondition is a mock function.
func (m *RecipeStepConditionDataManager) CreateRecipeStepCondition(ctx context.Context, input *types.RecipeStepConditionDatabaseCreationInput) (*types.RecipeStepCondition, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*types.RecipeStepCondition), args.Error(1)
}

// UpdateRecipeStepCondition is a mock function.
func (m *RecipeStepConditionDataManager) UpdateRecipeStepCondition(ctx context.Context, updated *types.RecipeStepCondition) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveRecipeStepCondition is a mock function.
func (m *RecipeStepConditionDataManager) ArchiveRecipeStepCondition(ctx context.Context, recipeStepID, recipeStepIngredientID string) error {
	return m.Called(ctx, recipeStepID, recipeStepIngredientID).Error(0)
}
