package recipeanalysis

import (
	"context"
	"image"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/pkg/types"
)

var _ RecipeAnalyzer = (*MockRecipeAnalyzer)(nil)

// MockRecipeAnalyzer is a mock RecipeAnalyzer.
type MockRecipeAnalyzer struct {
	mock.Mock
}

// GenerateMealPlanTasksForRecipe implements our interface.
func (m *MockRecipeAnalyzer) GenerateMealPlanTasksForRecipe(ctx context.Context, startsAt time.Time, mealPlanOptionID string, recipe *types.Recipe) ([]*types.MealPlanTaskDatabaseCreationInput, error) {
	returnArgs := m.Called(ctx, startsAt, mealPlanOptionID, recipe)

	return returnArgs.Get(0).([]*types.MealPlanTaskDatabaseCreationInput), returnArgs.Error(1)
}

// FindStepsEligibleForMealPlanTasks implements our interface.
func (m *MockRecipeAnalyzer) FindStepsEligibleForMealPlanTasks(ctx context.Context, recipe *types.Recipe) ([]*types.RecipeStep, error) {
	returnArgs := m.Called(ctx, recipe)

	return returnArgs.Get(0).([]*types.RecipeStep), returnArgs.Error(1)
}

// GenerateDAGDiagramForRecipe implements our interface.
func (m *MockRecipeAnalyzer) GenerateDAGDiagramForRecipe(ctx context.Context, recipe *types.Recipe) (image.Image, error) {
	returnArgs := m.Called(ctx, recipe)

	return returnArgs.Get(0).(image.Image), returnArgs.Error(1)
}
