package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dinnerdonebetter/backend/cmd/services/admin/components"
	"github.com/dinnerdonebetter/backend/cmd/services/admin/design"
	"github.com/dinnerdonebetter/backend/internal/authentication/cookies"
	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/domain/auth"
	authsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/auth"
	identitysvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/identity"
	"github.com/dinnerdonebetter/backend/internal/localdev"
	"github.com/dinnerdonebetter/backend/internal/platform/encoding"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/routing"
	phttp "github.com/dinnerdonebetter/backend/internal/platform/server/http"
	"github.com/dinnerdonebetter/backend/pkg/client"
	"google.golang.org/protobuf/types/known/timestamppb"

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

	r.Get("/users", ghttp.Adapt(s.UsersList))
	r.Get("/api/users/search", ghttp.Adapt(s.UsersSearch))

	router.Get("/login", ghttp.Adapt(s.LoginPage))
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
		if err != nil {
			logger.Error("building client", err)
			http.Redirect(res, req, "/login", http.StatusFound)
			return
		}

		handler.ServeHTTP(res, req.WithContext(context.WithValue(ctx, apiClientContextKey, c)))
	})
}

func fetchClientFromContext(ctx context.Context) (client.Client, error) {
	c, ok := ctx.Value(apiClientContextKey).(client.Client)
	if !ok {
		return nil, errors.New("no api client found")
	}

	return c, nil
}

func (s *AdminFrontendServer) homeRoute(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	_, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	return s.HomePage(""), nil
}

func (s *AdminFrontendServer) HomePage(title string) g.Node {
	if title == "" {
		title = "Dashboard"
	}

	return page(title,
		components.ContentContainer(&components.ContentContainerProps{
			Title:    title,
			Subtitle: "Welcome to the admin dashboard",
			Palette:  &design.StandardPalette,
		},
			components.Card(&design.StandardPalette,
				ghtml.H2(
					ghtml.Class(fmt.Sprintf("text-lg font-medium %s mb-4", design.TextColor(design.StandardPalette.Primary))),
					g.Text("Quick Stats"),
				),
				ghtml.Div(
					ghtml.Class("grid grid-cols-1 md:grid-cols-3 gap-4"),
					statCard("Total Users", "1,234", &design.StandardPalette),
					statCard("Active Sessions", "89", &design.StandardPalette),
					statCard("System Status", "Healthy", &design.StandardPalette),
				),
			),
		),
	)
}

func statCard(title, value string, palette *design.Palette) g.Node {
	return ghtml.Div(
		ghtml.Class(fmt.Sprintf("p-4 %s rounded-lg border %s",
			design.Background(design.Color{Value: "gray-50"}),
			design.BorderColor(palette.Background),
		)),
		ghtml.Div(
			ghtml.Class(fmt.Sprintf("text-sm font-medium %s", design.TextColor(palette.Text))),
			g.Text(title),
		),
		ghtml.Div(
			ghtml.Class(fmt.Sprintf("mt-1 text-2xl font-bold %s", design.TextColor(palette.Primary))),
			g.Text(value),
		),
	)
}

func (s *AdminFrontendServer) LoginPage(_ http.ResponseWriter, _ *http.Request) (g.Node, error) {
	return page("Login",
		s.componentRenderer.LoginForm(&components.LoginFormProps{}),
	), nil
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

func renderTimestamp(value any) g.Node {
	if value == nil {
		return g.Text("-")
	}

	switch v := value.(type) {
	case *timestamppb.Timestamp:
		if v == nil {
			return g.Text("-")
		}
		return g.Text(v.AsTime().Format("2006-01-02 15:04:05"))
	case timestamppb.Timestamp:
		return g.Text(v.AsTime().Format("2006-01-02 15:04:05"))
	default:
		return g.Text(fmt.Sprintf("%v", v))
	}
}

func (s *AdminFrontendServer) UsersList(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return page("Users", s.renderUsersError("Error: No API client available")), nil
	}

	usersRes, err := c.GetUsers(ctx, &identitysvc.GetUsersRequest{})
	if err != nil {
		return page("Users", s.renderUsersError(fmt.Sprintf("Error loading users: %v", err))), nil
	}

	// Use the new integrated TablePage component
	tablePageResult, err := components.TablePage(&components.TablePageProps[*identitysvc.User]{
		Title:             "Users",
		BaseSubtitle:      "Manage user accounts",
		Palette:           &design.StandardPalette,
		ShowSearch:        true,
		SearchPlaceholder: "Search users...",
		HTMXSearchTarget:  "/api/users/search",
		Data:              usersRes.Result,
		Actions:           []g.Node{},
		TableOptions: &components.TableOptions[*identitysvc.User]{
			TableID: "users-table",
			Palette: &design.StandardPalette,
			Fields: []string{
				"ID",
				"Username",
				"FirstName",
				"LastName",
				"EmailAddress",
				"ServiceRole",
				"AccountStatus",
				"AccountStatusExplanation",
				"Birthday",
				"PasswordLastChangedAt",
				"LastAcceptedTermsOfService",
				"LastAcceptedPrivacyPolicy",
				"TwoFactorSecretVerifiedAt",
				"EmailAddressVerifiedAt",
				"CreatedAt",
				"LastUpdatedAt",
				"ArchivedAt",
			},
			FieldReplacements: map[string]string{
				"EmailAddressVerifiedAt": "Email Verified At",
			},
			FieldRenderers: map[string]components.FieldRenderer{
				"Birthday":                  renderTimestamp,
				"CreatedAt":                 renderTimestamp,
				"TwoFactorSecretVerifiedAt": renderTimestamp,
				"LastUpdatedAt":             renderTimestamp,
				"ArchivedAt":                renderTimestamp,
			},
		},
		EmptyStateTitle:       "No users found",
		EmptyStateDescription: "Get started by creating your first user account.",
		EmptyStateActions:     []g.Node{},
		SubtitleGenerator: func(metadata components.TablePageMetadata) string {
			if metadata.EmptyState {
				return "Manage user accounts"
			}
			return fmt.Sprintf("Manage %d user accounts", metadata.TotalCount)
		},
	})
	if err != nil {
		return page("Users", s.renderUsersError(fmt.Sprintf("Error creating table: %v", err))), nil
	}

	return page("Users", tablePageResult.Node), nil
}

