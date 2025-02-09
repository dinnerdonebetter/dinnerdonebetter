package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func (s *Server) ArchiveWebhook(ctx context.Context, request *messages.ArchiveWebhookRequest) (*messages.ArchiveWebhookResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveWebhookTriggerEvent(ctx context.Context, request *messages.ArchiveWebhookTriggerEventRequest) (*messages.ArchiveWebhookTriggerEventResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateWebhook(ctx context.Context, request *messages.CreateWebhookRequest) (*messages.CreateWebhookResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateWebhookTriggerEvent(ctx context.Context, request *messages.CreateWebhookTriggerEventRequest) (*messages.CreateWebhookTriggerEventResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetWebhook(ctx context.Context, request *messages.GetWebhookRequest) (*messages.GetWebhookResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetWebhooks(ctx context.Context, request *messages.GetWebhooksRequest) (*messages.GetWebhooksResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}
