package grpc

import (
	"context"
	"errors"

	authentication2 "github.com/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/backend/internal/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/domain/auth/managers"
	identitymanager "github.com/dinnerdonebetter/backend/internal/domain/identity/manager"
	authsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/auth"
	"github.com/dinnerdonebetter/backend/internal/platform/featureflags"
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
		identityDataManager   identitymanager.IdentityDataManager
		authenticationManager authentication2.Manager
		authManager           managers.AuthManagerInterface
		featureFlagManager    featureflags.FeatureFlagManager
	}
)

func NewAuthService(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	identityDataManager identitymanager.IdentityDataManager,
	authManager managers.AuthManagerInterface,
	authenticationManager authentication2.Manager,
	featureFlagManager featureflags.FeatureFlagManager,
) authsvc.AuthServiceServer {
	return &serviceImpl{
		logger:                logging.EnsureLogger(logger).WithName(o11yName),
		tracer:                tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		identityDataManager:   identityDataManager,
		authManager:           authManager,
		authenticationManager: authenticationManager,
		featureFlagManager:    featureFlagManager,
	}
}

func (s *serviceImpl) fetchSessionContext(ctx context.Context) (*sessions.ContextData, error) {
	sessionContext, ok := ctx.Value(sessions.SessionContextDataKey).(*sessions.ContextData)
	if !ok {
		return nil, errors.New("session context not found")
	}

	return sessionContext, nil
}
