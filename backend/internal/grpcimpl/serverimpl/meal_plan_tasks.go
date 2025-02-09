package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func (s *Server) CreateMealPlanTask(ctx context.Context, request *messages.CreateMealPlanTaskRequest) (*messages.CreateMealPlanTaskResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetMealPlanTask(ctx context.Context, request *messages.GetMealPlanTaskRequest) (*messages.GetMealPlanTaskResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetMealPlanTasks(ctx context.Context, request *messages.GetMealPlanTasksRequest) (*messages.GetMealPlanTasksResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipeMealPlanTasks(ctx context.Context, request *messages.GetRecipeMealPlanTasksRequest) (*messages.GetRecipeMealPlanTasksResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) RunMealPlanTaskCreatorWorker(ctx context.Context, request *messages.RunMealPlanTaskCreatorWorkerRequest) (*messages.RunMealPlanTaskCreatorWorkerResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateMealPlanTaskStatus(ctx context.Context, request *messages.UpdateMealPlanTaskStatusRequest) (*messages.UpdateMealPlanTaskStatusResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}