func (s *AdminFrontendServer) UsersSearch(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return g.El("div",
			g.Attr("class", "overflow-x-auto"),
			ghtml.P(
				ghtml.Class(fmt.Sprintf("text-center py-8 %s", design.TextColor(design.StandardPalette.Warning))),
				g.Text("Error: No API client available"),
			),
		), nil
	}

	// Get search query from request
	searchQuery := req.URL.Query().Get("search")

	usersRes, err := c.GetUsers(ctx, &identitysvc.GetUsersRequest{})
	if err != nil {
		return g.El("div",
			g.Attr("class", "overflow-x-auto"),
			ghtml.P(
				ghtml.Class(fmt.Sprintf("text-center py-8 %s", design.TextColor(design.StandardPalette.Warning))),
				g.Text(fmt.Sprintf("Error loading users: %v", err)),
			),
		), nil
	}

	// Filter users based on search query
	var filteredUsers []*identitysvc.User
	if searchQuery == "" {
		// No search query, return all users
		filteredUsers = usersRes.Result
	} else {
		// Filter users by search query (case insensitive)
		searchQueryLower := strings.ToLower(searchQuery)
		for _, user := range usersRes.Result {
			if strings.Contains(strings.ToLower(user.Username), searchQueryLower) ||
				strings.Contains(strings.ToLower(user.FirstName), searchQueryLower) ||
				strings.Contains(strings.ToLower(user.LastName), searchQueryLower) ||
				strings.Contains(strings.ToLower(user.EmailAddress), searchQueryLower) {
				filteredUsers = append(filteredUsers, user)
			}
		}
	}

	// Generate just the table (not the full page)
	if len(filteredUsers) == 0 {
		return g.El("div",
			g.Attr("class", "overflow-x-auto"),
			components.EmptyState(
				"No users found",
				fmt.Sprintf("No users match the search term '%s'.", searchQuery),
				&design.StandardPalette,
				[]g.Node{
					components.ActionButton("Add User", "/users/new", &design.StandardPalette, true),
				},
			),
		), nil
	}

	table, err := components.Table(filteredUsers, &components.TableOptions[*identitysvc.User]{
		TableID: "users-table",
		Palette: &design.StandardPalette,
		Fields: []string{
			"ID",
			"Username",
			"FirstName",
			"LastName",
			"EmailAddress",
			"ServiceRole",
			"AccountStatus",
			"AccountStatusExplanation",
			"Birthday",
			"PasswordLastChangedAt",
			"LastAcceptedTermsOfService",
			"LastAcceptedPrivacyPolicy",
			"TwoFactorSecretVerifiedAt",
			"EmailAddressVerifiedAt",
			"CreatedAt",
			"LastUpdatedAt",
			"ArchivedAt",
		},
		FieldReplacements: map[string]string{
			"EmailAddressVerifiedAt": "Email Verified At",
		},
		FieldRenderers: map[string]components.FieldRenderer{
			"Birthday":                  renderTimestamp,
			"CreatedAt":                 renderTimestamp,
			"TwoFactorSecretVerifiedAt": renderTimestamp,
			"LastUpdatedAt":             renderTimestamp,
			"ArchivedAt":                renderTimestamp,
		},
	})
	if err != nil {
		return g.El("div",
			g.Attr("class", "overflow-x-auto"),
			ghtml.P(
				ghtml.Class(fmt.Sprintf("text-center py-8 %s", design.TextColor(design.StandardPalette.Warning))),
				g.Text(fmt.Sprintf("Error creating table: %v", err)),
			),
		), nil
	}

	// Wrap table in the same scrollable container structure for consistency
	return g.El("div",
		g.Attr("class", "overflow-x-auto"),
		table,
	), nil
}

// renderUsersError creates a consistent error display for the users page
func (s *AdminFrontendServer) renderUsersError(errorMsg string) g.Node {
	return components.ContentContainer(&components.ContentContainerProps{
		Title:    "Users",
		Subtitle: "Manage user accounts",
		Palette:  &design.StandardPalette,
	},
		components.Card(&design.StandardPalette,
			ghtml.P(
				ghtml.Class(fmt.Sprintf("text-center py-8 %s", design.TextColor(design.StandardPalette.Warning))),
				g.Text(errorMsg),
			),
		),
	)
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

	return s.HomePage(""), nil
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
