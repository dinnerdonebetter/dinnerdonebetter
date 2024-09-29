package authentication

import (
	"context"
	"github.com/dinnerdonebetter/backend/internal/authentication"
	"net/http"
	"strings"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
)

func ProvideOAuth2ServerImplementation(
	_ context.Context,
	logger logging.Logger,
	tracer tracing.Tracer,
	cfg *OAuth2Config,
	dataManager database.DataManager,
	authenticator authentication.Authenticator,
	jwtSigner authentication.JWTSigner,
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
			oauth2.CodeChallengeS256,
		},
	}

	oauth2Server := server.NewServer(oauth2ServerConfig, manager)

	oauth2Server.UserAuthorizationHandler = func(res http.ResponseWriter, req *http.Request) (userID string, err error) {
		ctx, span := tracer.StartCustomSpan(req.Context(), "oauth2_server.UserAuthorizationHandler")
		defer span.End()

		rawToken := req.Header.Get("Authorization")
		token := strings.TrimPrefix(rawToken, "Bearer ")

		parsedToken, err := jwtSigner.ParseJWT(ctx, token)
		if err != nil {
			return "", errors.ErrAccessDenied
		}

		subject, err := parsedToken.Claims.GetSubject()
		if err != nil {
			return "", errors.ErrAccessDenied
		}

		return subject, nil
	}

	oauth2Server.PasswordAuthorizationHandler = func(ctx context.Context, clientID, username, password string) (userID string, err error) {
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
			return "", errors.New("invalid username or password")
		}

		if !valid {
			return "", errors.New("invalid username or password")
		}

		return user.ID, nil
	}

	// this allows GET requests to retrieve tokens
	oauth2Server.SetAllowGetAccessRequest(true)

	// this determines how we identify clients from HTTP requests
	oauth2Server.SetClientInfoHandler(func(req *http.Request) (string, string, error) {
		clientID, clientSecret := req.Form.Get("client_id"), req.Form.Get("client_secret")
		if clientID == "" || clientSecret == "" {
			username, password, ok := req.BasicAuth()
			if !ok {
				return "", "", errors.ErrInvalidClient
			}

			return username, password, nil
		}

		return clientID, clientSecret, nil
	})

	oauth2Server.SetInternalErrorHandler(func(err error) *errors.Response {
		observability.AcknowledgeError(err, logger, nil, "internal oauth2 error")
		return &errors.Response{
			Error:       err,
			ErrorCode:   -1,
			Description: err.Error(),
			URI:         "",
			StatusCode:  http.StatusInternalServerError,
			Header:      nil,
		}
	})

	oauth2Server.SetResponseErrorHandler(func(res *errors.Response) {
		observability.AcknowledgeError(res.Error, logger, nil, "oauth2 response error")
	})

	return oauth2Server
}
