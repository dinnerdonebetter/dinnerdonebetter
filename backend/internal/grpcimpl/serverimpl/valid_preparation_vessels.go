package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func (s *Server) CreateValidPreparationVessel(ctx context.Context, request *messages.CreateValidPreparationVesselRequest) (*messages.CreateValidPreparationVesselResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidPreparationVessel(ctx context.Context, request *messages.GetValidPreparationVesselRequest) (*messages.GetValidPreparationVesselResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidPreparationVessels(ctx context.Context, request *messages.GetValidPreparationVesselsRequest) (*messages.GetValidPreparationVesselsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidPreparationVesselsByPreparation(ctx context.Context, request *messages.GetValidPreparationVesselsByPreparationRequest) (*messages.GetValidPreparationVesselsByPreparationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidPreparationVesselsByVessel(ctx context.Context, request *messages.GetValidPreparationVesselsByVesselRequest) (*messages.GetValidPreparationVesselsByVesselResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateValidPreparationVessel(ctx context.Context, request *messages.UpdateValidPreparationVesselRequest) (*messages.UpdateValidPreparationVesselResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveValidPreparationVessel(ctx context.Context, request *messages.ArchiveValidPreparationVesselRequest) (*messages.ArchiveValidPreparationVesselResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}
