package mock

import (
	"context"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/mock"
)

var _ models.RecipeStepDataManager = (*RecipeStepDataManager)(nil)

// RecipeStepDataManager is a mocked models.RecipeStepDataManager for testing
type RecipeStepDataManager struct {
	mock.Mock
}

// GetRecipeStep is a mock function
func (m *RecipeStepDataManager) GetRecipeStep(ctx context.Context, recipeStepID, userID uint64) (*models.RecipeStep, error) {
	args := m.Called(ctx, recipeStepID, userID)
	return args.Get(0).(*models.RecipeStep), args.Error(1)
}

// GetRecipeStepCount is a mock function
func (m *RecipeStepDataManager) GetRecipeStepCount(ctx context.Context, filter *models.QueryFilter, userID uint64) (uint64, error) {
	args := m.Called(ctx, filter, userID)
	return args.Get(0).(uint64), args.Error(1)
}

// GetAllRecipeStepsCount is a mock function
func (m *RecipeStepDataManager) GetAllRecipeStepsCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetRecipeSteps is a mock function
func (m *RecipeStepDataManager) GetRecipeSteps(ctx context.Context, filter *models.QueryFilter, userID uint64) (*models.RecipeStepList, error) {
	args := m.Called(ctx, filter, userID)
	return args.Get(0).(*models.RecipeStepList), args.Error(1)
}

// GetAllRecipeStepsForUser is a mock function
func (m *RecipeStepDataManager) GetAllRecipeStepsForUser(ctx context.Context, userID uint64) ([]models.RecipeStep, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.RecipeStep), args.Error(1)
}

// CreateRecipeStep is a mock function
func (m *RecipeStepDataManager) CreateRecipeStep(ctx context.Context, input *models.RecipeStepCreationInput) (*models.RecipeStep, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*models.RecipeStep), args.Error(1)
}

// UpdateRecipeStep is a mock function
func (m *RecipeStepDataManager) UpdateRecipeStep(ctx context.Context, updated *models.RecipeStep) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveRecipeStep is a mock function
func (m *RecipeStepDataManager) ArchiveRecipeStep(ctx context.Context, id, userID uint64) error {
	return m.Called(ctx, id, userID).Error(0)
}
