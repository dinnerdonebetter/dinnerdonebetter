package manager

import (
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/customroles"

	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing"

	"github.com/samber/do/v2"
)

// RegisterCustomRolesDataManager registers the custom roles data manager with the injector.
func RegisterCustomRolesDataManager(i do.Injector) {
	do.Provide[CustomRolesDataManager](i, func(i do.Injector) (CustomRolesDataManager, error) {
		return NewCustomRolesDataManager(
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[customroles.Repository](i),
			do.MustInvoke[*authorization.RolePermissionCache](i),
		), nil
	})
}
