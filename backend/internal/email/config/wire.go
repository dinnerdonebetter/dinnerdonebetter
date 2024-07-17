package config

import (
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/email"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/circuitbreaking"

	"github.com/google/wire"
)

var (
	// ProvidersEmail are what we provide to dependency injection.
	ProvidersEmail = wire.NewSet(
		ProvideEmailer,
	)
)

// ProvideEmailer provides an email.Emailer from a config.
func ProvideEmailer(cfg *Config, logger logging.Logger, tracerProvider tracing.TracerProvider, client *http.Client) (email.Emailer, error) {
	circuitBreaker := circuitbreaking.ProvideCircuitBreaker(cfg.CircuitBreakerConfig)
	return cfg.ProvideEmailer(logger, tracerProvider, client, circuitBreaker)
}
