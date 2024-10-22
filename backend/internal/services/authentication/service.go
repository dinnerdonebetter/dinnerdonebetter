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
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/routing"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
)

const (
	serviceName          = "auth_service"
	AuthProviderParamKey = "auth_provider"
)

// TODO: remove this.
var useProvidersMutex = sync.Mutex{}

type (
	// cookieEncoderDecoder is a stand-in interface for gorilla/securecookie.
	cookieEncoderDecoder interface {
		Encode(name string, value any) (string, error)
		Decode(name, value string, dst any) error
	}

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
) (types.AuthDataService, error) {
	dataChangesPublisher, publisherProviderErr := publisherProvider.ProvidePublisher(cfg.DataChangesTopicName)
	if publisherProviderErr != nil {
		return nil, fmt.Errorf("setting up %s data changes publisher: %w", serviceName, publisherProviderErr)
	}

	tracer := tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName))

	decryptedJWTSigningKey, err := base64.URLEncoding.DecodeString(cfg.JWTSigningKey)
	if err != nil {
		return nil, fmt.Errorf("decoding Token signing key: %w", err)
	}

	signer, err := authentication.NewJWTSigner(logger, tracerProvider, cfg.JWTAudience, decryptedJWTSigningKey)
	if err != nil {
		return nil, fmt.Errorf("creating Token signer: %w", err)
	}

	svc := &service{
		logger:                     logging.EnsureLogger(logger).WithName(serviceName),
		encoderDecoder:             encoder,
		config:                     cfg,
		userDataManager:            dataManager,
		householdMembershipManager: householdMembershipManager,
		authenticator:              authenticator,
		sessionContextDataFetcher:  FetchContextFromRequest,
		tracer:                     tracer,
		dataChangesPublisher:       dataChangesPublisher,
		featureFlagManager:         featureFlagManager,
		analyticsReporter:          analyticsReporter,
		jwtSigner:                  signer,
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
