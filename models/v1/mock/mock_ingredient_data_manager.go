package mock

import (
	"context"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/mock"
)

var _ models.IngredientDataManager = (*IngredientDataManager)(nil)

// IngredientDataManager is a mocked models.IngredientDataManager for testing
type IngredientDataManager struct {
	mock.Mock
}

// GetIngredient is a mock function
func (m *IngredientDataManager) GetIngredient(ctx context.Context, ingredientID uint64) (*models.Ingredient, error) {
	args := m.Called(ctx, ingredientID)
	return args.Get(0).(*models.Ingredient), args.Error(1)
}

// GetIngredientCount is a mock function
func (m *IngredientDataManager) GetIngredientCount(ctx context.Context, filter *models.QueryFilter) (uint64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(uint64), args.Error(1)
}

// GetAllIngredientsCount is a mock function
func (m *IngredientDataManager) GetAllIngredientsCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetIngredients is a mock function
func (m *IngredientDataManager) GetIngredients(ctx context.Context, filter *models.QueryFilter) (*models.IngredientList, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*models.IngredientList), args.Error(1)
}

// CreateIngredient is a mock function
func (m *IngredientDataManager) CreateIngredient(ctx context.Context, input *models.IngredientCreationInput) (*models.Ingredient, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*models.Ingredient), args.Error(1)
}

// UpdateIngredient is a mock function
func (m *IngredientDataManager) UpdateIngredient(ctx context.Context, updated *models.Ingredient) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveIngredient is a mock function
func (m *IngredientDataManager) ArchiveIngredient(ctx context.Context, id uint64) error {
	return m.Called(ctx, id).Error(0)
}
