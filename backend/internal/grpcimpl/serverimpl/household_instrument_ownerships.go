package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func (s *Server) CreateHouseholdInstrumentOwnership(ctx context.Context, request *messages.CreateHouseholdInstrumentOwnershipRequest) (*messages.CreateHouseholdInstrumentOwnershipResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetHouseholdInstrumentOwnership(ctx context.Context, request *messages.GetHouseholdInstrumentOwnershipRequest) (*messages.GetHouseholdInstrumentOwnershipResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetHouseholdInstrumentOwnerships(ctx context.Context, request *messages.GetHouseholdInstrumentOwnershipsRequest) (*messages.GetHouseholdInstrumentOwnershipsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateHouseholdInstrumentOwnership(ctx context.Context, request *messages.UpdateHouseholdInstrumentOwnershipRequest) (*messages.UpdateHouseholdInstrumentOwnershipResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveHouseholdInstrumentOwnership(ctx context.Context, request *messages.ArchiveHouseholdInstrumentOwnershipRequest) (*messages.ArchiveHouseholdInstrumentOwnershipResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}
