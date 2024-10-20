package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.RecipeRatingDataManager = (*RecipeRatingDataManagerMock)(nil)

// RecipeRatingDataManagerMock is a mocked types.RecipeRatingDataManager for testing.
type RecipeRatingDataManagerMock struct {
	mock.Mock
}

// RecipeRatingExists is a mock function.
func (m *RecipeRatingDataManagerMock) RecipeRatingExists(ctx context.Context, recipeID, recipeRatingID string) (bool, error) {
	args := m.Called(ctx, recipeID, recipeRatingID)
	return args.Bool(0), args.Error(1)
}

// GetRecipeRating is a mock function.
func (m *RecipeRatingDataManagerMock) GetRecipeRating(ctx context.Context, recipeID, recipeRatingID string) (*types.RecipeRating, error) {
	args := m.Called(ctx, recipeID, recipeRatingID)
	return args.Get(0).(*types.RecipeRating), args.Error(1)
}

// GetRecipeRatingsForRecipe is a mock function.
func (m *RecipeRatingDataManagerMock) GetRecipeRatingsForRecipe(ctx context.Context, recipeID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.RecipeRating], error) {
	args := m.Called(ctx, recipeID, filter)
	return args.Get(0).(*types.QueryFilteredResult[types.RecipeRating]), args.Error(1)
}

// GetRecipeRatingsForUser is a mock function.
func (m *RecipeRatingDataManagerMock) GetRecipeRatingsForUser(ctx context.Context, user string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.RecipeRating], error) {
	args := m.Called(ctx, user, filter)
	return args.Get(0).(*types.QueryFilteredResult[types.RecipeRating]), args.Error(1)
}

// CreateRecipeRating is a mock function.
func (m *RecipeRatingDataManagerMock) CreateRecipeRating(ctx context.Context, input *types.RecipeRatingDatabaseCreationInput) (*types.RecipeRating, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*types.RecipeRating), args.Error(1)
}

// UpdateRecipeRating is a mock function.
func (m *RecipeRatingDataManagerMock) UpdateRecipeRating(ctx context.Context, updated *types.RecipeRating) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveRecipeRating is a mock function.
func (m *RecipeRatingDataManagerMock) ArchiveRecipeRating(ctx context.Context, recipeID, recipeRatingID string) error {
	return m.Called(ctx, recipeID, recipeRatingID).Error(0)
}
