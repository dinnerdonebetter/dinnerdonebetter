package graphing

import (
	"context"
	"image"

	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/pkg/types"
)

var _ RecipeDAGDiagramGenerator = (*MockRecipeGrapher)(nil)

// MockRecipeGrapher is a mock RecipeDAGDiagramGenerator.
type MockRecipeGrapher struct {
	mock.Mock
}

// GenerateDAGDiagramForRecipe implements our interface.
func (m *MockRecipeGrapher) GenerateDAGDiagramForRecipe(ctx context.Context, recipe *types.Recipe) (image.Image, error) {
	returnArgs := m.Called(ctx, recipe)

	return returnArgs.Get(0).(image.Image), returnArgs.Error(1)
}
