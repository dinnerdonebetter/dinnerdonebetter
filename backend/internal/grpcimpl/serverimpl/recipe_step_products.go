package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func (s *Server) CreateRecipeStepProduct(ctx context.Context, request *messages.CreateRecipeStepProductRequest) (*messages.CreateRecipeStepProductResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipeStepProduct(ctx context.Context, request *messages.GetRecipeStepProductRequest) (*messages.GetRecipeStepProductResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipeStepProducts(ctx context.Context, request *messages.GetRecipeStepProductsRequest) (*messages.GetRecipeStepProductsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateRecipeStepProduct(ctx context.Context, request *messages.UpdateRecipeStepProductRequest) (*messages.UpdateRecipeStepProductResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveRecipeStepProduct(ctx context.Context, request *messages.ArchiveRecipeStepProductRequest) (*messages.ArchiveRecipeStepProductResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}
