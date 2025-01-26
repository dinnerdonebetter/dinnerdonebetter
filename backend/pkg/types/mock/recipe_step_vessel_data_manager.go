package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.RecipeStepVesselDataManager = (*RecipeStepVesselDataManagerMock)(nil)

// RecipeStepVesselDataManagerMock is a mocked types.RecipeStepVesselDataManager for testing.
type RecipeStepVesselDataManagerMock struct {
	mock.Mock
}

// RecipeStepVesselExists is a mock function.
func (m *RecipeStepVesselDataManagerMock) RecipeStepVesselExists(ctx context.Context, recipeID, recipeStepID, recipeStepVesselID string) (bool, error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID, recipeStepVesselID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetRecipeStepVessel is a mock function.
func (m *RecipeStepVesselDataManagerMock) GetRecipeStepVessel(ctx context.Context, recipeID, recipeStepID, recipeStepVesselID string) (*types.RecipeStepVessel, error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID, recipeStepVesselID)
	return returnValues.Get(0).(*types.RecipeStepVessel), returnValues.Error(1)
}

// GetRecipeStepVessels is a mock function.
func (m *RecipeStepVesselDataManagerMock) GetRecipeStepVessels(ctx context.Context, recipeID, recipeStepID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.RecipeStepVessel], error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[types.RecipeStepVessel]), returnValues.Error(1)
}

// CreateRecipeStepVessel is a mock function.
func (m *RecipeStepVesselDataManagerMock) CreateRecipeStepVessel(ctx context.Context, input *types.RecipeStepVesselDatabaseCreationInput) (*types.RecipeStepVessel, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*types.RecipeStepVessel), returnValues.Error(1)
}

// UpdateRecipeStepVessel is a mock function.
func (m *RecipeStepVesselDataManagerMock) UpdateRecipeStepVessel(ctx context.Context, updated *types.RecipeStepVessel) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveRecipeStepVessel is a mock function.
func (m *RecipeStepVesselDataManagerMock) ArchiveRecipeStepVessel(ctx context.Context, recipeStepID, recipeStepVesselID string) error {
	return m.Called(ctx, recipeStepID, recipeStepVesselID).Error(0)
}
