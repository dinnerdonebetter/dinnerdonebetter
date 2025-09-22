package mealplanning

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"

	"github.com/google/wire"
)

var (
	MPRepoProviders = wire.NewSet(
		ProvideMealPlanningRepository,
		ProvideValidEnumerationDataManager,
	)
)

func ProvideValidEnumerationDataManager(x mealplanning.Repository) mealplanning.ValidEnumerationDataManager {
	return x
}
