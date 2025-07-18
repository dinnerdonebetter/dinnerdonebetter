package grpc

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
)

func (s *ServiceImpl) ArchiveMeal(ctx context.Context, request *mealplanning.ArchiveMealRequest) (*mealplanning.ArchiveMealResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) ArchiveMealPlan(ctx context.Context, request *mealplanning.ArchiveMealPlanRequest) (*mealplanning.ArchiveMealPlanResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) ArchiveMealPlanEvent(ctx context.Context, request *mealplanning.ArchiveMealPlanEventRequest) (*mealplanning.ArchiveMealPlanEventResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	if err := s.mealPlanningManager.ArchiveMealPlanEvent(ctx, request.MealPlanID, request.MealPlanEventID); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "failed to archive meal plan event")
	}

	return &mealplanning.ArchiveMealPlanEventResponse{}, nil
}

func (s *ServiceImpl) ArchiveMealPlanGroceryListItem(ctx context.Context, request *mealplanning.ArchiveMealPlanGroceryListItemRequest) (*mealplanning.ArchiveMealPlanGroceryListItemResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) ArchiveMealPlanOption(ctx context.Context, request *mealplanning.ArchiveMealPlanOptionRequest) (*mealplanning.ArchiveMealPlanOptionResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) ArchiveMealPlanOptionVote(ctx context.Context, request *mealplanning.ArchiveMealPlanOptionVoteRequest) (*mealplanning.ArchiveMealPlanOptionVoteResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateMeal(ctx context.Context, request *mealplanning.CreateMealRequest) (*mealplanning.CreateMealResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateMealPlan(ctx context.Context, request *mealplanning.CreateMealPlanRequest) (*mealplanning.CreateMealPlanResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateMealPlanEvent(ctx context.Context, request *mealplanning.CreateMealPlanEventRequest) (*mealplanning.CreateMealPlanEventResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateMealPlanGroceryListItem(ctx context.Context, request *mealplanning.CreateMealPlanGroceryListItemRequest) (*mealplanning.CreateMealPlanGroceryListItemResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateMealPlanOption(ctx context.Context, request *mealplanning.CreateMealPlanOptionRequest) (*mealplanning.CreateMealPlanOptionResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateMealPlanOptionVote(ctx context.Context, request *mealplanning.CreateMealPlanOptionVoteRequest) (*mealplanning.CreateMealPlanOptionVoteResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateMealPlanTask(ctx context.Context, request *mealplanning.CreateMealPlanTaskRequest) (*mealplanning.CreateMealPlanTaskResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) FinalizeMealPlan(ctx context.Context, request *mealplanning.FinalizeMealPlanRequest) (*mealplanning.FinalizeMealPlanResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetMeal(ctx context.Context, request *mealplanning.GetMealRequest) (*mealplanning.GetMealResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetMealPlan(ctx context.Context, request *mealplanning.GetMealPlanRequest) (*mealplanning.GetMealPlanResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetMealPlanEvent(ctx context.Context, request *mealplanning.GetMealPlanEventRequest) (*mealplanning.GetMealPlanEventResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetMealPlanEvents(ctx context.Context, request *mealplanning.GetMealPlanEventsRequest) (*mealplanning.GetMealPlanEventsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetMealPlanGroceryListItem(ctx context.Context, request *mealplanning.GetMealPlanGroceryListItemRequest) (*mealplanning.GetMealPlanGroceryListItemResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetMealPlanGroceryListItemsForMealPlan(ctx context.Context, request *mealplanning.GetMealPlanGroceryListItemsForMealPlanRequest) (*mealplanning.GetMealPlanGroceryListItemsForMealPlanResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetMealPlanOption(ctx context.Context, request *mealplanning.GetMealPlanOptionRequest) (*mealplanning.GetMealPlanOptionResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetMealPlanOptionVote(ctx context.Context, request *mealplanning.GetMealPlanOptionVoteRequest) (*mealplanning.GetMealPlanOptionVoteResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetMealPlanOptionVotes(ctx context.Context, request *mealplanning.GetMealPlanOptionVotesRequest) (*mealplanning.GetMealPlanOptionVotesResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetMealPlanOptions(ctx context.Context, request *mealplanning.GetMealPlanOptionsRequest) (*mealplanning.GetMealPlanOptionsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetMealPlanTask(ctx context.Context, request *mealplanning.GetMealPlanTaskRequest) (*mealplanning.GetMealPlanTaskResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetMealPlanTasks(ctx context.Context, request *mealplanning.GetMealPlanTasksRequest) (*mealplanning.GetMealPlanTasksResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetMeals(ctx context.Context, request *mealplanning.GetMealsRequest) (*mealplanning.GetMealsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) RunFinalizeMealPlanWorker(ctx context.Context, request *mealplanning.RunFinalizeMealPlanWorkerRequest) (*mealplanning.RunFinalizeMealPlanWorkerResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) RunMealPlanGroceryListInitializerWorker(ctx context.Context, request *mealplanning.RunMealPlanGroceryListInitializerWorkerRequest) (*mealplanning.RunMealPlanGroceryListInitializerWorkerResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) RunMealPlanTaskCreatorWorker(ctx context.Context, request *mealplanning.RunMealPlanTaskCreatorWorkerRequest) (*mealplanning.RunMealPlanTaskCreatorWorkerResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) SearchForMeals(ctx context.Context, request *mealplanning.SearchForMealsRequest) (*mealplanning.SearchForMealsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateMealPlan(ctx context.Context, request *mealplanning.UpdateMealPlanRequest) (*mealplanning.UpdateMealPlanResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateMealPlanEvent(ctx context.Context, request *mealplanning.UpdateMealPlanEventRequest) (*mealplanning.UpdateMealPlanEventResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateMealPlanGroceryListItem(ctx context.Context, request *mealplanning.UpdateMealPlanGroceryListItemRequest) (*mealplanning.UpdateMealPlanGroceryListItemResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateMealPlanOption(ctx context.Context, request *mealplanning.UpdateMealPlanOptionRequest) (*mealplanning.UpdateMealPlanOptionResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateMealPlanOptionVote(ctx context.Context, request *mealplanning.UpdateMealPlanOptionVoteRequest) (*mealplanning.UpdateMealPlanOptionVoteResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateMealPlanTaskStatus(ctx context.Context, request *mealplanning.UpdateMealPlanTaskStatusRequest) (*mealplanning.UpdateMealPlanTaskStatusResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) ArchiveUserIngredientPreference(ctx context.Context, request *mealplanning.ArchiveUserIngredientPreferenceRequest) (*mealplanning.ArchiveUserIngredientPreferenceResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateUserIngredientPreference(ctx context.Context, request *mealplanning.CreateUserIngredientPreferenceRequest) (*mealplanning.CreateUserIngredientPreferenceResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetUserIngredientPreferences(ctx context.Context, request *mealplanning.GetUserIngredientPreferencesRequest) (*mealplanning.GetUserIngredientPreferencesResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateUserIngredientPreference(ctx context.Context, request *mealplanning.UpdateUserIngredientPreferenceRequest) (*mealplanning.UpdateUserIngredientPreferenceResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}
