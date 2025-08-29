package grpc

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/managers"
	mealplanningsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	mealplanfinalizer "github.com/dinnerdonebetter/backend/internal/services/mealplanning/workers/meal_plan_finalizer"
	mealplangrocerylistinitializer "github.com/dinnerdonebetter/backend/internal/services/mealplanning/workers/meal_plan_grocery_list_initializer"
	mealplantaskcreator "github.com/dinnerdonebetter/backend/internal/services/mealplanning/workers/meal_plan_task_creator"
)

var _ mealplanningsvc.MealPlanningServiceServer = (*serviceImpl)(nil)

const (
	authHeaderName = "Authorization"

	o11yName = "mealplanning_service"
)

type (
	serviceImpl struct {
		mealplanningsvc.UnimplementedMealPlanningServiceServer
		tracer                               tracing.Tracer
		logger                               logging.Logger
		sessionContextDataFetcher            func(context.Context) (sessions.ContextData, error)
		recipeManager                        managers.RecipeManager
		validEnumerationsManager             managers.ValidEnumerationsManager
		mealPlanningManager                  managers.MealPlanningManager
		mealPlanFinalizerWorker              *mealplanfinalizer.Worker
		mealPlanGroceryListInitializerWorker *mealplangrocerylistinitializer.Worker
		mealPlanTaskCreatorWorker            *mealplantaskcreator.Worker
	}
)

func NewService(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	recipeManager managers.RecipeManager,
	validEnumerationsManager managers.ValidEnumerationsManager,
	mealPlanningManager managers.MealPlanningManager,
	mealPlanFinalizerWorker *mealplanfinalizer.Worker,
	mealPlanGroceryListInitializerWorker *mealplangrocerylistinitializer.Worker,
	mealPlanTaskCreatorWorker *mealplantaskcreator.Worker,
) mealplanningsvc.MealPlanningServiceServer {
	return &serviceImpl{
		logger:                               logging.EnsureLogger(logger).WithName(o11yName),
		tracer:                               tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		recipeManager:                        recipeManager,
		validEnumerationsManager:             validEnumerationsManager,
		mealPlanningManager:                  mealPlanningManager,
		mealPlanFinalizerWorker:              mealPlanFinalizerWorker,
		mealPlanGroceryListInitializerWorker: mealPlanGroceryListInitializerWorker,
		mealPlanTaskCreatorWorker:            mealPlanTaskCreatorWorker,
	}
}
