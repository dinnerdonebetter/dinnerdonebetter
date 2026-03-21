package authentication

import (
	"github.com/samber/do/v2"
	"github.com/verygoodsoftwarenotvirus/platform/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/observability/tracing"
)

// RegisterAuth registers authentication providers with the injector.
func RegisterAuth(i do.Injector) {
	do.Provide[Authenticator](i, func(i do.Injector) (Authenticator, error) {
		return ProvideArgon2Authenticator(
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[tracing.TracerProvider](i),
		), nil
	})

	do.Provide[Hasher](i, func(i do.Injector) (Hasher, error) {
		return ProvideHasher(do.MustInvoke[Authenticator](i)), nil
	})
}

func ProvideHasher(authenticator Authenticator) Hasher {
	return authenticator
}
