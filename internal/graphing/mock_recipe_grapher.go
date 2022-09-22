package graphing

import (
	"context"
	"image"

	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/pkg/types"
)

var _ RecipeGrapher = (*MockRecipeGrapher)(nil)

// MockRecipeGrapher is a mock RecipeGrapher.
type MockRecipeGrapher struct {
	mock.Mock
}

// FindStepsEligibleForAdvancedCreation implements our interface.
func (m *MockRecipeGrapher) FindStepsEligibleForAdvancedCreation(ctx context.Context, recipe *types.Recipe) ([]*types.RecipeStep, error) {
	returnArgs := m.Called(ctx, recipe)

	return returnArgs.Get(0).([]*types.RecipeStep), returnArgs.Error(1)
}

// GenerateDAGDiagramForRecipe implements our interface.
func (m *MockRecipeGrapher) GenerateDAGDiagramForRecipe(ctx context.Context, recipe *types.Recipe) (image.Image, error) {
	returnArgs := m.Called(ctx, recipe)

	return returnArgs.Get(0).(image.Image), returnArgs.Error(1)
}
