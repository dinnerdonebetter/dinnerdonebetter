package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func (s *Server) CreateRecipeStepVessel(ctx context.Context, request *messages.CreateRecipeStepVesselRequest) (*messages.CreateRecipeStepVesselResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipeStepVessel(ctx context.Context, request *messages.GetRecipeStepVesselRequest) (*messages.GetRecipeStepVesselResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipeStepVessels(ctx context.Context, request *messages.GetRecipeStepVesselsRequest) (*messages.GetRecipeStepVesselsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateRecipeStepVessel(ctx context.Context, request *messages.UpdateRecipeStepVesselRequest) (*messages.UpdateRecipeStepVesselResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveRecipeStepVessel(ctx context.Context, request *messages.ArchiveRecipeStepVesselRequest) (*messages.ArchiveRecipeStepVesselResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}
