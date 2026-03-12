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
		log.Info("no analytics proxy sources configured, multisource reporter will be empty")
		return NewMultiSourceEventReporter(reporters, logger, tracerProvider), nil
	}

	for source, sourceCfg := range proxySources {
		log.WithValue("source", source).WithValue("provider", sourceCfg.Provider).Info("configuring analytics reporter for proxy source")
		r, err := sourceCfg.ProvideCollector(ctx, log, tracerProvider, metricsProvider)
		if err != nil {
			log.WithValue("source", source).WithValue("reason", err.Error()).Error("failed to create reporter for proxy source, using noop", err)
			reporters[source] = analytics.NewNoopEventReporter()
			continue
		}
		if r == nil {
			log.WithValue("source", source).WithValue("provider", sourceCfg.Provider).Info("ProvideCollector returned nil reporter, using noop")
			reporters[source] = analytics.NewNoopEventReporter()
			continue
		}
		log.WithValue("source", source).WithValue("provider", sourceCfg.Provider).Info("analytics reporter configured for proxy source")
		reporters[source] = r
	}

	return NewMultiSourceEventReporter(reporters, logger, tracerProvider), nil
}
