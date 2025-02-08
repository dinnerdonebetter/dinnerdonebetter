package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func (s *Server) ArchiveMeal(ctx context.Context, request *messages.ArchiveMealRequest) (*messages.Meal, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveMealPlan(ctx context.Context, request *messages.ArchiveMealPlanRequest) (*messages.MealPlan, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveMealPlanEvent(ctx context.Context, request *messages.ArchiveMealPlanEventRequest) (*messages.MealPlanEvent, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveMealPlanGroceryListItem(ctx context.Context, request *messages.ArchiveMealPlanGroceryListItemRequest) (*messages.MealPlanGroceryListItem, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveMealPlanOption(ctx context.Context, request *messages.ArchiveMealPlanOptionRequest) (*messages.MealPlanOption, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveMealPlanOptionVote(ctx context.Context, request *messages.ArchiveMealPlanOptionVoteRequest) (*messages.MealPlanOptionVote, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateMeal(ctx context.Context, input *messages.MealCreationRequestInput) (*messages.Meal, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateMealPlan(ctx context.Context, input *messages.MealPlanCreationRequestInput) (*messages.MealPlan, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateMealPlanEvent(ctx context.Context, request *messages.CreateMealPlanEventRequest) (*messages.MealPlanEvent, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateMealPlanGroceryListItem(ctx context.Context, request *messages.CreateMealPlanGroceryListItemRequest) (*messages.MealPlanGroceryListItem, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateMealPlanOption(ctx context.Context, request *messages.CreateMealPlanOptionRequest) (*messages.MealPlanOption, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateMealPlanOptionVote(ctx context.Context, request *messages.CreateMealPlanOptionVoteRequest) (*messages.MealPlanOptionVote, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateMealPlanTask(ctx context.Context, request *messages.CreateMealPlanTaskRequest) (*messages.MealPlanTask, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) FinalizeMealPlan(ctx context.Context, request *messages.FinalizeMealPlanRequest) (*messages.FinalizeMealPlansResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetMeal(ctx context.Context, request *messages.GetMealRequest) (*messages.Meal, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetMealPlan(ctx context.Context, request *messages.GetMealPlanRequest) (*messages.MealPlan, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetMealPlanEvent(ctx context.Context, request *messages.GetMealPlanEventRequest) (*messages.MealPlanEvent, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetMealPlanEvents(ctx context.Context, request *messages.GetMealPlanEventsRequest) (*messages.MealPlanEvent, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetMealPlanGroceryListItem(ctx context.Context, request *messages.GetMealPlanGroceryListItemRequest) (*messages.MealPlanGroceryListItem, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetMealPlanGroceryListItemsForMealPlan(ctx context.Context, request *messages.GetMealPlanGroceryListItemsForMealPlanRequest) (*messages.MealPlanGroceryListItem, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetMealPlanOption(ctx context.Context, request *messages.GetMealPlanOptionRequest) (*messages.MealPlanOption, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetMealPlanOptionVote(ctx context.Context, request *messages.GetMealPlanOptionVoteRequest) (*messages.MealPlanOptionVote, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetMealPlanOptionVotes(ctx context.Context, request *messages.GetMealPlanOptionVotesRequest) (*messages.MealPlanOptionVote, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetMealPlanOptions(ctx context.Context, request *messages.GetMealPlanOptionsRequest) (*messages.MealPlanOption, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetMealPlanTask(ctx context.Context, request *messages.GetMealPlanTaskRequest) (*messages.MealPlanTask, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetMealPlanTasks(ctx context.Context, request *messages.GetMealPlanTasksRequest) (*messages.MealPlanTask, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetMealPlansForHousehold(ctx context.Context, request *messages.GetMealPlansForHouseholdRequest) (*messages.MealPlan, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetMeals(ctx context.Context, request *messages.GetMealsRequest) (*messages.Meal, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateMealPlan(ctx context.Context, request *messages.UpdateMealPlanRequest) (*messages.MealPlan, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateMealPlanEvent(ctx context.Context, request *messages.UpdateMealPlanEventRequest) (*messages.MealPlanEvent, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateMealPlanGroceryListItem(ctx context.Context, request *messages.UpdateMealPlanGroceryListItemRequest) (*messages.MealPlanGroceryListItem, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateMealPlanOption(ctx context.Context, request *messages.UpdateMealPlanOptionRequest) (*messages.MealPlanOption, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateMealPlanOptionVote(ctx context.Context, request *messages.UpdateMealPlanOptionVoteRequest) (*messages.MealPlanOptionVote, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateMealPlanTaskStatus(ctx context.Context, request *messages.UpdateMealPlanTaskStatusRequest) (*messages.MealPlanTask, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}
