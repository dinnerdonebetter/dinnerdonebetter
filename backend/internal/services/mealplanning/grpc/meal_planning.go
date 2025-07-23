package grpc

import (
	"context"

	grpcconverters "github.com/dinnerdonebetter/backend/internal/grpc/converters"
	"github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/grpc/generated/types"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	converters "github.com/dinnerdonebetter/backend/internal/services/mealplanning/grpc/converters"

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

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to get session context data")
	}

	if err = s.mealPlanningManager.ArchiveMealPlan(ctx, request.MealPlanID, sessionContextData.GetActiveAccountID()); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to archive meal plan")
	}

	x := &mealplanning.ArchiveMealPlanResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
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

	x := &mealplanning.ArchiveMealPlanEventResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *ServiceImpl) ArchiveMealPlanGroceryListItem(ctx context.Context, request *mealplanning.ArchiveMealPlanGroceryListItemRequest) (*mealplanning.ArchiveMealPlanGroceryListItemResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey:                request.MealPlanID,
		keys.MealPlanGroceryListItemIDKey: request.MealPlanGroceryListItemID,
	}, span, s.logger)

	if err := s.mealPlanningManager.ArchiveMealPlanGroceryListItem(ctx, request.MealPlanID, request.MealPlanGroceryListItemID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to archive meal plan grocery list item")
	}

	x := &mealplanning.ArchiveMealPlanGroceryListItemResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *ServiceImpl) ArchiveMealPlanOption(ctx context.Context, request *mealplanning.ArchiveMealPlanOptionRequest) (*mealplanning.ArchiveMealPlanOptionResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey:       request.MealPlanID,
		keys.MealPlanEventIDKey:  request.MealPlanEventID,
		keys.MealPlanOptionIDKey: request.MealPlanOptionID,
	}, span, s.logger)

	if err := s.mealPlanningManager.ArchiveMealPlanOption(ctx, request.MealPlanID, request.MealPlanEventID, request.MealPlanOptionID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to archive meal plan option")
	}

	x := &mealplanning.ArchiveMealPlanOptionResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
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

	if err := s.mealPlanningManager.ArchiveMealPlanOptionVote(ctx, request.MealPlanID, request.MealPlanEventID, request.MealPlanOptionID, request.MealPlanOptionVoteID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to archive meal plan option vote")
	}

	x := &mealplanning.ArchiveMealPlanOptionVoteResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *ServiceImpl) ArchiveUserIngredientPreference(ctx context.Context, request *mealplanning.ArchiveUserIngredientPreferenceRequest) (*mealplanning.ArchiveUserIngredientPreferenceResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.UserIngredientPreferenceIDKey: request.UserIngredientPreferenceID,
	}, span, s.logger)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to get session context data")
	}

	if err = s.mealPlanningManager.ArchiveUserIngredientPreference(ctx, sessionContextData.GetUserID(), request.UserIngredientPreferenceID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to archive meal plan option vote")
	}

	x := &mealplanning.ArchiveUserIngredientPreferenceResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
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

func (s *ServiceImpl) CreateUserIngredientPreference(ctx context.Context, request *mealplanning.CreateUserIngredientPreferenceRequest) (*mealplanning.CreateUserIngredientPreferenceResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) FinalizeMealPlan(ctx context.Context, request *mealplanning.FinalizeMealPlanRequest) (*mealplanning.FinalizeMealPlanResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey: request.MealPlanID,
	}, span, s.logger)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to get session context data")
	}

	finalized, err := s.mealPlanningManager.FinalizeMealPlan(ctx, request.MealPlanID, sessionContextData.GetActiveAccountID())
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to finalize meal plan")
	}

	x := &mealplanning.FinalizeMealPlanResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Finalized: finalized,
	}

	return x, nil
}

func (s *ServiceImpl) GetMeal(ctx context.Context, request *mealplanning.GetMealRequest) (*mealplanning.GetMealResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealIDKey: request.MealID,
	}, span, s.logger)

	meal, err := s.mealPlanningManager.ReadMeal(ctx, request.MealID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to read meal")
	}

	x := &mealplanning.GetMealResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Result: converters.ConvertMealToGRPCMeal(meal),
	}

	return x, nil
}

func (s *ServiceImpl) GetMealPlan(ctx context.Context, request *mealplanning.GetMealPlanRequest) (*mealplanning.GetMealPlanResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey: request.MealPlanID,
	}, span, s.logger)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to get session context data")
	}

	mealPlan, err := s.mealPlanningManager.ReadMealPlan(ctx, request.MealPlanID, sessionContextData.GetActiveAccountID())
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to read meal plan")
	}

	x := &mealplanning.GetMealPlanResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Result: converters.ConvertMealPlanToGRPCMealPlan(mealPlan),
	}

	return x, nil
}

