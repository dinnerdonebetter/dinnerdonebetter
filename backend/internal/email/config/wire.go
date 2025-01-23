package emailcfg

import (
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/email"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"

	"github.com/google/wire"
)

var (
	// ProvidersEmail are what we provide to dependency injection.
	ProvidersEmail = wire.NewSet(
		ProvideEmailer,
	)
)

// ProvideEmailer provides an email.Emailer from a config.
func ProvideEmailer(cfg *Config, logger logging.Logger, tracerProvider tracing.TracerProvider, metricsProvider metrics.Provider, client *http.Client) (email.Emailer, error) {
	circuitBreaker, err := cfg.CircuitBreaker.ProvideCircuitBreaker(logger, metricsProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize email circuit breaker: %w", err)
	}

	return cfg.ProvideEmailer(logger, tracerProvider, client, circuitBreaker)
}
