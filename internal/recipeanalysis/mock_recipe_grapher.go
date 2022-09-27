package recipeanalysis

import (
	"context"
	"image"

	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/pkg/types"
)

var _ RecipeAnalyzer = (*MockRecipeAnalyzer)(nil)

// MockRecipeAnalyzer is a mock RecipeAnalyzer.
type MockRecipeAnalyzer struct {
	mock.Mock
}

// GenerateAdvancedStepCreationForRecipe implements our interface.
func (m *MockRecipeAnalyzer) GenerateAdvancedStepCreationForRecipe(ctx context.Context, mealPlanEvent *types.MealPlanEvent, mealPlanOptionID string, recipe *types.Recipe) ([]*types.MealPlanTaskDatabaseCreationInput, error) {
	returnArgs := m.Called(ctx, mealPlanEvent, mealPlanOptionID, recipe)

	return returnArgs.Get(0).([]*types.MealPlanTaskDatabaseCreationInput), returnArgs.Error(1)
}

// FindStepsEligibleForAdvancedCreation implements our interface.
func (m *MockRecipeAnalyzer) FindStepsEligibleForAdvancedCreation(ctx context.Context, recipe *types.Recipe) ([]*types.RecipeStep, error) {
	returnArgs := m.Called(ctx, recipe)

	return returnArgs.Get(0).([]*types.RecipeStep), returnArgs.Error(1)
}

// GenerateDAGDiagramForRecipe implements our interface.
func (m *MockRecipeAnalyzer) GenerateDAGDiagramForRecipe(ctx context.Context, recipe *types.Recipe) (image.Image, error) {
	returnArgs := m.Called(ctx, recipe)

	return returnArgs.Get(0).(image.Image), returnArgs.Error(1)
}
