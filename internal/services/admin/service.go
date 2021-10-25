package admin

import (
	"net/http"

	"github.com/alexedwards/scs/v2"

	"gitlab.com/prixfixe/prixfixe/internal/authentication"
	"gitlab.com/prixfixe/prixfixe/internal/encoding"
	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/internal/routing"
	authservice "gitlab.com/prixfixe/prixfixe/internal/services/authentication"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

const (
	serviceName = "auth_service"
)

type (
	// service handles passwords service-wide.
	service struct {
		config                    *authservice.Config
		logger                    logging.Logger
		authenticator             authentication.Authenticator
		userDB                    types.AdminUserDataManager
		encoderDecoder            encoding.ServerEncoderDecoder
		sessionManager            *scs.SessionManager
		sessionContextDataFetcher func(*http.Request) (*types.SessionContextData, error)
		userIDFetcher             func(*http.Request) string
		tracer                    tracing.Tracer
	}
)

// ProvideService builds a new AuthService.
func ProvideService(
	logger logging.Logger,
	cfg *authservice.Config,
	authenticator authentication.Authenticator,
	userDataManager types.AdminUserDataManager,
	sessionManager *scs.SessionManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
) types.AdminService {
	svc := &service{
		logger:                    logging.EnsureLogger(logger).WithName(serviceName),
		encoderDecoder:            encoder,
		config:                    cfg,
		userDB:                    userDataManager,
		authenticator:             authenticator,
		sessionManager:            sessionManager,
		sessionContextDataFetcher: authservice.FetchContextFromRequest,
		userIDFetcher:             routeParamManager.BuildRouteParamStringIDFetcher(UserIDURIParamKey),
		tracer:                    tracing.NewTracer(serviceName),
	}
	svc.sessionManager.Lifetime = cfg.Cookies.Lifetime

	return svc
}
