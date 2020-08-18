package mock

import (
	"context"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/mock"
)

var _ models.RecipeStepInstrumentDataManager = (*RecipeStepInstrumentDataManager)(nil)

// RecipeStepInstrumentDataManager is a mocked models.RecipeStepInstrumentDataManager for testing.
type RecipeStepInstrumentDataManager struct {
	mock.Mock
}

// RecipeStepInstrumentExists is a mock function.
func (m *RecipeStepInstrumentDataManager) RecipeStepInstrumentExists(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID uint64) (bool, error) {
	args := m.Called(ctx, recipeID, recipeStepID, recipeStepInstrumentID)
	return args.Bool(0), args.Error(1)
}

// GetRecipeStepInstrument is a mock function.
func (m *RecipeStepInstrumentDataManager) GetRecipeStepInstrument(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID uint64) (*models.RecipeStepInstrument, error) {
	args := m.Called(ctx, recipeID, recipeStepID, recipeStepInstrumentID)
	return args.Get(0).(*models.RecipeStepInstrument), args.Error(1)
}

// GetAllRecipeStepInstrumentsCount is a mock function.
func (m *RecipeStepInstrumentDataManager) GetAllRecipeStepInstrumentsCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetAllRecipeStepInstruments is a mock function.
func (m *RecipeStepInstrumentDataManager) GetAllRecipeStepInstruments(ctx context.Context, results chan []models.RecipeStepInstrument) error {
	args := m.Called(ctx, results)
	return args.Error(0)
}

// GetRecipeStepInstruments is a mock function.
func (m *RecipeStepInstrumentDataManager) GetRecipeStepInstruments(ctx context.Context, recipeID, recipeStepID uint64, filter *models.QueryFilter) (*models.RecipeStepInstrumentList, error) {
	args := m.Called(ctx, recipeID, recipeStepID, filter)
	return args.Get(0).(*models.RecipeStepInstrumentList), args.Error(1)
}

// GetRecipeStepInstrumentsWithIDs is a mock function.
func (m *RecipeStepInstrumentDataManager) GetRecipeStepInstrumentsWithIDs(ctx context.Context, recipeID, recipeStepID uint64, limit uint8, ids []uint64) ([]models.RecipeStepInstrument, error) {
	args := m.Called(ctx, recipeID, recipeStepID, limit, ids)
	return args.Get(0).([]models.RecipeStepInstrument), args.Error(1)
}

// CreateRecipeStepInstrument is a mock function.
func (m *RecipeStepInstrumentDataManager) CreateRecipeStepInstrument(ctx context.Context, input *models.RecipeStepInstrumentCreationInput) (*models.RecipeStepInstrument, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*models.RecipeStepInstrument), args.Error(1)
}

// UpdateRecipeStepInstrument is a mock function.
func (m *RecipeStepInstrumentDataManager) UpdateRecipeStepInstrument(ctx context.Context, updated *models.RecipeStepInstrument) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveRecipeStepInstrument is a mock function.
func (m *RecipeStepInstrumentDataManager) ArchiveRecipeStepInstrument(ctx context.Context, recipeStepID, recipeStepInstrumentID uint64) error {
	return m.Called(ctx, recipeStepID, recipeStepInstrumentID).Error(0)
}
