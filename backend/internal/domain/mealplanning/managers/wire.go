package managers

import "github.com/google/wire"

var (
	MPManagerProviders = wire.NewSet(
		NewMealPlanningManager,
		NewRecipeManager,
		NewValidEnumerationsManager,
	)
)
