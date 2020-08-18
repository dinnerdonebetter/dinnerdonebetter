package mock

import (
	"context"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/mock"
)

var _ models.RecipeStepProductDataManager = (*RecipeStepProductDataManager)(nil)

// RecipeStepProductDataManager is a mocked models.RecipeStepProductDataManager for testing.
type RecipeStepProductDataManager struct {
	mock.Mock
}

// RecipeStepProductExists is a mock function.
func (m *RecipeStepProductDataManager) RecipeStepProductExists(ctx context.Context, recipeID, recipeStepID, recipeStepProductID uint64) (bool, error) {
	args := m.Called(ctx, recipeID, recipeStepID, recipeStepProductID)
	return args.Bool(0), args.Error(1)
}

// GetRecipeStepProduct is a mock function.
func (m *RecipeStepProductDataManager) GetRecipeStepProduct(ctx context.Context, recipeID, recipeStepID, recipeStepProductID uint64) (*models.RecipeStepProduct, error) {
	args := m.Called(ctx, recipeID, recipeStepID, recipeStepProductID)
	return args.Get(0).(*models.RecipeStepProduct), args.Error(1)
}

// GetAllRecipeStepProductsCount is a mock function.
func (m *RecipeStepProductDataManager) GetAllRecipeStepProductsCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetAllRecipeStepProducts is a mock function.
func (m *RecipeStepProductDataManager) GetAllRecipeStepProducts(ctx context.Context, results chan []models.RecipeStepProduct) error {
	args := m.Called(ctx, results)
	return args.Error(0)
}

// GetRecipeStepProducts is a mock function.
func (m *RecipeStepProductDataManager) GetRecipeStepProducts(ctx context.Context, recipeID, recipeStepID uint64, filter *models.QueryFilter) (*models.RecipeStepProductList, error) {
	args := m.Called(ctx, recipeID, recipeStepID, filter)
	return args.Get(0).(*models.RecipeStepProductList), args.Error(1)
}

// GetRecipeStepProductsWithIDs is a mock function.
func (m *RecipeStepProductDataManager) GetRecipeStepProductsWithIDs(ctx context.Context, recipeID, recipeStepID uint64, limit uint8, ids []uint64) ([]models.RecipeStepProduct, error) {
	args := m.Called(ctx, recipeID, recipeStepID, limit, ids)
	return args.Get(0).([]models.RecipeStepProduct), args.Error(1)
}

// CreateRecipeStepProduct is a mock function.
func (m *RecipeStepProductDataManager) CreateRecipeStepProduct(ctx context.Context, input *models.RecipeStepProductCreationInput) (*models.RecipeStepProduct, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*models.RecipeStepProduct), args.Error(1)
}

// UpdateRecipeStepProduct is a mock function.
func (m *RecipeStepProductDataManager) UpdateRecipeStepProduct(ctx context.Context, updated *models.RecipeStepProduct) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveRecipeStepProduct is a mock function.
func (m *RecipeStepProductDataManager) ArchiveRecipeStepProduct(ctx context.Context, recipeStepID, recipeStepProductID uint64) error {
	return m.Called(ctx, recipeStepID, recipeStepProductID).Error(0)
}
