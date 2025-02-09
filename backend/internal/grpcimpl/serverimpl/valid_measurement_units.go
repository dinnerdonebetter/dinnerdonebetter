package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func (s *Server) CreateValidMeasurementUnit(ctx context.Context, request *messages.CreateValidMeasurementUnitRequest) (*messages.CreateValidMeasurementUnitResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidMeasurementUnit(ctx context.Context, request *messages.GetValidMeasurementUnitRequest) (*messages.GetValidMeasurementUnitResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidMeasurementUnits(ctx context.Context, request *messages.GetValidMeasurementUnitsRequest) (*messages.GetValidMeasurementUnitsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) SearchForValidMeasurementUnits(ctx context.Context, request *messages.SearchForValidMeasurementUnitsRequest) (*messages.SearchForValidMeasurementUnitsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) SearchValidMeasurementUnitsByIngredient(ctx context.Context, request *messages.SearchValidMeasurementUnitsByIngredientRequest) (*messages.SearchValidMeasurementUnitsByIngredientResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateValidMeasurementUnit(ctx context.Context, request *messages.UpdateValidMeasurementUnitRequest) (*messages.UpdateValidMeasurementUnitResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveValidMeasurementUnit(ctx context.Context, request *messages.ArchiveValidMeasurementUnitRequest) (*messages.ArchiveValidMeasurementUnitResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}
