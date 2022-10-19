package mocktypes

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/pkg/types"
)

var _ types.MealPlanDataManager = (*MealPlanDataManager)(nil)

// MealPlanDataManager is a mocked types.MealPlanDataManager for testing.
type MealPlanDataManager struct {
	mock.Mock
}

// MealPlanExists is a mock function.
func (m *MealPlanDataManager) MealPlanExists(ctx context.Context, mealPlanID, householdID string) (bool, error) {
	args := m.Called(ctx, mealPlanID, householdID)
	return args.Bool(0), args.Error(1)
}

// GetMealPlan is a mock function.
func (m *MealPlanDataManager) GetMealPlan(ctx context.Context, mealPlanID, householdID string) (*types.MealPlan, error) {
	args := m.Called(ctx, mealPlanID, householdID)
	return args.Get(0).(*types.MealPlan), args.Error(1)
}

// GetMealPlans is a mock function.
func (m *MealPlanDataManager) GetMealPlans(ctx context.Context, householdID string, filter *types.QueryFilter) (*types.MealPlanList, error) {
	args := m.Called(ctx, householdID, filter)
	return args.Get(0).(*types.MealPlanList), args.Error(1)
}

// CreateMealPlan is a mock function.
func (m *MealPlanDataManager) CreateMealPlan(ctx context.Context, input *types.MealPlanDatabaseCreationInput) (*types.MealPlan, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*types.MealPlan), args.Error(1)
}

// UpdateMealPlan is a mock function.
func (m *MealPlanDataManager) UpdateMealPlan(ctx context.Context, updated *types.MealPlan) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveMealPlan is a mock function.
func (m *MealPlanDataManager) ArchiveMealPlan(ctx context.Context, mealPlanID, householdID string) error {
	return m.Called(ctx, mealPlanID, householdID).Error(0)
}

// AttemptToFinalizeMealPlan is a mock function.
func (m *MealPlanDataManager) AttemptToFinalizeMealPlan(ctx context.Context, mealPlanID, householdID string) (changed bool, err error) {
	args := m.Called(ctx, mealPlanID, householdID)
	return args.Bool(0), args.Error(1)
}

// GetUnfinalizedMealPlansWithExpiredVotingPeriods is a mock function.
func (m *MealPlanDataManager) GetUnfinalizedMealPlansWithExpiredVotingPeriods(ctx context.Context) ([]*types.MealPlan, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*types.MealPlan), args.Error(1)
}

// GetFinalizedMealPlanIDsForTheNextWeek is a mock function.
func (m *MealPlanDataManager) GetFinalizedMealPlanIDsForTheNextWeek(ctx context.Context) ([]*types.FinalizedMealPlanDatabaseResult, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*types.FinalizedMealPlanDatabaseResult), args.Error(1)
}

// GetFinalizedMealPlansWithUninitializedGroceryLists is a mock function.
func (m *MealPlanDataManager) GetFinalizedMealPlansWithUninitializedGroceryLists(ctx context.Context) ([]*types.MealPlan, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*types.MealPlan), args.Error(1)
}
