package config

import (
	"github.com/google/wire"

	"github.com/prixfixeco/backend/internal/analytics"
	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/tracing"
)

var (
	// Providers are what we provide to dependency injection.
	Providers = wire.NewSet(
		ProvideEventReporter,
	)
)

// ProvideEventReporter provides a analytics.EventReporter from a config.
func ProvideEventReporter(cfg *Config, logger logging.Logger, tracerProvider tracing.TracerProvider) (analytics.EventReporter, error) {
	return cfg.ProvideCollector(logger, tracerProvider)
}
