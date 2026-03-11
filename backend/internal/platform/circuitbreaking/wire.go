package circuitbreaking

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/metrics"

	"github.com/google/wire"
)

var (
	Providers = wire.NewSet(
		ProvideCircuitBreaker,
	)
)

func ProvideCircuitBreaker(ctx context.Context, cfg *Config, logger logging.Logger, metricsProvider metrics.Provider) (CircuitBreaker, error) {
	return cfg.ProvideCircuitBreaker(ctx, logger, metricsProvider)
}
