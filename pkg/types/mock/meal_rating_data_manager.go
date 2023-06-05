package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.MealRatingDataManager = (*MealRatingDataManager)(nil)

// MealRatingDataManager is a mocked types.MealRatingDataManager for testing.
type MealRatingDataManager struct {
	mock.Mock
}

// MealRatingExists is a mock function.
func (m *MealRatingDataManager) MealRatingExists(ctx context.Context, mealRatingID string) (bool, error) {
	args := m.Called(ctx, mealRatingID)
	return args.Bool(0), args.Error(1)
}

// GetMealRating is a mock function.
func (m *MealRatingDataManager) GetMealRating(ctx context.Context, mealRatingID string) (*types.MealRating, error) {
	args := m.Called(ctx, mealRatingID)
	return args.Get(0).(*types.MealRating), args.Error(1)
}

// GetMealRatings is a mock function.
func (m *MealRatingDataManager) GetMealRatings(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.MealRating], error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*types.QueryFilteredResult[types.MealRating]), args.Error(1)
}

// CreateMealRating is a mock function.
func (m *MealRatingDataManager) CreateMealRating(ctx context.Context, input *types.MealRatingDatabaseCreationInput) (*types.MealRating, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*types.MealRating), args.Error(1)
}

// UpdateMealRating is a mock function.
func (m *MealRatingDataManager) UpdateMealRating(ctx context.Context, updated *types.MealRating) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveMealRating is a mock function.
func (m *MealRatingDataManager) ArchiveMealRating(ctx context.Context, mealRatingID string) error {
	return m.Called(ctx, mealRatingID).Error(0)
}
