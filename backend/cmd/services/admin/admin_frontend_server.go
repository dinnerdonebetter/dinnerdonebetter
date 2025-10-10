package main

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/dinnerdonebetter/backend/cmd/services/admin/components"
	"github.com/dinnerdonebetter/backend/internal/authentication/cookies"
	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/domain/auth"
	authsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/auth"
	mealplanninggrpc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/localdev"
	"github.com/dinnerdonebetter/backend/internal/platform/encoding"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/routing"
	phttp "github.com/dinnerdonebetter/backend/internal/platform/server/http"
	"github.com/dinnerdonebetter/backend/pkg/client"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	g "maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
	ghttp "maragu.dev/gomponents/http"
)

type ContextKey string

const (
	apiClientContextKey ContextKey = "api_client"
)

type AdminFrontendServer struct {
	tracer            tracing.Tracer
	logger            logging.Logger
	encoder           encoding.ServerEncoderDecoder
	cookieManager     cookies.Manager
	config            *config.AdminWebappConfig
	server            phttp.Server
	componentRenderer *components.ComponentRenderer
}

func NewAdminFrontendServer(
	ctx context.Context,
	apiClient client.Client,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	encoder encoding.ServerEncoderDecoder,
	cfg *config.AdminWebappConfig,
) (*AdminFrontendServer, error) {
	cookieMan, err := cookies.NewCookieManager(&cfg.Cookies, tracerProvider)
	if err != nil {
		return nil, err
	}

	metricsProvider, err := cfg.Observability.Metrics.ProvideMetricsProvider(ctx, logger)
	if err != nil {
		return nil, err
	}

	router, err := cfg.Routing.ProvideRouter(logger, tracerProvider, metricsProvider)
	if err != nil {
		return nil, err
	}

	server, err := phttp.ProvideHTTPServer(cfg.HTTPServer, logger, router, tracerProvider)
	if err != nil {
		return nil, err
	}

	s := &AdminFrontendServer{
		tracer:            tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		logger:            logging.EnsureLogger(logger).WithName(o11yName),
		componentRenderer: components.NewComponentRenderer(),
		cookieManager:     cookieMan,
		encoder:           encoder,
		config:            cfg,
		server:            server,
	}

	s.setupRoutes(router)

	return s, nil
}

func (s *AdminFrontendServer) Serve() {
	s.server.Serve()
}

func (s *AdminFrontendServer) setupRoutes(router routing.Router) {
	r := router.WithMiddleware(s.authMiddleware)

	r.Get("/", ghttp.Adapt(s.homeRoute))

	router.Get("/login", ghttp.Adapt(func(res http.ResponseWriter, req *http.Request) (g.Node, error) {
		return s.LoginPage(), nil
	}))

	router.Post("/login/submit", ghttp.Adapt(s.LoginSubmission))
}

func (s *AdminFrontendServer) authMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ctx, span := s.tracer.StartSpan(req.Context())
		defer span.End()

		logger := s.logger.WithRequest(req)

		cookieName := s.config.Cookies.CookieName
		cookie, err := req.Cookie(cookieName)
		if err != nil {
			logger.Error("no cookie found", err)
			http.Redirect(res, req, "/login", http.StatusFound)
			return
		} else if cookie == nil {
			logger.Debug("no cookie found")
			http.Redirect(res, req, "/login", http.StatusFound)
			return
		}

		var payload authPayload
		if err = s.cookieManager.Decode(ctx, cookieName, cookie.Value, &payload); err != nil {
			logger.Error("decoding cookie", err)
			http.Redirect(res, req, "/login", http.StatusFound)
			return
		}

		c, err := localdev.BuildInsecureOAuthedGRPCClient(
			ctx,
			s.config.APIServiceConnection.OAuth2APIClientID,
			s.config.APIServiceConnection.OAuth2APIClientSecret,
			s.config.APIServiceConnection.HTTPAPIServerURL,
			s.config.APIServiceConnection.GRPCAPIServerURL,
			payload.AccessToken,
		)

		handler.ServeHTTP(res, req.WithContext(context.WithValue(ctx, apiClientContextKey, c)))
	})
}

