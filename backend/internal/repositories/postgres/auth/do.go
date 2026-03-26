package auth

import (
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/audit"
	domainauth "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/auth"

	"github.com/samber/do/v2"
	"github.com/verygoodsoftwarenotvirus/platform/v3/database"
	"github.com/verygoodsoftwarenotvirus/platform/v3/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v3/observability/tracing"
)

// RegisterAuthRepository registers the auth repository with the injector.
func RegisterAuthRepository(i do.Injector) {
	do.Provide[domainauth.Repository](i, func(i do.Injector) (domainauth.Repository, error) {
		return ProvideAuthRepository(
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[audit.Repository](i),
			do.MustInvoke[database.Client](i),
		), nil
	})

	do.Provide[domainauth.PasswordResetTokenDataManager](i, func(i do.Injector) (domainauth.PasswordResetTokenDataManager, error) {
		return ProvidePasswordResetTokenDataManager(do.MustInvoke[domainauth.Repository](i)), nil
	})
}

func ProvidePasswordResetTokenDataManager(r domainauth.Repository) domainauth.PasswordResetTokenDataManager {
	return r
}
