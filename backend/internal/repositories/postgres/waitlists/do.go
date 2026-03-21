package waitlists

import (
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/audit"

	"github.com/samber/do/v2"
	"github.com/verygoodsoftwarenotvirus/platform/database"
	"github.com/verygoodsoftwarenotvirus/platform/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/observability/tracing"
)

// RegisterWaitlistsRepository registers the waitlists repository with the injector.
func RegisterWaitlistsRepository(i do.Injector) {
	do.Provide[*Repository](i, func(i do.Injector) (*Repository, error) {
		return ProvideWaitlistsRepository(
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[audit.Repository](i),
			do.MustInvoke[database.Client](i),
		), nil
	})
}
