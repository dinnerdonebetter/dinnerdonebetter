package grpc

import (
	"context"

	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/authentication/sessions"
	commentsmanager "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/comments/manager"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/mealplanning/managers"
	uploadedmediamanager "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/uploadedmedia/manager"
	mealplanningsvc "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	_ "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/mealplanning/errors"
	mealplanfinalizer "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/mealplanning/workers/meal_plan_finalizer"
	mealplangrocerylistinitializer "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/mealplanning/workers/meal_plan_grocery_list_initializer"
	mealplantaskcreator "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/mealplanning/workers/meal_plan_task_creator"

	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing"
	"github.com/verygoodsoftwarenotvirus/platform/v4/uploads"
)

var _ mealplanningsvc.MealPlanningServiceServer = (*serviceImpl)(nil)

const (
	o11yName = "mealplanning_service"
)

type (
	serviceImpl struct {
		mealplanningsvc.UnimplementedMealPlanningServiceServer
		tracer                               tracing.Tracer
		logger                               logging.Logger
		sessionContextDataFetcher            func(context.Context) (*sessions.ContextData, error)
		recipeManager                        managers.RecipeManager
		validEnumerationsManager             managers.ValidEnumerationsManager
		mealPlanningManager                  managers.MealPlanningManager
		mealPlanFinalizerWorker              *mealplanfinalizer.Worker
		mealPlanGroceryListInitializerWorker *mealplangrocerylistinitializer.Worker
		mealPlanTaskCreatorWorker            *mealplantaskcreator.Worker
		commentsManager                      commentsmanager.CommentsDataManager
		uploadedMediaManager                 uploadedmediamanager.UploadedMediaManager
		uploadManager                        uploads.UploadManager
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
	commentsManager commentsmanager.CommentsDataManager,
	uploadedMediaManager uploadedmediamanager.UploadedMediaManager,
	uploadManager uploads.UploadManager,
) mealplanningsvc.MealPlanningServiceServer {
	return &serviceImpl{
		logger:                               logging.NewNamedLogger(logger, o11yName),
		tracer:                               tracing.NewNamedTracer(tracerProvider, o11yName),
		recipeManager:                        recipeManager,
		validEnumerationsManager:             validEnumerationsManager,
		mealPlanningManager:                  mealPlanningManager,
		mealPlanFinalizerWorker:              mealPlanFinalizerWorker,
		mealPlanGroceryListInitializerWorker: mealPlanGroceryListInitializerWorker,
		mealPlanTaskCreatorWorker:            mealPlanTaskCreatorWorker,
		commentsManager:                      commentsManager,
		uploadedMediaManager:                 uploadedMediaManager,
		uploadManager:                        uploadManager,
		sessionContextDataFetcher:            sessions.FetchContextDataFromContext,
	}
}
