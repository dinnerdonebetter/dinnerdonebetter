package grpc

import (
	"context"
	"errors"

	"github.com/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/backend/internal/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/domain/auth/managers"
	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	authsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/auth"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
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
		authenticationManager authentication.Manager
		authManager           *managers.AuthManager
	}
)

func NewService(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	identityRepository identity.Repository,
	// what the actual fuck are we even doing here?
	authManager *managers.AuthManager,
	authenticationManager authentication.Manager,
) authsvc.AuthServiceServer {
	return &serviceImpl{
		logger:                logging.EnsureLogger(logger).WithName(o11yName),
		tracer:                tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		identityRepository:    identityRepository,
		authManager:           authManager,
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
