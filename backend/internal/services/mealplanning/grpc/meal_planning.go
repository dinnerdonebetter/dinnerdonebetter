package grpc

import (
	"context"

	grpcconverters "github.com/dinnerdonebetter/backend/internal/grpc/converters"
	mealplanningsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/grpc/generated/types"
	"github.com/dinnerdonebetter/backend/internal/platform/internalerrors"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	converters "github.com/dinnerdonebetter/backend/internal/services/mealplanning/grpc/converters"

	"google.golang.org/grpc/codes"
)

func (s *serviceImpl) ArchiveMeal(ctx context.Context, request *mealplanningsvc.ArchiveMealRequest) (*mealplanningsvc.ArchiveMealResponse, error) {
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

	x := &mealplanningsvc.ArchiveMealResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *serviceImpl) ArchiveMealPlan(ctx context.Context, request *mealplanningsvc.ArchiveMealPlanRequest) (*mealplanningsvc.ArchiveMealPlanResponse, error) {
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

	x := &mealplanningsvc.ArchiveMealPlanResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *serviceImpl) ArchiveMealPlanEvent(ctx context.Context, request *mealplanningsvc.ArchiveMealPlanEventRequest) (*mealplanningsvc.ArchiveMealPlanEventResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey:      request.MealPlanID,
		keys.MealPlanEventIDKey: request.MealPlanEventID,
	}, span, s.logger)

	if err := s.mealPlanningManager.ArchiveMealPlanEvent(ctx, request.MealPlanID, request.MealPlanEventID); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "failed to archive meal plan event")
	}

	x := &mealplanningsvc.ArchiveMealPlanEventResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *serviceImpl) ArchiveMealPlanGroceryListItem(ctx context.Context, request *mealplanningsvc.ArchiveMealPlanGroceryListItemRequest) (*mealplanningsvc.ArchiveMealPlanGroceryListItemResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey:                request.MealPlanID,
		keys.MealPlanGroceryListItemIDKey: request.MealPlanGroceryListItemID,
	}, span, s.logger)

	if err := s.mealPlanningManager.ArchiveMealPlanGroceryListItem(ctx, request.MealPlanID, request.MealPlanGroceryListItemID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to archive meal plan grocery list item")
	}

	x := &mealplanningsvc.ArchiveMealPlanGroceryListItemResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *serviceImpl) ArchiveMealPlanOption(ctx context.Context, request *mealplanningsvc.ArchiveMealPlanOptionRequest) (*mealplanningsvc.ArchiveMealPlanOptionResponse, error) {
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

	x := &mealplanningsvc.ArchiveMealPlanOptionResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *serviceImpl) ArchiveMealPlanOptionVote(ctx context.Context, request *mealplanningsvc.ArchiveMealPlanOptionVoteRequest) (*mealplanningsvc.ArchiveMealPlanOptionVoteResponse, error) {
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

	x := &mealplanningsvc.ArchiveMealPlanOptionVoteResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *serviceImpl) ArchiveUserIngredientPreference(ctx context.Context, request *mealplanningsvc.ArchiveUserIngredientPreferenceRequest) (*mealplanningsvc.ArchiveUserIngredientPreferenceResponse, error) {
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
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to archive user ingredient preference")
	}

	x := &mealplanningsvc.ArchiveUserIngredientPreferenceResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *serviceImpl) CreateMeal(ctx context.Context, request *mealplanningsvc.CreateMealRequest) (*mealplanningsvc.CreateMealResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to get session context data")
	}

	if request.Input == nil {
		return nil, internalerrors.ErrEmptyInputParameter
	}

	input := converters.ConvertGRPCMealCreationRequestInputToMealCreationRequestInput(request.Input)

	created, err := s.mealPlanningManager.CreateMeal(ctx, sessionContextData.GetUserID(), input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to create meal")
	}

	x := &mealplanningsvc.CreateMealResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Created: converters.ConvertMealToGRPCMeal(created),
	}

	return x, nil
}

