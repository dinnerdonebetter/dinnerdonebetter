package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func (s *Server) Ping(ctx context.Context, _ *messages.PingRequest) (*messages.PingResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return &messages.PingResponse{Meta: &messages.ResponseMeta{TraceID: span.SpanContext().TraceID().String()}}, nil
}
