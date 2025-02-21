package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types"

	"github.com/stretchr/testify/mock"
)

var _ types.RecipeStepProductDataManager = (*RecipeStepProductDataManagerMock)(nil)

// RecipeStepProductDataManagerMock is a mocked types.RecipeStepProductDataManager for testing.
type RecipeStepProductDataManagerMock struct {
	mock.Mock
}

// RecipeStepProductExists is a mock function.
func (m *RecipeStepProductDataManagerMock) RecipeStepProductExists(ctx context.Context, recipeID, recipeStepID, recipeStepProductID string) (bool, error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID, recipeStepProductID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetRecipeStepProduct is a mock function.
func (m *RecipeStepProductDataManagerMock) GetRecipeStepProduct(ctx context.Context, recipeID, recipeStepID, recipeStepProductID string) (*types.RecipeStepProduct, error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID, recipeStepProductID)
	return returnValues.Get(0).(*types.RecipeStepProduct), returnValues.Error(1)
}

// GetRecipeStepProducts is a mock function.
func (m *RecipeStepProductDataManagerMock) GetRecipeStepProducts(ctx context.Context, recipeID, recipeStepID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.RecipeStepProduct], error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[types.RecipeStepProduct]), returnValues.Error(1)
}

// CreateRecipeStepProduct is a mock function.
func (m *RecipeStepProductDataManagerMock) CreateRecipeStepProduct(ctx context.Context, input *types.RecipeStepProductDatabaseCreationInput) (*types.RecipeStepProduct, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*types.RecipeStepProduct), returnValues.Error(1)
}

// UpdateRecipeStepProduct is a mock function.
func (m *RecipeStepProductDataManagerMock) UpdateRecipeStepProduct(ctx context.Context, updated *types.RecipeStepProduct) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveRecipeStepProduct is a mock function.
func (m *RecipeStepProductDataManagerMock) ArchiveRecipeStepProduct(ctx context.Context, recipeStepID, recipeStepProductID string) error {
	return m.Called(ctx, recipeStepID, recipeStepProductID).Error(0)
}
