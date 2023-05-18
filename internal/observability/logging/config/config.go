package config

import (
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/logging/zap"
	"github.com/dinnerdonebetter/backend/internal/observability/logging/zerolog"
)

const (
	// ProviderZerolog indicates you'd like to use the zerolog logger.
	ProviderZerolog = "zerolog"
	// ProviderZap indicates you'd like to use the zap logger.
	ProviderZap = "zap"
)

type (
	// Config configures a zerologLogger.
	Config struct {
		_ struct{}

		Level    logging.Level `json:"level,omitempty"  mapstructure:"level" toml:"level"`
		Provider string        `json:"provider,omitempty" mapstructure:"provider" toml:"provider"`
	}
)

// ProvideLogger builds a logger according to the provided config.
func (cfg *Config) ProvideLogger() (logging.Logger, error) {
	var l logging.Logger

	switch cfg.Provider {
	case ProviderZerolog:
		l = zerolog.NewZerologLogger(cfg.Level)
	case ProviderZap:
		l = zap.NewZapLogger(cfg.Level)
	default:
		l = logging.NewNoopLogger()
	}

	return l, nil
}
