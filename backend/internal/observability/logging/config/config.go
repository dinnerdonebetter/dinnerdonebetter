package config

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

		Level          logging.Level `json:"level,omitempty"          toml:"level"`
		Provider       string        `json:"provider,omitempty"       toml:"provider"`
		OutputFilepath string        `json:"outputFilepath,omitempty" toml:"output_filepath"`
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
