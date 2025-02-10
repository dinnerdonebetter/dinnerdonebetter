package serverimpl

import (
	"context"
	"github.com/dinnerdonebetter/backend/internal/lib/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func (s *Server) buildResponseMeta(span tracing.Span) *messages.ResponseMeta {
	return &messages.ResponseMeta{TraceID: span.SpanContext().TraceID().String()}
}

func (s *Server) buildResponseMetaForUser(span tracing.Span, sessionContextData sessions.ContextData) *messages.ResponseMeta {
	return &messages.ResponseMeta{
		CurrentHouseholdID: sessionContextData.ActiveHouseholdID,
		TraceID:            span.SpanContext().TraceID().String(),
	}
}

func (s *Server) Ping(ctx context.Context, _ *messages.PingRequest) (*messages.PingResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return &messages.PingResponse{Meta: &messages.ResponseMeta{TraceID: span.SpanContext().TraceID().String()}}, nil
}

func (s *Server) CheckForReadiness(ctx context.Context, request *messages.CheckForReadinessRequest) (*messages.CheckForReadinessResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) PublishArbitraryQueueMessage(ctx context.Context, request *messages.PublishArbitraryQueueMessageRequest) (*messages.PublishArbitraryQueueMessageResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}
