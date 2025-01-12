package householdinvitations

import (
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/messagequeue"
	msgconfig "github.com/dinnerdonebetter/backend/internal/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/random"
	"github.com/dinnerdonebetter/backend/internal/routing"
	householdsservice "github.com/dinnerdonebetter/backend/internal/services/core/households"
	"github.com/dinnerdonebetter/backend/pkg/types"
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
	}
)

// ProvideHouseholdInvitationsService builds a new HouseholdInvitationDataService.
func ProvideHouseholdInvitationsService(
	logger logging.Logger,
	userDataManager types.UserDataManager,
	householdInvitationDataManager types.HouseholdInvitationDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
	secretGenerator random.Generator,
	queueConfig *msgconfig.QueuesConfig,
) (types.HouseholdInvitationDataService, error) {
	if queueConfig == nil {
		return nil, fmt.Errorf("nil queue config provided")
	}

	dataChangesPublisher, err := publisherProvider.ProvidePublisher(queueConfig.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up %s data changes publisher: %w", serviceName, err)
	}

	s := &service{
		logger:                         logging.EnsureLogger(logger).WithName(serviceName),
		userDataManager:                userDataManager,
		householdInvitationDataManager: householdInvitationDataManager,
		encoderDecoder:                 encoder,
		dataChangesPublisher:           dataChangesPublisher,
		secretGenerator:                secretGenerator,
		sessionContextDataFetcher:      authentication.FetchContextFromRequest,
		householdIDFetcher:             routeParamManager.BuildRouteParamStringIDFetcher(householdsservice.HouseholdIDURIParamKey),
		householdInvitationIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(HouseholdInvitationIDURIParamKey),
		tracer:                         tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
	}

	return s, nil
}
