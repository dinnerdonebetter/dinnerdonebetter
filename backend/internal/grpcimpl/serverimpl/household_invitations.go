package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func (s *Server) CreateHouseholdInvitation(ctx context.Context, request *messages.CreateHouseholdInvitationRequest) (*messages.CreateHouseholdInvitationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) AcceptHouseholdInvitation(ctx context.Context, request *messages.AcceptHouseholdInvitationRequest) (*messages.AcceptHouseholdInvitationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CancelHouseholdInvitation(ctx context.Context, request *messages.CancelHouseholdInvitationRequest) (*messages.CancelHouseholdInvitationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) RejectHouseholdInvitation(ctx context.Context, request *messages.RejectHouseholdInvitationRequest) (*messages.RejectHouseholdInvitationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetHouseholdInvitation(ctx context.Context, request *messages.GetHouseholdInvitationRequest) (*messages.GetHouseholdInvitationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetHouseholdInvitationByID(ctx context.Context, request *messages.GetHouseholdInvitationByIDRequest) (*messages.GetHouseholdInvitationByIDResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetReceivedHouseholdInvitations(ctx context.Context, request *messages.GetReceivedHouseholdInvitationsRequest) (*messages.GetReceivedHouseholdInvitationsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetSentHouseholdInvitations(ctx context.Context, request *messages.GetSentHouseholdInvitationsRequest) (*messages.GetSentHouseholdInvitationsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}
