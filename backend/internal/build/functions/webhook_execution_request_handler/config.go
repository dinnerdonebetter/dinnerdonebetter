package webhookexecutionrequesthandler

import (
	"github.com/dinnerdonebetter/backend/internal/config"

	"github.com/google/wire"
)

var (
	ConfigProviders = wire.NewSet(
		wire.FieldsOf(
			new(*config.WebhookExecutionRequestHandlerConfig),
			"Queues",
			"Encoding",
			"Events",
			"Observability",
			"Database",
		),
	)
)
