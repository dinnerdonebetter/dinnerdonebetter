package mealplangrocerylistinitializer

import (
	"github.com/google/wire"
)

var (
	// ProvidersMealPlanGroceryListInitializer are what we provide to dependency injection.
	ProvidersMealPlanGroceryListInitializer = wire.NewSet(
		NewMealPlanGroceryListInitializer,
	)
)
