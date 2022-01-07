package householdinvitations

import (
	"fmt"
	"net/http"

	"github.com/prixfixeco/api_server/internal/messagequeue"

	"github.com/prixfixeco/api_server/internal/customerdata"
	"github.com/prixfixeco/api_server/internal/encoding"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/internal/random"
	"github.com/prixfixeco/api_server/internal/routing"
	authservice "github.com/prixfixeco/api_server/internal/services/authentication"
	householdsservice "github.com/prixfixeco/api_server/internal/services/households"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	serviceName string = "household_invitations_service"
)

var (
	_ types.HouseholdInvitationDataService = (*service)(nil)
)

type (
	// service handles webhooks.
	service struct {
		logger                         logging.Logger
		userDataManager                types.UserDataManager
		householdInvitationDataManager types.HouseholdInvitationDataManager
		tracer                         tracing.Tracer
		encoderDecoder                 encoding.ServerEncoderDecoder
		secretGenerator                random.Generator
		dataChangesPublisher           messagequeue.Publisher
		householdIDFetcher             func(*http.Request) string
		householdInvitationIDFetcher   func(*http.Request) string
		sessionContextDataFetcher      func(*http.Request) (*types.SessionContextData, error)
		customerDataCollector          customerdata.Collector
	}
)

// ProvideHouseholdInvitationsService builds a new HouseholdInvitationDataService.
func ProvideHouseholdInvitationsService(
	logger logging.Logger,
	cfg *Config,
	userDataManager types.UserDataManager,
	householdInvitationDataManager types.HouseholdInvitationDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	customerDataCollector customerdata.Collector,
	tracerProvider tracing.TracerProvider,
) (types.HouseholdInvitationDataService, error) {
	dataChangesPublisher, err := publisherProvider.ProviderPublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up household invitations service data changes publisher: %w", err)
	}

	s := &service{
		logger:                         logging.EnsureLogger(logger).WithName(serviceName),
		userDataManager:                userDataManager,
		householdInvitationDataManager: householdInvitationDataManager,
		encoderDecoder:                 encoder,
		dataChangesPublisher:           dataChangesPublisher,
		secretGenerator:                random.NewGenerator(logger, tracerProvider),
		sessionContextDataFetcher:      authservice.FetchContextFromRequest,
		householdIDFetcher:             routeParamManager.BuildRouteParamStringIDFetcher(householdsservice.HouseholdIDURIParamKey),
		householdInvitationIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(HouseholdInvitationIDURIParamKey),
		tracer:                         tracing.NewTracer(tracerProvider.Tracer(serviceName)),
		customerDataCollector:          customerDataCollector,
	}

	return s, nil
}
