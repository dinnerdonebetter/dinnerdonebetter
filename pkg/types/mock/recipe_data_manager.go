package mock

import (
	"context"

	"gitlab.com/prixfixe/prixfixe/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.RecipeDataManager = (*RecipeDataManager)(nil)

// RecipeDataManager is a mocked types.RecipeDataManager for testing.
type RecipeDataManager struct {
	mock.Mock
}

// RecipeExists is a mock function.
func (m *RecipeDataManager) RecipeExists(ctx context.Context, recipeID uint64) (bool, error) {
	args := m.Called(ctx, recipeID)
	return args.Bool(0), args.Error(1)
}

// GetRecipe is a mock function.
func (m *RecipeDataManager) GetRecipe(ctx context.Context, recipeID uint64) (*types.Recipe, error) {
	args := m.Called(ctx, recipeID)
	return args.Get(0).(*types.Recipe), args.Error(1)
}

// GetAllRecipesCount is a mock function.
func (m *RecipeDataManager) GetAllRecipesCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetAllRecipes is a mock function.
func (m *RecipeDataManager) GetAllRecipes(ctx context.Context, results chan []*types.Recipe, bucketSize uint16) error {
	args := m.Called(ctx, results, bucketSize)
	return args.Error(0)
}

// GetRecipes is a mock function.
func (m *RecipeDataManager) GetRecipes(ctx context.Context, filter *types.QueryFilter) (*types.RecipeList, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*types.RecipeList), args.Error(1)
}

// GetRecipesWithIDs is a mock function.
func (m *RecipeDataManager) GetRecipesWithIDs(ctx context.Context, accountID uint64, limit uint8, ids []uint64) ([]*types.Recipe, error) {
	args := m.Called(ctx, accountID, limit, ids)
	return args.Get(0).([]*types.Recipe), args.Error(1)
}

// CreateRecipe is a mock function.
func (m *RecipeDataManager) CreateRecipe(ctx context.Context, input *types.RecipeCreationInput, createdByUser uint64) (*types.Recipe, error) {
	args := m.Called(ctx, input, createdByUser)
	return args.Get(0).(*types.Recipe), args.Error(1)
}

// UpdateRecipe is a mock function.
func (m *RecipeDataManager) UpdateRecipe(ctx context.Context, updated *types.Recipe, changedByUser uint64, changes []*types.FieldChangeSummary) error {
	return m.Called(ctx, updated, changedByUser, changes).Error(0)
}

// ArchiveRecipe is a mock function.
func (m *RecipeDataManager) ArchiveRecipe(ctx context.Context, recipeID, accountID, archivedBy uint64) error {
	return m.Called(ctx, recipeID, accountID, archivedBy).Error(0)
}

// GetAuditLogEntriesForRecipe is a mock function.
func (m *RecipeDataManager) GetAuditLogEntriesForRecipe(ctx context.Context, recipeID uint64) ([]*types.AuditLogEntry, error) {
	args := m.Called(ctx, recipeID)
	return args.Get(0).([]*types.AuditLogEntry), args.Error(1)
}
