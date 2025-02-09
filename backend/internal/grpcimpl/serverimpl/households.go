package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func (s *Server) CreateHousehold(ctx context.Context, request *messages.CreateHouseholdRequest) (*messages.CreateHouseholdResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetActiveHousehold(ctx context.Context, request *messages.GetActiveHouseholdRequest) (*messages.GetActiveHouseholdResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetHousehold(ctx context.Context, request *messages.GetHouseholdRequest) (*messages.GetHouseholdResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetHouseholds(ctx context.Context, request *messages.GetHouseholdsRequest) (*messages.GetHouseholdsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) SetDefaultHousehold(ctx context.Context, request *messages.SetDefaultHouseholdRequest) (*messages.SetDefaultHouseholdResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) TransferHouseholdOwnership(ctx context.Context, request *messages.TransferHouseholdOwnershipRequest) (*messages.TransferHouseholdOwnershipResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateHousehold(ctx context.Context, request *messages.UpdateHouseholdRequest) (*messages.UpdateHouseholdResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateHouseholdMemberPermissions(ctx context.Context, request *messages.UpdateHouseholdMemberPermissionsRequest) (*messages.UpdateHouseholdMemberPermissionsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveHousehold(ctx context.Context, request *messages.ArchiveHouseholdRequest) (*messages.ArchiveHouseholdResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveUserMembership(ctx context.Context, request *messages.ArchiveUserMembershipRequest) (*messages.ArchiveUserMembershipResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}
