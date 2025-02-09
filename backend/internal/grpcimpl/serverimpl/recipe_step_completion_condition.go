package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func (s *Server) CreateRecipeStepCompletionCondition(ctx context.Context, request *messages.CreateRecipeStepCompletionConditionRequest) (*messages.CreateRecipeStepCompletionConditionResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipeStepCompletionCondition(ctx context.Context, request *messages.GetRecipeStepCompletionConditionRequest) (*messages.GetRecipeStepCompletionConditionResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipeStepCompletionConditions(ctx context.Context, request *messages.GetRecipeStepCompletionConditionsRequest) (*messages.GetRecipeStepCompletionConditionsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateRecipeStepCompletionCondition(ctx context.Context, request *messages.UpdateRecipeStepCompletionConditionRequest) (*messages.UpdateRecipeStepCompletionConditionResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveRecipeStepCompletionCondition(ctx context.Context, request *messages.ArchiveRecipeStepCompletionConditionRequest) (*messages.ArchiveRecipeStepCompletionConditionResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}
