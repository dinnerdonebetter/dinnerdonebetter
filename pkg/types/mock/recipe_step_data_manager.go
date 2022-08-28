package mocktypes

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/pkg/types"
)

var _ types.RecipeStepDataManager = (*RecipeStepDataManager)(nil)

// RecipeStepDataManager is a mocked types.RecipeStepDataManager for testing.
type RecipeStepDataManager struct {
	mock.Mock
}

// RecipeStepExists is a mock function.
func (m *RecipeStepDataManager) RecipeStepExists(ctx context.Context, recipeID, recipeStepID string) (bool, error) {
	args := m.Called(ctx, recipeID, recipeStepID)
	return args.Bool(0), args.Error(1)
}

// GetRecipeStep is a mock function.
func (m *RecipeStepDataManager) GetRecipeStep(ctx context.Context, recipeID, recipeStepID string) (*types.RecipeStep, error) {
	args := m.Called(ctx, recipeID, recipeStepID)
	return args.Get(0).(*types.RecipeStep), args.Error(1)
}

// GetRecipeSteps is a mock function.
func (m *RecipeStepDataManager) GetRecipeSteps(ctx context.Context, recipeID string, filter *types.QueryFilter) (*types.RecipeStepList, error) {
	args := m.Called(ctx, recipeID, filter)
	return args.Get(0).(*types.RecipeStepList), args.Error(1)
}

// CreateRecipeStep is a mock function.
func (m *RecipeStepDataManager) CreateRecipeStep(ctx context.Context, input *types.RecipeStepDatabaseCreationInput) (*types.RecipeStep, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*types.RecipeStep), args.Error(1)
}

// UpdateRecipeStep is a mock function.
func (m *RecipeStepDataManager) UpdateRecipeStep(ctx context.Context, updated *types.RecipeStep) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveRecipeStep is a mock function.
func (m *RecipeStepDataManager) ArchiveRecipeStep(ctx context.Context, recipeID, recipeStepID string) error {
	return m.Called(ctx, recipeID, recipeStepID).Error(0)
}
