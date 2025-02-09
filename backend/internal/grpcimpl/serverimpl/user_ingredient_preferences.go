package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func (s *Server) ArchiveUserIngredientPreference(ctx context.Context, request *messages.ArchiveUserIngredientPreferenceRequest) (*messages.ArchiveUserIngredientPreferenceResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateUserIngredientPreference(ctx context.Context, request *messages.CreateUserIngredientPreferenceRequest) (*messages.CreateUserIngredientPreferenceResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetUserIngredientPreferences(ctx context.Context, request *messages.GetUserIngredientPreferencesRequest) (*messages.GetUserIngredientPreferencesResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}
