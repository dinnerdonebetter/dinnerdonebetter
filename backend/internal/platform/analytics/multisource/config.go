package multisource

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/platform/analytics"
	analyticscfg "github.com/dinnerdonebetter/backend/internal/platform/analytics/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

// ProvideMultiSourceEventReporter builds a MultiSourceEventReporter from proxy sources config.
// For each source, attempts to create an EventReporter via ProvideCollector.
// If creation fails (e.g. missing credentials) or provider is unset, uses Noop for that source.
func ProvideMultiSourceEventReporter(
	ctx context.Context,
	proxySources map[string]*analyticscfg.SourceConfig,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	metricsProvider metrics.Provider,
) (*MultiSourceEventReporter, error) {
	reporters := make(map[string]analytics.EventReporter)
	log := logging.EnsureLogger(logger).WithName(name)

	if len(proxySources) == 0 {
		return NewMultiSourceEventReporter(reporters), nil
	}

	for source, sourceCfg := range proxySources {
		r, err := sourceCfg.ProvideCollector(ctx, log, tracerProvider, metricsProvider)
		if err != nil {
			log.WithValue("source", source).WithValue("reason", err.Error()).Info("using noop for source")
			reporters[source] = analytics.NewNoopEventReporter()
			continue
		}
		if r == nil {
			reporters[source] = analytics.NewNoopEventReporter()
			continue
		}
		reporters[source] = r
	}

	return NewMultiSourceEventReporter(reporters), nil
}
