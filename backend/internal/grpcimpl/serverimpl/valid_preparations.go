package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) CreateValidPreparation(ctx context.Context, input *messages.ValidPreparationCreationRequestInput) (*messages.ValidPreparation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidPreparation(ctx context.Context, request *messages.GetValidPreparationRequest) (*messages.ValidPreparation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidPreparations(ctx context.Context, request *messages.GetValidPreparationsRequest) (*messages.ValidPreparation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRandomValidPreparation(ctx context.Context, _ *emptypb.Empty) (*messages.ValidPreparation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) SearchForValidPreparations(ctx context.Context, request *messages.SearchForValidPreparationsRequest) (*messages.ValidPreparation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateValidPreparation(ctx context.Context, request *messages.UpdateValidPreparationRequest) (*messages.ValidPreparation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveValidPreparation(ctx context.Context, request *messages.ArchiveValidPreparationRequest) (*messages.ValidPreparation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}
