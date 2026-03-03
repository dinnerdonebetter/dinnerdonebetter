package managers

import (
	mealplangrocerylistinitializer "github.com/dinnerdonebetter/backend/internal/services/mealplanning/workers/meal_plan_grocery_list_initializer"
	mealplantaskcreator "github.com/dinnerdonebetter/backend/internal/services/mealplanning/workers/meal_plan_task_creator"

	"github.com/google/wire"
)

var (
	MPManagerProviders = wire.NewSet(
		NewMealPlanningManager,
		NewRecipeManager,
		NewValidEnumerationsManager,
		BuildMealPlanTaskCreatorWorker,
		BuildMealPlanGroceryListInitializerWorker,
	)
)

func BuildMealPlanTaskCreatorWorker(x *mealplantaskcreator.Worker) mealPlanTaskCreatorWorker {
	return x
}

func BuildMealPlanGroceryListInitializerWorker(x *mealplangrocerylistinitializer.Worker) mealPlanGroceryListInitializerWorker {
	return x
}
