package config

import (
	"github.com/google/wire"
	"github.com/prixfixeco/backend/internal/observability/logging"
)

var (
	Providers = wire.NewSet(
		ProvideLogger,
	)
)

func ProvideLogger(cfg *Config) (logging.Logger, error) {
	return cfg.ProvideLogger()
}
