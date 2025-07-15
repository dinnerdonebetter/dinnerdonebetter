package grpc

import (
	"context"
	"github.com/dinnerdonebetter/backend/internal/domain/oauth"
	"github.com/dinnerdonebetter/backend/internal/domain/oauth/converters"
	oauthsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/oauth"
	grpctypes "github.com/dinnerdonebetter/backend/internal/grpc/generated/types"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/random"

	"google.golang.org/grpc/codes"
)

const (
	clientIDSize     = 16
	clientSecretSize = 16
)

var _ oauthsvc.OAuthServiceServer = (*ServiceImpl)(nil)

type (
	ServiceImpl struct {
		oauthsvc.UnimplementedOAuthServiceServer
		tracer          tracing.Tracer
		secretGenerator random.Generator
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

	logger := s.logger.WithValue(keys.OAuth2ClientIDKey, request.OAuth2ClientID)

	if err := s.oauthRepository.ArchiveOAuth2Client(ctx, request.OAuth2ClientID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving oauth2 client")
	}

	x := &oauthsvc.ArchiveOAuth2ClientResponse{
		ResponseDetails: &grpctypes.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *ServiceImpl) CreateOAuth2Client(ctx context.Context, request *oauthsvc.CreateOAuth2ClientRequest) (*oauthsvc.CreateOAuth2ClientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.Clone()

	dbInput := converters.ConvertOAuth2ClientCreationRequestInputToOAuth2ClientDatabaseCreationInput(&oauth.OAuth2ClientCreationRequestInput{
		Name:        request.Name,
		Description: request.Description,
	})
	dbInput.ID = identifiers.New()

	var err error
	if dbInput.ClientID, err = s.secretGenerator.GenerateHexEncodedString(ctx, clientIDSize); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "generating client id")

	}

	if dbInput.ClientSecret, err = s.secretGenerator.GenerateHexEncodedString(ctx, clientSecretSize); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "generating client secret")

	}

	created, err := s.oauthRepository.CreateOAuth2Client(ctx, dbInput)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating oauth2 client")
	}

	x := &oauthsvc.CreateOAuth2ClientResponse{
		ResponseDetails: &grpctypes.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Created: ConvertOAuth2ClientToGRPCOAuth2Client(created),
	}

	return x, nil
}

func (s *ServiceImpl) GetOAuth2Client(ctx context.Context, request *oauthsvc.GetOAuth2ClientRequest) (*oauthsvc.GetOAuth2ClientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithValue(keys.OAuth2ClientIDKey, request.OAuth2ClientID)

	oauth2Client, err := s.oauthRepository.GetOAuth2ClientByDatabaseID(ctx, request.OAuth2ClientID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "getting oauth2 client by database ID")
	}

	x := &oauthsvc.GetOAuth2ClientResponse{
		ResponseDetails: &grpctypes.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Result: ConvertOAuth2ClientToGRPCOAuth2Client(oauth2Client),
	}

	return x, nil
}

func (s *ServiceImpl) GetOAuth2Clients(ctx context.Context, request *oauthsvc.GetOAuth2ClientsRequest) (*oauthsvc.GetOAuth2ClientsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &oauthsvc.GetOAuth2ClientsResponse{
		ResponseDetails: &grpctypes.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}
