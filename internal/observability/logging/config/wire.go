package config

import (
	"github.com/prixfixeco/backend/internal/observability/logging"

	"github.com/google/wire"
)

var (
	Providers = wire.NewSet(
		ProvideLogger,
	)
)

func ProvideLogger(cfg *Config) (logging.Logger, error) {
	return cfg.ProvideLogger()
}
