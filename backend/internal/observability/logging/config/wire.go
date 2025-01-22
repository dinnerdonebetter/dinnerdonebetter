package loggingcfg

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability/logging"

	"github.com/google/wire"
)

var (
	ProvidersLogConfig = wire.NewSet(
		ProvideLogger,
	)
)

func ProvideLogger(ctx context.Context, cfg *Config) (logging.Logger, error) {
	return cfg.ProvideLogger(ctx)
}
