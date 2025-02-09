package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func (s *Server) CreateRecipeRating(ctx context.Context, request *messages.CreateRecipeRatingRequest) (*messages.CreateRecipeRatingResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipeRating(ctx context.Context, request *messages.GetRecipeRatingRequest) (*messages.GetRecipeRatingResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipeRatingsForRecipe(ctx context.Context, request *messages.GetRecipeRatingsForRecipeRequest) (*messages.GetRecipeRatingsForRecipeResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateRecipeRating(ctx context.Context, request *messages.UpdateRecipeRatingRequest) (*messages.UpdateRecipeRatingResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveRecipeRating(ctx context.Context, request *messages.ArchiveRecipeRatingRequest) (*messages.ArchiveRecipeRatingResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}
