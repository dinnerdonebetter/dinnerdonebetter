package grpc

import (
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/oauth/manager"
	oauthsvc "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/grpc/generated/services/oauth"

	"github.com/samber/do/v2"
	"github.com/verygoodsoftwarenotvirus/platform/v3/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v3/observability/tracing"
)

// RegisterOAuthService registers the OAuth gRPC service with the injector.
func RegisterOAuthService(i do.Injector) {
	do.Provide[OAuthMethodPermissions](i, func(i do.Injector) (OAuthMethodPermissions, error) {
		return ProvideMethodPermissions(), nil
	})

	do.Provide[oauthsvc.OAuthServiceServer](i, func(i do.Injector) (oauthsvc.OAuthServiceServer, error) {
		return NewService(
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[manager.OAuth2Manager](i),
		), nil
	})
}