func (s *serviceImpl) CreateMealPlan(ctx context.Context, request *mealplanningsvc.CreateMealPlanRequest) (*mealplanningsvc.CreateMealPlanResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to get session context data")
	}

	input := converters.ConvertGRPCMealPlanCreationRequestInputToMealPlanCreationRequestInput(request.Input)

	created, err := s.mealPlanningManager.CreateMealPlan(ctx, sessionContextData.GetActiveAccountID(), sessionContextData.GetUserID(), input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to create meal plan")
	}

	x := &mealplanningsvc.CreateMealPlanResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Created: converters.ConvertMealPlanToGRPCMealPlan(created),
	}

	return x, nil
}

func (s *serviceImpl) CreateMealPlanEvent(ctx context.Context, request *mealplanningsvc.CreateMealPlanEventRequest) (*mealplanningsvc.CreateMealPlanEventResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey: request.MealPlanID,
	}, span, s.logger)

	input := converters.ConvertGRPCMealPlanEventCreationRequestInputToMealPlanEventCreationRequestInput(request.Input)

	created, err := s.mealPlanningManager.CreateMealPlanEvent(ctx, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to create meal plan event")
	}

	x := &mealplanningsvc.CreateMealPlanEventResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Created: converters.ConvertMealPlanEventToGRPCMealPlanEvent(created),
	}

	return x, nil
}

func (s *serviceImpl) CreateMealPlanGroceryListItem(ctx context.Context, request *mealplanningsvc.CreateMealPlanGroceryListItemRequest) (*mealplanningsvc.CreateMealPlanGroceryListItemResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey: request.MealPlanID,
	}, span, s.logger)

	input := converters.ConvertGRPCMealPlanGroceryListItemCreationRequestInputToMealPlanGroceryListItemCreationRequestInput(request.Input)

	created, err := s.mealPlanningManager.CreateMealPlanGroceryListItem(ctx, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to create meal plan grocery list item")
	}

	x := &mealplanningsvc.CreateMealPlanGroceryListItemResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Created: converters.ConvertMealPlanGroceryListItemToGRPCMealPlanGroceryListItem(created),
	}

	return x, nil
}

func (s *serviceImpl) CreateMealPlanOption(ctx context.Context, request *mealplanningsvc.CreateMealPlanOptionRequest) (*mealplanningsvc.CreateMealPlanOptionResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey: request.MealPlanID,
	}, span, s.logger)

	input := converters.ConvertGRPCMealPlanOptionCreationRequestInputToMealPlanOptionCreationRequestInput(request.Input)

	created, err := s.mealPlanningManager.CreateMealPlanOptionWithEventID(ctx, request.MealPlanEventID, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to create meal plan option")
	}

	x := &mealplanningsvc.CreateMealPlanOptionResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Created: converters.ConvertMealPlanOptionToGRPCMealPlanOption(created),
	}

	return x, nil
}

func (s *serviceImpl) CreateMealPlanOptionVote(ctx context.Context, request *mealplanningsvc.CreateMealPlanOptionVoteRequest) (*mealplanningsvc.CreateMealPlanOptionVoteResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey: request.MealPlanID,
	}, span, s.logger)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to get session context data")
	}

	input := converters.ConvertGRPCMealPlanOptionVoteCreationRequestInputToMealPlanOptionVoteCreationRequestInput(request.Input)
	for i := range input.Votes {
		input.Votes[i].ByUser = sessionContextData.GetUserID()
	}

	created, err := s.mealPlanningManager.CreateMealPlanOptionVotes(ctx, sessionContextData.GetUserID(), input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to create meal plan option vote")
	}

	x := &mealplanningsvc.CreateMealPlanOptionVoteResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	for _, creation := range created {
		x.Created = append(x.Created, converters.ConvertMealPlanOptionVoteToGRPCMealPlanOptionVote(creation))
	}

	return x, nil
}

