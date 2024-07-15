package users

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/analytics"
	"github.com/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/featureflags"
	"github.com/dinnerdonebetter/backend/internal/messagequeue"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/random"
	"github.com/dinnerdonebetter/backend/internal/routing"
	authservice "github.com/dinnerdonebetter/backend/internal/services/authentication"
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
		householdUserMembershipDataManager types.HouseholdUserMembershipDataManager
		householdInvitationDataManager     types.HouseholdInvitationDataManager
		passwordResetTokenDataManager      types.PasswordResetTokenDataManager
		tracer                             tracing.Tracer
		authenticator                      authentication.Authenticator
		logger                             logging.Logger
		encoderDecoder                     encoding.ServerEncoderDecoder
		dataChangesPublisher               messagequeue.Publisher
		analyticsReporter                  analytics.EventReporter
		userDataManager                    types.UserDataManager
		secretGenerator                    random.Generator
		userIDFetcher                      func(*http.Request) string
		authSettings                       *authservice.Config
		sessionContextDataFetcher          func(*http.Request) (*types.SessionContextData, error)
		featureFlagManager                 featureflags.FeatureFlagManager
	}
)

// ErrNilConfig is returned when you provide a nil configuration to the users service constructor.
var ErrNilConfig = errors.New("nil config provided")

// ProvideUsersService builds a new UsersService.
func ProvideUsersService(
	_ context.Context,
	cfg *Config,
	authSettings *authservice.Config,
	logger logging.Logger,
	userDataManager types.UserDataManager,
	householdInvitationDataManager types.HouseholdInvitationDataManager,
	householdUserMembershipDataManager types.HouseholdUserMembershipDataManager,
	authenticator authentication.Authenticator,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	tracerProvider tracing.TracerProvider,
	publisherProvider messagequeue.PublisherProvider,
	secretGenerator random.Generator,
	passwordResetTokenDataManager types.PasswordResetTokenDataManager,
	featureFlagManager featureflags.FeatureFlagManager,
	analyticsReporter analytics.EventReporter,
) (types.UserDataService, error) {
	if cfg == nil {
		return nil, ErrNilConfig
	}

	dataChangesPublisher, err := publisherProvider.ProvidePublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up users service data changes publisher: %w", err)
	}

	s := &service{
		logger:                             logging.EnsureLogger(logger).WithName(serviceName),
		userDataManager:                    userDataManager,
		householdInvitationDataManager:     householdInvitationDataManager,
		authenticator:                      authenticator,
		userIDFetcher:                      routeParamManager.BuildRouteParamStringIDFetcher(UserIDURIParamKey),
		sessionContextDataFetcher:          authservice.FetchContextFromRequest,
		encoderDecoder:                     encoder,
		authSettings:                       authSettings,
		secretGenerator:                    secretGenerator,
		householdUserMembershipDataManager: householdUserMembershipDataManager,
		tracer:                             tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
		dataChangesPublisher:               dataChangesPublisher,
		passwordResetTokenDataManager:      passwordResetTokenDataManager,
		featureFlagManager:                 featureFlagManager,
		analyticsReporter:                  analyticsReporter,
	}

	return s, nil
}
