package auditlogentries

import (
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/lib/authentication"
	"github.com/dinnerdonebetter/backend/internal/lib/encoding"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/lib/routing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	serviceName string = "audit_log_entries_service"
)

var _ types.AuditLogEntryDataService = (*service)(nil)

type (
	// service handles audit log entries.
	service struct {
		logger                    logging.Logger
		auditLogEntryDataManager  types.AuditLogEntryDataManager
		auditLogEntryIDFetcher    func(*http.Request) string
		sessionContextDataFetcher func(*http.Request) (*types.SessionContextData, error)
		encoderDecoder            encoding.ServerEncoderDecoder
		tracer                    tracing.Tracer
	}
)

// ProvideService builds a new AuditLogEntriesService.
func ProvideService(
	logger logging.Logger,
	auditLogEntryDataManager types.AuditLogEntryDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	tracerProvider tracing.TracerProvider,
) (types.AuditLogEntryDataService, error) {
	svc := &service{
		logger:                    logging.EnsureLogger(logger).WithName(serviceName),
		auditLogEntryIDFetcher:    routeParamManager.BuildRouteParamStringIDFetcher(AuditLogEntryIDURIParamKey),
		sessionContextDataFetcher: authentication.FetchContextFromRequest,
		auditLogEntryDataManager:  auditLogEntryDataManager,
		encoderDecoder:            encoder,
		tracer:                    tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
	}

	return svc, nil
}
