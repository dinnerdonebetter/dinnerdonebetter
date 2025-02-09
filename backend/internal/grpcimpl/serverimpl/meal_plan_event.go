package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func (s *Server) CreateMealPlanEvent(ctx context.Context, request *messages.CreateMealPlanEventRequest) (*messages.CreateMealPlanEventResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetMealPlanEvent(ctx context.Context, request *messages.GetMealPlanEventRequest) (*messages.GetMealPlanEventResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetMealPlanEvents(ctx context.Context, request *messages.GetMealPlanEventsRequest) (*messages.GetMealPlanEventsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateMealPlanEvent(ctx context.Context, request *messages.UpdateMealPlanEventRequest) (*messages.UpdateMealPlanEventResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveMealPlanEvent(ctx context.Context, request *messages.ArchiveMealPlanEventRequest) (*messages.ArchiveMealPlanEventResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}
