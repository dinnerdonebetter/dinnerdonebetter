package searchdataindexscheduler

import (
	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/search/text/indexing"
	coreindexing "github.com/dinnerdonebetter/backend/internal/services/core/indexing"
	eatingindexing "github.com/dinnerdonebetter/backend/internal/services/eating/indexing"
)

func ProvideIndexFunctions(identityRepo identity.Repository, mealPlanningRepo mealplanning.Repository) map[string]indexing.Function {
	outputMap := map[string]indexing.Function{}
	coreMap := coreindexing.BuildCoreDataIndexingFunctions(identityRepo)
	eatingMap := eatingindexing.BuildEatingDataIndexingFunctions(mealPlanningRepo)

	for k, v := range coreMap {
		outputMap[k] = v
	}

	for k, v := range eatingMap {
		outputMap[k] = v
	}

	return outputMap
}
