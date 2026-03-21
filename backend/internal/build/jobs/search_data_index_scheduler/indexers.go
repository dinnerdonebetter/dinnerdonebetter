package searchdataindexscheduler

import (
	"maps"

	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	coreindexing "github.com/dinnerdonebetter/backend/internal/services/identity/indexing"
	eatingindexing "github.com/dinnerdonebetter/backend/internal/services/mealplanning/indexing"

	"github.com/verygoodsoftwarenotvirus/platform/search/text/indexing"
)

func ProvideIndexFunctions(identityRepo identity.Repository, mealPlanningRepo mealplanning.Repository) map[string]indexing.Function {
	outputMap := map[string]indexing.Function{}
	coreMap := coreindexing.BuildCoreDataIndexingFunctions(identityRepo)
	eatingMap := eatingindexing.BuildEatingDataIndexingFunctions(mealPlanningRepo)

	maps.Copy(outputMap, coreMap)

	maps.Copy(outputMap, eatingMap)

	return outputMap
}
