package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func (s *Server) CreateMealPlanOptionVote(ctx context.Context, request *messages.CreateMealPlanOptionVoteRequest) (*messages.CreateMealPlanOptionVoteResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetMealPlanOptionVote(ctx context.Context, request *messages.GetMealPlanOptionVoteRequest) (*messages.GetMealPlanOptionVoteResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetMealPlanOptionVotes(ctx context.Context, request *messages.GetMealPlanOptionVotesRequest) (*messages.GetMealPlanOptionVotesResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateMealPlanOptionVote(ctx context.Context, request *messages.UpdateMealPlanOptionVoteRequest) (*messages.UpdateMealPlanOptionVoteResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveMealPlanOptionVote(ctx context.Context, request *messages.ArchiveMealPlanOptionVoteRequest) (*messages.ArchiveMealPlanOptionVoteResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}
