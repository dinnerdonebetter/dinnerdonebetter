package mealplantaskcreator

import (
	"github.com/google/wire"
)

var (
	// ProvidersMealPlanTaskCreator are what we provide to dependency injection.
	ProvidersMealPlanTaskCreator = wire.NewSet(
		NewMealPlanTaskCreator,
	)
)
