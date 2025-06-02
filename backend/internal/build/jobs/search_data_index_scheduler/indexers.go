package searchdataindexscheduler

import (
	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/lib/search/text/indexing"
	coreindexing "github.com/dinnerdonebetter/backend/internal/services/core/indexing"
	eatingindexing "github.com/dinnerdonebetter/backend/internal/services/eating/indexing"
)

func ProvideIndexFunctions(dataManager database.DataManager) map[string]indexing.Function {
	outputMap := map[string]indexing.Function{}
	coreMap := coreindexing.BuildCoreDataIndexingFunctions(dataManager)
	eatingMap := eatingindexing.BuildEatingDataIndexingFunctions(dataManager)

	for k, v := range coreMap {
		outputMap[k] = v
	}

	for k, v := range eatingMap {
		outputMap[k] = v
	}

	return outputMap
}
