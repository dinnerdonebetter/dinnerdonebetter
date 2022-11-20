package mocktypes

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/backend/pkg/types"
)

var _ types.MealPlanOptionVoteDataManager = (*MealPlanOptionVoteDataManager)(nil)

// MealPlanOptionVoteDataManager is a mocked types.MealPlanOptionVoteDataManager for testing.
type MealPlanOptionVoteDataManager struct {
	mock.Mock
}

// MealPlanOptionVoteExists is a mock function.
func (m *MealPlanOptionVoteDataManager) MealPlanOptionVoteExists(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID string) (bool, error) {
	args := m.Called(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID)
	return args.Bool(0), args.Error(1)
}

// GetMealPlanOptionVote is a mock function.
func (m *MealPlanOptionVoteDataManager) GetMealPlanOptionVote(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID string) (*types.MealPlanOptionVote, error) {
	args := m.Called(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID)
	return args.Get(0).(*types.MealPlanOptionVote), args.Error(1)
}

// GetMealPlanOptionVotes is a mock function.
func (m *MealPlanOptionVoteDataManager) GetMealPlanOptionVotes(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.MealPlanOptionVote], error) {
	args := m.Called(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, filter)
	return args.Get(0).(*types.QueryFilteredResult[types.MealPlanOptionVote]), args.Error(1)
}

// CreateMealPlanOptionVote is a mock function.
func (m *MealPlanOptionVoteDataManager) CreateMealPlanOptionVote(ctx context.Context, input *types.MealPlanOptionVotesDatabaseCreationInput) ([]*types.MealPlanOptionVote, error) {
	args := m.Called(ctx, input)
	return args.Get(0).([]*types.MealPlanOptionVote), args.Error(1)
}

// UpdateMealPlanOptionVote is a mock function.
func (m *MealPlanOptionVoteDataManager) UpdateMealPlanOptionVote(ctx context.Context, updated *types.MealPlanOptionVote) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveMealPlanOptionVote is a mock function.
func (m *MealPlanOptionVoteDataManager) ArchiveMealPlanOptionVote(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID string) error {
	return m.Called(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID).Error(0)
}

// GetMealPlanOptionVotesForMealPlanOption is a mock function.
func (m *MealPlanOptionVoteDataManager) GetMealPlanOptionVotesForMealPlanOption(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string) ([]*types.MealPlanOptionVote, error) {
	args := m.Called(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID)
	return args.Get(0).([]*types.MealPlanOptionVote), args.Error(1)
}
