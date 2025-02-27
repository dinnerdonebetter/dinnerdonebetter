package searchdataindexscheduler

import (
	"maps"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/lib/search/text/indexing"
	coreindexing "github.com/dinnerdonebetter/backend/internal/services/core/indexing"
	eatingindexing "github.com/dinnerdonebetter/backend/internal/services/eating/indexing"
)

func ProvideIndexFunctions(dataManager database.DataManager) map[string]indexing.Function {
	outputMap := map[string]indexing.Function{}
	coreMap := coreindexing.BuildCoreDataIndexingFunctions(dataManager)
	eatingMap := eatingindexing.BuildEatingDataIndexingFunctions(dataManager)

	maps.Copy(outputMap, coreMap)

	maps.Copy(outputMap, eatingMap)

	return outputMap
}
