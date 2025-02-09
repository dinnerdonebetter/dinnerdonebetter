package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func (s *Server) CreateValidIngredientState(ctx context.Context, request *messages.CreateValidIngredientStateRequest) (*messages.CreateValidIngredientStateResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredientState(ctx context.Context, request *messages.GetValidIngredientStateRequest) (*messages.GetValidIngredientStateResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredientStates(ctx context.Context, request *messages.GetValidIngredientStatesRequest) (*messages.GetValidIngredientStatesResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) SearchForValidIngredientStates(ctx context.Context, request *messages.SearchForValidIngredientStatesRequest) (*messages.SearchForValidIngredientStatesResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateValidIngredientState(ctx context.Context, request *messages.UpdateValidIngredientStateRequest) (*messages.UpdateValidIngredientStateResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveValidIngredientState(ctx context.Context, request *messages.ArchiveValidIngredientStateRequest) (*messages.ArchiveValidIngredientStateResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}
