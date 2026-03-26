package manager

import (
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/audit"

	"github.com/samber/do/v2"
	"github.com/verygoodsoftwarenotvirus/platform/v3/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v3/observability/tracing"
)

// RegisterAuditDataManager registers the audit data manager with the injector.
func RegisterAuditDataManager(i do.Injector) {
	do.Provide[AuditDataManager](i, func(i do.Injector) (AuditDataManager, error) {
		return NewAuditDataManager(
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[audit.Repository](i),
		), nil
	})
}
