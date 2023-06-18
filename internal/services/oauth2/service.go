package oauth2

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/messagequeue"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/routing"
	authservice "github.com/dinnerdonebetter/backend/internal/services/authentication"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
)

const (
	serviceName = "oauth2_service"
)

type (
	// Service handles oauth2.
	Service struct {
		logger                    logging.Logger
		tracer                    tracing.Tracer
		encoderDecoder            encoding.ServerEncoderDecoder
		dataChangesPublisher      messagequeue.Publisher
		oauth2Server              *server.Server
		sessionContextDataFetcher func(*http.Request) (*types.SessionContextData, error)
	}
)

// ProvideOAuth2Service builds a new oauth2 Service.
func ProvideOAuth2Service(
	_ context.Context,
	logger logging.Logger,
	cfg *Config,
	encoder encoding.ServerEncoderDecoder,
	_ routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
) (*Service, error) {
	dataChangesPublisher, err := publisherProvider.ProvidePublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up data changes publisher: %w", err)
	}

	manager := manage.NewManager()

	// token memory store
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	// client memory store
	manager.MapClientStorage(&oauth2ClientStoreImpl{domain: cfg.Domain})

	oauth2ServerConfig := &server.Config{
		TokenType: "Bearer",
		AllowedResponseTypes: []oauth2.ResponseType{
			oauth2.Token,
		},
		AllowedGrantTypes: []oauth2.GrantType{
			oauth2.PasswordCredentials,
			oauth2.Refreshing,
		},
		AllowedCodeChallengeMethods: []oauth2.CodeChallengeMethod{
			oauth2.CodeChallengeS256,
		},
	}

	oauth2Server := server.NewServer(oauth2ServerConfig, manager)

	// this allows GET requests to retrieve tokens
	oauth2Server.SetAllowGetAccessRequest(true)

	// this determines how we identify clients from HTTP requests
	oauth2Server.SetClientInfoHandler(server.ClientFormHandler)

	oauth2Server.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		observability.AcknowledgeError(err, logger, nil, "internal oauth2 error")
		return nil
	})

	oauth2Server.SetResponseErrorHandler(func(res *errors.Response) {
		observability.AcknowledgeError(err, logger, nil, "oauth2 response error")
		return
	})

	s := &Service{
		logger:                    logging.EnsureLogger(logger).WithName(serviceName),
		encoderDecoder:            encoder,
		dataChangesPublisher:      dataChangesPublisher,
		oauth2Server:              oauth2Server,
		sessionContextDataFetcher: authservice.FetchContextFromRequest,
		tracer:                    tracing.NewTracer(tracerProvider.Tracer(serviceName)),
	}

	return s, nil
}
