package grpc

import (
	auditmanager "github.com/dinnerdonebetter/backend/internal/domain/audit/manager"
	auditsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/audit"

	"github.com/samber/do/v2"
	"github.com/verygoodsoftwarenotvirus/platform/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/observability/tracing"
)

// RegisterAuditService registers the audit gRPC service with the injector.
func RegisterAuditService(i do.Injector) {
	do.Provide[AuditMethodPermissions](i, func(i do.Injector) (AuditMethodPermissions, error) {
		return ProvideMethodPermissions(), nil
	})

	do.Provide[auditsvc.AuditServiceServer](i, func(i do.Injector) (auditsvc.AuditServiceServer, error) {
		return NewService(
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[auditmanager.AuditDataManager](i),
		), nil
	})
}
