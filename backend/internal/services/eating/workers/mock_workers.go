package workers

import (
	"context"

	"github.com/stretchr/testify/mock"
)

var _ MealPlanFinalizationWorker = (*MockMealPlanFinalizationWorker)(nil)

type MockMealPlanFinalizationWorker struct {
	mock.Mock
}

// FinalizeExpiredMealPlans implement our interface.
func (m *MockMealPlanFinalizationWorker) FinalizeExpiredMealPlans(ctx context.Context) (int, error) {
	returnValues := m.Called(ctx)
	return returnValues.Get(0).(int), returnValues.Error(1)
}

// FinalizeExpiredMealPlansWithoutReturningCount implement our interface.
func (m *MockMealPlanFinalizationWorker) FinalizeExpiredMealPlansWithoutReturningCount(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}
