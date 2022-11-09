package config

import (
	"github.com/google/wire"

	"github.com/prixfixeco/backend/internal/customerdata"
	"github.com/prixfixeco/backend/internal/observability/logging"
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
