package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.MealPlanDataManager = (*MealPlanDataManagerMock)(nil)

// MealPlanDataManagerMock is a mocked types.MealPlanDataManager for testing.
type MealPlanDataManagerMock struct {
	mock.Mock
}

// MealPlanExists is a mock function.
func (m *MealPlanDataManagerMock) MealPlanExists(ctx context.Context, mealPlanID, householdID string) (bool, error) {
	returnValues := m.Called(ctx, mealPlanID, householdID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetMealPlan is a mock function.
func (m *MealPlanDataManagerMock) GetMealPlan(ctx context.Context, mealPlanID, householdID string) (*types.MealPlan, error) {
	returnValues := m.Called(ctx, mealPlanID, householdID)
	return returnValues.Get(0).(*types.MealPlan), returnValues.Error(1)
}

// GetMealPlans is a mock function.
func (m *MealPlanDataManagerMock) GetMealPlansForHousehold(ctx context.Context, householdID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.MealPlan], error) {
	returnValues := m.Called(ctx, householdID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[types.MealPlan]), returnValues.Error(1)
}

// CreateMealPlan is a mock function.
func (m *MealPlanDataManagerMock) CreateMealPlan(ctx context.Context, input *types.MealPlanDatabaseCreationInput) (*types.MealPlan, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*types.MealPlan), returnValues.Error(1)
}

// UpdateMealPlan is a mock function.
func (m *MealPlanDataManagerMock) UpdateMealPlan(ctx context.Context, updated *types.MealPlan) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveMealPlan is a mock function.
func (m *MealPlanDataManagerMock) ArchiveMealPlan(ctx context.Context, mealPlanID, householdID string) error {
	return m.Called(ctx, mealPlanID, householdID).Error(0)
}

// AttemptToFinalizeMealPlan is a mock function.
func (m *MealPlanDataManagerMock) AttemptToFinalizeMealPlan(ctx context.Context, mealPlanID, householdID string) (changed bool, err error) {
	returnValues := m.Called(ctx, mealPlanID, householdID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetUnfinalizedMealPlansWithExpiredVotingPeriods is a mock function.
func (m *MealPlanDataManagerMock) GetUnfinalizedMealPlansWithExpiredVotingPeriods(ctx context.Context) ([]*types.MealPlan, error) {
	returnValues := m.Called(ctx)
	return returnValues.Get(0).([]*types.MealPlan), returnValues.Error(1)
}

// GetFinalizedMealPlanIDsForTheNextWeek is a mock function.
func (m *MealPlanDataManagerMock) GetFinalizedMealPlanIDsForTheNextWeek(ctx context.Context) ([]*types.FinalizedMealPlanDatabaseResult, error) {
	returnValues := m.Called(ctx)
	return returnValues.Get(0).([]*types.FinalizedMealPlanDatabaseResult), returnValues.Error(1)
}

// GetFinalizedMealPlansWithUninitializedGroceryLists is a mock function.
func (m *MealPlanDataManagerMock) GetFinalizedMealPlansWithUninitializedGroceryLists(ctx context.Context) ([]*types.MealPlan, error) {
	returnValues := m.Called(ctx)
	return returnValues.Get(0).([]*types.MealPlan), returnValues.Error(1)
}