func (s *serviceImpl) CreateMealPlanTask(ctx context.Context, request *mealplanningsvc.CreateMealPlanTaskRequest) (*mealplanningsvc.CreateMealPlanTaskResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey: request.MealPlanID,
	}, span, s.logger)

	input := converters.ConvertGRPCMealPlanTaskCreationRequestInputToMealPlanTaskCreationRequestInput(request.Input)

	created, err := s.mealPlanningManager.CreateMealPlanTask(ctx, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to create meal plan task")
	}

	x := &mealplanningsvc.CreateMealPlanTaskResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Created: converters.ConvertMealPlanTaskToGRPCMealPlanTask(created),
	}

	return x, nil
}

func (s *serviceImpl) CreateUserIngredientPreference(ctx context.Context, request *mealplanningsvc.CreateUserIngredientPreferenceRequest) (*mealplanningsvc.CreateUserIngredientPreferenceResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to get session context data")
	}

	input := converters.ConvertGRPCUserIngredientPreferenceCreationRequestInputToUserIngredientPreferenceCreationRequestInput(request.Input)

	created, err := s.mealPlanningManager.CreateUserIngredientPreference(ctx, sessionContextData.GetUserID(), input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to create user ingredient preference")
	}

	x := &mealplanningsvc.CreateUserIngredientPreferenceResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	for _, creation := range created {
		x.Created = append(x.Created, converters.ConvertUserIngredientPreferenceToGRPCUserIngredientPreference(creation))
	}

	return x, nil
}

func (s *serviceImpl) FinalizeMealPlan(ctx context.Context, request *mealplanningsvc.FinalizeMealPlanRequest) (*mealplanningsvc.FinalizeMealPlanResponse, error) {
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

	x := &mealplanningsvc.FinalizeMealPlanResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Finalized: finalized,
	}

	return x, nil
}

func (s *serviceImpl) GetMeal(ctx context.Context, request *mealplanningsvc.GetMealRequest) (*mealplanningsvc.GetMealResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealIDKey: request.MealID,
	}, span, s.logger)

	meal, err := s.mealPlanningManager.ReadMeal(ctx, request.MealID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to read meal")
	}

	x := &mealplanningsvc.GetMealResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Result: converters.ConvertMealToGRPCMeal(meal),
	}

	return x, nil
}

func (s *serviceImpl) GetMealPlan(ctx context.Context, request *mealplanningsvc.GetMealPlanRequest) (*mealplanningsvc.GetMealPlanResponse, error) {
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

	x := &mealplanningsvc.GetMealPlanResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Result: converters.ConvertMealPlanToGRPCMealPlan(mealPlan),
	}

	return x, nil
}

func (s *serviceImpl) GetMealPlansForAccount(ctx context.Context, request *mealplanningsvc.GetMealPlansForAccountRequest) (*mealplanningsvc.GetMealPlansForAccountResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, s.logger, span, codes.Unauthenticated, "failed to get session context data")
	}

	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)

	logger := observability.ObserveValues(nil, span, s.logger)

	mealPlans, _, err := s.mealPlanningManager.ListMealPlans(ctx, sessionContextData.GetActiveAccountID(), filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to list meal plans")
	}

	x := &mealplanningsvc.GetMealPlansForAccountResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Filter: request.Filter,
	}

	for _, mealPlan := range mealPlans {
		x.Results = append(x.Results, converters.ConvertMealPlanToGRPCMealPlan(mealPlan))
	}

	return x, nil
}

func (s *serviceImpl) GetMealPlanEvent(ctx context.Context, request *mealplanningsvc.GetMealPlanEventRequest) (*mealplanningsvc.GetMealPlanEventResponse, error) {
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

	x := &mealplanningsvc.GetMealPlanEventResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Result: converters.ConvertMealPlanEventToGRPCMealPlanEvent(mealPlanEvent),
	}

	return x, nil
}

func (s *serviceImpl) GetMealPlanEvents(ctx context.Context, request *mealplanningsvc.GetMealPlanEventsRequest) (*mealplanningsvc.GetMealPlanEventsResponse, error) {
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

	x := &mealplanningsvc.GetMealPlanEventsResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	for _, mealPlanEvent := range mealPlanEvents {
		x.Results = append(x.Results, converters.ConvertMealPlanEventToGRPCMealPlanEvent(mealPlanEvent))
	}

	return x, nil
}

