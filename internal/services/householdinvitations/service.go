package householdinvitations

import (
	"fmt"
	"net/http"

	"github.com/prixfixeco/backend/internal/email"
	"github.com/prixfixeco/backend/internal/encoding"
	"github.com/prixfixeco/backend/internal/messagequeue"
	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/internal/random"
	"github.com/prixfixeco/backend/internal/routing"
	authservice "github.com/prixfixeco/backend/internal/services/authentication"
	householdsservice "github.com/prixfixeco/backend/internal/services/households"
	"github.com/prixfixeco/backend/pkg/types"
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
		emailer                        email.Emailer
		secretGenerator                random.Generator
		outboundEmailsPublisher        messagequeue.Publisher
		dataChangesPublisher           messagequeue.Publisher
		householdIDFetcher             func(*http.Request) string
		householdInvitationIDFetcher   func(*http.Request) string
		sessionContextDataFetcher      func(*http.Request) (*types.SessionContextData, error)
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
	tracerProvider tracing.TracerProvider,
	emailer email.Emailer,
	secretGenerator random.Generator,
) (types.HouseholdInvitationDataService, error) {
	dataChangesPublisher, err := publisherProvider.ProviderPublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up household invitations service data changes publisher: %w", err)
	}

	outboundEmailsPublisher, err := publisherProvider.ProviderPublisher(cfg.OutboundEmailsTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up household invitations service data changes publisher: %w", err)
	}

	s := &service{
		logger:                         logging.EnsureLogger(logger).WithName(serviceName),
		userDataManager:                userDataManager,
		householdInvitationDataManager: householdInvitationDataManager,
		encoderDecoder:                 encoder,
		outboundEmailsPublisher:        outboundEmailsPublisher,
		dataChangesPublisher:           dataChangesPublisher,
		emailer:                        emailer,
		secretGenerator:                secretGenerator,
		sessionContextDataFetcher:      authservice.FetchContextFromRequest,
		householdIDFetcher:             routeParamManager.BuildRouteParamStringIDFetcher(householdsservice.HouseholdIDURIParamKey),
		householdInvitationIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(HouseholdInvitationIDURIParamKey),
		tracer:                         tracing.NewTracer(tracerProvider.Tracer(serviceName)),
	}

	return s, nil
}
