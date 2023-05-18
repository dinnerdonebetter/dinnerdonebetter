package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.RecipeStepVesselDataManager = (*RecipeStepVesselDataManager)(nil)

// RecipeStepVesselDataManager is a mocked types.RecipeStepVesselDataManager for testing.
type RecipeStepVesselDataManager struct {
	mock.Mock
}

// RecipeStepVesselExists is a mock function.
func (m *RecipeStepVesselDataManager) RecipeStepVesselExists(ctx context.Context, recipeID, recipeStepID, recipeStepVesselID string) (bool, error) {
	args := m.Called(ctx, recipeID, recipeStepID, recipeStepVesselID)
	return args.Bool(0), args.Error(1)
}

// GetRecipeStepVessel is a mock function.
func (m *RecipeStepVesselDataManager) GetRecipeStepVessel(ctx context.Context, recipeID, recipeStepID, recipeStepVesselID string) (*types.RecipeStepVessel, error) {
	args := m.Called(ctx, recipeID, recipeStepID, recipeStepVesselID)
	return args.Get(0).(*types.RecipeStepVessel), args.Error(1)
}

// GetRecipeStepVessels is a mock function.
func (m *RecipeStepVesselDataManager) GetRecipeStepVessels(ctx context.Context, recipeID, recipeStepID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.RecipeStepVessel], error) {
	args := m.Called(ctx, recipeID, recipeStepID, filter)
	return args.Get(0).(*types.QueryFilteredResult[types.RecipeStepVessel]), args.Error(1)
}

// CreateRecipeStepVessel is a mock function.
func (m *RecipeStepVesselDataManager) CreateRecipeStepVessel(ctx context.Context, input *types.RecipeStepVesselDatabaseCreationInput) (*types.RecipeStepVessel, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*types.RecipeStepVessel), args.Error(1)
}

// UpdateRecipeStepVessel is a mock function.
func (m *RecipeStepVesselDataManager) UpdateRecipeStepVessel(ctx context.Context, updated *types.RecipeStepVessel) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveRecipeStepVessel is a mock function.
func (m *RecipeStepVesselDataManager) ArchiveRecipeStepVessel(ctx context.Context, recipeStepID, recipeStepVesselID string) error {
	return m.Called(ctx, recipeStepID, recipeStepVesselID).Error(0)
}