func (s *serviceImpl) GetMealPlanGroceryListItem(ctx context.Context, request *mealplanningsvc.GetMealPlanGroceryListItemRequest) (*mealplanningsvc.GetMealPlanGroceryListItemResponse, error) {
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

	x := &mealplanningsvc.GetMealPlanGroceryListItemResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Result: converters.ConvertMealPlanGroceryListItemToGRPCMealPlanGroceryListItem(mealPlanGroceryListItem),
	}

	return x, nil
}

func (s *serviceImpl) GetMealPlanGroceryListItemsForMealPlan(ctx context.Context, request *mealplanningsvc.GetMealPlanGroceryListItemsForMealPlanRequest) (*mealplanningsvc.GetMealPlanGroceryListItemsForMealPlanResponse, error) {
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

	x := &mealplanningsvc.GetMealPlanGroceryListItemsForMealPlanResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	for _, mealPlanGroceryListItem := range mealPlanGroceryListItems {
		x.Results = append(x.Results, converters.ConvertMealPlanGroceryListItemToGRPCMealPlanGroceryListItem(mealPlanGroceryListItem))
	}

	return x, nil
}

func (s *serviceImpl) GetMealPlanOption(ctx context.Context, request *mealplanningsvc.GetMealPlanOptionRequest) (*mealplanningsvc.GetMealPlanOptionResponse, error) {
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

	x := &mealplanningsvc.GetMealPlanOptionResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Result: converters.ConvertMealPlanOptionToGRPCMealPlanOption(mealPlanOption),
	}

	return x, nil
}

func (s *serviceImpl) GetMealPlanOptionVote(ctx context.Context, request *mealplanningsvc.GetMealPlanOptionVoteRequest) (*mealplanningsvc.GetMealPlanOptionVoteResponse, error) {
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

	x := &mealplanningsvc.GetMealPlanOptionVoteResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Result: converters.ConvertMealPlanOptionVoteToGRPCMealPlanOptionVote(mealPlanOptionVote),
	}

	return x, nil
}

func (s *serviceImpl) GetMealPlanOptionVotes(ctx context.Context, request *mealplanningsvc.GetMealPlanOptionVotesRequest) (*mealplanningsvc.GetMealPlanOptionVotesResponse, error) {
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

	x := &mealplanningsvc.GetMealPlanOptionVotesResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	for _, mealPlanOptionVote := range mealPlanOptionVotes {
		x.Results = append(x.Results, converters.ConvertMealPlanOptionVoteToGRPCMealPlanOptionVote(mealPlanOptionVote))
	}

	return x, nil
}

func (s *serviceImpl) GetMealPlanOptions(ctx context.Context, request *mealplanningsvc.GetMealPlanOptionsRequest) (*mealplanningsvc.GetMealPlanOptionsResponse, error) {
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

	x := &mealplanningsvc.GetMealPlanOptionsResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	for _, mealPlanOption := range mealPlanOptions {
		x.Results = append(x.Results, converters.ConvertMealPlanOptionToGRPCMealPlanOption(mealPlanOption))
	}

	return x, nil
}

func (s *serviceImpl) GetMealPlanTask(ctx context.Context, request *mealplanningsvc.GetMealPlanTaskRequest) (*mealplanningsvc.GetMealPlanTaskResponse, error) {
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

	x := &mealplanningsvc.GetMealPlanTaskResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Result: converters.ConvertMealPlanTaskToGRPCMealPlanTask(mealPlanTask),
	}

	return x, nil
}

func (s *serviceImpl) GetMealPlanTasks(ctx context.Context, request *mealplanningsvc.GetMealPlanTasksRequest) (*mealplanningsvc.GetMealPlanTasksResponse, error) {
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

	x := &mealplanningsvc.GetMealPlanTasksResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	for _, mealPlanTask := range mealPlanTasks {
		x.Results = append(x.Results, converters.ConvertMealPlanTaskToGRPCMealPlanTask(mealPlanTask))
	}

	return x, nil
}

