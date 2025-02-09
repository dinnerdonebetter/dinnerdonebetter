package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func (s *Server) CreateValidIngredientGroup(ctx context.Context, request *messages.CreateValidIngredientGroupRequest) (*messages.CreateValidIngredientGroupResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredientGroup(ctx context.Context, request *messages.GetValidIngredientGroupRequest) (*messages.GetValidIngredientGroupResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredientGroups(ctx context.Context, request *messages.GetValidIngredientGroupsRequest) (*messages.GetValidIngredientGroupsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) SearchForValidIngredientGroups(ctx context.Context, request *messages.SearchForValidIngredientGroupsRequest) (*messages.SearchForValidIngredientGroupsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateValidIngredientGroup(ctx context.Context, request *messages.UpdateValidIngredientGroupRequest) (*messages.UpdateValidIngredientGroupResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveValidIngredientGroup(ctx context.Context, request *messages.ArchiveValidIngredientGroupRequest) (*messages.ArchiveValidIngredientGroupResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}