func (s *ServiceImpl) GetMealPlanEvent(ctx context.Context, request *mealplanning.GetMealPlanEventRequest) (*mealplanning.GetMealPlanEventResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey:      request.MealPlanID,
		keys.MealPlanEventIDKey: request.MealPlanEventID,
	}, span, s.logger)

	mealPlanEvent, err := s.mealPlanningManager.ReadMealPlanEvent(ctx, request.MealPlanID, request.MealPlanEventID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to read meal plan")
	}

	x := &mealplanning.GetMealPlanEventResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Result: converters.ConvertMealPlanEventToGRPCMealPlanEvent(mealPlanEvent),
	}

	return x, nil
}

func (s *ServiceImpl) GetMealPlanEvents(ctx context.Context, request *mealplanning.GetMealPlanEventsRequest) (*mealplanning.GetMealPlanEventsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey: request.MealPlanID,
	}, span, s.logger)

	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)

	mealPlanEvents, _, err := s.mealPlanningManager.ListMealPlanEvents(ctx, request.MealPlanID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to fetch list of meal plan events")
	}

	x := &mealplanning.GetMealPlanEventsResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	for _, mealPlanEvent := range mealPlanEvents {
		x.Results = append(x.Results, converters.ConvertMealPlanEventToGRPCMealPlanEvent(mealPlanEvent))
	}

	return x, nil
}

func (s *ServiceImpl) GetMealPlanGroceryListItem(ctx context.Context, request *mealplanning.GetMealPlanGroceryListItemRequest) (*mealplanning.GetMealPlanGroceryListItemResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey:                request.MealPlanID,
		keys.MealPlanGroceryListItemIDKey: request.MealPlanGroceryListItemID,
	}, span, s.logger)

	mealPlanGroceryListItem, err := s.mealPlanningManager.ReadMealPlanGroceryListItem(ctx, request.MealPlanID, request.MealPlanGroceryListItemID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to read meal plan grocery list item")
	}

	x := &mealplanning.GetMealPlanGroceryListItemResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Result: converters.ConvertMealPlanGroceryListItemToGRPCMealPlanGroceryListItem(mealPlanGroceryListItem),
	}

	return x, nil
}

func (s *ServiceImpl) GetMealPlanGroceryListItemsForMealPlan(ctx context.Context, request *mealplanning.GetMealPlanGroceryListItemsForMealPlanRequest) (*mealplanning.GetMealPlanGroceryListItemsForMealPlanResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey: request.MealPlanID,
	}, span, s.logger)

	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)

	mealPlanGroceryListItems, _, err := s.mealPlanningManager.ListMealPlanGroceryListItemsByMealPlan(ctx, request.MealPlanID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to fetch list of meal plan grocery list items")
	}

	x := &mealplanning.GetMealPlanGroceryListItemsForMealPlanResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	for _, mealPlanGroceryListItem := range mealPlanGroceryListItems {
		x.Results = append(x.Results, converters.ConvertMealPlanGroceryListItemToGRPCMealPlanGroceryListItem(mealPlanGroceryListItem))
	}

	return x, nil
}

func (s *ServiceImpl) GetMealPlanOption(ctx context.Context, request *mealplanning.GetMealPlanOptionRequest) (*mealplanning.GetMealPlanOptionResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey:       request.MealPlanID,
		keys.MealPlanEventIDKey:  request.MealPlanEventID,
		keys.MealPlanOptionIDKey: request.MealPlanOptionID,
	}, span, s.logger)

	mealPlanOption, err := s.mealPlanningManager.ReadMealPlanOption(ctx, request.MealPlanID, request.MealPlanEventID, request.MealPlanOptionID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to read meal plan grocery list item")
	}

	x := &mealplanning.GetMealPlanOptionResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Result: converters.ConvertMealPlanOptionToGRPCMealPlanOption(mealPlanOption),
	}

	return x, nil
}

func (s *ServiceImpl) GetMealPlanOptionVote(ctx context.Context, request *mealplanning.GetMealPlanOptionVoteRequest) (*mealplanning.GetMealPlanOptionVoteResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey:           request.MealPlanID,
		keys.MealPlanOptionIDKey:     request.MealPlanOptionID,
		keys.MealPlanEventIDKey:      request.MealPlanEventID,
		keys.MealPlanOptionVoteIDKey: request.MealPlanOptionVoteID,
	}, span, s.logger)

	mealPlanOptionVote, err := s.mealPlanningManager.ReadMealPlanOptionVote(ctx, request.MealPlanID, request.MealPlanEventID, request.MealPlanOptionID, request.MealPlanOptionVoteID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to read meal plan grocery list item")
	}

	x := &mealplanning.GetMealPlanOptionVoteResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Result: converters.ConvertMealPlanOptionVoteToGRPCMealPlanOptionVote(mealPlanOptionVote),
	}

	return x, nil
}

