package grpc

import (
	commentsmanager "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/comments/manager"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/mealplanning/managers"
	uploadedmediamanager "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/uploadedmedia/manager"
	mealplanningsvc "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	mealplanfinalizer "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/mealplanning/workers/meal_plan_finalizer"
	mealplangrocerylistinitializer "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/mealplanning/workers/meal_plan_grocery_list_initializer"
	mealplantaskcreator "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/mealplanning/workers/meal_plan_task_creator"

	"github.com/samber/do/v2"
	"github.com/verygoodsoftwarenotvirus/platform/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/observability/tracing"
	"github.com/verygoodsoftwarenotvirus/platform/uploads"
)

// RegisterMealPlanningService registers the meal planning gRPC service with the injector.
func RegisterMealPlanningService(i do.Injector) {
	do.Provide[MealPlanningMethodPermissions](i, func(i do.Injector) (MealPlanningMethodPermissions, error) {
		return ProvideMethodPermissions(), nil
	})

	do.Provide[mealplanningsvc.MealPlanningServiceServer](i, func(i do.Injector) (mealplanningsvc.MealPlanningServiceServer, error) {
		return NewService(
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[managers.RecipeManager](i),
			do.MustInvoke[managers.ValidEnumerationsManager](i),
			do.MustInvoke[managers.MealPlanningManager](i),
			do.MustInvoke[*mealplanfinalizer.Worker](i),
			do.MustInvoke[*mealplangrocerylistinitializer.Worker](i),
			do.MustInvoke[*mealplantaskcreator.Worker](i),
			do.MustInvoke[commentsmanager.CommentsDataManager](i),
			do.MustInvoke[uploadedmediamanager.UploadedMediaManager](i),
			do.MustInvoke[uploads.UploadManager](i),
		), nil
	})
}
