package grpc

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/grpc/generated/types"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"

	"google.golang.org/grpc/codes"
)

func (s *ServiceImpl) ArchiveMeal(ctx context.Context, request *mealplanning.ArchiveMealRequest) (*mealplanning.ArchiveMealResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealIDKey: request.MealID,
	}, span, s.logger)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to fetch session context data")
	}

	if err = s.mealPlanningManager.ArchiveMeal(ctx, request.MealID, sessionContextData.GetUserID()); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to archive meal")
	}

	x := &mealplanning.ArchiveMealResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *ServiceImpl) ArchiveMealPlan(ctx context.Context, request *mealplanning.ArchiveMealPlanRequest) (*mealplanning.ArchiveMealPlanResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey: request.MealPlanID,
	}, span, s.logger)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) ArchiveMealPlanEvent(ctx context.Context, request *mealplanning.ArchiveMealPlanEventRequest) (*mealplanning.ArchiveMealPlanEventResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey:      request.MealPlanID,
		keys.MealPlanEventIDKey: request.MealPlanEventID,
	}, span, s.logger)

	if err := s.mealPlanningManager.ArchiveMealPlanEvent(ctx, request.MealPlanID, request.MealPlanEventID); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "failed to archive meal plan event")
	}

	return &mealplanning.ArchiveMealPlanEventResponse{}, nil
}