func (s *serviceImpl) GetMeals(ctx context.Context, request *mealplanningsvc.GetMealsRequest) (*mealplanningsvc.GetMealsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)

	meals, _, err := s.mealPlanningManager.ListMeals(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to fetch list of meals")
	}

	x := &mealplanningsvc.GetMealsResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	for _, meal := range meals {
		x.Results = append(x.Results, converters.ConvertMealToGRPCMeal(meal))
	}

	return x, nil
}

func (s *serviceImpl) GetUserIngredientPreference(ctx context.Context, request *mealplanningsvc.GetUserIngredientPreferenceRequest) (*mealplanningsvc.GetUserIngredientPreferenceResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to get session context data")
	}

	userIngredientPreference, err := s.mealPlanningManager.ReadUserIngredientPreference(ctx, sessionContextData.GetUserID(), request.UserIngredientPreferenceID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to fetch list of meals")
	}

	x := &mealplanningsvc.GetUserIngredientPreferenceResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Result: converters.ConvertUserIngredientPreferenceToGRPCUserIngredientPreference(userIngredientPreference),
	}

	return x, nil
}

func (s *serviceImpl) GetUserIngredientPreferences(ctx context.Context, request *mealplanningsvc.GetUserIngredientPreferencesRequest) (*mealplanningsvc.GetUserIngredientPreferencesResponse, error) {
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

	x := &mealplanningsvc.GetUserIngredientPreferencesResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	for _, userIngredientPreference := range userIngredientPreferences {
		x.Results = append(x.Results, converters.ConvertUserIngredientPreferenceToGRPCUserIngredientPreference(userIngredientPreference))
	}

	return x, nil
}

func (s *serviceImpl) RunFinalizeMealPlanWorker(ctx context.Context, _ *mealplanningsvc.RunFinalizeMealPlanWorkerRequest) (*mealplanningsvc.RunFinalizeMealPlanWorkerResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	finalized, err := s.mealPlanFinalizerWorker.Work(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "running meal plan finalizer worker")
	}

	x := &mealplanningsvc.RunFinalizeMealPlanWorkerResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Finalized: uint32(finalized),
	}

	return x, nil
}

func (s *serviceImpl) RunMealPlanGroceryListInitializerWorker(ctx context.Context, _ *mealplanningsvc.RunMealPlanGroceryListInitializerWorkerRequest) (*mealplanningsvc.RunMealPlanGroceryListInitializerWorkerResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	if err := s.mealPlanGroceryListInitializerWorker.Work(ctx); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "running meal plan grocery list initializer worker")
	}

	x := &mealplanningsvc.RunMealPlanGroceryListInitializerWorkerResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *serviceImpl) RunMealPlanTaskCreatorWorker(ctx context.Context, _ *mealplanningsvc.RunMealPlanTaskCreatorWorkerRequest) (*mealplanningsvc.RunMealPlanTaskCreatorWorkerResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	if err := s.mealPlanTaskCreatorWorker.Work(ctx); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "running meal plan task creator worker")
	}

	x := &mealplanningsvc.RunMealPlanTaskCreatorWorkerResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *serviceImpl) SearchForMeals(ctx context.Context, request *mealplanningsvc.SearchForMealsRequest) (*mealplanningsvc.SearchForMealsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)

	meals, err := s.mealPlanningManager.SearchMeals(ctx, request.Query, !request.UseSearchService, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to search for meals")
	}

	x := &mealplanningsvc.SearchForMealsResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	for _, meal := range meals {
		x.Results = append(x.Results, converters.ConvertMealToGRPCMeal(meal))
	}

	return x, nil
}

