package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func (s *Server) CreateValidMeasurementUnitConversion(ctx context.Context, request *messages.CreateValidMeasurementUnitConversionRequest) (*messages.CreateValidMeasurementUnitConversionResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidMeasurementUnitConversion(ctx context.Context, request *messages.GetValidMeasurementUnitConversionRequest) (*messages.GetValidMeasurementUnitConversionResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidMeasurementUnitConversionsFromUnit(ctx context.Context, request *messages.GetValidMeasurementUnitConversionsFromUnitRequest) (*messages.GetValidMeasurementUnitConversionsFromUnitResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidMeasurementUnitConversionsToUnit(ctx context.Context, request *messages.GetValidMeasurementUnitConversionsToUnitRequest) (*messages.GetValidMeasurementUnitConversionsToUnitResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateValidMeasurementUnitConversion(ctx context.Context, request *messages.UpdateValidMeasurementUnitConversionRequest) (*messages.UpdateValidMeasurementUnitConversionResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveValidMeasurementUnitConversion(ctx context.Context, request *messages.ArchiveValidMeasurementUnitConversionRequest) (*messages.ArchiveValidMeasurementUnitConversionResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}
