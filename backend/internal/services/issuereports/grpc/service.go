package grpc

import (
	"context"

	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/authentication/sessions"
	issuereportsmanager "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/issuereports/manager"
	issuereportssvc "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/grpc/generated/services/issue_reports"

	"github.com/verygoodsoftwarenotvirus/platform/v3/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v3/observability/tracing"
)

const (
	o11yName = "issue_reports_service"
)

var _ issuereportssvc.IssueReportsServiceServer = (*serviceImpl)(nil)

type (
	serviceImpl struct {
		issuereportssvc.UnimplementedIssueReportsServiceServer
		tracer                    tracing.Tracer
		logger                    logging.Logger
		sessionContextDataFetcher func(context.Context) (*sessions.ContextData, error)
		issueReportsManager       issuereportsmanager.IssueReportsDataManager
	}
)

func NewService(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	issueReportsManager issuereportsmanager.IssueReportsDataManager,
) issuereportssvc.IssueReportsServiceServer {
	return &serviceImpl{
		logger:                    logging.EnsureLogger(logger).WithName(o11yName),
		tracer:                    tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		sessionContextDataFetcher: sessions.FetchContextDataFromContext,
		issueReportsManager:       issueReportsManager,
	}
}
