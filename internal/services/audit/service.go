package audit

import (
	"net/http"

	"gitlab.com/prixfixe/prixfixe/internal/encoding"
	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/internal/routing"
	authservice "gitlab.com/prixfixe/prixfixe/internal/services/authentication"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

const (
	serviceName = "audit_log_entries_service"
)

var (
	_ types.AuditLogEntryDataService = (*service)(nil)
)

type (
	// service handles audit log entries.
	service struct {
		logger                    logging.Logger
		auditLog                  types.AuditLogEntryDataManager
		auditLogEntryIDFetcher    func(*http.Request) uint64
		sessionContextDataFetcher func(*http.Request) (*types.SessionContextData, error)
		encoderDecoder            encoding.ServerEncoderDecoder
		tracer                    tracing.Tracer
	}
)

// ProvideService builds a new service.
func ProvideService(
	logger logging.Logger,
	auditLog types.AuditLogEntryDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
) types.AuditLogEntryDataService {
	return &service{
		logger:                    logging.EnsureLogger(logger).WithName(serviceName),
		auditLog:                  auditLog,
		auditLogEntryIDFetcher:    routeParamManager.BuildRouteParamIDFetcher(logger, LogEntryURIParamKey, "audit log entry"),
		sessionContextDataFetcher: authservice.FetchContextFromRequest,
		encoderDecoder:            encoder,
		tracer:                    tracing.NewTracer(serviceName),
	}
}
