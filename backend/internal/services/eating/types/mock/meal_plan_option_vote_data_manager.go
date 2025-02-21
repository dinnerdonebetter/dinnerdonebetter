package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types"

	"github.com/stretchr/testify/mock"
)

var _ types.MealPlanOptionVoteDataManager = (*MealPlanOptionVoteDataManagerMock)(nil)

// MealPlanOptionVoteDataManagerMock is a mocked types.MealPlanOptionVoteDataManager for testing.
type MealPlanOptionVoteDataManagerMock struct {
	mock.Mock
}

// MealPlanOptionVoteExists is a mock function.
func (m *MealPlanOptionVoteDataManagerMock) MealPlanOptionVoteExists(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID string) (bool, error) {
	returnValues := m.Called(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetMealPlanOptionVote is a mock function.
func (m *MealPlanOptionVoteDataManagerMock) GetMealPlanOptionVote(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID string) (*types.MealPlanOptionVote, error) {
	returnValues := m.Called(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID)
	return returnValues.Get(0).(*types.MealPlanOptionVote), returnValues.Error(1)
}

// GetMealPlanOptionVotes is a mock function.
func (m *MealPlanOptionVoteDataManagerMock) GetMealPlanOptionVotes(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.MealPlanOptionVote], error) {
	returnValues := m.Called(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[types.MealPlanOptionVote]), returnValues.Error(1)
}

// CreateMealPlanOptionVote is a mock function.
func (m *MealPlanOptionVoteDataManagerMock) CreateMealPlanOptionVote(ctx context.Context, input *types.MealPlanOptionVotesDatabaseCreationInput) ([]*types.MealPlanOptionVote, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).([]*types.MealPlanOptionVote), returnValues.Error(1)
}

// UpdateMealPlanOptionVote is a mock function.
func (m *MealPlanOptionVoteDataManagerMock) UpdateMealPlanOptionVote(ctx context.Context, updated *types.MealPlanOptionVote) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveMealPlanOptionVote is a mock function.
func (m *MealPlanOptionVoteDataManagerMock) ArchiveMealPlanOptionVote(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID string) error {
	return m.Called(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID).Error(0)
}

// GetMealPlanOptionVotesForMealPlanOption is a mock function.
func (m *MealPlanOptionVoteDataManagerMock) GetMealPlanOptionVotesForMealPlanOption(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string) ([]*types.MealPlanOptionVote, error) {
	returnValues := m.Called(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID)
	return returnValues.Get(0).([]*types.MealPlanOptionVote), returnValues.Error(1)
}
