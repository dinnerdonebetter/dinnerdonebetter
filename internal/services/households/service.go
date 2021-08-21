package households

import (
	"net/http"

	"gitlab.com/prixfixe/prixfixe/internal/encoding"
	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
	"gitlab.com/prixfixe/prixfixe/internal/observability/metrics"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/internal/routing"
	"gitlab.com/prixfixe/prixfixe/internal/search"
	authservice "gitlab.com/prixfixe/prixfixe/internal/services/authentication"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

const (
	counterName        metrics.CounterName = "households"
	counterDescription string              = "the number of households managed by the households service"
	serviceName        string              = "households_service"
)

var _ types.HouseholdDataService = (*service)(nil)

type (
	// SearchIndex is a type alias for dependency injection's sake.
	SearchIndex search.IndexManager

	// service handles to-do list households.
	service struct {
		logger                         logging.Logger
		householdDataManager           types.HouseholdDataManager
		householdMembershipDataManager types.HouseholdUserMembershipDataManager
		householdIDFetcher             func(*http.Request) uint64
		userIDFetcher                  func(*http.Request) uint64
		sessionContextDataFetcher      func(*http.Request) (*types.SessionContextData, error)
		householdCounter               metrics.UnitCounter
		encoderDecoder                 encoding.ServerEncoderDecoder
		tracer                         tracing.Tracer
	}
)

// ProvideService builds a new HouseholdsService.
func ProvideService(
	logger logging.Logger,
	householdDataManager types.HouseholdDataManager,
	householdMembershipDataManager types.HouseholdUserMembershipDataManager,
	encoder encoding.ServerEncoderDecoder,
	counterProvider metrics.UnitCounterProvider,
	routeParamManager routing.RouteParamManager,
) types.HouseholdDataService {
	return &service{
		logger:                         logging.EnsureLogger(logger).WithName(serviceName),
		householdIDFetcher:             routeParamManager.BuildRouteParamIDFetcher(logger, HouseholdIDURIParamKey, "household"),
		userIDFetcher:                  routeParamManager.BuildRouteParamIDFetcher(logger, UserIDURIParamKey, "user"),
		sessionContextDataFetcher:      authservice.FetchContextFromRequest,
		householdDataManager:           householdDataManager,
		householdMembershipDataManager: householdMembershipDataManager,
		encoderDecoder:                 encoder,
		householdCounter:               metrics.EnsureUnitCounter(counterProvider, logger, counterName, counterDescription),
		tracer:                         tracing.NewTracer(serviceName),
	}
}
