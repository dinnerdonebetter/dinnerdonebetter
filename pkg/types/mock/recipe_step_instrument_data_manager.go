package mocktypes

import (
	"context"

	"github.com/stretchr/testify/mock"

	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

var _ types.RecipeStepInstrumentDataManager = (*RecipeStepInstrumentDataManager)(nil)

// RecipeStepInstrumentDataManager is a mocked types.RecipeStepInstrumentDataManager for testing.
type RecipeStepInstrumentDataManager struct {
	mock.Mock
}

// RecipeStepInstrumentExists is a mock function.
func (m *RecipeStepInstrumentDataManager) RecipeStepInstrumentExists(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string) (bool, error) {
	args := m.Called(ctx, recipeID, recipeStepID, recipeStepInstrumentID)
	return args.Bool(0), args.Error(1)
}

// GetRecipeStepInstrument is a mock function.
func (m *RecipeStepInstrumentDataManager) GetRecipeStepInstrument(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string) (*types.RecipeStepInstrument, error) {
	args := m.Called(ctx, recipeID, recipeStepID, recipeStepInstrumentID)
	return args.Get(0).(*types.RecipeStepInstrument), args.Error(1)
}

// GetTotalRecipeStepInstrumentCount is a mock function.
func (m *RecipeStepInstrumentDataManager) GetTotalRecipeStepInstrumentCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetRecipeStepInstruments is a mock function.
func (m *RecipeStepInstrumentDataManager) GetRecipeStepInstruments(ctx context.Context, recipeID, recipeStepID string, filter *types.QueryFilter) (*types.RecipeStepInstrumentList, error) {
	args := m.Called(ctx, recipeID, recipeStepID, filter)
	return args.Get(0).(*types.RecipeStepInstrumentList), args.Error(1)
}

// GetRecipeStepInstrumentsWithIDs is a mock function.
func (m *RecipeStepInstrumentDataManager) GetRecipeStepInstrumentsWithIDs(ctx context.Context, recipeStepID string, limit uint8, ids []string) ([]*types.RecipeStepInstrument, error) {
	args := m.Called(ctx, recipeStepID, limit, ids)
	return args.Get(0).([]*types.RecipeStepInstrument), args.Error(1)
}

// CreateRecipeStepInstrument is a mock function.
func (m *RecipeStepInstrumentDataManager) CreateRecipeStepInstrument(ctx context.Context, input *types.RecipeStepInstrumentDatabaseCreationInput) (*types.RecipeStepInstrument, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*types.RecipeStepInstrument), args.Error(1)
}

// UpdateRecipeStepInstrument is a mock function.
func (m *RecipeStepInstrumentDataManager) UpdateRecipeStepInstrument(ctx context.Context, updated *types.RecipeStepInstrument) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveRecipeStepInstrument is a mock function.
func (m *RecipeStepInstrumentDataManager) ArchiveRecipeStepInstrument(ctx context.Context, recipeStepID, recipeStepInstrumentID string) error {
	return m.Called(ctx, recipeStepID, recipeStepInstrumentID).Error(0)
}