func (s *ServiceImpl) ArchiveMealPlanGroceryListItem(ctx context.Context, request *mealplanning.ArchiveMealPlanGroceryListItemRequest) (*mealplanning.ArchiveMealPlanGroceryListItemResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey:                request.MealPlanID,
		keys.MealPlanGroceryListItemIDKey: request.MealPlanGroceryListItemID,
	}, span, s.logger)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) ArchiveMealPlanOption(ctx context.Context, request *mealplanning.ArchiveMealPlanOptionRequest) (*mealplanning.ArchiveMealPlanOptionResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey:       request.MealPlanID,
		keys.MealPlanEventIDKey:  request.MealPlanEventID,
		keys.MealPlanOptionIDKey: request.MealPlanOptionID,
	}, span, s.logger)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) ArchiveMealPlanOptionVote(ctx context.Context, request *mealplanning.ArchiveMealPlanOptionVoteRequest) (*mealplanning.ArchiveMealPlanOptionVoteResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanOptionVoteIDKey: request.MealPlanOptionVoteID,
		keys.MealPlanOptionIDKey:     request.MealPlanOptionID,
		keys.MealPlanEventIDKey:      request.MealPlanEventID,
		keys.MealPlanIDKey:           request.MealPlanID,
	}, span, s.logger)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateMeal(ctx context.Context, request *mealplanning.CreateMealRequest) (*mealplanning.CreateMealResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateMealPlan(ctx context.Context, request *mealplanning.CreateMealPlanRequest) (*mealplanning.CreateMealPlanResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateMealPlanEvent(ctx context.Context, request *mealplanning.CreateMealPlanEventRequest) (*mealplanning.CreateMealPlanEventResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey: request.MealPlanID,
	}, span, s.logger)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateMealPlanGroceryListItem(ctx context.Context, request *mealplanning.CreateMealPlanGroceryListItemRequest) (*mealplanning.CreateMealPlanGroceryListItemResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey: request.MealPlanID,
	}, span, s.logger)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateMealPlanOption(ctx context.Context, request *mealplanning.CreateMealPlanOptionRequest) (*mealplanning.CreateMealPlanOptionResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey: request.MealPlanID,
	}, span, s.logger)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateMealPlanOptionVote(ctx context.Context, request *mealplanning.CreateMealPlanOptionVoteRequest) (*mealplanning.CreateMealPlanOptionVoteResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey: request.MealPlanID,
	}, span, s.logger)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateMealPlanTask(ctx context.Context, request *mealplanning.CreateMealPlanTaskRequest) (*mealplanning.CreateMealPlanTaskResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey: request.MealPlanID,
	}, span, s.logger)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) FinalizeMealPlan(ctx context.Context, request *mealplanning.FinalizeMealPlanRequest) (*mealplanning.FinalizeMealPlanResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey: request.MealPlanID,
	}, span, s.logger)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetMeal(ctx context.Context, request *mealplanning.GetMealRequest) (*mealplanning.GetMealResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealIDKey: request.MealID,
	}, span, s.logger)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetMealPlan(ctx context.Context, request *mealplanning.GetMealPlanRequest) (*mealplanning.GetMealPlanResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey: request.MealPlanID,
	}, span, s.logger)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetMealPlanEvent(ctx context.Context, request *mealplanning.GetMealPlanEventRequest) (*mealplanning.GetMealPlanEventResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey:      request.MealPlanID,
		keys.MealPlanEventIDKey: request.MealPlanEventID,
	}, span, s.logger)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetMealPlanEvents(ctx context.Context, request *mealplanning.GetMealPlanEventsRequest) (*mealplanning.GetMealPlanEventsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey: request.MealPlanID,
	}, span, s.logger)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetMealPlanGroceryListItem(ctx context.Context, request *mealplanning.GetMealPlanGroceryListItemRequest) (*mealplanning.GetMealPlanGroceryListItemResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey:                request.MealPlanID,
		keys.MealPlanGroceryListItemIDKey: request.MealPlanGroceryListItemID,
	}, span, s.logger)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetMealPlanGroceryListItemsForMealPlan(ctx context.Context, request *mealplanning.GetMealPlanGroceryListItemsForMealPlanRequest) (*mealplanning.GetMealPlanGroceryListItemsForMealPlanResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey: request.MealPlanID,
	}, span, s.logger)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetMealPlanOption(ctx context.Context, request *mealplanning.GetMealPlanOptionRequest) (*mealplanning.GetMealPlanOptionResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey:       request.MealPlanID,
		keys.MealPlanOptionIDKey: request.MealPlanOptionID,
	}, span, s.logger)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetMealPlanOptionVote(ctx context.Context, request *mealplanning.GetMealPlanOptionVoteRequest) (*mealplanning.GetMealPlanOptionVoteResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey:           request.MealPlanID,
		keys.MealPlanOptionIDKey:     request.MealPlanOptionID,
		keys.MealPlanOptionVoteIDKey: request.MealPlanOptionVoteID,
	}, span, s.logger)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetMealPlanOptionVotes(ctx context.Context, request *mealplanning.GetMealPlanOptionVotesRequest) (*mealplanning.GetMealPlanOptionVotesResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey:       request.MealPlanID,
		keys.MealPlanOptionIDKey: request.MealPlanOptionID,
	}, span, s.logger)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetMealPlanOptions(ctx context.Context, request *mealplanning.GetMealPlanOptionsRequest) (*mealplanning.GetMealPlanOptionsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey: request.MealPlanID,
	}, span, s.logger)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetMealPlanTask(ctx context.Context, request *mealplanning.GetMealPlanTaskRequest) (*mealplanning.GetMealPlanTaskResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey:     request.MealPlanID,
		keys.MealPlanTaskIDKey: request.MealPlanTaskID,
	}, span, s.logger)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetMealPlanTasks(ctx context.Context, request *mealplanning.GetMealPlanTasksRequest) (*mealplanning.GetMealPlanTasksResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey: request.MealPlanID,
	}, span, s.logger)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetMeals(ctx context.Context, request *mealplanning.GetMealsRequest) (*mealplanning.GetMealsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) RunFinalizeMealPlanWorker(ctx context.Context, request *mealplanning.RunFinalizeMealPlanWorkerRequest) (*mealplanning.RunFinalizeMealPlanWorkerResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) RunMealPlanGroceryListInitializerWorker(ctx context.Context, request *mealplanning.RunMealPlanGroceryListInitializerWorkerRequest) (*mealplanning.RunMealPlanGroceryListInitializerWorkerResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) RunMealPlanTaskCreatorWorker(ctx context.Context, request *mealplanning.RunMealPlanTaskCreatorWorkerRequest) (*mealplanning.RunMealPlanTaskCreatorWorkerResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) SearchForMeals(ctx context.Context, request *mealplanning.SearchForMealsRequest) (*mealplanning.SearchForMealsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateMealPlan(ctx context.Context, request *mealplanning.UpdateMealPlanRequest) (*mealplanning.UpdateMealPlanResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey: request.MealPlanID,
	}, span, s.logger)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateMealPlanEvent(ctx context.Context, request *mealplanning.UpdateMealPlanEventRequest) (*mealplanning.UpdateMealPlanEventResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey:      request.MealPlanID,
		keys.MealPlanEventIDKey: request.MealPlanEventID,
	}, span, s.logger)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateMealPlanGroceryListItem(ctx context.Context, request *mealplanning.UpdateMealPlanGroceryListItemRequest) (*mealplanning.UpdateMealPlanGroceryListItemResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey:                request.MealPlanID,
		keys.MealPlanGroceryListItemIDKey: request.MealPlanGroceryListItemID,
	}, span, s.logger)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateMealPlanOption(ctx context.Context, request *mealplanning.UpdateMealPlanOptionRequest) (*mealplanning.UpdateMealPlanOptionResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey:       request.MealPlanID,
		keys.MealPlanOptionIDKey: request.MealPlanOptionID,
	}, span, s.logger)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateMealPlanOptionVote(ctx context.Context, request *mealplanning.UpdateMealPlanOptionVoteRequest) (*mealplanning.UpdateMealPlanOptionVoteResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey:           request.MealPlanID,
		keys.MealPlanOptionIDKey:     request.MealPlanOptionID,
		keys.MealPlanOptionVoteIDKey: request.MealPlanOptionVoteID,
	}, span, s.logger)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateMealPlanTaskStatus(ctx context.Context, request *mealplanning.UpdateMealPlanTaskStatusRequest) (*mealplanning.UpdateMealPlanTaskStatusResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey:     request.MealPlanID,
		keys.MealPlanTaskIDKey: request.MealPlanTaskID,
	}, span, s.logger)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) ArchiveUserIngredientPreference(ctx context.Context, request *mealplanning.ArchiveUserIngredientPreferenceRequest) (*mealplanning.ArchiveUserIngredientPreferenceResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.UserIngredientPreferenceIDKey: request.UserIngredientPreferenceID,
	}, span, s.logger)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateUserIngredientPreference(ctx context.Context, request *mealplanning.CreateUserIngredientPreferenceRequest) (*mealplanning.CreateUserIngredientPreferenceResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetUserIngredientPreferences(ctx context.Context, request *mealplanning.GetUserIngredientPreferencesRequest) (*mealplanning.GetUserIngredientPreferencesResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateUserIngredientPreference(ctx context.Context, request *mealplanning.UpdateUserIngredientPreferenceRequest) (*mealplanning.UpdateUserIngredientPreferenceResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.UserIngredientPreferenceIDKey: request.UserIngredientPreferenceID,
	}, span, s.logger)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}
