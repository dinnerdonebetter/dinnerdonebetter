package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func (s *Server) CreateRecipeStepIngredient(ctx context.Context, request *messages.CreateRecipeStepIngredientRequest) (*messages.CreateRecipeStepIngredientResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipeStepIngredient(ctx context.Context, request *messages.GetRecipeStepIngredientRequest) (*messages.GetRecipeStepIngredientResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipeStepIngredients(ctx context.Context, request *messages.GetRecipeStepIngredientsRequest) (*messages.GetRecipeStepIngredientsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateRecipeStepIngredient(ctx context.Context, request *messages.UpdateRecipeStepIngredientRequest) (*messages.UpdateRecipeStepIngredientResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveRecipeStepIngredient(ctx context.Context, request *messages.ArchiveRecipeStepIngredientRequest) (*messages.ArchiveRecipeStepIngredientResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}