func (s *serviceImpl) UpdateMealPlan(ctx context.Context, request *mealplanningsvc.UpdateMealPlanRequest) (*mealplanningsvc.UpdateMealPlanResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey: request.MealPlanID,
	}, span, s.logger)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to get session context data")
	}

	input := converters.ConvertGRPCMealPlanUpdateRequestInputToMealPlanUpdateRequestInput(request.Input)

	if err = s.mealPlanningManager.UpdateMealPlan(ctx, request.MealPlanID, sessionContextData.GetActiveAccountID(), input); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to update meal plan")
	}

	updated, err := s.mealPlanningManager.ReadMealPlan(ctx, request.MealPlanID, sessionContextData.GetActiveAccountID())
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to fetch updated meal plan")
	}

	x := &mealplanningsvc.UpdateMealPlanResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Updated: converters.ConvertMealPlanToGRPCMealPlan(updated),
	}

	return x, nil
}

func (s *serviceImpl) UpdateMealPlanEvent(ctx context.Context, request *mealplanningsvc.UpdateMealPlanEventRequest) (*mealplanningsvc.UpdateMealPlanEventResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey:      request.MealPlanID,
		keys.MealPlanEventIDKey: request.MealPlanEventID,
	}, span, s.logger)

	input := converters.ConvertGRPCMealPlanEventUpdateRequestInputToMealPlanEventUpdateRequestInput(request.Input)

	if err := s.mealPlanningManager.UpdateMealPlanEvent(ctx, request.MealPlanID, request.MealPlanEventID, input); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to update meal plan event")
	}

	updated, err := s.mealPlanningManager.ReadMealPlanEvent(ctx, request.MealPlanID, request.MealPlanEventID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to fetch updated meal plan event")
	}

	x := &mealplanningsvc.UpdateMealPlanEventResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Updated: converters.ConvertMealPlanEventToGRPCMealPlanEvent(updated),
	}

	return x, nil
}

func (s *serviceImpl) UpdateMealPlanGroceryListItem(ctx context.Context, request *mealplanningsvc.UpdateMealPlanGroceryListItemRequest) (*mealplanningsvc.UpdateMealPlanGroceryListItemResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey:                request.MealPlanID,
		keys.MealPlanGroceryListItemIDKey: request.MealPlanGroceryListItemID,
	}, span, s.logger)

	input := converters.ConvertGRPCMealPlanGroceryListItemUpdateRequestInputToMealPlanGroceryListItemUpdateRequestInput(request.Input)

	if err := s.mealPlanningManager.UpdateMealPlanGroceryListItem(ctx, request.MealPlanID, request.MealPlanGroceryListItemID, input); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to update meal plan grocery list item")
	}

	updated, err := s.mealPlanningManager.ReadMealPlanGroceryListItem(ctx, request.MealPlanID, request.MealPlanGroceryListItemID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to fetch updated meal plan grocery list item")
	}

	x := &mealplanningsvc.UpdateMealPlanGroceryListItemResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Updated: converters.ConvertMealPlanGroceryListItemToGRPCMealPlanGroceryListItem(updated),
	}

	return x, nil
}

func (s *serviceImpl) UpdateMealPlanOption(ctx context.Context, request *mealplanningsvc.UpdateMealPlanOptionRequest) (*mealplanningsvc.UpdateMealPlanOptionResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey:       request.MealPlanID,
		keys.MealPlanOptionIDKey: request.MealPlanOptionID,
		keys.MealPlanEventIDKey:  request.MealPlanEventID,
	}, span, s.logger)

	input := converters.ConvertGRPCMealPlanOptionUpdateRequestInputToMealPlanOptionUpdateRequestInput(request.Input)

	if err := s.mealPlanningManager.UpdateMealPlanOption(ctx, request.MealPlanID, request.MealPlanEventID, request.MealPlanOptionID, input); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to update meal plan option")
	}

	updated, err := s.mealPlanningManager.ReadMealPlanOption(ctx, request.MealPlanID, request.MealPlanEventID, request.MealPlanOptionID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to fetch updated meal plan option")
	}

	x := &mealplanningsvc.UpdateMealPlanOptionResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Updated: converters.ConvertMealPlanOptionToGRPCMealPlanOption(updated),
	}

	return x, nil
}

