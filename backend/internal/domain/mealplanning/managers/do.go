package managers

import (
	"context"

	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/mealplanning/recipeanalysis"
	mealplangrocerylistinitializer "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/mealplanning/workers/meal_plan_grocery_list_initializer"
	mealplantaskcreator "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/mealplanning/workers/meal_plan_task_creator"

	"github.com/samber/do/v2"
	"github.com/verygoodsoftwarenotvirus/platform/v2/messagequeue"
	msgconfig "github.com/verygoodsoftwarenotvirus/platform/v2/messagequeue/config"
	"github.com/verygoodsoftwarenotvirus/platform/v2/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v2/observability/metrics"
	"github.com/verygoodsoftwarenotvirus/platform/v2/observability/tracing"
	textsearchcfg "github.com/verygoodsoftwarenotvirus/platform/v2/search/text/config"
)

// RegisterManagers registers all meal planning managers with the injector.
func RegisterManagers(i do.Injector) {
	do.Provide[mealPlanTaskCreatorWorker](i, func(i do.Injector) (mealPlanTaskCreatorWorker, error) {
		return BuildMealPlanTaskCreatorWorker(do.MustInvoke[*mealplantaskcreator.Worker](i)), nil
	})

	do.Provide[mealPlanGroceryListInitializerWorker](i, func(i do.Injector) (mealPlanGroceryListInitializerWorker, error) {
		return BuildMealPlanGroceryListInitializerWorker(do.MustInvoke[*mealplangrocerylistinitializer.Worker](i)), nil
	})

	do.Provide[MealPlanningManager](i, func(i do.Injector) (MealPlanningManager, error) {
		return NewMealPlanningManager(
			do.MustInvoke[context.Context](i),
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[mealplanning.Repository](i),
			do.MustInvoke[*msgconfig.QueuesConfig](i),
			do.MustInvoke[messagequeue.PublisherProvider](i),
			do.MustInvoke[*textsearchcfg.Config](i),
			do.MustInvoke[metrics.Provider](i),
			do.MustInvoke[mealPlanGroceryListInitializerWorker](i),
			do.MustInvoke[mealPlanTaskCreatorWorker](i),
		)
	})

	do.Provide[RecipeManager](i, func(i do.Injector) (RecipeManager, error) {
		return NewRecipeManager(
			do.MustInvoke[context.Context](i),
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[mealplanning.Repository](i),
			do.MustInvoke[*msgconfig.QueuesConfig](i),
			do.MustInvoke[messagequeue.PublisherProvider](i),
			do.MustInvoke[recipeanalysis.RecipeAnalyzer](i),
			do.MustInvoke[*textsearchcfg.Config](i),
			do.MustInvoke[metrics.Provider](i),
		)
	})

	do.Provide[ValidEnumerationsManager](i, func(i do.Injector) (ValidEnumerationsManager, error) {
		return NewValidEnumerationsManager(
			do.MustInvoke[context.Context](i),
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[mealplanning.ValidEnumerationDataManager](i),
			do.MustInvoke[*msgconfig.QueuesConfig](i),
			do.MustInvoke[messagequeue.PublisherProvider](i),
			do.MustInvoke[*textsearchcfg.Config](i),
			do.MustInvoke[metrics.Provider](i),
		)
	})
}

func BuildMealPlanTaskCreatorWorker(x *mealplantaskcreator.Worker) mealPlanTaskCreatorWorker {
	return x
}

func BuildMealPlanGroceryListInitializerWorker(x *mealplangrocerylistinitializer.Worker) mealPlanGroceryListInitializerWorker {
	return x
}
