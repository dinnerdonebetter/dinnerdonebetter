package mocktypes

import (
	"context"

	"github.com/stretchr/testify/mock"

	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

var _ types.MealPlanDataManager = (*MealPlanDataManager)(nil)

// MealPlanDataManager is a mocked types.MealPlanDataManager for testing.
type MealPlanDataManager struct {
	mock.Mock
}

// MealPlanExists is a mock function.
func (m *MealPlanDataManager) MealPlanExists(ctx context.Context, mealPlanID string) (bool, error) {
	args := m.Called(ctx, mealPlanID)
	return args.Bool(0), args.Error(1)
}

// GetMealPlan is a mock function.
func (m *MealPlanDataManager) GetMealPlan(ctx context.Context, mealPlanID string) (*types.MealPlan, error) {
	args := m.Called(ctx, mealPlanID)
	return args.Get(0).(*types.MealPlan), args.Error(1)
}

// GetTotalMealPlanCount is a mock function.
func (m *MealPlanDataManager) GetTotalMealPlanCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetMealPlans is a mock function.
func (m *MealPlanDataManager) GetMealPlans(ctx context.Context, filter *types.QueryFilter) (*types.MealPlanList, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*types.MealPlanList), args.Error(1)
}

// GetMealPlansWithIDs is a mock function.
func (m *MealPlanDataManager) GetMealPlansWithIDs(ctx context.Context, householdID string, limit uint8, ids []string) ([]*types.MealPlan, error) {
	args := m.Called(ctx, householdID, limit, ids)
	return args.Get(0).([]*types.MealPlan), args.Error(1)
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
