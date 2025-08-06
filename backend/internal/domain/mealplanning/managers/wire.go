package managers

import "github.com/google/wire"

var (
	ProvidersManagers = wire.NewSet(
		NewMealPlanningManager,
		NewRecipeManager,
		NewValidEnumerationsManager,
	)
)
