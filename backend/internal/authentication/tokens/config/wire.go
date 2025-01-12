package tokenscfg

import (
	"github.com/dinnerdonebetter/backend/internal/authentication/tokens"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"

	"github.com/google/wire"
)

var (
	// ProvidersTokenIssuers are what we provide to dependency injection.
	ProvidersTokenIssuers = wire.NewSet(
		ProvideTokenIssuer,
	)
)

// ProvideTokenIssuer provides a tokens.Issuer from a config.
func ProvideTokenIssuer(cfg *Config, logger logging.Logger, tracerProvider tracing.TracerProvider) (tokens.Issuer, error) {
	return cfg.ProvideTokenIssuer(logger, tracerProvider)
}
