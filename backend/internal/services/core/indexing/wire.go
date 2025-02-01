package indexing

import (
	"github.com/google/wire"
)

var (
	ProvidersCoreSearchDataIndexScheduler = wire.NewSet(
		BuildCoreDataIndexingFunctions,
	)
)
