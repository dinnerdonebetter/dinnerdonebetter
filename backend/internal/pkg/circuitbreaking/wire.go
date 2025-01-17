package circuitbreaking

import (
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/metrics"

	"github.com/google/wire"
)

var (
	Providers = wire.NewSet(
		ProvideCircuitBreaker,
	)
)

func ProvideCircuitBreaker(cfg *Config, logger logging.Logger, metricsProvider metrics.Provider) (CircuitBreaker, error) {
	return cfg.ProvideCircuitBreaker(logger, metricsProvider)
}
