package mock

import (
	"context"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/mock"
)

var _ models.ValidIngredientTagDataManager = (*ValidIngredientTagDataManager)(nil)

// ValidIngredientTagDataManager is a mocked models.ValidIngredientTagDataManager for testing.
type ValidIngredientTagDataManager struct {
	mock.Mock
}

// ValidIngredientTagExists is a mock function.
func (m *ValidIngredientTagDataManager) ValidIngredientTagExists(ctx context.Context, validIngredientTagID uint64) (bool, error) {
	args := m.Called(ctx, validIngredientTagID)
	return args.Bool(0), args.Error(1)
}

// GetValidIngredientTag is a mock function.
func (m *ValidIngredientTagDataManager) GetValidIngredientTag(ctx context.Context, validIngredientTagID uint64) (*models.ValidIngredientTag, error) {
	args := m.Called(ctx, validIngredientTagID)
	return args.Get(0).(*models.ValidIngredientTag), args.Error(1)
}

// GetAllValidIngredientTagsCount is a mock function.
func (m *ValidIngredientTagDataManager) GetAllValidIngredientTagsCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetValidIngredientTags is a mock function.
func (m *ValidIngredientTagDataManager) GetValidIngredientTags(ctx context.Context, filter *models.QueryFilter) (*models.ValidIngredientTagList, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*models.ValidIngredientTagList), args.Error(1)
}

// CreateValidIngredientTag is a mock function.
func (m *ValidIngredientTagDataManager) CreateValidIngredientTag(ctx context.Context, input *models.ValidIngredientTagCreationInput) (*models.ValidIngredientTag, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*models.ValidIngredientTag), args.Error(1)
}

// UpdateValidIngredientTag is a mock function.
func (m *ValidIngredientTagDataManager) UpdateValidIngredientTag(ctx context.Context, updated *models.ValidIngredientTag) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveValidIngredientTag is a mock function.
func (m *ValidIngredientTagDataManager) ArchiveValidIngredientTag(ctx context.Context, validIngredientTagID uint64) error {
	return m.Called(ctx, validIngredientTagID).Error(0)
}
