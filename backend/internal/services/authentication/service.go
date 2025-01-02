package authentication

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"sync"

	"github.com/dinnerdonebetter/backend/internal/analytics"
	"github.com/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/featureflags"
	"github.com/dinnerdonebetter/backend/internal/messagequeue"
	msgconfig "github.com/dinnerdonebetter/backend/internal/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/internalerrors"
	"github.com/dinnerdonebetter/backend/internal/routing"
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
		sessionContextDataFetcher  func(*http.Request) (*types.SessionContextData, error)
		authProviderFetcher        func(*http.Request) string
		tracer                     tracing.Tracer
		dataChangesPublisher       messagequeue.Publisher
		oauth2Server               *server.Server
		jwtSigner                  authentication.JWTSigner
		rejectedRequestCounter     metrics.Int64Counter
	}
)

// ProvideService builds a new AuthDataService.
func ProvideService(
	ctx context.Context,
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

	decryptedJWTSigningKey, err := base64.URLEncoding.DecodeString(cfg.JWTSigningKey)
	if err != nil {
		return nil, fmt.Errorf("decoding json web token signing key: %w", err)
	}

	signer, err := authentication.NewJWTSigner(logger, tracerProvider, cfg.JWTAudience, decryptedJWTSigningKey)
	if err != nil {
		metrics.NewNoopMetricsProvider()
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
		sessionContextDataFetcher:  authentication.FetchContextFromRequest,
		tracer:                     tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
		dataChangesPublisher:       dataChangesPublisher,
		featureFlagManager:         featureFlagManager,
		analyticsReporter:          analyticsReporter,
		jwtSigner:                  signer,
		rejectedRequestCounter:     rejectedRequestCounter,
		authProviderFetcher:        routeParamManager.BuildRouteParamStringIDFetcher(AuthProviderParamKey),
		oauth2Server:               ProvideOAuth2ServerImplementation(ctx, logger, tracer, &cfg.OAuth2, dataManager, authenticator, signer),
	}

	useProvidersMutex.Lock()
	goth.UseProviders(
		google.New(svc.config.SSO.Google.ClientID, svc.config.SSO.Google.ClientID, svc.config.SSO.Google.CallbackURL),
	)
	useProvidersMutex.Unlock()

	return svc, nil
}
