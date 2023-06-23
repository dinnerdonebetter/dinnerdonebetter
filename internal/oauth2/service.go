package oauth2

import (
	"context"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
)

const (
	serviceName = "oauth2_service"
)

type (
	// Service handles oauth2.
	Service struct {
		logger         logging.Logger
		tracer         tracing.Tracer
		encoderDecoder encoding.ServerEncoderDecoder
		oauth2Server   *server.Server
	}
)

// ProvideOAuth2Service builds a new oauth2 Service.
func ProvideOAuth2Service(
	_ context.Context,
	logger logging.Logger,
	cfg *Config,
	dataManager database.DataManager,
	encoder encoding.ServerEncoderDecoder,
	tracerProvider tracing.TracerProvider,
) (*Service, error) {
	manager := manage.NewManager()
	manager.MapAuthorizeGenerate(generates.NewAuthorizeGenerate())
	manager.MapAccessGenerate(generates.NewAccessGenerate())

	// token memory store
	manager.MapTokenStorage(&oauth2TokenStoreImpl{
		tracer:      tracing.NewTracer(tracerProvider.Tracer("oauth2_token_store")),
		logger:      logging.EnsureLogger(logger),
		dataManager: dataManager,
	})

	// client memory store
	manager.MapClientStorage(newOAuth2ClientStore(cfg.Domain, logger, tracerProvider, dataManager))

	oauth2ServerConfig := &server.Config{
		TokenType: "Bearer",
		AllowedResponseTypes: []oauth2.ResponseType{
			oauth2.Token,
		},
		AllowedGrantTypes: []oauth2.GrantType{
			// oauth2.ClientCredentials,
			oauth2.PasswordCredentials,
			oauth2.Refreshing,
		},
		AllowedCodeChallengeMethods: []oauth2.CodeChallengeMethod{
			oauth2.CodeChallengeS256,
		},
	}

	oauth2Server := server.NewServer(oauth2ServerConfig, manager)

	oauth2Server.UserAuthorizationHandler = func(res http.ResponseWriter, req *http.Request) (userID string, err error) {
		return "", errors.ErrAccessDenied
	}

	oauth2Server.PasswordAuthorizationHandler = func(ctx context.Context, clientID, username, password string) (userID string, err error) {
		user, err := dataManager.GetUserByUsername(ctx, username)
		if err != nil {
			return "", errors.New("invalid username or password")
		}

		// TODO: validate password here, duh

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

	s := &Service{
		logger:         logging.EnsureLogger(logger).WithName(serviceName),
		encoderDecoder: encoder,
		oauth2Server:   oauth2Server,
		tracer:         tracing.NewTracer(tracerProvider.Tracer(serviceName)),
	}

	return s, nil
}
