package households

import (
	"fmt"
	"net/http"

	"github.com/prixfixeco/api_server/internal/messagequeue"

	"github.com/prixfixeco/api_server/internal/customerdata"
	"github.com/prixfixeco/api_server/internal/encoding"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/metrics"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/internal/routing"
	"github.com/prixfixeco/api_server/internal/search"
	authservice "github.com/prixfixeco/api_server/internal/services/authentication"
	"github.com/prixfixeco/api_server/pkg/types"
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
		householdInvitationDataManager types.HouseholdInvitationDataManager
		householdMembershipDataManager types.HouseholdUserMembershipDataManager
		tracer                         tracing.Tracer
		householdCounter               metrics.UnitCounter
		encoderDecoder                 encoding.ServerEncoderDecoder
		preWritesPublisher             messagequeue.Publisher
		sessionContextDataFetcher      func(*http.Request) (*types.SessionContextData, error)
		userIDFetcher                  func(*http.Request) string
		householdIDFetcher             func(*http.Request) string
		customerDataCollector          customerdata.Collector
	}
)

// ProvideService builds a new HouseholdsService.
func ProvideService(
	logger logging.Logger,
	cfg Config,
	householdDataManager types.HouseholdDataManager,
	householdInvitationDataManager types.HouseholdInvitationDataManager,
	householdMembershipDataManager types.HouseholdUserMembershipDataManager,
	encoder encoding.ServerEncoderDecoder,
	counterProvider metrics.UnitCounterProvider,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	customerDataCollector customerdata.Collector,
	tracerProvider tracing.TracerProvider,
) (types.HouseholdDataService, error) {
	preWritesPublisher, err := publisherProvider.ProviderPublisher(cfg.PreWritesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up event publisher: %w", err)
	}

	s := &service{
		logger:                         logging.EnsureLogger(logger).WithName(serviceName),
		householdIDFetcher:             routeParamManager.BuildRouteParamStringIDFetcher(HouseholdIDURIParamKey),
		userIDFetcher:                  routeParamManager.BuildRouteParamStringIDFetcher(UserIDURIParamKey),
		sessionContextDataFetcher:      authservice.FetchContextFromRequest,
		householdDataManager:           householdDataManager,
		householdInvitationDataManager: householdInvitationDataManager,
		householdMembershipDataManager: householdMembershipDataManager,
		encoderDecoder:                 encoder,
		preWritesPublisher:             preWritesPublisher,
		householdCounter:               metrics.EnsureUnitCounter(counterProvider, logger, counterName, counterDescription),
		tracer:                         tracing.NewTracer(tracerProvider.Tracer(serviceName)),
		customerDataCollector:          customerDataCollector,
	}

	return s, nil
}
