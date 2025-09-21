package metricscfg

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/metrics"

	"github.com/google/wire"
)

var (
	// MetricsProviders is a Wire provider set that provides a tracing.TracerProvider.
	MetricsProviders = wire.NewSet(
		ProvideMetricsProvider,
	)
)

func ProvideMetricsProvider(ctx context.Context, logger logging.Logger, c *Config) (metrics.Provider, error) {
	return c.ProvideMetricsProvider(ctx, logger)
}