func (s *AdminFrontendServer) homeRoute(res http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)

	c, ok := req.Context().Value(apiClientContextKey).(client.Client)
	if !ok {
		return nil, errors.New("no api client found")
	}

	results, err := c.GetRecipes(ctx, &mealplanninggrpc.GetRecipesRequest{})
	if err != nil {
		return nil, err
	}

	logger.WithValue("count", len(results.Results)).Info("got recipes")

	return s.HomePage(), nil
}

func (s *AdminFrontendServer) HomePage() g.Node {
	return page("Home",
		ghtml.H1(g.Text("Home")),
	)
}

func (s *AdminFrontendServer) LoginPage() g.Node {
	return page("Login",
		s.componentRenderer.LoginForm(&components.LoginFormProps{}),
	)
}

func fetchErrorString(err error, key string) string {
	var validErr validation.Errors
	if errors.As(err, &validErr) {
		if validationErr := validErr[key]; validationErr != nil {
			var validationLibError validation.ErrorObject
			if errors.As(validationErr, &validationLibError) {
				return validationLibError.Error()
			}
		}
	}

	return ""
}

type authPayload struct {
	AccessToken string
}

func (s *AdminFrontendServer) LoginSubmission(res http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)

	var loginInput *auth.UserLoginInput
	if err := s.encoder.DecodeRequest(ctx, req, &loginInput); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request")
		return s.componentRenderer.LoginForm(&components.LoginFormProps{GeneralError: err.Error()}), nil
	}

	var (
		usernameError, passwordError, totpError string
	)

	if err := loginInput.ValidateWithContext(ctx); err != nil {
		usernameError = fetchErrorString(err, "username")
		passwordError = fetchErrorString(err, "password")
		totpError = fetchErrorString(err, "totpToken")
	}

	if usernameError != "" || passwordError != "" || totpError != "" {
		return s.componentRenderer.LoginForm(&components.LoginFormProps{
			UsernameError: usernameError,
			PasswordError: passwordError,
			TOTPError:     totpError,
		}), nil
	}

	unauthedClient, err := client.BuildUnauthenticatedGRPCClient(s.config.APIServiceConnection.GRPCAPIServerURL)
	if err != nil {
		return s.componentRenderer.LoginForm(&components.LoginFormProps{
			GeneralError: err.Error(),
		}), nil
	}

	tokenRes, err := unauthedClient.AdminLoginForToken(ctx, &authsvc.AdminLoginForTokenRequest{
		Input: &authsvc.UserLoginInput{
			Username:  loginInput.Username,
			Password:  loginInput.Password,
			TOTPToken: loginInput.TOTPToken,
		},
	})
	if err != nil {
		return s.componentRenderer.LoginForm(&components.LoginFormProps{
			GeneralError: err.Error(),
		}), nil
	}

	encodedCookie, err := s.cookieManager.Encode(ctx, s.config.Cookies.CookieName, &authPayload{AccessToken: tokenRes.Result.AccessToken})
	if err != nil {
		return s.componentRenderer.LoginForm(&components.LoginFormProps{
			GeneralError: err.Error(),
		}), nil
	}

	http.SetCookie(res, s.buildCookie(ctx, encodedCookie))

	return s.HomePage(), nil
}

// buildCookie provides a consistent way of constructing an HTTP cookie.
func (s *AdminFrontendServer) buildCookie(ctx context.Context, value string) *http.Cookie {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	expiry := time.Now().Add(s.config.Cookies.Lifetime)

	// https://www.calhoun.io/securing-cookies-in-go/
	cookie := &http.Cookie{
		Name:     s.config.Cookies.CookieName,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   s.config.Cookies.SecureOnly,
		// Domain:   s.config.Cookies.Domain,
		Expires:  expiry,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(time.Until(expiry).Seconds()),
	}

	return cookie
}
