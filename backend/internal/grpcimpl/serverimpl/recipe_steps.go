package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func (s *Server) CreateRecipeStep(ctx context.Context, request *messages.CreateRecipeStepRequest) (*messages.CreateRecipeStepResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipeStep(ctx context.Context, request *messages.GetRecipeStepRequest) (*messages.GetRecipeStepResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipeSteps(ctx context.Context, request *messages.GetRecipeStepsRequest) (*messages.GetRecipeStepsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateRecipeStep(ctx context.Context, request *messages.UpdateRecipeStepRequest) (*messages.UpdateRecipeStepResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveRecipeStep(ctx context.Context, request *messages.ArchiveRecipeStepRequest) (*messages.ArchiveRecipeStepResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}
