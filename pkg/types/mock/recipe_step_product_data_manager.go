package mock

import (
	"context"

	"gitlab.com/prixfixe/prixfixe/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.RecipeStepProductDataManager = (*RecipeStepProductDataManager)(nil)

// RecipeStepProductDataManager is a mocked types.RecipeStepProductDataManager for testing.
type RecipeStepProductDataManager struct {
	mock.Mock
}

// RecipeStepProductExists is a mock function.
func (m *RecipeStepProductDataManager) RecipeStepProductExists(ctx context.Context, recipeID, recipeStepID, recipeStepProductID uint64) (bool, error) {
	args := m.Called(ctx, recipeID, recipeStepID, recipeStepProductID)
	return args.Bool(0), args.Error(1)
}

// GetRecipeStepProduct is a mock function.
func (m *RecipeStepProductDataManager) GetRecipeStepProduct(ctx context.Context, recipeID, recipeStepID, recipeStepProductID uint64) (*types.RecipeStepProduct, error) {
	args := m.Called(ctx, recipeID, recipeStepID, recipeStepProductID)
	return args.Get(0).(*types.RecipeStepProduct), args.Error(1)
}

// GetAllRecipeStepProductsCount is a mock function.
func (m *RecipeStepProductDataManager) GetAllRecipeStepProductsCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetAllRecipeStepProducts is a mock function.
func (m *RecipeStepProductDataManager) GetAllRecipeStepProducts(ctx context.Context, results chan []*types.RecipeStepProduct, bucketSize uint16) error {
	args := m.Called(ctx, results, bucketSize)
	return args.Error(0)
}

// GetRecipeStepProducts is a mock function.
func (m *RecipeStepProductDataManager) GetRecipeStepProducts(ctx context.Context, recipeID, recipeStepID uint64, filter *types.QueryFilter) (*types.RecipeStepProductList, error) {
	args := m.Called(ctx, recipeID, recipeStepID, filter)
	return args.Get(0).(*types.RecipeStepProductList), args.Error(1)
}

// GetRecipeStepProductsWithIDs is a mock function.
func (m *RecipeStepProductDataManager) GetRecipeStepProductsWithIDs(ctx context.Context, recipeStepID uint64, limit uint8, ids []uint64) ([]*types.RecipeStepProduct, error) {
	args := m.Called(ctx, recipeStepID, limit, ids)
	return args.Get(0).([]*types.RecipeStepProduct), args.Error(1)
}

// CreateRecipeStepProduct is a mock function.
func (m *RecipeStepProductDataManager) CreateRecipeStepProduct(ctx context.Context, input *types.RecipeStepProductCreationInput, createdByUser uint64) (*types.RecipeStepProduct, error) {
	args := m.Called(ctx, input, createdByUser)
	return args.Get(0).(*types.RecipeStepProduct), args.Error(1)
}

// UpdateRecipeStepProduct is a mock function.
func (m *RecipeStepProductDataManager) UpdateRecipeStepProduct(ctx context.Context, updated *types.RecipeStepProduct, changedByUser uint64, changes []*types.FieldChangeSummary) error {
	return m.Called(ctx, updated, changedByUser, changes).Error(0)
}

// ArchiveRecipeStepProduct is a mock function.
func (m *RecipeStepProductDataManager) ArchiveRecipeStepProduct(ctx context.Context, recipeStepID, recipeStepProductID, archivedBy uint64) error {
	return m.Called(ctx, recipeStepID, recipeStepProductID, archivedBy).Error(0)
}

// GetAuditLogEntriesForRecipeStepProduct is a mock function.
func (m *RecipeStepProductDataManager) GetAuditLogEntriesForRecipeStepProduct(ctx context.Context, recipeStepProductID uint64) ([]*types.AuditLogEntry, error) {
	args := m.Called(ctx, recipeStepProductID)
	return args.Get(0).([]*types.AuditLogEntry), args.Error(1)
}
