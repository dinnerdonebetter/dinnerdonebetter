package config

import (
	"context"

	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/logging/zap"
	"github.com/prixfixeco/backend/internal/observability/logging/zerolog"
)

const (
	// ProviderZerolog indicates you'd like to use the zerolog Logger.
	ProviderZerolog = "zerolog"
	// ProviderZap indicates you'd like to use the zap Logger.
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

// ProvideLogger builds a Logger according to the provided config.
func (cfg *Config) ProvideLogger(_ context.Context) (logging.Logger, error) {
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
