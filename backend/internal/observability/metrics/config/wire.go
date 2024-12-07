package config

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability/metrics"

	"github.com/google/wire"
)

var (
	// ProvidersMetrics is a Wire provider set that provides a tracing.TracerProvider.
	ProvidersMetrics = wire.NewSet(
		ProvideMetricsProvider,
	)
)

func ProvideMetricsProvider(ctx context.Context, c *Config) (metrics.Provider, error) {
	return c.ProvideMetricsProvider(ctx)
}
