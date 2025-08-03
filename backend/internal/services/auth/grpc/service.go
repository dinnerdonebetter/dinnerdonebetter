package grpc

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	"github.com/dinnerdonebetter/backend/internal/domain/identity/managers"
	authsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/auth"
	"github.com/dinnerdonebetter/backend/internal/platform/authentication"
	"github.com/dinnerdonebetter/backend/internal/platform/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"

	oauth2server "github.com/go-oauth2/oauth2/v4"
)

const (
	o11yName = "auth_service"
)

var _ authsvc.AuthServiceServer = (*serviceImpl)(nil)

type (
	serviceImpl struct {
		authsvc.UnimplementedAuthServiceServer
		tracer                    tracing.Tracer
		logger                    logging.Logger
		identityRepository        identity.Repository
		oauth2ClientManager       oauth2server.Manager
		authenticationManager     authentication.Manager
		authenticator             authentication.Authenticator
		sessionContextDataFetcher func(context.Context) (*sessions.ContextData, error)
		authManager               managers.AuthManager
	}
)

func NewService(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	identityRepository identity.Repository,
	oauth2ClientManager oauth2server.Manager,
	authManager managers.AuthManager,
	authenticator authentication.Authenticator,
) authsvc.AuthServiceServer {
	return &serviceImpl{
		logger:              logging.EnsureLogger(logger).WithName(o11yName),
		tracer:              tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		identityRepository:  identityRepository,
		oauth2ClientManager: oauth2ClientManager,
		authManager:         authManager,
		authenticator:       authenticator,
	}
}
