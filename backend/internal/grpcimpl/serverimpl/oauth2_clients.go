package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func (s *Server) ArchiveOAuth2Client(ctx context.Context, request *messages.ArchiveOAuth2ClientRequest) (*messages.ArchiveOAuth2ClientResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateOAuth2Client(ctx context.Context, request *messages.CreateOAuth2ClientRequest) (*messages.CreateOAuth2ClientResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetOAuth2Client(ctx context.Context, request *messages.GetOAuth2ClientRequest) (*messages.GetOAuth2ClientResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetOAuth2Clients(ctx context.Context, request *messages.GetOAuth2ClientsRequest) (*messages.GetOAuth2ClientsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}