func (s *serviceImpl) UpdateMealPlanOptionVote(ctx context.Context, request *mealplanningsvc.UpdateMealPlanOptionVoteRequest) (*mealplanningsvc.UpdateMealPlanOptionVoteResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey:           request.MealPlanID,
		keys.MealPlanOptionIDKey:     request.MealPlanOptionID,
		keys.MealPlanOptionVoteIDKey: request.MealPlanOptionVoteID,
	}, span, s.logger)

	input := converters.ConvertGRPCMealPlanOptionVoteUpdateRequestInputToMealPlanOptionVoteUpdateRequestInput(request.Input)

	if err := s.mealPlanningManager.UpdateMealPlanOptionVote(ctx, request.MealPlanID, request.MealPlanEventID, request.MealPlanOptionID, request.MealPlanOptionVoteID, input); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to update meal plan option vote")
	}

	updated, err := s.mealPlanningManager.ReadMealPlanOptionVote(ctx, request.MealPlanID, request.MealPlanEventID, request.MealPlanOptionID, request.MealPlanOptionVoteID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to fetch updated meal plan option vote")
	}

	x := &mealplanningsvc.UpdateMealPlanOptionVoteResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Updated: converters.ConvertMealPlanOptionVoteToGRPCMealPlanOptionVote(updated),
	}

	return x, nil
}

func (s *serviceImpl) UpdateMealPlanTaskStatus(ctx context.Context, request *mealplanningsvc.UpdateMealPlanTaskStatusRequest) (*mealplanningsvc.UpdateMealPlanTaskStatusResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey:     request.MealPlanID,
		keys.MealPlanTaskIDKey: request.MealPlanTaskID,
	}, span, s.logger)

	input := converters.ConvertGRPCMealPlanTaskStatusChangeRequestInputToMealPlanTaskStatusChangeRequestInput(request.Input)

	if err := s.mealPlanningManager.MealPlanTaskStatusChange(ctx, input); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to update meal plan task status")
	}

	updated, err := s.mealPlanningManager.ReadMealPlanTask(ctx, request.MealPlanID, request.MealPlanTaskID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to fetch updated meal plan task status")
	}

	x := &mealplanningsvc.UpdateMealPlanTaskStatusResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Updated: converters.ConvertMealPlanTaskToGRPCMealPlanTask(updated),
	}

	return x, nil
}

func (s *serviceImpl) UpdateUserIngredientPreference(ctx context.Context, request *mealplanningsvc.UpdateUserIngredientPreferenceRequest) (*mealplanningsvc.UpdateUserIngredientPreferenceResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.UserIngredientPreferenceIDKey: request.UserIngredientPreferenceID,
	}, span, s.logger)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to fetch session context data")
	}

	input := converters.ConvertGRPCUserIngredientPreferenceUpdateRequestInputToUserIngredientPreferenceUpdateRequestInput(request.Input)

	if err = s.mealPlanningManager.UpdateUserIngredientPreference(ctx, request.UserIngredientPreferenceID, sessionContextData.GetUserID(), input); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to update meal plan task status")
	}

	updated, err := s.mealPlanningManager.ReadUserIngredientPreference(ctx, sessionContextData.GetUserID(), request.UserIngredientPreferenceID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to fetch updated meal plan task status")
	}

	x := &mealplanningsvc.UpdateUserIngredientPreferenceResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Updated: converters.ConvertUserIngredientPreferenceToGRPCUserIngredientPreference(updated),
	}

	return x, nil
}

func (s *serviceImpl) CreateAccountInstrumentOwnership(ctx context.Context, request *mealplanningsvc.CreateAccountInstrumentOwnershipRequest) (*mealplanningsvc.CreateAccountInstrumentOwnershipResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, s.logger, span, codes.Unauthenticated, "failed to fetch session context data")
	}

	logger := observability.ObserveValues(nil, span, s.logger)

	input := converters.ConvertGRPCAccountInstrumentOwnershipCreationRequestInputToAccountInstrumentOwnershipCreationRequestInput(request.Input)
	if err = input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.InvalidArgument, "failed to validate account instrument ownership creation request input")
	}

	created, err := s.mealPlanningManager.CreateAccountInstrumentOwnership(ctx, sessionContextData.GetActiveAccountID(), input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to create account instrument ownership")
	}

	x := &mealplanningsvc.CreateAccountInstrumentOwnershipResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Created: converters.ConvertAccountInstrumentOwnershipToGRPCAccountInstrumentOwnership(created),
	}

	return x, nil
}

