package tokenscfg

import (
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/authentication/tokens"

	"github.com/samber/do/v2"
	"github.com/verygoodsoftwarenotvirus/platform/v2/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v2/observability/tracing"
)

// ProvideTokenIssuer provides a tokens.Issuer from a config.
func ProvideTokenIssuer(cfg *Config, logger logging.Logger, tracerProvider tracing.TracerProvider) (tokens.Issuer, error) {
	return cfg.ProvideTokenIssuer(logger, tracerProvider)
}

// RegisterTokenIssuer registers the token issuer with the injector.
func RegisterTokenIssuer(i do.Injector) {
	do.Provide[tokens.Issuer](i, func(i do.Injector) (tokens.Issuer, error) {
		return ProvideTokenIssuer(
			do.MustInvoke[*Config](i),
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[tracing.TracerProvider](i),
		)
	})
}
