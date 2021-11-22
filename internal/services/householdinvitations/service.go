package householdinvitations

import (
	"fmt"
	"net/http"

	"github.com/prixfixeco/api_server/internal/customerdata"
	"github.com/prixfixeco/api_server/internal/encoding"
	"github.com/prixfixeco/api_server/internal/messagequeue/publishers"
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
		preWritesPublisher             publishers.Publisher
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
	publisherProvider publishers.PublisherProvider,
	customerDataCollector customerdata.Collector,
) (types.HouseholdInvitationDataService, error) {
	preWritesPublisher, err := publisherProvider.ProviderPublisher(cfg.PreWritesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up pre-writes producer: %w", err)
	}

	s := &service{
		logger:                         logging.EnsureLogger(logger).WithName(serviceName),
		userDataManager:                userDataManager,
		householdInvitationDataManager: householdInvitationDataManager,
		encoderDecoder:                 encoder,
		preWritesPublisher:             preWritesPublisher,
		secretGenerator:                random.NewGenerator(logger),
		sessionContextDataFetcher:      authservice.FetchContextFromRequest,
		householdIDFetcher:             routeParamManager.BuildRouteParamStringIDFetcher(householdsservice.HouseholdIDURIParamKey),
		householdInvitationIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(HouseholdInvitationIDURIParamKey),
		tracer:                         tracing.NewTracer(serviceName),
		customerDataCollector:          customerDataCollector,
	}

	return s, nil
}
