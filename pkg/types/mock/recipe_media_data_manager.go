package mocktypes

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/pkg/types"
)

var _ types.RecipeMediaDataManager = (*MockRecipeMediaDataManager)(nil)

// MockRecipeMediaDataManager is a mocked types.RecipeMediaDataManager for testing.
type MockRecipeMediaDataManager struct {
	mock.Mock
}

// RecipeMediaExists is a mock function.
func (m *MockRecipeMediaDataManager) RecipeMediaExists(ctx context.Context, validPreparationID string) (bool, error) {
	args := m.Called(ctx, validPreparationID)
	return args.Bool(0), args.Error(1)
}

// GetRecipeMedia is a mock function.
func (m *MockRecipeMediaDataManager) GetRecipeMedia(ctx context.Context, validPreparationID string) (*types.RecipeMedia, error) {
	args := m.Called(ctx, validPreparationID)
	return args.Get(0).(*types.RecipeMedia), args.Error(1)
}

// CreateRecipeMedia is a mock function.
func (m *MockRecipeMediaDataManager) CreateRecipeMedia(ctx context.Context, input *types.RecipeMediaDatabaseCreationInput) (*types.RecipeMedia, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*types.RecipeMedia), args.Error(1)
}

// UpdateRecipeMedia is a mock function.
func (m *MockRecipeMediaDataManager) UpdateRecipeMedia(ctx context.Context, updated *types.RecipeMedia) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveRecipeMedia is a mock function.
func (m *MockRecipeMediaDataManager) ArchiveRecipeMedia(ctx context.Context, validPreparationID string) error {
	return m.Called(ctx, validPreparationID).Error(0)
}
