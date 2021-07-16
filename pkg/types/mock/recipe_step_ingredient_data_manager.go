package mock

import (
	"context"

	"gitlab.com/prixfixe/prixfixe/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.RecipeStepIngredientDataManager = (*RecipeStepIngredientDataManager)(nil)

// RecipeStepIngredientDataManager is a mocked types.RecipeStepIngredientDataManager for testing.
type RecipeStepIngredientDataManager struct {
	mock.Mock
}

// RecipeStepIngredientExists is a mock function.
func (m *RecipeStepIngredientDataManager) RecipeStepIngredientExists(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID uint64) (bool, error) {
	args := m.Called(ctx, recipeID, recipeStepID, recipeStepIngredientID)
	return args.Bool(0), args.Error(1)
}

// GetRecipeStepIngredient is a mock function.
func (m *RecipeStepIngredientDataManager) GetRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID uint64) (*types.RecipeStepIngredient, error) {
	args := m.Called(ctx, recipeID, recipeStepID, recipeStepIngredientID)
	return args.Get(0).(*types.RecipeStepIngredient), args.Error(1)
}

// GetAllRecipeStepIngredientsCount is a mock function.
func (m *RecipeStepIngredientDataManager) GetAllRecipeStepIngredientsCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetAllRecipeStepIngredients is a mock function.
func (m *RecipeStepIngredientDataManager) GetAllRecipeStepIngredients(ctx context.Context, results chan []*types.RecipeStepIngredient, bucketSize uint16) error {
	args := m.Called(ctx, results, bucketSize)
	return args.Error(0)
}

// GetRecipeStepIngredients is a mock function.
func (m *RecipeStepIngredientDataManager) GetRecipeStepIngredients(ctx context.Context, recipeID, recipeStepID uint64, filter *types.QueryFilter) (*types.RecipeStepIngredientList, error) {
	args := m.Called(ctx, recipeID, recipeStepID, filter)
	return args.Get(0).(*types.RecipeStepIngredientList), args.Error(1)
}

// GetRecipeStepIngredientsWithIDs is a mock function.
func (m *RecipeStepIngredientDataManager) GetRecipeStepIngredientsWithIDs(ctx context.Context, recipeStepID uint64, limit uint8, ids []uint64) ([]*types.RecipeStepIngredient, error) {
	args := m.Called(ctx, recipeStepID, limit, ids)
	return args.Get(0).([]*types.RecipeStepIngredient), args.Error(1)
}

// CreateRecipeStepIngredient is a mock function.
func (m *RecipeStepIngredientDataManager) CreateRecipeStepIngredient(ctx context.Context, input *types.RecipeStepIngredientCreationInput, createdByUser uint64) (*types.RecipeStepIngredient, error) {
	args := m.Called(ctx, input, createdByUser)
	return args.Get(0).(*types.RecipeStepIngredient), args.Error(1)
}

// UpdateRecipeStepIngredient is a mock function.
func (m *RecipeStepIngredientDataManager) UpdateRecipeStepIngredient(ctx context.Context, updated *types.RecipeStepIngredient, changedByUser uint64, changes []*types.FieldChangeSummary) error {
	return m.Called(ctx, updated, changedByUser, changes).Error(0)
}

// ArchiveRecipeStepIngredient is a mock function.
func (m *RecipeStepIngredientDataManager) ArchiveRecipeStepIngredient(ctx context.Context, recipeStepID, recipeStepIngredientID, archivedBy uint64) error {
	return m.Called(ctx, recipeStepID, recipeStepIngredientID, archivedBy).Error(0)
}

// GetAuditLogEntriesForRecipeStepIngredient is a mock function.
func (m *RecipeStepIngredientDataManager) GetAuditLogEntriesForRecipeStepIngredient(ctx context.Context, recipeStepIngredientID uint64) ([]*types.AuditLogEntry, error) {
	args := m.Called(ctx, recipeStepIngredientID)
	return args.Get(0).([]*types.AuditLogEntry), args.Error(1)
}
