package authentication

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/backend/internal/authentication/tokens"
	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
)

func ProvideOAuth2ServerImplementation(
	logger logging.Logger,
	tracer tracing.Tracer,
	cfg *OAuth2Config,
	dataManager database.DataManager,
	authenticator authentication.Authenticator,
	tokenIssuer tokens.Issuer,
) *server.Server {
	manager := manage.NewManager()

	// we don't care at the moment
	manager.SetValidateURIHandler(func(_, _ string) error {
		return nil
	})
	manager.MapAuthorizeGenerate(generates.NewAuthorizeGenerate())
	manager.MapAccessGenerate(generates.NewAccessGenerate())
	manager.MapClientStorage(newOAuth2ClientStore(cfg.Domain, logger, tracer, dataManager))
	manager.MapTokenStorage(&oauth2TokenStoreImpl{
		tracer:      tracer,
		logger:      logging.EnsureLogger(logger),
		dataManager: dataManager,
	})

	oauth2ServerConfig := &server.Config{
		TokenType: "Bearer",
		AllowedResponseTypes: []oauth2.ResponseType{
			oauth2.Code,
		},
		AllowedGrantTypes: []oauth2.GrantType{
			oauth2.AuthorizationCode,
			oauth2.Refreshing,
		},
		AllowedCodeChallengeMethods: []oauth2.CodeChallengeMethod{
			oauth2.CodeChallengePlain,
		},
	}

	oauth2Server := server.NewServer(oauth2ServerConfig, manager)

	oauth2Server.AuthorizeScopeHandler = AuthorizeScopeHandler(logger)
	oauth2Server.AccessTokenExpHandler = AccessTokenExpHandler(logger)
	oauth2Server.ClientScopeHandler = ClientScopeHandler(logger)
	oauth2Server.UserAuthorizationHandler = buildUserAuthorizationHandler(tracer, logger, tokenIssuer)
	oauth2Server.PasswordAuthorizationHandler = buildPasswordAuthorizationHandler(logger, authenticator, dataManager)
	// this allows GET requests to retrieve tokens
	oauth2Server.SetAllowGetAccessRequest(true)
	oauth2Server.ClientInfoHandler = buildClientInfoHandler()
	oauth2Server.InternalErrorHandler = buildInternalErrorHandler(logger)
	oauth2Server.ResponseErrorHandler = buildOAuth2ErrorHandler(logger)

	return oauth2Server
}

func buildOAuth2ErrorHandler(logger logging.Logger) func(*errors.Response) {
	return func(res *errors.Response) {
		observability.AcknowledgeError(res.Error, logger, nil, "oauth2 response error")
	}
}

func buildInternalErrorHandler(logger logging.Logger) func(error) *errors.Response {
	return func(err error) *errors.Response {
		observability.AcknowledgeError(err, logger, nil, "internal oauth2 error")
		return &errors.Response{
			Error:       err,
			ErrorCode:   -1,
			Description: err.Error(),
			URI:         "",
			StatusCode:  http.StatusInternalServerError,
			Header:      nil,
		}
	}
}

// this determines how we identify clients from HTTP requests.
func buildClientInfoHandler() func(*http.Request) (string, string, error) {
	return func(req *http.Request) (string, string, error) {
		clientID, clientSecret := req.Form.Get("client_id"), req.Form.Get("client_secret")
		if clientID == "" || clientSecret == "" {
			username, password, ok := req.BasicAuth()
			if !ok {
				return "", "", errors.ErrInvalidClient
			}

			return username, password, nil
		}

		return clientID, clientSecret, nil
	}
}

func buildPasswordAuthorizationHandler(logger logging.Logger, authenticator authentication.Authenticator, dataManager database.DataManager) func(context.Context, string, string, string) (string, error) {
	return func(ctx context.Context, clientID, username, password string) (userID string, err error) {
		l := logger.WithValue(keys.OAuth2ClientIDKey, clientID).WithValue(keys.UsernameKey, username)
		l.Info("PasswordAuthorizationHandler invoked")

		user, err := dataManager.GetUserByUsername(ctx, username)
		if err != nil {
			return "", errors.New("invalid username or password")
		}

		valid, err := authenticator.CredentialsAreValid(
			ctx,
			user.HashedPassword,
			password,
			"",
			"",
		)
		if err != nil {
			l.Error("validating credentials", err)
			return "", errors.New("invalid username or password")
		}

		if !valid {
			l.Info("invalid credentials")
			return "", errors.New("invalid username or password")
		}

		return user.ID, nil
	}
}

func buildUserAuthorizationHandler(tracer tracing.Tracer, logger logging.Logger, tokenIssuer tokens.Issuer) func(http.ResponseWriter, *http.Request) (string, error) {
	return func(res http.ResponseWriter, req *http.Request) (userID string, err error) {
		ctx, span := tracer.StartCustomSpan(req.Context(), "oauth2_server.UserAuthorizationHandler")
		defer span.End()

		l := logger.WithRequest(req)
		l.Info("UserAuthorizationHandler invoked")

		rawToken := req.Header.Get("Authorization")
		token := strings.TrimPrefix(rawToken, "Bearer ")

		subject, err := tokenIssuer.ParseUserIDFromToken(ctx, token)
		if err != nil {
			l.Error("parsing token in UserAuthorizationHandler", err)
			return "", errors.ErrAccessDenied
		}

		return subject, nil
	}
}

func AuthorizeScopeHandler(_ logging.Logger) func(http.ResponseWriter, *http.Request) (string, error) {
	return func(_ http.ResponseWriter, req *http.Request) (scope string, err error) {
		return req.URL.Query().Get("scope"), nil
	}
}

func AccessTokenExpHandler(_ logging.Logger) func(http.ResponseWriter, *http.Request) (time.Duration, error) {
	return func(_ http.ResponseWriter, _ *http.Request) (time.Duration, error) {
		return 24 * time.Hour, nil
	}
}

func ClientScopeHandler(_ logging.Logger) func(_ *oauth2.TokenGenerateRequest) (allowed bool, err error) {
	return func(_ *oauth2.TokenGenerateRequest) (allowed bool, err error) {
		return true, nil
	}
}
