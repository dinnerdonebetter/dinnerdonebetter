package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func (s *Server) CreateValidVessel(ctx context.Context, request *messages.CreateValidVesselRequest) (*messages.CreateValidVesselResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidVessel(ctx context.Context, request *messages.GetValidVesselRequest) (*messages.GetValidVesselResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRandomValidVessel(ctx context.Context, request *messages.GetRandomValidVesselRequest) (*messages.GetRandomValidVesselResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidVessels(ctx context.Context, request *messages.GetValidVesselsRequest) (*messages.GetValidVesselsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) SearchForValidVessels(ctx context.Context, request *messages.SearchForValidVesselsRequest) (*messages.SearchForValidVesselsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateValidVessel(ctx context.Context, request *messages.UpdateValidVesselRequest) (*messages.UpdateValidVesselResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveValidVessel(ctx context.Context, request *messages.ArchiveValidVesselRequest) (*messages.ArchiveValidVesselResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}
