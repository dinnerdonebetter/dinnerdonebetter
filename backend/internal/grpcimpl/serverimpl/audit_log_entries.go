package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func (s *Server) GetAuditLogEntriesForHousehold(ctx context.Context, request *messages.GetAuditLogEntriesForHouseholdRequest) (*messages.GetAuditLogEntriesForHouseholdResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetAuditLogEntriesForUser(ctx context.Context, request *messages.GetAuditLogEntriesForUserRequest) (*messages.GetAuditLogEntriesForUserResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetAuditLogEntryByID(ctx context.Context, request *messages.GetAuditLogEntryByIDRequest) (*messages.GetAuditLogEntryByIDResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}
