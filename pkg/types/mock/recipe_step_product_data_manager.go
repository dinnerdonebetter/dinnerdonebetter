package mocktypes

import (
	"context"

	"github.com/stretchr/testify/mock"

	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

var _ types.RecipeStepProductDataManager = (*RecipeStepProductDataManager)(nil)

// RecipeStepProductDataManager is a mocked types.RecipeStepProductDataManager for testing.
type RecipeStepProductDataManager struct {
	mock.Mock
}

// RecipeStepProductExists is a mock function.
func (m *RecipeStepProductDataManager) RecipeStepProductExists(ctx context.Context, recipeID, recipeStepID, recipeStepProductID string) (bool, error) {
	args := m.Called(ctx, recipeID, recipeStepID, recipeStepProductID)
	return args.Bool(0), args.Error(1)
}

// GetRecipeStepProduct is a mock function.
func (m *RecipeStepProductDataManager) GetRecipeStepProduct(ctx context.Context, recipeID, recipeStepID, recipeStepProductID string) (*types.RecipeStepProduct, error) {
	args := m.Called(ctx, recipeID, recipeStepID, recipeStepProductID)
	return args.Get(0).(*types.RecipeStepProduct), args.Error(1)
}

// GetTotalRecipeStepProductCount is a mock function.
func (m *RecipeStepProductDataManager) GetTotalRecipeStepProductCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetRecipeStepProducts is a mock function.
func (m *RecipeStepProductDataManager) GetRecipeStepProducts(ctx context.Context, recipeID, recipeStepID string, filter *types.QueryFilter) (*types.RecipeStepProductList, error) {
	args := m.Called(ctx, recipeID, recipeStepID, filter)
	return args.Get(0).(*types.RecipeStepProductList), args.Error(1)
}

// GetRecipeStepProductsWithIDs is a mock function.
func (m *RecipeStepProductDataManager) GetRecipeStepProductsWithIDs(ctx context.Context, recipeStepID string, limit uint8, ids []string) ([]*types.RecipeStepProduct, error) {
	args := m.Called(ctx, recipeStepID, limit, ids)
	return args.Get(0).([]*types.RecipeStepProduct), args.Error(1)
}

// CreateRecipeStepProduct is a mock function.
func (m *RecipeStepProductDataManager) CreateRecipeStepProduct(ctx context.Context, input *types.RecipeStepProductDatabaseCreationInput) (*types.RecipeStepProduct, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*types.RecipeStepProduct), args.Error(1)
}

// UpdateRecipeStepProduct is a mock function.
func (m *RecipeStepProductDataManager) UpdateRecipeStepProduct(ctx context.Context, updated *types.RecipeStepProduct) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveRecipeStepProduct is a mock function.
func (m *RecipeStepProductDataManager) ArchiveRecipeStepProduct(ctx context.Context, recipeStepID, recipeStepProductID string) error {
	return m.Called(ctx, recipeStepID, recipeStepProductID).Error(0)
}
