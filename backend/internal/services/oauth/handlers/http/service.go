package http

import (
	"fmt"
	"github.com/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/backend/internal/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/authentication/tokens"
	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	"github.com/dinnerdonebetter/backend/internal/domain/oauth"
	"github.com/dinnerdonebetter/backend/internal/platform/encoding"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"net/http"

	"github.com/go-oauth2/oauth2/v4/server"
)

const (
	serviceName                = "auth_service"
	AuthProviderParamKey       = "auth_provider"
	rejectedRequestCounterName = "auth_service.rejected_requests"
)

type (
	// service handles passwords service-wide.
	service struct {
		logger                    logging.Logger
		authenticator             authentication.Authenticator
		userDataManager           identity.UserDataManager
		encoderDecoder            encoding.ServerEncoderDecoder
		sessionContextDataFetcher func(*http.Request) (*sessions.ContextData, error)
		tracer                    tracing.Tracer
		oauth2Server              *server.Server
		tokenIssuer               tokens.Issuer
		rejectedRequestCounter    metrics.Int64Counter
	}
)

// ProvideService builds a new AuthDataService.
func ProvideService(
	logger logging.Logger,
	cfg *Config,
	authenticator authentication.Authenticator,
	oauthRepo oauth.Repository,
	identityRepo identity.Repository,
	encoder encoding.ServerEncoderDecoder,
	tracerProvider tracing.TracerProvider,
	metricsProvider metrics.Provider,
) (oauth.OAuth2Service, error) {
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
		logger:                    logging.EnsureLogger(logger).WithName(serviceName),
		encoderDecoder:            encoder,
		userDataManager:           identityRepo,
		authenticator:             authenticator,
		sessionContextDataFetcher: sessions.FetchContextDataFromRequest,
		tracer:                    tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
		tokenIssuer:               signer,
		rejectedRequestCounter:    rejectedRequestCounter,
		oauth2Server:              ProvideOAuth2ServerImplementation(logger, tracer, cfg, oauthRepo, identityRepo, authenticator, signer),
	}

	return svc, nil
}
