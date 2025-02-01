package indexing

import (
	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/lib/search/text/indexing"
)

const (
	// IndexTypeUsers represents the users index.
	IndexTypeUsers = "users"
)

func BuildCoreDataIndexingFunctions(dataManager database.DataManager) map[string]indexing.Function {
	return map[string]indexing.Function{
		IndexTypeUsers: dataManager.GetUserIDsThatNeedSearchIndexing,
	}
}
