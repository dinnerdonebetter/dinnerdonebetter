package mock

import (
	"context"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/mock"
)

var _ models.RecipeStepInstrumentDataManager = (*RecipeStepInstrumentDataManager)(nil)

// RecipeStepInstrumentDataManager is a mocked models.RecipeStepInstrumentDataManager for testing
type RecipeStepInstrumentDataManager struct {
	mock.Mock
}

// GetRecipeStepInstrument is a mock function
func (m *RecipeStepInstrumentDataManager) GetRecipeStepInstrument(ctx context.Context, recipeStepInstrumentID, userID uint64) (*models.RecipeStepInstrument, error) {
	args := m.Called(ctx, recipeStepInstrumentID, userID)
	return args.Get(0).(*models.RecipeStepInstrument), args.Error(1)
}

// GetRecipeStepInstrumentCount is a mock function
func (m *RecipeStepInstrumentDataManager) GetRecipeStepInstrumentCount(ctx context.Context, filter *models.QueryFilter, userID uint64) (uint64, error) {
	args := m.Called(ctx, filter, userID)
	return args.Get(0).(uint64), args.Error(1)
}

// GetAllRecipeStepInstrumentsCount is a mock function
func (m *RecipeStepInstrumentDataManager) GetAllRecipeStepInstrumentsCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetRecipeStepInstruments is a mock function
func (m *RecipeStepInstrumentDataManager) GetRecipeStepInstruments(ctx context.Context, filter *models.QueryFilter, userID uint64) (*models.RecipeStepInstrumentList, error) {
	args := m.Called(ctx, filter, userID)
	return args.Get(0).(*models.RecipeStepInstrumentList), args.Error(1)
}

// GetAllRecipeStepInstrumentsForUser is a mock function
func (m *RecipeStepInstrumentDataManager) GetAllRecipeStepInstrumentsForUser(ctx context.Context, userID uint64) ([]models.RecipeStepInstrument, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.RecipeStepInstrument), args.Error(1)
}

// CreateRecipeStepInstrument is a mock function
func (m *RecipeStepInstrumentDataManager) CreateRecipeStepInstrument(ctx context.Context, input *models.RecipeStepInstrumentCreationInput) (*models.RecipeStepInstrument, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*models.RecipeStepInstrument), args.Error(1)
}

// UpdateRecipeStepInstrument is a mock function
func (m *RecipeStepInstrumentDataManager) UpdateRecipeStepInstrument(ctx context.Context, updated *models.RecipeStepInstrument) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveRecipeStepInstrument is a mock function
func (m *RecipeStepInstrumentDataManager) ArchiveRecipeStepInstrument(ctx context.Context, id, userID uint64) error {
	return m.Called(ctx, id, userID).Error(0)
}
