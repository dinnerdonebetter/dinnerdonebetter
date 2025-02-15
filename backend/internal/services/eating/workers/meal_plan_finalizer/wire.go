package mealplanfinalizer

import (
	"github.com/google/wire"
)

var (
	// ProvidersMealPlanFinalizer are what we provide to dependency injection.
	ProvidersMealPlanFinalizer = wire.NewSet(
		NewMealPlanFinalizer,
	)
)
