package mock

import (
	"context"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/mock"
)

var _ models.ValidIngredientDataManager = (*ValidIngredientDataManager)(nil)

// ValidIngredientDataManager is a mocked models.ValidIngredientDataManager for testing.
type ValidIngredientDataManager struct {
	mock.Mock
}

// ValidIngredientExists is a mock function.
func (m *ValidIngredientDataManager) ValidIngredientExists(ctx context.Context, validIngredientID uint64) (bool, error) {
	args := m.Called(ctx, validIngredientID)
	return args.Bool(0), args.Error(1)
}

// GetValidIngredient is a mock function.
func (m *ValidIngredientDataManager) GetValidIngredient(ctx context.Context, validIngredientID uint64) (*models.ValidIngredient, error) {
	args := m.Called(ctx, validIngredientID)
	return args.Get(0).(*models.ValidIngredient), args.Error(1)
}

// GetAllValidIngredientsCount is a mock function.
func (m *ValidIngredientDataManager) GetAllValidIngredientsCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetAllValidIngredients is a mock function.
func (m *ValidIngredientDataManager) GetAllValidIngredients(ctx context.Context, results chan []models.ValidIngredient) error {
	args := m.Called(ctx, results)
	return args.Error(0)
}

// GetValidIngredients is a mock function.
func (m *ValidIngredientDataManager) GetValidIngredients(ctx context.Context, filter *models.QueryFilter) (*models.ValidIngredientList, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*models.ValidIngredientList), args.Error(1)
}

// GetValidIngredientsWithIDs is a mock function.
func (m *ValidIngredientDataManager) GetValidIngredientsWithIDs(ctx context.Context, limit uint8, ids []uint64) ([]models.ValidIngredient, error) {
	args := m.Called(ctx, limit, ids)
	return args.Get(0).([]models.ValidIngredient), args.Error(1)
}

// CreateValidIngredient is a mock function.
func (m *ValidIngredientDataManager) CreateValidIngredient(ctx context.Context, input *models.ValidIngredientCreationInput) (*models.ValidIngredient, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*models.ValidIngredient), args.Error(1)
}

// UpdateValidIngredient is a mock function.
func (m *ValidIngredientDataManager) UpdateValidIngredient(ctx context.Context, updated *models.ValidIngredient) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveValidIngredient is a mock function.
func (m *ValidIngredientDataManager) ArchiveValidIngredient(ctx context.Context, validIngredientID uint64) error {
	return m.Called(ctx, validIngredientID).Error(0)
}
