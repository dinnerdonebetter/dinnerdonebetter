package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/dinnerdonebetter/backend/cmd/services/admin/components"
	"github.com/dinnerdonebetter/backend/internal/authentication/cookies"
	"github.com/dinnerdonebetter/backend/internal/domain/auth"
	"github.com/dinnerdonebetter/backend/internal/platform/encoding"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/client"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	g "maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
	ghttp "maragu.dev/gomponents/http"
)

const (
	grpcServerAddress = ""
	clientID          = ""
	clientSecret      = ""
	authToken         = ""
)

type AdminFrontendServer struct {
	tracer            tracing.Tracer
	logger            logging.Logger
	mux               *http.ServeMux
	encoder           encoding.ServerEncoderDecoder
	componentRenderer *components.ComponentRenderer
	apiClient         client.Client
	cookieManager     cookies.Manager
}

func NewAdminFrontendServer(
	ctx context.Context,
	tracerProvider tracing.TracerProvider,
	logger logging.Logger,
	encoder encoding.ServerEncoderDecoder,
	mux *http.ServeMux,
) (*AdminFrontendServer, error) {
	cookieMan, err := cookies.NewCookieManager(&cookies.Config{}, tracerProvider)
	if err != nil {
		return nil, err
	}

	opt, err := client.WithOAuth2Credentials(
		ctx,
		grpcServerAddress,
		clientID,
		clientSecret,
		authToken,
	)
	if err != nil {
		return nil, err
	}

	apiClient, err := client.BuildClient(grpcServerAddress, opt)
	if err != nil {
		return nil, err
	}

	s := &AdminFrontendServer{
		apiClient:         apiClient,
		mux:               mux,
		tracer:            tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		logger:            logging.EnsureLogger(logger).WithName(o11yName),
		componentRenderer: components.NewComponentRenderer(),
		cookieManager:     cookieMan,
		encoder:           encoder,
	}

	s.setupRoutes(mux)

	return s, nil
}

func (s *AdminFrontendServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	s.mux.ServeHTTP(res, req)
}

func (s *AdminFrontendServer) setupRoutes(mux *http.ServeMux) {
	mux.Handle("GET /", ghttp.Adapt(func(res http.ResponseWriter, req *http.Request) (g.Node, error) {
		return s.HomePage(), nil
	}))

	mux.Handle("GET /login", ghttp.Adapt(func(res http.ResponseWriter, req *http.Request) (g.Node, error) {
		return s.LoginPage(), nil
	}))

	mux.Handle("POST /login/submit", ghttp.Adapt(s.LoginSubmission))
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

	return s.HomePage(), nil
}
