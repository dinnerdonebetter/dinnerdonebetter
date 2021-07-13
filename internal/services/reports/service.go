package reports

import (
	"net/http"

	"gitlab.com/prixfixe/prixfixe/internal/encoding"
	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
	"gitlab.com/prixfixe/prixfixe/internal/observability/metrics"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	routing "gitlab.com/prixfixe/prixfixe/internal/routing"
	"gitlab.com/prixfixe/prixfixe/internal/search"
	authservice "gitlab.com/prixfixe/prixfixe/internal/services/authentication"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

const (
	counterName        metrics.CounterName = "reports"
	counterDescription string              = "the number of reports managed by the reports service"
	serviceName        string              = "reports_service"
)

var _ types.ReportDataService = (*service)(nil)

type (
	// SearchIndex is a type alias for dependency injection's sake.
	SearchIndex search.IndexManager

	// service handles reports.
	service struct {
		logger                    logging.Logger
		reportDataManager         types.ReportDataManager
		reportIDFetcher           func(*http.Request) uint64
		sessionContextDataFetcher func(*http.Request) (*types.SessionContextData, error)
		reportCounter             metrics.UnitCounter
		encoderDecoder            encoding.ServerEncoderDecoder
		tracer                    tracing.Tracer
	}
)

// ProvideService builds a new ReportsService.
func ProvideService(
	logger logging.Logger,
	cfg Config,
	reportDataManager types.ReportDataManager,
	encoder encoding.ServerEncoderDecoder,
	counterProvider metrics.UnitCounterProvider,
	routeParamManager routing.RouteParamManager,
) (types.ReportDataService, error) {
	svc := &service{
		logger:                    logging.EnsureLogger(logger).WithName(serviceName),
		reportIDFetcher:           routeParamManager.BuildRouteParamIDFetcher(logger, ReportIDURIParamKey, "report"),
		sessionContextDataFetcher: authservice.FetchContextFromRequest,
		reportDataManager:         reportDataManager,
		encoderDecoder:            encoder,
		reportCounter:             metrics.EnsureUnitCounter(counterProvider, logger, counterName, counterDescription),
		tracer:                    tracing.NewTracer(serviceName),
	}

	return svc, nil
}
