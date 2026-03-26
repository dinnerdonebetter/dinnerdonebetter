package managers

import (
	"context"

	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/auth"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/identity"

	"github.com/samber/do/v2"
	"github.com/verygoodsoftwarenotvirus/platform/v3/messagequeue"
	msgconfig "github.com/verygoodsoftwarenotvirus/platform/v3/messagequeue/config"
	"github.com/verygoodsoftwarenotvirus/platform/v3/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v3/observability/tracing"
	"github.com/verygoodsoftwarenotvirus/platform/v3/qrcodes"
	"github.com/verygoodsoftwarenotvirus/platform/v3/random"
)

// RegisterAuthManager registers the auth manager with the injector.
func RegisterAuthManager(i do.Injector) {
	do.Provide[AuthManagerInterface](i, func(i do.Injector) (AuthManagerInterface, error) {
		return ProvideAuthManager(
			do.MustInvoke[context.Context](i),
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[auth.PasswordResetTokenDataManager](i),
			do.MustInvoke[identity.UserDataManager](i),
			do.MustInvoke[authentication.Authenticator](i),
			do.MustInvoke[messagequeue.PublisherProvider](i),
			do.MustInvoke[random.Generator](i),
			do.MustInvoke[qrcodes.Builder](i),
			do.MustInvoke[*msgconfig.QueuesConfig](i),
		)
	})
}
