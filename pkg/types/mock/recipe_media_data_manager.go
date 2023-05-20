package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.RecipeMediaDataManager = (*RecipeMediaDataManager)(nil)

// RecipeMediaDataManager is a mocked types.RecipeMediaDataManager for testing.
type RecipeMediaDataManager struct {
	mock.Mock
}

// RecipeMediaExists is a mock function.
func (m *RecipeMediaDataManager) RecipeMediaExists(ctx context.Context, validPreparationID string) (bool, error) {
	args := m.Called(ctx, validPreparationID)
	return args.Bool(0), args.Error(1)
}

// GetRecipeMedia is a mock function.
func (m *RecipeMediaDataManager) GetRecipeMedia(ctx context.Context, validPreparationID string) (*types.RecipeMedia, error) {
	args := m.Called(ctx, validPreparationID)
	return args.Get(0).(*types.RecipeMedia), args.Error(1)
}

// CreateRecipeMedia is a mock function.
func (m *RecipeMediaDataManager) CreateRecipeMedia(ctx context.Context, input *types.RecipeMediaDatabaseCreationInput) (*types.RecipeMedia, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*types.RecipeMedia), args.Error(1)
}

// UpdateRecipeMedia is a mock function.
func (m *RecipeMediaDataManager) UpdateRecipeMedia(ctx context.Context, updated *types.RecipeMedia) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveRecipeMedia is a mock function.
func (m *RecipeMediaDataManager) ArchiveRecipeMedia(ctx context.Context, validPreparationID string) error {
	return m.Called(ctx, validPreparationID).Error(0)
}
