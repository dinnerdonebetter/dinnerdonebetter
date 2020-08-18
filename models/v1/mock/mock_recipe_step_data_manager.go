package mock

import (
	"context"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/mock"
)

var _ models.RecipeStepDataManager = (*RecipeStepDataManager)(nil)

// RecipeStepDataManager is a mocked models.RecipeStepDataManager for testing.
type RecipeStepDataManager struct {
	mock.Mock
}

// RecipeStepExists is a mock function.
func (m *RecipeStepDataManager) RecipeStepExists(ctx context.Context, recipeID, recipeStepID uint64) (bool, error) {
	args := m.Called(ctx, recipeID, recipeStepID)
	return args.Bool(0), args.Error(1)
}

// GetRecipeStep is a mock function.
func (m *RecipeStepDataManager) GetRecipeStep(ctx context.Context, recipeID, recipeStepID uint64) (*models.RecipeStep, error) {
	args := m.Called(ctx, recipeID, recipeStepID)
	return args.Get(0).(*models.RecipeStep), args.Error(1)
}

// GetAllRecipeStepsCount is a mock function.
func (m *RecipeStepDataManager) GetAllRecipeStepsCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetAllRecipeSteps is a mock function.
func (m *RecipeStepDataManager) GetAllRecipeSteps(ctx context.Context, results chan []models.RecipeStep) error {
	args := m.Called(ctx, results)
	return args.Error(0)
}

// GetRecipeSteps is a mock function.
func (m *RecipeStepDataManager) GetRecipeSteps(ctx context.Context, recipeID uint64, filter *models.QueryFilter) (*models.RecipeStepList, error) {
	args := m.Called(ctx, recipeID, filter)
	return args.Get(0).(*models.RecipeStepList), args.Error(1)
}

// GetRecipeStepsWithIDs is a mock function.
func (m *RecipeStepDataManager) GetRecipeStepsWithIDs(ctx context.Context, recipeID uint64, limit uint8, ids []uint64) ([]models.RecipeStep, error) {
	args := m.Called(ctx, recipeID, limit, ids)
	return args.Get(0).([]models.RecipeStep), args.Error(1)
}

// CreateRecipeStep is a mock function.
func (m *RecipeStepDataManager) CreateRecipeStep(ctx context.Context, input *models.RecipeStepCreationInput) (*models.RecipeStep, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*models.RecipeStep), args.Error(1)
}

// UpdateRecipeStep is a mock function.
func (m *RecipeStepDataManager) UpdateRecipeStep(ctx context.Context, updated *models.RecipeStep) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveRecipeStep is a mock function.
func (m *RecipeStepDataManager) ArchiveRecipeStep(ctx context.Context, recipeID, recipeStepID uint64) error {
	return m.Called(ctx, recipeID, recipeStepID).Error(0)
}
