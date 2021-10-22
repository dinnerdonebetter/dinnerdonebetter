package mocktypes

import (
	"context"

	"github.com/stretchr/testify/mock"

	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

var _ types.MealPlanOptionVoteDataManager = (*MealPlanOptionVoteDataManager)(nil)

// MealPlanOptionVoteDataManager is a mocked types.MealPlanOptionVoteDataManager for testing.
type MealPlanOptionVoteDataManager struct {
	mock.Mock
}

// MealPlanOptionVoteExists is a mock function.
func (m *MealPlanOptionVoteDataManager) MealPlanOptionVoteExists(ctx context.Context, mealPlanOptionVoteID string) (bool, error) {
	args := m.Called(ctx, mealPlanOptionVoteID)
	return args.Bool(0), args.Error(1)
}

// GetMealPlanOptionVote is a mock function.
func (m *MealPlanOptionVoteDataManager) GetMealPlanOptionVote(ctx context.Context, mealPlanOptionVoteID string) (*types.MealPlanOptionVote, error) {
	args := m.Called(ctx, mealPlanOptionVoteID)
	return args.Get(0).(*types.MealPlanOptionVote), args.Error(1)
}

// GetTotalMealPlanOptionVoteCount is a mock function.
func (m *MealPlanOptionVoteDataManager) GetTotalMealPlanOptionVoteCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetMealPlanOptionVotes is a mock function.
func (m *MealPlanOptionVoteDataManager) GetMealPlanOptionVotes(ctx context.Context, filter *types.QueryFilter) (*types.MealPlanOptionVoteList, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*types.MealPlanOptionVoteList), args.Error(1)
}

// GetMealPlanOptionVotesWithIDs is a mock function.
func (m *MealPlanOptionVoteDataManager) GetMealPlanOptionVotesWithIDs(ctx context.Context, householdID string, limit uint8, ids []string) ([]*types.MealPlanOptionVote, error) {
	args := m.Called(ctx, householdID, limit, ids)
	return args.Get(0).([]*types.MealPlanOptionVote), args.Error(1)
}

// CreateMealPlanOptionVote is a mock function.
func (m *MealPlanOptionVoteDataManager) CreateMealPlanOptionVote(ctx context.Context, input *types.MealPlanOptionVoteDatabaseCreationInput) (*types.MealPlanOptionVote, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*types.MealPlanOptionVote), args.Error(1)
}

// UpdateMealPlanOptionVote is a mock function.
func (m *MealPlanOptionVoteDataManager) UpdateMealPlanOptionVote(ctx context.Context, updated *types.MealPlanOptionVote) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveMealPlanOptionVote is a mock function.
func (m *MealPlanOptionVoteDataManager) ArchiveMealPlanOptionVote(ctx context.Context, mealPlanOptionVoteID, householdID string) error {
	return m.Called(ctx, mealPlanOptionVoteID, householdID).Error(0)
}
