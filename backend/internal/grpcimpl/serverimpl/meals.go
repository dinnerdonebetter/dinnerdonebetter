package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func (s *Server) CreateMeal(ctx context.Context, request *messages.CreateMealRequest) (*messages.CreateMealResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetMeal(ctx context.Context, request *messages.GetMealRequest) (*messages.GetMealResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetMeals(ctx context.Context, request *messages.GetMealsRequest) (*messages.GetMealsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) SearchForMeals(ctx context.Context, request *messages.SearchForMealsRequest) (*messages.SearchForMealsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveMeal(ctx context.Context, request *messages.ArchiveMealRequest) (*messages.ArchiveMealResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}
