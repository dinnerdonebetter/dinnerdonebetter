package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func (s *Server) CreateMealPlanOption(ctx context.Context, request *messages.CreateMealPlanOptionRequest) (*messages.CreateMealPlanOptionResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetMealPlanOption(ctx context.Context, request *messages.GetMealPlanOptionRequest) (*messages.GetMealPlanOptionResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetMealPlanOptions(ctx context.Context, request *messages.GetMealPlanOptionsRequest) (*messages.GetMealPlanOptionsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateMealPlanOption(ctx context.Context, request *messages.UpdateMealPlanOptionRequest) (*messages.UpdateMealPlanOptionResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveMealPlanOption(ctx context.Context, request *messages.ArchiveMealPlanOptionRequest) (*messages.ArchiveMealPlanOptionResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}
