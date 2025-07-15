package grpc

import (
	"github.com/dinnerdonebetter/backend/internal/domain/oauth"
	oauthsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/oauth"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

var _ oauthsvc.OAuthServiceServer = (*ServiceImpl)(nil)

type (
	ServiceImpl struct {
		oauthsvc.UnimplementedOAuthServiceServer
		tracer          tracing.Tracer
		logger          logging.Logger
		oauthRepository oauth.Repository
	}
)

func NewService(
	oauthRepository oauth.Repository,
) oauthsvc.OAuthServiceServer {
	return &ServiceImpl{
		oauthRepository: oauthRepository,
	}
}
