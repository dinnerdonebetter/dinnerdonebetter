package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func (s *Server) AdminUpdateUserStatus(ctx context.Context, request *messages.AdminUpdateUserStatusRequest) (*messages.AdminUpdateUserStatusResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}
