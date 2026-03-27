package interceptors

import (
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/authentication/tokens"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/auth"
	identitymanager "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/identity/manager"

	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/samber/do/v2"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing"
)

// RegisterAuthInterceptor registers the auth interceptor with the injector.
func RegisterAuthInterceptor(i do.Injector) {
	do.Provide[*AuthInterceptor](i, func(i do.Injector) (*AuthInterceptor, error) {
		return ProvideAuthInterceptor(
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[identitymanager.IdentityDataManager](i),
			do.MustInvoke[auth.Repository](i),
			do.MustInvoke[*manage.Manager](i),
			do.MustInvoke[tokens.Issuer](i),
			do.MustInvoke[MethodPermissionsMap](i),
		), nil
	})
}
