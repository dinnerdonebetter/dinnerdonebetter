package config

import (
	"github.com/google/wire"

	"github.com/prixfixeco/api_server/internal/customerdata"
	"github.com/prixfixeco/api_server/internal/observability/logging"
)

var (
	// Providers are what we provide to dependency injection.
	Providers = wire.NewSet(
		ProvideCollector,
	)
)

// ProvideCollector provides a customerdata.Collector from a config.
func ProvideCollector(cfg *Config, logger logging.Logger) (customerdata.Collector, error) {
	return cfg.ProvideCollector(logger)
}
