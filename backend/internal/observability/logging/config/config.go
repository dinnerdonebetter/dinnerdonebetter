package loggingcfg

import (
	"context"
	"strings"

	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/logging/otelslog"
	"github.com/dinnerdonebetter/backend/internal/observability/logging/slog"
	"github.com/dinnerdonebetter/backend/internal/observability/logging/zap"
	"github.com/dinnerdonebetter/backend/internal/observability/logging/zerolog"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// ProviderZerolog indicates you'd like to use the zerolog logger.
	ProviderZerolog = "zerolog"
	// ProviderZap indicates you'd like to use the zap logger.
	ProviderZap = "zap"
	// ProviderSlog indicates you'd like to use the slog logger.
	ProviderSlog = "slog"
	// ProviderOtelSlog indicates you'd like to use the otel-enabled slog logger.
	ProviderOtelSlog = "otelslog"
)

type (
	// Config configures a Logger.
	Config struct {
		_ struct{} `json:"-"`

		ServiceName string           `env:"SERVICE_NAME" json:"serviceName"`
		Level       logging.Level    `env:"LEVEL"        json:"level,omitempty"`
		OtelSlog    *otelslog.Config `env:"init"         envPrefix:"OTEL_SLOG_"    json:"otelslog,omitempty"`
		Provider    string           `env:"PROVIDER"     json:"provider,omitempty"`
	}
)

func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, &cfg,
		validation.Field(&cfg.ServiceName, validation.Required),
		validation.Field(&cfg.Level, validation.In(logging.AllLevels())),
		validation.Field(&cfg.Provider, validation.In(ProviderZerolog, ProviderZap, ProviderSlog, ProviderOtelSlog)),
		validation.Field(&cfg.OtelSlog, validation.When(cfg.Provider == ProviderOtelSlog, validation.Required)),
	)
}

// ProvideLogger builds a logger according to the provided config.
func (cfg *Config) ProvideLogger(ctx context.Context) (logger logging.Logger, err error) {
	switch strings.TrimSpace(strings.ToLower(cfg.Provider)) {
	case ProviderZerolog:
		logger = zerolog.NewZerologLogger(cfg.Level)
	case ProviderZap:
		logger = zap.NewZapLogger(cfg.Level)
	case ProviderSlog:
		logger = slog.NewSlogLogger(cfg.Level)
	case ProviderOtelSlog:
		logger, err = otelslog.NewOtelSlogLogger(ctx, cfg.Level, cfg.ServiceName, cfg.OtelSlog)
	default:
		logger = logging.NewNoopLogger()
	}

	return logger, err
}
