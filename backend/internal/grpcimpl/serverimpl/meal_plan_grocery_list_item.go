package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func (s *Server) CreateMealPlanGroceryListItem(ctx context.Context, request *messages.CreateMealPlanGroceryListItemRequest) (*messages.CreateMealPlanGroceryListItemResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetMealPlanGroceryListItem(ctx context.Context, request *messages.GetMealPlanGroceryListItemRequest) (*messages.GetMealPlanGroceryListItemResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetMealPlanGroceryListItemsForMealPlan(ctx context.Context, request *messages.GetMealPlanGroceryListItemsForMealPlanRequest) (*messages.GetMealPlanGroceryListItemsForMealPlanResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateMealPlanGroceryListItem(ctx context.Context, request *messages.UpdateMealPlanGroceryListItemRequest) (*messages.UpdateMealPlanGroceryListItemResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveMealPlanGroceryListItem(ctx context.Context, request *messages.ArchiveMealPlanGroceryListItemRequest) (*messages.ArchiveMealPlanGroceryListItemResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) RunMealPlanGroceryListInitializerWorker(ctx context.Context, request *messages.RunMealPlanGroceryListInitializerWorkerRequest) (*messages.RunMealPlanGroceryListInitializerWorkerResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}
