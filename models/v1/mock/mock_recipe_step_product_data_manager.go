package mock

import (
	"context"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/mock"
)

var _ models.RecipeStepProductDataManager = (*RecipeStepProductDataManager)(nil)

// RecipeStepProductDataManager is a mocked models.RecipeStepProductDataManager for testing
type RecipeStepProductDataManager struct {
	mock.Mock
}

// GetRecipeStepProduct is a mock function
func (m *RecipeStepProductDataManager) GetRecipeStepProduct(ctx context.Context, recipeStepProductID, userID uint64) (*models.RecipeStepProduct, error) {
	args := m.Called(ctx, recipeStepProductID, userID)
	return args.Get(0).(*models.RecipeStepProduct), args.Error(1)
}

// GetRecipeStepProductCount is a mock function
func (m *RecipeStepProductDataManager) GetRecipeStepProductCount(ctx context.Context, filter *models.QueryFilter, userID uint64) (uint64, error) {
	args := m.Called(ctx, filter, userID)
	return args.Get(0).(uint64), args.Error(1)
}

// GetAllRecipeStepProductsCount is a mock function
func (m *RecipeStepProductDataManager) GetAllRecipeStepProductsCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetRecipeStepProducts is a mock function
func (m *RecipeStepProductDataManager) GetRecipeStepProducts(ctx context.Context, filter *models.QueryFilter, userID uint64) (*models.RecipeStepProductList, error) {
	args := m.Called(ctx, filter, userID)
	return args.Get(0).(*models.RecipeStepProductList), args.Error(1)
}

// GetAllRecipeStepProductsForUser is a mock function
func (m *RecipeStepProductDataManager) GetAllRecipeStepProductsForUser(ctx context.Context, userID uint64) ([]models.RecipeStepProduct, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.RecipeStepProduct), args.Error(1)
}

// CreateRecipeStepProduct is a mock function
func (m *RecipeStepProductDataManager) CreateRecipeStepProduct(ctx context.Context, input *models.RecipeStepProductCreationInput) (*models.RecipeStepProduct, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*models.RecipeStepProduct), args.Error(1)
}

// UpdateRecipeStepProduct is a mock function
func (m *RecipeStepProductDataManager) UpdateRecipeStepProduct(ctx context.Context, updated *models.RecipeStepProduct) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveRecipeStepProduct is a mock function
func (m *RecipeStepProductDataManager) ArchiveRecipeStepProduct(ctx context.Context, id, userID uint64) error {
	return m.Called(ctx, id, userID).Error(0)
}
