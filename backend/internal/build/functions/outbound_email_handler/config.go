package outboundemailhandler

import (
	"github.com/dinnerdonebetter/backend/internal/config"

	"github.com/google/wire"
)

var (
	ConfigProviders = wire.NewSet(
		wire.FieldsOf(
			new(*config.OutboundEmailHandlerConfig),
			"Queues",
			"Email",
			"Analytics",
			"Events",
			"Observability",
		),
	)
)
