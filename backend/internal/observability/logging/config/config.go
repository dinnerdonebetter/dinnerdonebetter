package loggingcfg

import (
	"strings"

	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/logging/slog"
	"github.com/dinnerdonebetter/backend/internal/observability/logging/zap"
	"github.com/dinnerdonebetter/backend/internal/observability/logging/zerolog"
)

const (
	// ProviderZerolog indicates you'd like to use the zerolog logger.
	ProviderZerolog = "zerolog"
	// ProviderZap indicates you'd like to use the zap logger.
	ProviderZap = "zap"
	// ProviderSlog indicates you'd like to use the slog logger.
	ProviderSlog = "slog"
)

type (
	// Config configures a Logger.
	Config struct {
		_ struct{} `json:"-"`

		Level          logging.Level `env:"LEVEL"           json:"level,omitempty"`
		Provider       string        `env:"PROVIDER"        json:"provider,omitempty"`
		OutputFilepath string        `env:"OUTPUT_FILEPATH" json:"outputFilepath,omitempty"`
	}
)

// ProvideLogger builds a logger according to the provided config.
func (cfg *Config) ProvideLogger() logging.Logger {
	var logger logging.Logger

	switch strings.TrimSpace(strings.ToLower(cfg.Provider)) {
	case ProviderZerolog:
		logger = zerolog.NewZerologLogger(cfg.Level)
	case ProviderZap:
		logger = zap.NewZapLogger(cfg.Level)
	case ProviderSlog:
		logger = slog.NewSlogLogger(cfg.Level, cfg.OutputFilepath)
	default:
		logger = logging.NewNoopLogger()
	}

	return logger
}
