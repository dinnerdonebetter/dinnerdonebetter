package analyticscfg

import (
	"github.com/dinnerdonebetter/backend/internal/lib/analytics"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"

	"github.com/google/wire"
)

var (
	// ProvidersAnalytics are what we provide to dependency injection.
	ProvidersAnalytics = wire.NewSet(
		ProvideEventReporter,
	)
)

// ProvideEventReporter provides an analytics.EventReporter from a config.
func ProvideEventReporter(cfg *Config, logger logging.Logger, tracerProvider tracing.TracerProvider, metricsProvider metrics.Provider) (analytics.EventReporter, error) {
	return cfg.ProvideCollector(logger, tracerProvider, metricsProvider)
}
