package routingcfg

import (
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/routing"

	"github.com/google/wire"
)

var (
	// RoutingConfigProviders are what we provide to the dependency injector.
	RoutingConfigProviders = wire.NewSet(
		// ProvideRouterViaConfig,
		ProvideRouteParamManager,
	)
)

func ProvideRouterViaConfig(
	cfg *Config,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	metricProvider metrics.Provider,
) (routing.Router, error) {
	return cfg.ProvideRouter(logger, tracerProvider, metricProvider)
}
