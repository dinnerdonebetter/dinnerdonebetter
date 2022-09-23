package mocktypes

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/pkg/types"
)

var _ types.AdvancedPrepStepDataManager = (*AdvancedPrepStepDataManager)(nil)

// AdvancedPrepStepDataManager is a mocked types.AdvancedPrepStepDataManager for testing.
type AdvancedPrepStepDataManager struct {
	mock.Mock
}

// AdvancedPrepStepExists is a mock function.
func (m *AdvancedPrepStepDataManager) AdvancedPrepStepExists(ctx context.Context, advancedPrepStepID string) (bool, error) {
	args := m.Called(ctx, advancedPrepStepID)
	return args.Bool(0), args.Error(1)
}

// GetAdvancedPrepStep is a mock function.
func (m *AdvancedPrepStepDataManager) GetAdvancedPrepStep(ctx context.Context, advancedPrepStepID string) (*types.AdvancedPrepStep, error) {
	args := m.Called(ctx, advancedPrepStepID)
	return args.Get(0).(*types.AdvancedPrepStep), args.Error(1)
}

// GetAdvancedPrepStepsForMealPlan is a mock function.
func (m *AdvancedPrepStepDataManager) GetAdvancedPrepStepsForMealPlan(ctx context.Context, mealPlanID string) ([]*types.AdvancedPrepStep, error) {
	args := m.Called(ctx, mealPlanID)
	return args.Get(0).([]*types.AdvancedPrepStep), args.Error(1)
}

// CreateAdvancedPrepStepsForMealPlanOption is a mock function.
func (m *AdvancedPrepStepDataManager) CreateAdvancedPrepStepsForMealPlanOption(ctx context.Context, mealPlanOptionID string, inputs []*types.AdvancedPrepStepDatabaseCreationInput) ([]*types.AdvancedPrepStep, error) {
	args := m.Called(ctx, mealPlanOptionID, inputs)
	return args.Get(0).([]*types.AdvancedPrepStep), args.Error(1)
}

// MarkAdvancedPrepStepAsComplete is a mock function.
func (m *AdvancedPrepStepDataManager) MarkAdvancedPrepStepAsComplete(ctx context.Context, advancedPrepStepID string) error {
	return m.Called(ctx, advancedPrepStepID).Error(0)
}
