package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.RecipeMediaDataManager = (*RecipeMediaDataManagerMock)(nil)

// RecipeMediaDataManagerMock is a mocked types.RecipeMediaDataManager for testing.
type RecipeMediaDataManagerMock struct {
	mock.Mock
}

// RecipeMediaExists is a mock function.
func (m *RecipeMediaDataManagerMock) RecipeMediaExists(ctx context.Context, validPreparationID string) (bool, error) {
	returnValues := m.Called(ctx, validPreparationID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetRecipeMedia is a mock function.
func (m *RecipeMediaDataManagerMock) GetRecipeMedia(ctx context.Context, validPreparationID string) (*types.RecipeMedia, error) {
	returnValues := m.Called(ctx, validPreparationID)
	return returnValues.Get(0).(*types.RecipeMedia), returnValues.Error(1)
}

// CreateRecipeMedia is a mock function.
func (m *RecipeMediaDataManagerMock) CreateRecipeMedia(ctx context.Context, input *types.RecipeMediaDatabaseCreationInput) (*types.RecipeMedia, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*types.RecipeMedia), returnValues.Error(1)
}

// UpdateRecipeMedia is a mock function.
func (m *RecipeMediaDataManagerMock) UpdateRecipeMedia(ctx context.Context, updated *types.RecipeMedia) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveRecipeMedia is a mock function.
func (m *RecipeMediaDataManagerMock) ArchiveRecipeMedia(ctx context.Context, validPreparationID string) error {
	return m.Called(ctx, validPreparationID).Error(0)
}
