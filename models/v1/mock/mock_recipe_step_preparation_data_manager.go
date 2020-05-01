package mock

import (
	"context"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/mock"
)

var _ models.RecipeStepPreparationDataManager = (*RecipeStepPreparationDataManager)(nil)

// RecipeStepPreparationDataManager is a mocked models.RecipeStepPreparationDataManager for testing.
type RecipeStepPreparationDataManager struct {
	mock.Mock
}

// RecipeStepPreparationExists is a mock function.
func (m *RecipeStepPreparationDataManager) RecipeStepPreparationExists(ctx context.Context, recipeID, recipeStepID, recipeStepPreparationID uint64) (bool, error) {
	args := m.Called(ctx, recipeID, recipeStepID, recipeStepPreparationID)
	return args.Bool(0), args.Error(1)
}

// GetRecipeStepPreparation is a mock function.
func (m *RecipeStepPreparationDataManager) GetRecipeStepPreparation(ctx context.Context, recipeID, recipeStepID, recipeStepPreparationID uint64) (*models.RecipeStepPreparation, error) {
	args := m.Called(ctx, recipeID, recipeStepID, recipeStepPreparationID)
	return args.Get(0).(*models.RecipeStepPreparation), args.Error(1)
}

// GetAllRecipeStepPreparationsCount is a mock function.
func (m *RecipeStepPreparationDataManager) GetAllRecipeStepPreparationsCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetRecipeStepPreparations is a mock function.
func (m *RecipeStepPreparationDataManager) GetRecipeStepPreparations(ctx context.Context, recipeID, recipeStepID uint64, filter *models.QueryFilter) (*models.RecipeStepPreparationList, error) {
	args := m.Called(ctx, recipeID, recipeStepID, filter)
	return args.Get(0).(*models.RecipeStepPreparationList), args.Error(1)
}

// CreateRecipeStepPreparation is a mock function.
func (m *RecipeStepPreparationDataManager) CreateRecipeStepPreparation(ctx context.Context, input *models.RecipeStepPreparationCreationInput) (*models.RecipeStepPreparation, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*models.RecipeStepPreparation), args.Error(1)
}

// UpdateRecipeStepPreparation is a mock function.
func (m *RecipeStepPreparationDataManager) UpdateRecipeStepPreparation(ctx context.Context, updated *models.RecipeStepPreparation) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveRecipeStepPreparation is a mock function.
func (m *RecipeStepPreparationDataManager) ArchiveRecipeStepPreparation(ctx context.Context, recipeStepID, recipeStepPreparationID uint64) error {
	return m.Called(ctx, recipeStepID, recipeStepPreparationID).Error(0)
}
