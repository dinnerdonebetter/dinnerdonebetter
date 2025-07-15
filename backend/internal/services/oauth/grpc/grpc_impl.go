package grpc

import (
	"context"

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

func (s *ServiceImpl) ArchiveOAuth2Client(ctx context.Context, request *oauthsvc.ArchiveOAuth2ClientRequest) (*oauthsvc.ArchiveOAuth2ClientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &oauthsvc.ArchiveOAuth2ClientResponse{}

	return x, nil
}

func (s *ServiceImpl) CreateOAuth2Client(ctx context.Context, request *oauthsvc.CreateOAuth2ClientRequest) (*oauthsvc.CreateOAuth2ClientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &oauthsvc.CreateOAuth2ClientResponse{}

	return x, nil
}

func (s *ServiceImpl) GetOAuth2Client(ctx context.Context, request *oauthsvc.GetOAuth2ClientRequest) (*oauthsvc.GetOAuth2ClientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &oauthsvc.GetOAuth2ClientResponse{}

	return x, nil
}

func (s *ServiceImpl) GetOAuth2Clients(ctx context.Context, request *oauthsvc.GetOAuth2ClientsRequest) (*oauthsvc.GetOAuth2ClientsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &oauthsvc.GetOAuth2ClientsResponse{}

	return x, nil
}
