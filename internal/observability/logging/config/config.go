package config

import (
	"context"

	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/logging/googlecloud"
	"github.com/prixfixeco/backend/internal/observability/logging/zerolog"
)

const (
	// ProviderZerolog indicates you'd like to use the zerolog Logger.
	ProviderZerolog = "zerolog"
	// ProviderGoogleCloud indicates you'd like to use the googlecloud Logger.
	ProviderGoogleCloud = "googlecloud"
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
func (cfg *Config) ProvideLogger(ctx context.Context) (logging.Logger, error) {
	var (
		l   logging.Logger
		err error
	)

	switch cfg.Provider {
	case ProviderZerolog:
		l = zerolog.NewZerologLogger()
	case ProviderGoogleCloud:
		l, err = googlecloud.NewLogger(ctx)
		if err != nil {
			return nil, err
		}
	default:
		l = logging.NewNoopLogger()
	}

	l.SetLevel(cfg.Level)

	return l, nil
}
