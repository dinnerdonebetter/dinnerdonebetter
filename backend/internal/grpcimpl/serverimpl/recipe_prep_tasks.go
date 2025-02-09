package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func (s *Server) CreateRecipePrepTask(ctx context.Context, request *messages.CreateRecipePrepTaskRequest) (*messages.CreateRecipePrepTaskResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipePrepTask(ctx context.Context, request *messages.GetRecipePrepTaskRequest) (*messages.GetRecipePrepTaskResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipePrepTasks(ctx context.Context, request *messages.GetRecipePrepTasksRequest) (*messages.GetRecipePrepTasksResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateRecipePrepTask(ctx context.Context, request *messages.UpdateRecipePrepTaskRequest) (*messages.UpdateRecipePrepTaskResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveRecipePrepTask(ctx context.Context, request *messages.ArchiveRecipePrepTaskRequest) (*messages.ArchiveRecipePrepTaskResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}
