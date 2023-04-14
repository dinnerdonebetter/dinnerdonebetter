package config

import (
	"net/http"

	"github.com/prixfixeco/backend/internal/featureflags"
	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/tracing"

	"github.com/google/wire"
)

var (
	Providers = wire.NewSet(
		ProvideFeatureFlagManager,
	)
)

func ProvideFeatureFlagManager(c *Config, logger logging.Logger, tracerProvider tracing.TracerProvider, httpClient *http.Client) (featureflags.FeatureFlagManager, error) {
	return c.ProvideFeatureFlagManager(logger, tracerProvider, httpClient)
}
