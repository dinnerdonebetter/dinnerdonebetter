package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func (s *Server) CreateValidInstrument(ctx context.Context, request *messages.CreateValidInstrumentRequest) (*messages.CreateValidInstrumentResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidInstrument(ctx context.Context, request *messages.GetValidInstrumentRequest) (*messages.GetValidInstrumentResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRandomValidInstrument(ctx context.Context, request *messages.GetRandomValidInstrumentRequest) (*messages.GetRandomValidInstrumentResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidInstruments(ctx context.Context, request *messages.GetValidInstrumentsRequest) (*messages.GetValidInstrumentsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) SearchForValidInstruments(ctx context.Context, request *messages.SearchForValidInstrumentsRequest) (*messages.SearchForValidInstrumentsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateValidInstrument(ctx context.Context, request *messages.UpdateValidInstrumentRequest) (*messages.UpdateValidInstrumentResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveValidInstrument(ctx context.Context, request *messages.ArchiveValidInstrumentRequest) (*messages.ArchiveValidInstrumentResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}
