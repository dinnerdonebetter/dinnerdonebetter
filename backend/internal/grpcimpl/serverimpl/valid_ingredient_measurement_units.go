package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func (s *Server) CreateValidIngredientMeasurementUnit(ctx context.Context, request *messages.CreateValidIngredientMeasurementUnitRequest) (*messages.CreateValidIngredientMeasurementUnitResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredientMeasurementUnit(ctx context.Context, request *messages.GetValidIngredientMeasurementUnitRequest) (*messages.GetValidIngredientMeasurementUnitResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredientMeasurementUnits(ctx context.Context, request *messages.GetValidIngredientMeasurementUnitsRequest) (*messages.GetValidIngredientMeasurementUnitsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredientMeasurementUnitsByIngredient(ctx context.Context, request *messages.GetValidIngredientMeasurementUnitsByIngredientRequest) (*messages.GetValidIngredientMeasurementUnitsByIngredientResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredientMeasurementUnitsByMeasurementUnit(ctx context.Context, request *messages.GetValidIngredientMeasurementUnitsByMeasurementUnitRequest) (*messages.GetValidIngredientMeasurementUnitsByMeasurementUnitResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateValidIngredientMeasurementUnit(ctx context.Context, request *messages.UpdateValidIngredientMeasurementUnitRequest) (*messages.UpdateValidIngredientMeasurementUnitResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveValidIngredientMeasurementUnit(ctx context.Context, request *messages.ArchiveValidIngredientMeasurementUnitRequest) (*messages.ArchiveValidIngredientMeasurementUnitResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}
