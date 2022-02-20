package authentication

import (
	"fmt"
	"net/http"

	"github.com/prixfixeco/api_server/internal/messagequeue"

	"github.com/alexedwards/scs/v2"
	"github.com/gorilla/securecookie"

	"github.com/prixfixeco/api_server/internal/authentication"
	"github.com/prixfixeco/api_server/internal/encoding"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	serviceName           = "auth_service"
	userIDContextKey      = string(types.UserIDContextKey)
	householdIDContextKey = string(types.HouseholdIDContextKey)
	cookieErrorLogName    = "_COOKIE_CONSTRUCTION_ERROR_"
	cookieSecretSize      = 64
)

type (
	// cookieEncoderDecoder is a stand-in interface for gorilla/securecookie.
	cookieEncoderDecoder interface {
		Encode(name string, value interface{}) (string, error)
		Decode(name, value string, dst interface{}) error
	}

	// service handles passwords service-wide.
	service struct {
		config                     *Config
		logger                     logging.Logger
		authenticator              authentication.Authenticator
		userDataManager            types.UserDataManager
		apiClientManager           types.APIClientDataManager
		householdMembershipManager types.HouseholdUserMembershipDataManager
		encoderDecoder             encoding.ServerEncoderDecoder
		cookieManager              cookieEncoderDecoder
		sessionManager             sessionManager
		sessionContextDataFetcher  func(*http.Request) (*types.SessionContextData, error)
		tracer                     tracing.Tracer
		dataChangesPublisher       messagequeue.Publisher
	}
)

// ProvideService builds a new AuthService.
func ProvideService(
	logger logging.Logger,
	cfg *Config,
	authenticator authentication.Authenticator,
	userDataManager types.UserDataManager,
	apiClientsService types.APIClientDataManager,
	householdMembershipManager types.HouseholdUserMembershipDataManager,
	sessionManager *scs.SessionManager,
	encoder encoding.ServerEncoderDecoder,
	tracerProvider tracing.TracerProvider,
	publisherProvider messagequeue.PublisherProvider,
) (types.AuthService, error) {
	hashKey := []byte(cfg.Cookies.HashKey)
	if len(hashKey) == 0 {
		hashKey = securecookie.GenerateRandomKey(cookieSecretSize)
	}

	dataChangesPublisher, err := publisherProvider.ProviderPublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up auth service data changes publisher: %w", err)
	}

	svc := &service{
		logger:                     logging.EnsureLogger(logger).WithName(serviceName),
		encoderDecoder:             encoder,
		config:                     cfg,
		userDataManager:            userDataManager,
		apiClientManager:           apiClientsService,
		householdMembershipManager: householdMembershipManager,
		authenticator:              authenticator,
		sessionManager:             sessionManager,
		sessionContextDataFetcher:  FetchContextFromRequest,
		cookieManager:              securecookie.New(hashKey, []byte(cfg.Cookies.BlockKey)),
		tracer:                     tracing.NewTracer(tracerProvider.Tracer(serviceName)),
		dataChangesPublisher:       dataChangesPublisher,
	}

	if _, err := svc.cookieManager.Encode(cfg.Cookies.Name, "blah"); err != nil {
		logger.WithValue("cookie_signing_key_length", len(cfg.Cookies.BlockKey)).Error(err, "building test cookie")
		return nil, fmt.Errorf("building test cookie: %w", err)
	}

	return svc, nil
}
