package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func (s *Server) CreateValidPreparationInstrument(ctx context.Context, request *messages.CreateValidPreparationInstrumentRequest) (*messages.CreateValidPreparationInstrumentResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidPreparationInstrument(ctx context.Context, request *messages.GetValidPreparationInstrumentRequest) (*messages.GetValidPreparationInstrumentResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidPreparationInstruments(ctx context.Context, request *messages.GetValidPreparationInstrumentsRequest) (*messages.GetValidPreparationInstrumentsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidPreparationInstrumentsByInstrument(ctx context.Context, request *messages.GetValidPreparationInstrumentsByInstrumentRequest) (*messages.GetValidPreparationInstrumentsByInstrumentResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidPreparationInstrumentsByPreparation(ctx context.Context, request *messages.GetValidPreparationInstrumentsByPreparationRequest) (*messages.GetValidPreparationInstrumentsByPreparationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateValidPreparationInstrument(ctx context.Context, request *messages.UpdateValidPreparationInstrumentRequest) (*messages.UpdateValidPreparationInstrumentResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveValidPreparationInstrument(ctx context.Context, request *messages.ArchiveValidPreparationInstrumentRequest) (*messages.ArchiveValidPreparationInstrumentResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}
