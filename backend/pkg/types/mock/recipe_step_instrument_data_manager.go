package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.RecipeStepInstrumentDataManager = (*RecipeStepInstrumentDataManagerMock)(nil)

// RecipeStepInstrumentDataManagerMock is a mocked types.RecipeStepInstrumentDataManager for testing.
type RecipeStepInstrumentDataManagerMock struct {
	mock.Mock
}

// RecipeStepInstrumentExists is a mock function.
func (m *RecipeStepInstrumentDataManagerMock) RecipeStepInstrumentExists(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string) (bool, error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID, recipeStepInstrumentID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetRecipeStepInstrument is a mock function.
func (m *RecipeStepInstrumentDataManagerMock) GetRecipeStepInstrument(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string) (*types.RecipeStepInstrument, error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID, recipeStepInstrumentID)
	return returnValues.Get(0).(*types.RecipeStepInstrument), returnValues.Error(1)
}

// GetRecipeStepInstruments is a mock function.
func (m *RecipeStepInstrumentDataManagerMock) GetRecipeStepInstruments(ctx context.Context, recipeID, recipeStepID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.RecipeStepInstrument], error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID, filter)
	return returnValues.Get(0).(*types.QueryFilteredResult[types.RecipeStepInstrument]), returnValues.Error(1)
}

// CreateRecipeStepInstrument is a mock function.
func (m *RecipeStepInstrumentDataManagerMock) CreateRecipeStepInstrument(ctx context.Context, input *types.RecipeStepInstrumentDatabaseCreationInput) (*types.RecipeStepInstrument, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*types.RecipeStepInstrument), returnValues.Error(1)
}

// UpdateRecipeStepInstrument is a mock function.
func (m *RecipeStepInstrumentDataManagerMock) UpdateRecipeStepInstrument(ctx context.Context, updated *types.RecipeStepInstrument) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveRecipeStepInstrument is a mock function.
func (m *RecipeStepInstrumentDataManagerMock) ArchiveRecipeStepInstrument(ctx context.Context, recipeStepID, recipeStepInstrumentID string) error {
	return m.Called(ctx, recipeStepID, recipeStepInstrumentID).Error(0)
}
