package searchindexrequesthandler

import (
	"github.com/dinnerdonebetter/backend/internal/config"

	"github.com/google/wire"
)

var (
	ConfigProviders = wire.NewSet(
		wire.FieldsOf(
			new(*config.SearchIndexRequestHandlerConfig),
			"Queues",
			"Search",
			"Events",
			"Observability",
			"Database",
		),
	)
)
