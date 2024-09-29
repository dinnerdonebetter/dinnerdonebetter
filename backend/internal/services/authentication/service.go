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

	"github.com/alexedwards/scs/v2"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/gorilla/securecookie"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
)

const (
	serviceName           = "auth_service"
	userIDContextKey      = string(types.UserIDContextKey)
	householdIDContextKey = string(types.HouseholdIDContextKey)
	cookieErrorLogName    = "_COOKIE_CONSTRUCTION_ERROR_"
	cookieSecretSize      = 64
	AuthProviderParamKey  = "auth_provider"
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
		cookieManager              cookieEncoderDecoder
		sessionManager             sessionManager
		sessionContextDataFetcher  func(*http.Request) (*types.SessionContextData, error)
		authProviderFetcher        func(*http.Request) string
		tracer                     tracing.Tracer
		dataChangesPublisher       messagequeue.Publisher
		oauth2Server               *server.Server
		jwtSigner                  authentication.JWTSigner
	}
)

// ProvideService builds a new AuthService.
func ProvideService(
	ctx context.Context,
	logger logging.Logger,
	cfg *Config,
	authenticator authentication.Authenticator,
	dataManager database.DataManager,
	householdMembershipManager types.HouseholdUserMembershipDataManager,
	sessionManager *scs.SessionManager,
	encoder encoding.ServerEncoderDecoder,
	tracerProvider tracing.TracerProvider,
	publisherProvider messagequeue.PublisherProvider,
	featureFlagManager featureflags.FeatureFlagManager,
	analyticsReporter analytics.EventReporter,
	routeParamManager routing.RouteParamManager,
) (types.AuthService, error) {
	hashKey := []byte(cfg.Cookies.HashKey)
	if len(hashKey) == 0 {
		hashKey = securecookie.GenerateRandomKey(cookieSecretSize)
	}

	dataChangesPublisher, publisherProviderErr := publisherProvider.ProvidePublisher(cfg.DataChangesTopicName)
	if publisherProviderErr != nil {
		return nil, fmt.Errorf("setting up auth service data changes publisher: %w", publisherProviderErr)
	}

	tracer := tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName))

	decryptedJWTSigningKey, err := base64.URLEncoding.DecodeString(cfg.JWTSigningKey)
	if err != nil {
		return nil, fmt.Errorf("decoding Token signing key: %w", err)
	}

	svc := &service{
		logger:                     logging.EnsureLogger(logger).WithName(serviceName),
		encoderDecoder:             encoder,
		config:                     cfg,
		userDataManager:            dataManager,
		householdMembershipManager: householdMembershipManager,
		authenticator:              authenticator,
		sessionManager:             sessionManager,
		sessionContextDataFetcher:  FetchContextFromRequest,
		cookieManager:              securecookie.New(hashKey, []byte(cfg.Cookies.BlockKey)),
		tracer:                     tracer,
		dataChangesPublisher:       dataChangesPublisher,
		featureFlagManager:         featureFlagManager,
		analyticsReporter:          analyticsReporter,
		authProviderFetcher:        routeParamManager.BuildRouteParamStringIDFetcher(AuthProviderParamKey),
		oauth2Server:               ProvideOAuth2ServerImplementation(ctx, logger, tracer, &cfg.OAuth2, dataManager),
	}

	signer, err := authentication.NewJWTSigner(logger, tracerProvider, cfg.JWTAudience, decryptedJWTSigningKey)
	if err != nil {
		return nil, fmt.Errorf("creating Token signer: %w", err)
	}
	svc.jwtSigner = signer

	if _, err = svc.cookieManager.Encode(cfg.Cookies.Name, "blah"); err != nil {
		logger.WithValue("cookie_signing_key_length", len(cfg.Cookies.BlockKey)).Error(err, "building test cookie")
		return nil, fmt.Errorf("building test cookie: %w", err)
	}

	useProvidersMutex.Lock()
	goth.UseProviders(
		google.New(svc.config.SSO.Google.ClientID, svc.config.SSO.Google.ClientID, svc.config.SSO.Google.CallbackURL),
	)
	useProvidersMutex.Unlock()

	return svc, nil
}
