package auditlogentries

import (
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/audit"

	"github.com/samber/do/v2"
	"github.com/verygoodsoftwarenotvirus/platform/database"
	"github.com/verygoodsoftwarenotvirus/platform/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/observability/tracing"
)

// RegisterAuditLogRepository registers the audit log repository with the injector.
func RegisterAuditLogRepository(i do.Injector) {
	do.Provide[audit.Repository](i, func(i do.Injector) (audit.Repository, error) {
		return ProvideAuditLogRepository(
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[database.Client](i),
		), nil
	})
}