func (s *ServiceImpl) GetMealPlanOptionVotes(ctx context.Context, request *mealplanning.GetMealPlanOptionVotesRequest) (*mealplanning.GetMealPlanOptionVotesResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey:       request.MealPlanID,
		keys.MealPlanOptionIDKey: request.MealPlanOptionID,
		keys.MealPlanEventIDKey:  request.MealPlanEventID,
	}, span, s.logger)

	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)

	mealPlanOptionVotes, _, err := s.mealPlanningManager.ListMealPlanOptionVotes(ctx, request.MealPlanID, request.MealPlanEventID, request.MealPlanOptionID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to fetch list of meal plan option votes")
	}

	x := &mealplanning.GetMealPlanOptionVotesResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	for _, mealPlanOptionVote := range mealPlanOptionVotes {
		x.Results = append(x.Results, converters.ConvertMealPlanOptionVoteToGRPCMealPlanOptionVote(mealPlanOptionVote))
	}

	return x, nil
}

func (s *ServiceImpl) GetMealPlanOptions(ctx context.Context, request *mealplanning.GetMealPlanOptionsRequest) (*mealplanning.GetMealPlanOptionsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey: request.MealPlanID,
	}, span, s.logger)

	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)

	mealPlanOptions, _, err := s.mealPlanningManager.ListMealPlanOptions(ctx, request.MealPlanID, request.MealPlanEventID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to fetch list of meal plan options")
	}

	x := &mealplanning.GetMealPlanOptionsResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	for _, mealPlanOption := range mealPlanOptions {
		x.Results = append(x.Results, converters.ConvertMealPlanOptionToGRPCMealPlanOption(mealPlanOption))
	}

	return x, nil
}

func (s *ServiceImpl) GetMealPlanTask(ctx context.Context, request *mealplanning.GetMealPlanTaskRequest) (*mealplanning.GetMealPlanTaskResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey:     request.MealPlanID,
		keys.MealPlanTaskIDKey: request.MealPlanTaskID,
	}, span, s.logger)

	mealPlanTask, err := s.mealPlanningManager.ReadMealPlanTask(ctx, request.MealPlanID, request.MealPlanTaskID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to read meal plan grocery list item")
	}

	x := &mealplanning.GetMealPlanTaskResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Result: converters.ConvertMealPlanTaskToGRPCMealPlanTask(mealPlanTask),
	}

	return x, nil
}

func (s *ServiceImpl) GetMealPlanTasks(ctx context.Context, request *mealplanning.GetMealPlanTasksRequest) (*mealplanning.GetMealPlanTasksResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey: request.MealPlanID,
	}, span, s.logger)

	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)

	mealPlanTasks, _, err := s.mealPlanningManager.ListMealPlanTasksByMealPlan(ctx, request.MealPlanID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to fetch list of meal plan tasks")
	}

	x := &mealplanning.GetMealPlanTasksResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	for _, mealPlanTask := range mealPlanTasks {
		x.Results = append(x.Results, converters.ConvertMealPlanTaskToGRPCMealPlanTask(mealPlanTask))
	}

	return x, nil
}

func (s *ServiceImpl) GetMeals(ctx context.Context, request *mealplanning.GetMealsRequest) (*mealplanning.GetMealsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)

	meals, _, err := s.mealPlanningManager.ListMeals(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to fetch list of meals")
	}

	x := &mealplanning.GetMealsResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	for _, meal := range meals {
		x.Results = append(x.Results, converters.ConvertMealToGRPCMeal(meal))
	}

	return x, nil
}

func (s *ServiceImpl) GetUserIngredientPreferences(ctx context.Context, request *mealplanning.GetUserIngredientPreferencesRequest) (*mealplanning.GetUserIngredientPreferencesResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to get session context data")
	}

	userIngredientPreferences, _, err := s.mealPlanningManager.ListUserIngredientPreferences(ctx, sessionContextData.GetUserID(), filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to fetch list of meals")
	}

	x := &mealplanning.GetUserIngredientPreferencesResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	for _, userIngredientPreference := range userIngredientPreferences {
		x.Results = append(x.Results, converters.ConvertUserIngredientPreferenceToGRPCUserIngredientPreference(userIngredientPreference))
	}

	return x, nil
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

func (s *ServiceImpl) UpdateUserIngredientPreference(ctx context.Context, request *mealplanning.UpdateUserIngredientPreferenceRequest) (*mealplanning.UpdateUserIngredientPreferenceResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.UserIngredientPreferenceIDKey: request.UserIngredientPreferenceID,
	}, span, s.logger)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}
