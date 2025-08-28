package grpc

import (
	"context"
	"errors"
	authentication2 "github.com/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/backend/internal/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/domain/auth/managers"
	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	authsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/auth"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"

	"github.com/go-oauth2/oauth2/v4/manage"
)

const (
	o11yName = "auth_service"
)

var _ authsvc.AuthServiceServer = (*serviceImpl)(nil)

type (
	serviceImpl struct {
		authsvc.UnimplementedAuthServiceServer
		tracer                tracing.Tracer
		logger                logging.Logger
		identityRepository    identity.Repository
		oauth2ClientManager   *manage.Manager
		authenticationManager authentication2.Manager
		authenticator         authentication2.Authenticator
		authManager           *managers.AuthManager
	}
)

func NewService(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	identityRepository identity.Repository,
	oauth2ClientManager *manage.Manager,
	// bruh what the actual fuck are we even doing here
	authManager *managers.AuthManager,
	authenticator authentication2.Authenticator,
	authenticationManager authentication2.Manager,
) authsvc.AuthServiceServer {
	return &serviceImpl{
		logger:                logging.EnsureLogger(logger).WithName(o11yName),
		tracer:                tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		identityRepository:    identityRepository,
		authManager:           authManager,
		oauth2ClientManager:   oauth2ClientManager,
		authenticator:         authenticator,
		authenticationManager: authenticationManager,
	}
}

func (s *serviceImpl) fetchSessionContext(ctx context.Context) (*sessions.ContextData, error) {
	sessionContext, ok := ctx.Value(sessions.SessionContextDataKey).(*sessions.ContextData)
	if !ok {
		return nil, errors.New("session context not found")
	}

	return sessionContext, nil
}
