package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.RecipeStepCompletionConditionDataManager = (*RecipeStepCompletionConditionDataManagerMock)(nil)

// RecipeStepCompletionConditionDataManagerMock is a mocked types.RecipeStepCompletionConditionDataManager for testing.
type RecipeStepCompletionConditionDataManagerMock struct {
	mock.Mock
}

// RecipeStepCompletionConditionExists is a mock function.
func (m *RecipeStepCompletionConditionDataManagerMock) RecipeStepCompletionConditionExists(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) (bool, error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID, recipeStepIngredientID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetRecipeStepCompletionCondition is a mock function.
func (m *RecipeStepCompletionConditionDataManagerMock) GetRecipeStepCompletionCondition(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) (*types.RecipeStepCompletionCondition, error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID, recipeStepIngredientID)
	return returnValues.Get(0).(*types.RecipeStepCompletionCondition), returnValues.Error(1)
}

// GetRecipeStepCompletionConditions is a mock function.
func (m *RecipeStepCompletionConditionDataManagerMock) GetRecipeStepCompletionConditions(ctx context.Context, recipeID, recipeStepID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.RecipeStepCompletionCondition], error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[types.RecipeStepCompletionCondition]), returnValues.Error(1)
}

// CreateRecipeStepCompletionCondition is a mock function.
func (m *RecipeStepCompletionConditionDataManagerMock) CreateRecipeStepCompletionCondition(ctx context.Context, input *types.RecipeStepCompletionConditionDatabaseCreationInput) (*types.RecipeStepCompletionCondition, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*types.RecipeStepCompletionCondition), returnValues.Error(1)
}

// UpdateRecipeStepCompletionCondition is a mock function.
func (m *RecipeStepCompletionConditionDataManagerMock) UpdateRecipeStepCompletionCondition(ctx context.Context, updated *types.RecipeStepCompletionCondition) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveRecipeStepCompletionCondition is a mock function.
func (m *RecipeStepCompletionConditionDataManagerMock) ArchiveRecipeStepCompletionCondition(ctx context.Context, recipeStepID, recipeStepIngredientID string) error {
	return m.Called(ctx, recipeStepID, recipeStepIngredientID).Error(0)
}
