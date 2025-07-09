package eatinggrpc

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
)

func (s *ServiceImpl) ArchiveMeal(ctx context.Context, request *messages.ArchiveMealRequest) (*messages.ArchiveMealResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) ArchiveMealPlan(ctx context.Context, request *messages.ArchiveMealPlanRequest) (*messages.ArchiveMealPlanResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) ArchiveMealPlanEvent(ctx context.Context, request *messages.ArchiveMealPlanEventRequest) (*messages.ArchiveMealPlanEventResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	if err := s.mealPlanningManager.ArchiveMealPlanEvent(ctx, request.MealPlanID, request.MealPlanEventID); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "failed to archive meal plan event")
	}

	return &messages.ArchiveMealPlanEventResponse{}, nil
}

func (s *ServiceImpl) ArchiveMealPlanGroceryListItem(ctx context.Context, request *messages.ArchiveMealPlanGroceryListItemRequest) (*messages.ArchiveMealPlanGroceryListItemResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) ArchiveMealPlanOption(ctx context.Context, request *messages.ArchiveMealPlanOptionRequest) (*messages.ArchiveMealPlanOptionResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) ArchiveMealPlanOptionVote(ctx context.Context, request *messages.ArchiveMealPlanOptionVoteRequest) (*messages.ArchiveMealPlanOptionVoteResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateMeal(ctx context.Context, request *messages.CreateMealRequest) (*messages.CreateMealResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateMealPlan(ctx context.Context, request *messages.CreateMealPlanRequest) (*messages.CreateMealPlanResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateMealPlanEvent(ctx context.Context, request *messages.CreateMealPlanEventRequest) (*messages.CreateMealPlanEventResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateMealPlanGroceryListItem(ctx context.Context, request *messages.CreateMealPlanGroceryListItemRequest) (*messages.CreateMealPlanGroceryListItemResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateMealPlanOption(ctx context.Context, request *messages.CreateMealPlanOptionRequest) (*messages.CreateMealPlanOptionResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateMealPlanOptionVote(ctx context.Context, request *messages.CreateMealPlanOptionVoteRequest) (*messages.CreateMealPlanOptionVoteResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateMealPlanTask(ctx context.Context, request *messages.CreateMealPlanTaskRequest) (*messages.CreateMealPlanTaskResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) FinalizeMealPlan(ctx context.Context, request *messages.FinalizeMealPlanRequest) (*messages.FinalizeMealPlanResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetMeal(ctx context.Context, request *messages.GetMealRequest) (*messages.GetMealResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetMealPlan(ctx context.Context, request *messages.GetMealPlanRequest) (*messages.GetMealPlanResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetMealPlanEvent(ctx context.Context, request *messages.GetMealPlanEventRequest) (*messages.GetMealPlanEventResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetMealPlanEvents(ctx context.Context, request *messages.GetMealPlanEventsRequest) (*messages.GetMealPlanEventsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetMealPlanGroceryListItem(ctx context.Context, request *messages.GetMealPlanGroceryListItemRequest) (*messages.GetMealPlanGroceryListItemResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetMealPlanGroceryListItemsForMealPlan(ctx context.Context, request *messages.GetMealPlanGroceryListItemsForMealPlanRequest) (*messages.GetMealPlanGroceryListItemsForMealPlanResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetMealPlanOption(ctx context.Context, request *messages.GetMealPlanOptionRequest) (*messages.GetMealPlanOptionResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetMealPlanOptionVote(ctx context.Context, request *messages.GetMealPlanOptionVoteRequest) (*messages.GetMealPlanOptionVoteResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetMealPlanOptionVotes(ctx context.Context, request *messages.GetMealPlanOptionVotesRequest) (*messages.GetMealPlanOptionVotesResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetMealPlanOptions(ctx context.Context, request *messages.GetMealPlanOptionsRequest) (*messages.GetMealPlanOptionsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetMealPlanTask(ctx context.Context, request *messages.GetMealPlanTaskRequest) (*messages.GetMealPlanTaskResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetMealPlanTasks(ctx context.Context, request *messages.GetMealPlanTasksRequest) (*messages.GetMealPlanTasksResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetMealPlansForHousehold(ctx context.Context, request *messages.GetMealPlansForHouseholdRequest) (*messages.GetMealPlansForHouseholdResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetMeals(ctx context.Context, request *messages.GetMealsRequest) (*messages.GetMealsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) RunFinalizeMealPlanWorker(ctx context.Context, request *messages.RunFinalizeMealPlanWorkerRequest) (*messages.RunFinalizeMealPlanWorkerResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) RunMealPlanGroceryListInitializerWorker(ctx context.Context, request *messages.RunMealPlanGroceryListInitializerWorkerRequest) (*messages.RunMealPlanGroceryListInitializerWorkerResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) RunMealPlanTaskCreatorWorker(ctx context.Context, request *messages.RunMealPlanTaskCreatorWorkerRequest) (*messages.RunMealPlanTaskCreatorWorkerResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) SearchForMeals(ctx context.Context, request *messages.SearchForMealsRequest) (*messages.SearchForMealsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateMealPlan(ctx context.Context, request *messages.UpdateMealPlanRequest) (*messages.UpdateMealPlanResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateMealPlanEvent(ctx context.Context, request *messages.UpdateMealPlanEventRequest) (*messages.UpdateMealPlanEventResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateMealPlanGroceryListItem(ctx context.Context, request *messages.UpdateMealPlanGroceryListItemRequest) (*messages.UpdateMealPlanGroceryListItemResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateMealPlanOption(ctx context.Context, request *messages.UpdateMealPlanOptionRequest) (*messages.UpdateMealPlanOptionResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateMealPlanOptionVote(ctx context.Context, request *messages.UpdateMealPlanOptionVoteRequest) (*messages.UpdateMealPlanOptionVoteResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateMealPlanTaskStatus(ctx context.Context, request *messages.UpdateMealPlanTaskStatusRequest) (*messages.UpdateMealPlanTaskStatusResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}
