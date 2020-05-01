package mock

import (
	"context"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/mock"
)

var _ models.IngredientTagMappingDataManager = (*IngredientTagMappingDataManager)(nil)

// IngredientTagMappingDataManager is a mocked models.IngredientTagMappingDataManager for testing.
type IngredientTagMappingDataManager struct {
	mock.Mock
}

// IngredientTagMappingExists is a mock function.
func (m *IngredientTagMappingDataManager) IngredientTagMappingExists(ctx context.Context, validIngredientID, ingredientTagMappingID uint64) (bool, error) {
	args := m.Called(ctx, validIngredientID, ingredientTagMappingID)
	return args.Bool(0), args.Error(1)
}

// GetIngredientTagMapping is a mock function.
func (m *IngredientTagMappingDataManager) GetIngredientTagMapping(ctx context.Context, validIngredientID, ingredientTagMappingID uint64) (*models.IngredientTagMapping, error) {
	args := m.Called(ctx, validIngredientID, ingredientTagMappingID)
	return args.Get(0).(*models.IngredientTagMapping), args.Error(1)
}

// GetAllIngredientTagMappingsCount is a mock function.
func (m *IngredientTagMappingDataManager) GetAllIngredientTagMappingsCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetIngredientTagMappings is a mock function.
func (m *IngredientTagMappingDataManager) GetIngredientTagMappings(ctx context.Context, validIngredientID uint64, filter *models.QueryFilter) (*models.IngredientTagMappingList, error) {
	args := m.Called(ctx, validIngredientID, filter)
	return args.Get(0).(*models.IngredientTagMappingList), args.Error(1)
}

// CreateIngredientTagMapping is a mock function.
func (m *IngredientTagMappingDataManager) CreateIngredientTagMapping(ctx context.Context, input *models.IngredientTagMappingCreationInput) (*models.IngredientTagMapping, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*models.IngredientTagMapping), args.Error(1)
}

// UpdateIngredientTagMapping is a mock function.
func (m *IngredientTagMappingDataManager) UpdateIngredientTagMapping(ctx context.Context, updated *models.IngredientTagMapping) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveIngredientTagMapping is a mock function.
func (m *IngredientTagMappingDataManager) ArchiveIngredientTagMapping(ctx context.Context, validIngredientID, ingredientTagMappingID uint64) error {
	return m.Called(ctx, validIngredientID, ingredientTagMappingID).Error(0)
}
