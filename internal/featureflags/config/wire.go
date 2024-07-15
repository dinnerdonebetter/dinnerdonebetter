package config

import (
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/featureflags"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/circuitbreaking"

	"github.com/google/wire"
)

var (
	ProvidersFeatureFlags = wire.NewSet(
		ProvideFeatureFlagManager,
	)
)

func ProvideFeatureFlagManager(c *Config, logger logging.Logger, tracerProvider tracing.TracerProvider, httpClient *http.Client) (featureflags.FeatureFlagManager, error) {
	circuitBreaker, err := circuitbreaking.ProvideCircuitBreaker(c.CircuitBreakingConfig)
	if err != nil {
		return nil, fmt.Errorf("providing circuit breaker for feature flagging: %w", err)
	}

	return c.ProvideFeatureFlagManager(logger, tracerProvider, httpClient, circuitBreaker)
}
