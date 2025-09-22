package grpc

import (
	"github.com/dinnerdonebetter/backend/internal/domain/oauth/manager"
	oauthsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/oauth"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

const (
	o11yName = "oauth_service"
)

var _ oauthsvc.OAuthServiceServer = (*serviceImpl)(nil)

type (
	serviceImpl struct {
		oauthsvc.UnimplementedOAuthServiceServer
		tracer           tracing.Tracer
		logger           logging.Logger
		oauthDataManager manager.OAuth2Manager
	}
)

func NewService(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	oauthDataManager manager.OAuth2Manager,
) oauthsvc.OAuthServiceServer {
	return &serviceImpl{
		logger:           logging.EnsureLogger(logger).WithName(o11yName),
		tracer:           tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		oauthDataManager: oauthDataManager,
	}
}
