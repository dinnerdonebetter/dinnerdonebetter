package authentication

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/lib/analytics"
	"github.com/dinnerdonebetter/backend/internal/lib/authentication"
	"github.com/dinnerdonebetter/backend/internal/lib/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/lib/authentication/tokens"
	"github.com/dinnerdonebetter/backend/internal/lib/encoding"
	"github.com/dinnerdonebetter/backend/internal/lib/featureflags"
	"github.com/dinnerdonebetter/backend/internal/lib/internalerrors"
	"github.com/dinnerdonebetter/backend/internal/lib/messagequeue"
	msgconfig "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/lib/routing"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
)

const (
	serviceName                = "auth_service"
	AuthProviderParamKey       = "auth_provider"
	rejectedRequestCounterName = "auth_service.rejected_requests"
)

// TODO: remove this when Goth can handle concurrency.
var useProvidersMutex = sync.Mutex{}

type (
	// service handles passwords service-wide.
	service struct {
		config                     *Config
		logger                     logging.Logger
		authenticator              authentication.Authenticator
		analyticsReporter          analytics.EventReporter
		featureFlagManager         featureflags.FeatureFlagManager
		userDataManager            types.UserDataManager
		householdMembershipManager types.HouseholdUserMembershipDataManager
		encoderDecoder             encoding.ServerEncoderDecoder
		sessionContextDataFetcher  func(*http.Request) (*sessions.ContextData, error)
		authProviderFetcher        func(*http.Request) string
		tracer                     tracing.Tracer
		dataChangesPublisher       messagequeue.Publisher
		oauth2Server               *server.Server
		tokenIssuer                tokens.Issuer
		rejectedRequestCounter     metrics.Int64Counter
	}
)

// ProvideService builds a new AuthDataService.
func ProvideService(
	logger logging.Logger,
	cfg *Config,
	authenticator authentication.Authenticator,
	dataManager database.DataManager,
	householdMembershipManager types.HouseholdUserMembershipDataManager,
	encoder encoding.ServerEncoderDecoder,
	tracerProvider tracing.TracerProvider,
	publisherProvider messagequeue.PublisherProvider,
	featureFlagManager featureflags.FeatureFlagManager,
	analyticsReporter analytics.EventReporter,
	routeParamManager routing.RouteParamManager,
	metricsProvider metrics.Provider,
	queuesConfig *msgconfig.QueuesConfig,
) (types.AuthDataService, error) {
	if queuesConfig == nil {
		return nil, internalerrors.NilConfigError("queuesConfig for AuthDataService")
	}

	dataChangesPublisher, publisherProviderErr := publisherProvider.ProvidePublisher(queuesConfig.DataChangesTopicName)
	if publisherProviderErr != nil {
		return nil, fmt.Errorf("setting up %s data changes publisher: %w", serviceName, publisherProviderErr)
	}

	signer, err := cfg.Tokens.ProvideTokenIssuer(logger, tracerProvider)
	if err != nil {
		return nil, fmt.Errorf("creating json web token signer: %w", err)
	}

	rejectedRequestCounter, err := metricsProvider.NewInt64Counter(rejectedRequestCounterName)
	if err != nil {
		return nil, fmt.Errorf("creating rejected request counter: %w", err)
	}

	tracer := tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName))

	svc := &service{
		logger:                     logging.EnsureLogger(logger).WithName(serviceName),
		encoderDecoder:             encoder,
		config:                     cfg,
		userDataManager:            dataManager,
		householdMembershipManager: householdMembershipManager,
		authenticator:              authenticator,
		sessionContextDataFetcher:  sessions.FetchContextFromRequest,
		tracer:                     tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
		dataChangesPublisher:       dataChangesPublisher,
		featureFlagManager:         featureFlagManager,
		analyticsReporter:          analyticsReporter,
		tokenIssuer:                signer,
		rejectedRequestCounter:     rejectedRequestCounter,
		authProviderFetcher:        routeParamManager.BuildRouteParamStringIDFetcher(AuthProviderParamKey),
		oauth2Server:               ProvideOAuth2ServerImplementation(logger, tracer, &cfg.OAuth2, dataManager, authenticator, signer),
	}

	useProvidersMutex.Lock()
	goth.UseProviders(
		google.New(
			svc.config.SSO.Google.ClientID,
			svc.config.SSO.Google.ClientSecret,
			svc.config.SSO.Google.CallbackURL,
		),
	)
	useProvidersMutex.Unlock()

	return svc, nil
}
