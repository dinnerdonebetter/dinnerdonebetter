package grpc

import (
	issuereportsmanager "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/issuereports/manager"
	issuereportssvc "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/grpc/generated/services/issue_reports"

	"github.com/samber/do/v2"
	"github.com/verygoodsoftwarenotvirus/platform/v2/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v2/observability/tracing"
)

// RegisterIssueReportsService registers the issue reports gRPC service with the injector.
func RegisterIssueReportsService(i do.Injector) {
	do.Provide[IssueReportsMethodPermissions](i, func(i do.Injector) (IssueReportsMethodPermissions, error) {
		return ProvideMethodPermissions(), nil
	})

	do.Provide[issuereportssvc.IssueReportsServiceServer](i, func(i do.Injector) (issuereportssvc.IssueReportsServiceServer, error) {
		return NewService(
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[issuereportsmanager.IssueReportsDataManager](i),
		), nil
	})
}
