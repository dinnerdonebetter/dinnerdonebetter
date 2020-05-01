package mock

import (
	"context"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/mock"
)

var _ models.RecipeIterationDataManager = (*RecipeIterationDataManager)(nil)

// RecipeIterationDataManager is a mocked models.RecipeIterationDataManager for testing.
type RecipeIterationDataManager struct {
	mock.Mock
}

// RecipeIterationExists is a mock function.
func (m *RecipeIterationDataManager) RecipeIterationExists(ctx context.Context, recipeID, recipeIterationID uint64) (bool, error) {
	args := m.Called(ctx, recipeID, recipeIterationID)
	return args.Bool(0), args.Error(1)
}

// GetRecipeIteration is a mock function.
func (m *RecipeIterationDataManager) GetRecipeIteration(ctx context.Context, recipeID, recipeIterationID uint64) (*models.RecipeIteration, error) {
	args := m.Called(ctx, recipeID, recipeIterationID)
	return args.Get(0).(*models.RecipeIteration), args.Error(1)
}

// GetAllRecipeIterationsCount is a mock function.
func (m *RecipeIterationDataManager) GetAllRecipeIterationsCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetRecipeIterations is a mock function.
func (m *RecipeIterationDataManager) GetRecipeIterations(ctx context.Context, recipeID uint64, filter *models.QueryFilter) (*models.RecipeIterationList, error) {
	args := m.Called(ctx, recipeID, filter)
	return args.Get(0).(*models.RecipeIterationList), args.Error(1)
}

// CreateRecipeIteration is a mock function.
func (m *RecipeIterationDataManager) CreateRecipeIteration(ctx context.Context, input *models.RecipeIterationCreationInput) (*models.RecipeIteration, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*models.RecipeIteration), args.Error(1)
}

// UpdateRecipeIteration is a mock function.
func (m *RecipeIterationDataManager) UpdateRecipeIteration(ctx context.Context, updated *models.RecipeIteration) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveRecipeIteration is a mock function.
func (m *RecipeIterationDataManager) ArchiveRecipeIteration(ctx context.Context, recipeID, recipeIterationID uint64) error {
	return m.Called(ctx, recipeID, recipeIterationID).Error(0)
}