func (s *serviceImpl) GetAccountInstrumentOwnership(ctx context.Context, request *mealplanningsvc.GetAccountInstrumentOwnershipRequest) (*mealplanningsvc.GetAccountInstrumentOwnershipResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, s.logger, span, codes.Unauthenticated, "failed to fetch session context data")
	}

	logger := observability.ObserveValues(map[string]any{
		keys.AccountInstrumentOwnershipIDKey: request.AccountInstrumentOwnershipID,
	}, span, s.logger)

	result, err := s.mealPlanningManager.ReadAccountInstrumentOwnership(ctx, sessionContextData.GetActiveAccountID(), request.AccountInstrumentOwnershipID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to fetch account instrument ownership")
	}

	x := &mealplanningsvc.GetAccountInstrumentOwnershipResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Result: converters.ConvertAccountInstrumentOwnershipToGRPCAccountInstrumentOwnership(result),
	}

	return x, nil
}

func (s *serviceImpl) GetAccountInstrumentOwnerships(ctx context.Context, request *mealplanningsvc.GetAccountInstrumentOwnershipsRequest) (*mealplanningsvc.GetAccountInstrumentOwnershipsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, s.logger, span, codes.Unauthenticated, "failed to fetch session context data")
	}

	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)

	logger := observability.ObserveValues(nil, span, s.logger)

	results, _, err := s.mealPlanningManager.ListAccountInstrumentOwnerships(ctx, sessionContextData.GetActiveAccountID(), filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to list account instrument ownerships")
	}

	x := &mealplanningsvc.GetAccountInstrumentOwnershipsResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	for _, result := range results {
		x.Results = append(x.Results, converters.ConvertAccountInstrumentOwnershipToGRPCAccountInstrumentOwnership(result))
	}

	return x, nil
}

func (s *serviceImpl) UpdateAccountInstrumentOwnership(ctx context.Context, request *mealplanningsvc.UpdateAccountInstrumentOwnershipRequest) (*mealplanningsvc.UpdateAccountInstrumentOwnershipResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, s.logger, span, codes.Unauthenticated, "failed to fetch session context data")
	}

	logger := observability.ObserveValues(map[string]any{
		keys.AccountInstrumentOwnershipIDKey: request.AccountInstrumentOwnershipID,
	}, span, s.logger)

	input := converters.ConvertGRPCAccountInstrumentOwnershipUpdateRequestInputToAccountInstrumentOwnershipUpdateRequestInput(request.Input)

	accountInstrumentOwnership, err := s.mealPlanningManager.ReadAccountInstrumentOwnership(ctx, sessionContextData.GetActiveAccountID(), request.AccountInstrumentOwnershipID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to fetch account instrument ownership")
	}

	if err = s.mealPlanningManager.UpdateAccountInstrumentOwnership(ctx, accountInstrumentOwnership.ID, accountInstrumentOwnership.BelongsToAccount, input); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to update account instrument ownership")
	}

	x := &mealplanningsvc.UpdateAccountInstrumentOwnershipResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *serviceImpl) ArchiveAccountInstrumentOwnership(ctx context.Context, request *mealplanningsvc.ArchiveAccountInstrumentOwnershipRequest) (*mealplanningsvc.ArchiveAccountInstrumentOwnershipResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, s.logger, span, codes.Unauthenticated, "failed to fetch session context data")
	}

	logger := observability.ObserveValues(map[string]any{
		keys.AccountInstrumentOwnershipIDKey: request.AccountInstrumentOwnershipID,
	}, span, s.logger)

	if err = s.mealPlanningManager.ArchiveAccountInstrumentOwnership(ctx, sessionContextData.GetActiveAccountID(), request.AccountInstrumentOwnershipID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to archive account instrument ownership")
	}

	x := &mealplanningsvc.ArchiveAccountInstrumentOwnershipResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}
