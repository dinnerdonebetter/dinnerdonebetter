package mocktypes

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/pkg/types"
)

var _ types.MealPlanEventDataManager = (*MealPlanEventDataManager)(nil)

// MealPlanEventDataManager is a mocked types.MealPlanEventDataManager for testing.
type MealPlanEventDataManager struct {
	mock.Mock
}

// MealPlanEventExists is a mock function.
func (m *MealPlanEventDataManager) MealPlanEventExists(ctx context.Context, mealPlanEventID string) (bool, error) {
	args := m.Called(ctx, mealPlanEventID)
	return args.Bool(0), args.Error(1)
}

// GetMealPlanEvent is a mock function.
func (m *MealPlanEventDataManager) GetMealPlanEvent(ctx context.Context, mealPlanEventID string) (*types.MealPlanEvent, error) {
	args := m.Called(ctx, mealPlanEventID)
	return args.Get(0).(*types.MealPlanEvent), args.Error(1)
}

// GetMealPlanEvents is a mock function.
func (m *MealPlanEventDataManager) GetMealPlanEvents(ctx context.Context, filter *types.QueryFilter) (*types.MealPlanEventList, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*types.MealPlanEventList), args.Error(1)
}

// CreateMealPlanEvent is a mock function.
func (m *MealPlanEventDataManager) CreateMealPlanEvent(ctx context.Context, input *types.MealPlanEventDatabaseCreationInput) (*types.MealPlanEvent, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*types.MealPlanEvent), args.Error(1)
}

// UpdateMealPlanEvent is a mock function.
func (m *MealPlanEventDataManager) UpdateMealPlanEvent(ctx context.Context, updated *types.MealPlanEvent) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveMealPlanEvent is a mock function.
func (m *MealPlanEventDataManager) ArchiveMealPlanEvent(ctx context.Context, mealPlanEventID string) error {
	return m.Called(ctx, mealPlanEventID).Error(0)
}

// AttemptToFinalizeMealPlanEvent is a mock function.
func (m *MealPlanEventDataManager) AttemptToFinalizeMealPlanEvent(ctx context.Context, mealPlanEventID string) (changed bool, err error) {
	args := m.Called(ctx, mealPlanEventID)
	return args.Bool(0), args.Error(1)
}

// GetUnfinalizedMealPlanEventsWithExpiredVotingPeriods is a mock function.
func (m *MealPlanEventDataManager) GetUnfinalizedMealPlanEventsWithExpiredVotingPeriods(ctx context.Context) ([]*types.MealPlanEvent, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*types.MealPlanEvent), args.Error(1)
}

// GetFinalizedMealPlanEventIDsForTheNextWeek is a mock function.
func (m *MealPlanEventDataManager) GetFinalizedMealPlanEventIDsForTheNextWeek(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}
