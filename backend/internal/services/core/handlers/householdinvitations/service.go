package accountinvitations

import (
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/platform/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/platform/encoding"
	"github.com/dinnerdonebetter/backend/internal/platform/messagequeue"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/random"
	"github.com/dinnerdonebetter/backend/internal/platform/routing"
	accountsservice "github.com/dinnerdonebetter/backend/internal/services/core/handlers/accounts"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	serviceName string = "account_invitations_service"
)

var (
	_ types.AccountInvitationDataService = (*service)(nil)
)

type (
	// service handles webhooks.
	service struct {
		logger                       logging.Logger
		userDataManager              types.UserDataManager
		accountInvitationDataManager types.AccountInvitationDataManager
		tracer                       tracing.Tracer
		encoderDecoder               encoding.ServerEncoderDecoder
		secretGenerator              random.Generator
		dataChangesPublisher         messagequeue.Publisher
		accountIDFetcher             func(*http.Request) string
		accountInvitationIDFetcher   func(*http.Request) string
		sessionContextDataFetcher    func(*http.Request) (*sessions.ContextData, error)
	}
)

// ProvideAccountInvitationsService builds a new AccountInvitationDataService.
func ProvideAccountInvitationsService(
	logger logging.Logger,
	userDataManager types.UserDataManager,
	accountInvitationDataManager types.AccountInvitationDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
	secretGenerator random.Generator,
	queueConfig *msgconfig.QueuesConfig,
) (types.AccountInvitationDataService, error) {
	if queueConfig == nil {
		return nil, fmt.Errorf("nil queue config provided")
	}

	dataChangesPublisher, err := publisherProvider.ProvidePublisher(queueConfig.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up %s data changes publisher: %w", serviceName, err)
	}

	s := &service{
		logger:                       logging.EnsureLogger(logger).WithName(serviceName),
		userDataManager:              userDataManager,
		accountInvitationDataManager: accountInvitationDataManager,
		encoderDecoder:               encoder,
		dataChangesPublisher:         dataChangesPublisher,
		secretGenerator:              secretGenerator,
		sessionContextDataFetcher:    sessions.FetchContextFromRequest,
		accountIDFetcher:             routeParamManager.BuildRouteParamStringIDFetcher(accountsservice.AccountIDURIParamKey),
		accountInvitationIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(AccountInvitationIDURIParamKey),
		tracer:                       tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
	}

	return s, nil
}
