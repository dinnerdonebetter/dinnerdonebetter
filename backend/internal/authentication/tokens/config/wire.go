package tokenscfg

import (
	"github.com/dinnerdonebetter/backend/internal/authentication/tokens"

	"github.com/google/wire"
	"github.com/verygoodsoftwarenotvirus/platform/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/observability/tracing"
)

var (
	// TokenIssuerProviders are what we provide to dependency injection.
	TokenIssuerProviders = wire.NewSet(
		ProvideTokenIssuer,
	)
)

// ProvideTokenIssuer provides a tokens.Issuer from a config.
func ProvideTokenIssuer(cfg *Config, logger logging.Logger, tracerProvider tracing.TracerProvider) (tokens.Issuer, error) {
	return cfg.ProvideTokenIssuer(logger, tracerProvider)
}
