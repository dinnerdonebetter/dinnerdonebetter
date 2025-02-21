package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types"

	"github.com/stretchr/testify/mock"
)

var _ types.RecipeStepDataManager = (*RecipeStepDataManagerMock)(nil)

// RecipeStepDataManagerMock is a mocked types.RecipeStepDataManager for testing.
type RecipeStepDataManagerMock struct {
	mock.Mock
}

// RecipeStepExists is a mock function.
func (m *RecipeStepDataManagerMock) RecipeStepExists(ctx context.Context, recipeID, recipeStepID string) (bool, error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetRecipeStep is a mock function.
func (m *RecipeStepDataManagerMock) GetRecipeStep(ctx context.Context, recipeID, recipeStepID string) (*types.RecipeStep, error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID)
	return returnValues.Get(0).(*types.RecipeStep), returnValues.Error(1)
}

// GetRecipeSteps is a mock function.
func (m *RecipeStepDataManagerMock) GetRecipeSteps(ctx context.Context, recipeID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.RecipeStep], error) {
	returnValues := m.Called(ctx, recipeID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[types.RecipeStep]), returnValues.Error(1)
}

// CreateRecipeStep is a mock function.
func (m *RecipeStepDataManagerMock) CreateRecipeStep(ctx context.Context, input *types.RecipeStepDatabaseCreationInput) (*types.RecipeStep, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*types.RecipeStep), returnValues.Error(1)
}

// UpdateRecipeStep is a mock function.
func (m *RecipeStepDataManagerMock) UpdateRecipeStep(ctx context.Context, updated *types.RecipeStep) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveRecipeStep is a mock function.
func (m *RecipeStepDataManagerMock) ArchiveRecipeStep(ctx context.Context, recipeID, recipeStepID string) error {
	return m.Called(ctx, recipeID, recipeStepID).Error(0)
}
