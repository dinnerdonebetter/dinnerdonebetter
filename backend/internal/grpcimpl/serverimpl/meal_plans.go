package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func (s *Server) CreateMealPlan(ctx context.Context, request *messages.CreateMealPlanRequest) (*messages.CreateMealPlanResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetMealPlansForHousehold(ctx context.Context, request *messages.GetMealPlansForHouseholdRequest) (*messages.GetMealPlansForHouseholdResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) FinalizeMealPlan(ctx context.Context, request *messages.FinalizeMealPlanRequest) (*messages.FinalizeMealPlanResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetMealPlan(ctx context.Context, request *messages.GetMealPlanRequest) (*messages.GetMealPlanResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateMealPlan(ctx context.Context, request *messages.UpdateMealPlanRequest) (*messages.UpdateMealPlanResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) RunFinalizeMealPlanWorker(ctx context.Context, request *messages.RunFinalizeMealPlanWorkerRequest) (*messages.RunFinalizeMealPlanWorkerResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveMealPlan(ctx context.Context, request *messages.ArchiveMealPlanRequest) (*messages.ArchiveMealPlanResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}
