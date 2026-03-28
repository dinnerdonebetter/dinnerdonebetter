package customroles

import (
	domaincustomroles "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/customroles"

	"github.com/verygoodsoftwarenotvirus/platform/v4/database"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing"

	"github.com/samber/do/v2"
)

// RegisterCustomRolesRepository registers the custom roles repository with the injector.
func RegisterCustomRolesRepository(i do.Injector) {
	do.Provide[domaincustomroles.Repository](i, func(i do.Injector) (domaincustomroles.Repository, error) {
		return ProvideCustomRolesRepository(
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[database.Client](i),
		), nil
	})
}
