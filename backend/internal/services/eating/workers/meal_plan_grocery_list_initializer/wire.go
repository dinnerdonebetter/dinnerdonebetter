package mealplangrocerylistinitializer

import (
	"github.com/dinnerdonebetter/backend/internal/database"

	"github.com/google/wire"
)

var (
	// ProvidersMealPlanGroceryListInitializer are what we provide to dependency injection.
	ProvidersMealPlanGroceryListInitializer = wire.NewSet(
		NewMealPlanGroceryListInitializer,
		ProvideMealPlanGroceryListInitializerDataManager,
	)
)

func ProvideMealPlanGroceryListInitializerDataManager(dataManager database.DataManager) dataManager {
	return dataManager
}
