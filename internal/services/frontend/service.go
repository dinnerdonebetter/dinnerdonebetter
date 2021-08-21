package frontend

import (
	"context"
	"net/http"

	"gitlab.com/prixfixe/prixfixe/internal/authorization"
	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	panicking "gitlab.com/prixfixe/prixfixe/internal/panicking"
	routing "gitlab.com/prixfixe/prixfixe/internal/routing"
	authservice "gitlab.com/prixfixe/prixfixe/internal/services/authentication"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

const (
	serviceName string = "frontend_service"
)

type (
	// AuthService is a subset of the larger types.AuthService interface.
	AuthService interface {
		UserAttributionMiddleware(next http.Handler) http.Handler
		PermissionFilterMiddleware(permissions ...authorization.Permission) func(next http.Handler) http.Handler
		ServiceAdminMiddleware(next http.Handler) http.Handler

		AuthenticateUser(ctx context.Context, loginData *types.UserLoginInput) (*types.User, *http.Cookie, error)
		LogoutUser(ctx context.Context, sessionCtxData *types.SessionContextData, req *http.Request, res http.ResponseWriter) error
	}

	// UsersService is a subset of the larger types.UsersService interface.
	UsersService interface {
		RegisterUser(ctx context.Context, registrationInput *types.UserRegistrationInput) (*types.UserCreationResponse, error)
		VerifyUserTwoFactorSecret(ctx context.Context, input *types.TOTPSecretVerificationInput) error
	}

	// Service serves HTML.
	Service interface {
		SetupRoutes(router routing.Router)
	}

	service struct {
		logger                    logging.Logger
		tracer                    tracing.Tracer
		panicker                  panicking.Panicker
		authService               AuthService
		config                    *Config
		sessionContextDataFetcher func(*http.Request) (*types.SessionContextData, error)
	}
)

// ProvideService builds a new Service.
func ProvideService(
	cfg *Config,
	logger logging.Logger,
	authService AuthService,
) Service {
	svc := &service{
		config:                    cfg,
		logger:                    logging.EnsureLogger(logger).WithName(serviceName),
		tracer:                    tracing.NewTracer(serviceName),
		panicker:                  panicking.NewProductionPanicker(),
		sessionContextDataFetcher: authservice.FetchContextFromRequest,
		authService:               authService,
	}

	if cfg.Debug {
		svc.logger.SetLevel(logging.DebugLevel)
	}

	return svc
}
