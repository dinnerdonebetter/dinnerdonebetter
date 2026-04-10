package oauth

import (
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/audit"
	domainoauth "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/oauth"

	"github.com/verygoodsoftwarenotvirus/platform/v5/database"
	databasecfg "github.com/verygoodsoftwarenotvirus/platform/v5/database/config"
	"github.com/verygoodsoftwarenotvirus/platform/v5/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v5/observability/tracing"

	"github.com/samber/do/v2"
)

// RegisterOAuthRepository registers the OAuth repository with the injector.
func RegisterOAuthRepository(i do.Injector) {
	do.Provide[domainoauth.Repository](i, func(i do.Injector) (domainoauth.Repository, error) {
		return ProvideOAuthRepository(
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[audit.Repository](i),
			do.MustInvoke[*databasecfg.Config](i),
			do.MustInvoke[database.Client](i),
		), nil
	})
}
