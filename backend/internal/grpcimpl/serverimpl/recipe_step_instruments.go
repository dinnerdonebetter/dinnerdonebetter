package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func (s *Server) CreateRecipeStepInstrument(ctx context.Context, request *messages.CreateRecipeStepInstrumentRequest) (*messages.CreateRecipeStepInstrumentResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipeStepInstrument(ctx context.Context, request *messages.GetRecipeStepInstrumentRequest) (*messages.GetRecipeStepInstrumentResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipeStepInstruments(ctx context.Context, request *messages.GetRecipeStepInstrumentsRequest) (*messages.GetRecipeStepInstrumentsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateRecipeStepInstrument(ctx context.Context, request *messages.UpdateRecipeStepInstrumentRequest) (*messages.UpdateRecipeStepInstrumentResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveRecipeStepInstrument(ctx context.Context, request *messages.ArchiveRecipeStepInstrumentRequest) (*messages.ArchiveRecipeStepInstrumentResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}
