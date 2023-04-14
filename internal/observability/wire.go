package observability

import "github.com/google/wire"

var (
	Providers = wire.NewSet(
		wire.FieldsOf(
			new(*Config),
			"Logging",
			"Metrics",
			"Tracing",
		),
	)
)
