package users

import (
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/platform/analytics"
	"github.com/dinnerdonebetter/backend/internal/platform/authentication"
	"github.com/dinnerdonebetter/backend/internal/platform/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/platform/encoding"
	"github.com/dinnerdonebetter/backend/internal/platform/featureflags"
	"github.com/dinnerdonebetter/backend/internal/platform/messagequeue"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/random"
	"github.com/dinnerdonebetter/backend/internal/platform/routing"
	authservice "github.com/dinnerdonebetter/backend/internal/services/core/handlers/authentication"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	serviceName = "users_service"
)

var _ types.UserDataService = (*service)(nil)

type (
	// RequestValidator validates request.
	RequestValidator interface {
		Validate(req *http.Request) (bool, error)
	}

	// service handles our users.
	service struct {
		accountUserMembershipDataManager types.AccountUserMembershipDataManager
		accountInvitationDataManager     types.AccountInvitationDataManager
		passwordResetTokenDataManager    types.PasswordResetTokenDataManager
		tracer                           tracing.Tracer
		authenticator                    authentication.Authenticator
		logger                           logging.Logger
		encoderDecoder                   encoding.ServerEncoderDecoder
		dataChangesPublisher             messagequeue.Publisher
		analyticsReporter                analytics.EventReporter
		userDataManager                  types.UserDataManager
		secretGenerator                  random.Generator
		userIDFetcher                    func(*http.Request) string
		authSettings                     *authservice.Config
		sessionContextDataFetcher        func(*http.Request) (*sessions.ContextData, error)
		featureFlagManager               featureflags.FeatureFlagManager
	}
)

// ProvideUsersService builds a new UsersService.
func ProvideUsersService(
	authSettings *authservice.Config,
	logger logging.Logger,
	userDataManager types.UserDataManager,
	accountInvitationDataManager types.AccountInvitationDataManager,
	accountUserMembershipDataManager types.AccountUserMembershipDataManager,
	authenticator authentication.Authenticator,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	tracerProvider tracing.TracerProvider,
	publisherProvider messagequeue.PublisherProvider,
	secretGenerator random.Generator,
	passwordResetTokenDataManager types.PasswordResetTokenDataManager,
	featureFlagManager featureflags.FeatureFlagManager,
	analyticsReporter analytics.EventReporter,
	queueConfig *msgconfig.QueuesConfig,
) (types.UserDataService, error) {
	if queueConfig == nil {
		return nil, fmt.Errorf("nil queue config provided")
	}

	dataChangesPublisher, err := publisherProvider.ProvidePublisher(queueConfig.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up %s data changes publisher: %w", serviceName, err)
	}

	s := &service{
		logger:                           logging.EnsureLogger(logger).WithName(serviceName),
		userDataManager:                  userDataManager,
		accountInvitationDataManager:     accountInvitationDataManager,
		authenticator:                    authenticator,
		userIDFetcher:                    routeParamManager.BuildRouteParamStringIDFetcher(UserIDURIParamKey),
		sessionContextDataFetcher:        sessions.FetchContextFromRequest,
		encoderDecoder:                   encoder,
		authSettings:                     authSettings,
		secretGenerator:                  secretGenerator,
		accountUserMembershipDataManager: accountUserMembershipDataManager,
		tracer:                           tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
		dataChangesPublisher:             dataChangesPublisher,
		passwordResetTokenDataManager:    passwordResetTokenDataManager,
		featureFlagManager:               featureFlagManager,
		analyticsReporter:                analyticsReporter,
	}

	return s, nil
}
