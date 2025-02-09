package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func (s *Server) CreateValidIngredientStateIngredient(ctx context.Context, request *messages.CreateValidIngredientStateIngredientRequest) (*messages.CreateValidIngredientStateIngredientResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredientStateIngredient(ctx context.Context, request *messages.GetValidIngredientStateIngredientRequest) (*messages.GetValidIngredientStateIngredientResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredientStateIngredients(ctx context.Context, request *messages.GetValidIngredientStateIngredientsRequest) (*messages.GetValidIngredientStateIngredientsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredientStateIngredientsByIngredient(ctx context.Context, request *messages.GetValidIngredientStateIngredientsByIngredientRequest) (*messages.GetValidIngredientStateIngredientsByIngredientResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredientStateIngredientsByIngredientState(ctx context.Context, request *messages.GetValidIngredientStateIngredientsByIngredientStateRequest) (*messages.GetValidIngredientStateIngredientsByIngredientStateResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateValidIngredientStateIngredient(ctx context.Context, request *messages.UpdateValidIngredientStateIngredientRequest) (*messages.UpdateValidIngredientStateIngredientResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveValidIngredientStateIngredient(ctx context.Context, request *messages.ArchiveValidIngredientStateIngredientRequest) (*messages.ArchiveValidIngredientStateIngredientResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}
