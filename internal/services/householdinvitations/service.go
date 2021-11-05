package householdinvitations

import (
	"fmt"
	"net/http"

	"github.com/prixfixeco/api_server/internal/encoding"
	"github.com/prixfixeco/api_server/internal/messagequeue/publishers"
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

var _ types.HouseholdInvitationDataService = (*service)(nil)

type (
	// SearchIndex is a type alias for dependency injection's sake.
	SearchIndex search.IndexManager

	// service handles to-do list households.
	service struct {
		logger                         logging.Logger
		householdInvitationDataManager types.HouseholdInvitationDataManager
		tracer                         tracing.Tracer
		householdInvitationCounter     metrics.UnitCounter
		encoderDecoder                 encoding.ServerEncoderDecoder
		preWritesPublisher             publishers.Publisher
		sessionContextDataFetcher      func(*http.Request) (*types.SessionContextData, error)
		householdIDFetcher             func(*http.Request) string
		householdInvitationIDFetcher   func(*http.Request) string
	}
)

// ProvideService builds a new HouseholdInvitationsService.
func ProvideService(
	logger logging.Logger,
	cfg Config,
	householdDataManager types.HouseholdInvitationDataManager,
	encoder encoding.ServerEncoderDecoder,
	counterProvider metrics.UnitCounterProvider,
	routeParamManager routing.RouteParamManager,
	publisherProvider publishers.PublisherProvider,
) (types.HouseholdInvitationDataService, error) {
	preWritesPublisher, err := publisherProvider.ProviderPublisher(cfg.PreWritesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up event publisher: %w", err)
	}

	s := &service{
		logger:                         logging.EnsureLogger(logger).WithName(serviceName),
		householdIDFetcher:             routeParamManager.BuildRouteParamStringIDFetcher(HouseholdIDURIParamKey),
		householdInvitationIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(HouseholdInvitationIDURIParamKey),
		sessionContextDataFetcher:      authservice.FetchContextFromRequest,
		householdInvitationDataManager: householdDataManager,
		encoderDecoder:                 encoder,
		preWritesPublisher:             preWritesPublisher,
		householdInvitationCounter:     metrics.EnsureUnitCounter(counterProvider, logger, counterName, counterDescription),
		tracer:                         tracing.NewTracer(serviceName),
	}

	return s, nil
}
