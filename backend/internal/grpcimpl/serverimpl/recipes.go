package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func (s *Server) CreateRecipe(ctx context.Context, request *messages.CreateRecipeRequest) (*messages.CreateRecipeResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipe(ctx context.Context, request *messages.GetRecipeRequest) (*messages.GetRecipeResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipes(ctx context.Context, request *messages.GetRecipesRequest) (*messages.GetRecipesResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) SearchForRecipes(ctx context.Context, request *messages.SearchForRecipesRequest) (*messages.SearchForRecipesResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateRecipe(ctx context.Context, request *messages.UpdateRecipeRequest) (*messages.UpdateRecipeResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveRecipe(ctx context.Context, request *messages.ArchiveRecipeRequest) (*messages.ArchiveRecipeResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CloneRecipe(ctx context.Context, request *messages.CloneRecipeRequest) (*messages.CloneRecipeResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}
