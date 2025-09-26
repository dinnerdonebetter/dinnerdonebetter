package grpc

import (
	"testing"

	mockmanagers "github.com/dinnerdonebetter/backend/internal/domain/mealplanning/managers/mock"
	mealplanningsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	mealplanfinalizer "github.com/dinnerdonebetter/backend/internal/services/mealplanning/workers/meal_plan_finalizer"
	mealplangrocerylistinitializer "github.com/dinnerdonebetter/backend/internal/services/mealplanning/workers/meal_plan_grocery_list_initializer"
	mealplantaskcreator "github.com/dinnerdonebetter/backend/internal/services/mealplanning/workers/meal_plan_task_creator"

	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	t.Parallel()

	t.Run("standard", func(t *testing.T) {
		t.Parallel()

		logger := logging.NewNoopLogger()
		tracerProvider := tracing.NewNoopTracerProvider()
		recipeManager := &mockmanagers.MockRecipeManager{}
		validEnumerationsManager := &mockmanagers.MockValidEnumerationsManager{}
		mealPlanningManager := &mockmanagers.MockMealPlanningManager{}
		mealPlanFinalizerWorker := &mealplanfinalizer.Worker{}
		mealPlanGroceryListInitializerWorker := &mealplangrocerylistinitializer.Worker{}
		mealPlanTaskCreatorWorker := &mealplantaskcreator.Worker{}

		service := NewService(
			logger,
			tracerProvider,
			recipeManager,
			validEnumerationsManager,
			mealPlanningManager,
			mealPlanFinalizerWorker,
			mealPlanGroceryListInitializerWorker,
			mealPlanTaskCreatorWorker,
		)

		assert.NotNil(t, service)
		assert.Implements(t, (*mealplanningsvc.MealPlanningServiceServer)(nil), service)

		// Type assertion to ensure we get the correct implementation
		impl, ok := service.(*serviceImpl)
		assert.True(t, ok)
		assert.NotNil(t, impl.logger)
		assert.NotNil(t, impl.tracer)
		assert.Equal(t, recipeManager, impl.recipeManager)
		assert.Equal(t, validEnumerationsManager, impl.validEnumerationsManager)
		assert.Equal(t, mealPlanningManager, impl.mealPlanningManager)
		assert.Equal(t, mealPlanFinalizerWorker, impl.mealPlanFinalizerWorker)
		assert.Equal(t, mealPlanGroceryListInitializerWorker, impl.mealPlanGroceryListInitializerWorker)
		assert.Equal(t, mealPlanTaskCreatorWorker, impl.mealPlanTaskCreatorWorker)
		assert.NotNil(t, impl.sessionContextDataFetcher)
	})
}
