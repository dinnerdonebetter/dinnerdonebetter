package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func (s *Server) CreateValidIngredientPreparation(ctx context.Context, request *messages.CreateValidIngredientPreparationRequest) (*messages.CreateValidIngredientPreparationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredientPreparation(ctx context.Context, request *messages.GetValidIngredientPreparationRequest) (*messages.GetValidIngredientPreparationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredientPreparations(ctx context.Context, request *messages.GetValidIngredientPreparationsRequest) (*messages.GetValidIngredientPreparationsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredientPreparationsByIngredient(ctx context.Context, request *messages.GetValidIngredientPreparationsByIngredientRequest) (*messages.GetValidIngredientPreparationsByIngredientResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredientPreparationsByPreparation(ctx context.Context, request *messages.GetValidIngredientPreparationsByPreparationRequest) (*messages.GetValidIngredientPreparationsByPreparationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateValidIngredientPreparation(ctx context.Context, request *messages.UpdateValidIngredientPreparationRequest) (*messages.UpdateValidIngredientPreparationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveValidIngredientPreparation(ctx context.Context, request *messages.ArchiveValidIngredientPreparationRequest) (*messages.ArchiveValidIngredientPreparationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}
