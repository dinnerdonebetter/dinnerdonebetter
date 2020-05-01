package mock

import (
	"context"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/mock"
)

var _ models.RecipeIterationStepDataManager = (*RecipeIterationStepDataManager)(nil)

// RecipeIterationStepDataManager is a mocked models.RecipeIterationStepDataManager for testing.
type RecipeIterationStepDataManager struct {
	mock.Mock
}

// RecipeIterationStepExists is a mock function.
func (m *RecipeIterationStepDataManager) RecipeIterationStepExists(ctx context.Context, recipeID, recipeIterationStepID uint64) (bool, error) {
	args := m.Called(ctx, recipeID, recipeIterationStepID)
	return args.Bool(0), args.Error(1)
}

// GetRecipeIterationStep is a mock function.
func (m *RecipeIterationStepDataManager) GetRecipeIterationStep(ctx context.Context, recipeID, recipeIterationStepID uint64) (*models.RecipeIterationStep, error) {
	args := m.Called(ctx, recipeID, recipeIterationStepID)
	return args.Get(0).(*models.RecipeIterationStep), args.Error(1)
}

// GetAllRecipeIterationStepsCount is a mock function.
func (m *RecipeIterationStepDataManager) GetAllRecipeIterationStepsCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetRecipeIterationSteps is a mock function.
func (m *RecipeIterationStepDataManager) GetRecipeIterationSteps(ctx context.Context, recipeID uint64, filter *models.QueryFilter) (*models.RecipeIterationStepList, error) {
	args := m.Called(ctx, recipeID, filter)
	return args.Get(0).(*models.RecipeIterationStepList), args.Error(1)
}

// CreateRecipeIterationStep is a mock function.
func (m *RecipeIterationStepDataManager) CreateRecipeIterationStep(ctx context.Context, input *models.RecipeIterationStepCreationInput) (*models.RecipeIterationStep, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*models.RecipeIterationStep), args.Error(1)
}

// UpdateRecipeIterationStep is a mock function.
func (m *RecipeIterationStepDataManager) UpdateRecipeIterationStep(ctx context.Context, updated *models.RecipeIterationStep) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveRecipeIterationStep is a mock function.
func (m *RecipeIterationStepDataManager) ArchiveRecipeIterationStep(ctx context.Context, recipeID, recipeIterationStepID uint64) error {
	return m.Called(ctx, recipeID, recipeIterationStepID).Error(0)
}
