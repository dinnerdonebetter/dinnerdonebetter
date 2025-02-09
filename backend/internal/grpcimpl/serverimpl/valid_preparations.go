package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func (s *Server) CreateValidPreparation(ctx context.Context, request *messages.CreateValidPreparationRequest) (*messages.CreateValidPreparationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidPreparation(ctx context.Context, request *messages.GetValidPreparationRequest) (*messages.GetValidPreparationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRandomValidPreparation(ctx context.Context, request *messages.GetRandomValidPreparationRequest) (*messages.GetRandomValidPreparationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) SearchForValidPreparations(ctx context.Context, request *messages.SearchForValidPreparationsRequest) (*messages.SearchForValidPreparationsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidPreparations(ctx context.Context, request *messages.GetValidPreparationsRequest) (*messages.GetValidPreparationsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateValidPreparation(ctx context.Context, request *messages.UpdateValidPreparationRequest) (*messages.UpdateValidPreparationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveValidPreparation(ctx context.Context, request *messages.ArchiveValidPreparationRequest) (*messages.ArchiveValidPreparationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}
