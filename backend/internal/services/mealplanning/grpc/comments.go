package grpc

import (
	"context"
	"errors"

	"github.com/dinnerdonebetter/backend/internal/domain/comments"
	grpcconverters "github.com/dinnerdonebetter/backend/internal/grpc/converters"
	commentssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/comments"
	"github.com/dinnerdonebetter/backend/internal/grpc/generated/types"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	converters "github.com/dinnerdonebetter/backend/internal/services/comments/grpc/converters"

	"google.golang.org/grpc/codes"
)

func (s *serviceImpl) AddCommentToRecipe(ctx context.Context, request *commentssvc.AddCommentToRecipeRequest) (*commentssvc.AddCommentToRecipeResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.RecipeIDKey: request.RecipeId,
	}, span, s.logger)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "fetching session context data")
	}

	input := converters.ConvertProtoCommentCreationRequestInputToDomain(
		request.Input,
		comments.CommentTargetTypeRecipes,
		request.RecipeId,
		sessionContextData.GetUserID(),
	)
	if input == nil {
		return nil, observability.PrepareAndLogGRPCStatus(errors.New("input is required"), logger, span, codes.InvalidArgument, "input is required")
	}

	comment, err := s.commentsManager.CreateComment(ctx, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating comment")
	}

	return &commentssvc.AddCommentToRecipeResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceId: span.SpanContext().TraceID().String(),
		},
		Comment: converters.ConvertCommentToGRPCComment(comment),
	}, nil
}

func (s *serviceImpl) AddCommentToMeal(ctx context.Context, request *commentssvc.AddCommentToMealRequest) (*commentssvc.AddCommentToMealResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealIDKey: request.MealId,
	}, span, s.logger)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "fetching session context data")
	}

	input := converters.ConvertProtoCommentCreationRequestInputToDomain(
		request.Input,
		comments.CommentTargetTypeMeals,
		request.MealId,
		sessionContextData.GetUserID(),
	)
	if input == nil {
		return nil, observability.PrepareAndLogGRPCStatus(errors.New("input is required"), logger, span, codes.InvalidArgument, "input is required")
	}

	comment, err := s.commentsManager.CreateComment(ctx, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating comment")
	}

	return &commentssvc.AddCommentToMealResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceId: span.SpanContext().TraceID().String(),
		},
		Comment: converters.ConvertCommentToGRPCComment(comment),
	}, nil
}

func (s *serviceImpl) AddCommentToMealPlan(ctx context.Context, request *commentssvc.AddCommentToMealPlanRequest) (*commentssvc.AddCommentToMealPlanResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealPlanIDKey: request.MealPlanId,
	}, span, s.logger)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "fetching session context data")
	}

	if _, err = s.mealPlanningManager.ReadMealPlan(ctx, request.MealPlanId, sessionContextData.GetActiveAccountID()); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.InvalidArgument, "validating meal plan access")
	}

	input := converters.ConvertProtoCommentCreationRequestInputToDomain(
		request.Input,
		comments.CommentTargetTypeMealPlans,
		request.MealPlanId,
		sessionContextData.GetUserID(),
	)
	if input == nil {
		return nil, observability.PrepareAndLogGRPCStatus(errors.New("input is required"), logger, span, codes.InvalidArgument, "input is required")
	}

	comment, err := s.commentsManager.CreateComment(ctx, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating comment")
	}

	return &commentssvc.AddCommentToMealPlanResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceId: span.SpanContext().TraceID().String(),
		},
		Comment: converters.ConvertCommentToGRPCComment(comment),
	}, nil
}

func (s *serviceImpl) CreateComment(ctx context.Context, request *commentssvc.CreateCommentRequest) (*commentssvc.CreateCommentResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	if request.Input == nil {
		return nil, observability.PrepareAndLogGRPCStatus(errors.New("input is required"), s.logger, span, codes.InvalidArgument, "input is required")
	}

	logger := observability.ObserveValues(map[string]any{
		"target_type":   request.Input.GetTargetType(),
		"referenced_id": request.Input.GetReferencedId(),
	}, span, s.logger)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "fetching session context data")
	}

	input := converters.ConvertProtoCommentCreationRequestInputToDomain(
		request.Input,
		"",
		"",
		sessionContextData.GetUserID(),
	)
	if input == nil {
		return nil, observability.PrepareAndLogGRPCStatus(errors.New("input is required"), logger, span, codes.InvalidArgument, "input is required")
	}
	if input.TargetType == "" || input.ReferencedID == "" {
		return nil, observability.PrepareAndLogGRPCStatus(errors.New("target_type and referenced_id are required"), logger, span, codes.InvalidArgument, "target_type and referenced_id are required")
	}

	switch input.TargetType {
	case comments.CommentTargetTypeRecipes, comments.CommentTargetTypeMeals:
		// no access check
	case comments.CommentTargetTypeMealPlans:
		if _, err = s.mealPlanningManager.ReadMealPlan(ctx, input.ReferencedID, sessionContextData.GetActiveAccountID()); err != nil {
			return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.InvalidArgument, "validating meal plan access")
		}
	default:
		return nil, observability.PrepareAndLogGRPCStatus(errors.New("invalid target type"), logger, span, codes.InvalidArgument, "invalid target type")
	}

	comment, err := s.commentsManager.CreateComment(ctx, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating comment")
	}

	return &commentssvc.CreateCommentResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceId: span.SpanContext().TraceID().String(),
		},
		Comment: converters.ConvertCommentToGRPCComment(comment),
	}, nil
}

func (s *serviceImpl) GetCommentsForReference(ctx context.Context, request *commentssvc.GetCommentsForReferenceRequest) (*commentssvc.GetCommentsForReferenceResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.CommentIDKey: request.ReferencedId,
		"target_type":     request.TargetType,
	}, span, s.logger)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "fetching session context data")
	}

	switch request.TargetType {
	case comments.CommentTargetTypeRecipes, comments.CommentTargetTypeMeals:
		// no access check
	case comments.CommentTargetTypeMealPlans:
		if _, err = s.mealPlanningManager.ReadMealPlan(ctx, request.ReferencedId, sessionContextData.GetActiveAccountID()); err != nil {
			return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.InvalidArgument, "validating meal plan access")
		}
	default:
		return nil, observability.PrepareAndLogGRPCStatus(errors.New("invalid target type"), logger, span, codes.InvalidArgument, "invalid target type")
	}

	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	result, err := s.commentsManager.GetCommentsForReference(ctx, request.TargetType, request.ReferencedId, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching comments")
	}

	x := &commentssvc.GetCommentsForReferenceResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceId: span.SpanContext().TraceID().String(),
		},
		Pagination: grpcconverters.ConvertPaginationToGRPCPagination(result.Pagination, filter),
	}
	for _, c := range result.Data {
		x.Data = append(x.Data, converters.ConvertCommentToGRPCComment(c))
	}

	return x, nil
}

func (s *serviceImpl) UpdateComment(ctx context.Context, request *commentssvc.UpdateCommentRequest) (*commentssvc.UpdateCommentResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.CommentIDKey: request.CommentId,
	}, span, s.logger)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "fetching session context data")
	}

	comment, err := s.commentsManager.GetComment(ctx, request.CommentId)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching comment")
	}

	if comment.BelongsToUser != sessionContextData.GetUserID() {
		return nil, observability.PrepareAndLogGRPCStatus(errors.New("comment does not belong to user"), logger, span, codes.PermissionDenied, "comment does not belong to user")
	}

	updateInput := converters.ConvertProtoCommentUpdateRequestInputToDomain(request.Input)
	if updateInput == nil {
		return nil, observability.PrepareAndLogGRPCStatus(errors.New("input is required"), logger, span, codes.InvalidArgument, "input is required")
	}

	if err = s.commentsManager.UpdateComment(ctx, request.CommentId, sessionContextData.GetUserID(), updateInput); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "updating comment")
	}

	updated, err := s.commentsManager.GetComment(ctx, request.CommentId)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching updated comment")
	}

	return &commentssvc.UpdateCommentResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceId: span.SpanContext().TraceID().String(),
		},
		Comment: converters.ConvertCommentToGRPCComment(updated),
	}, nil
}

func (s *serviceImpl) ArchiveComment(ctx context.Context, request *commentssvc.ArchiveCommentRequest) (*commentssvc.ArchiveCommentResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.CommentIDKey: request.CommentId,
	}, span, s.logger)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "fetching session context data")
	}

	comment, err := s.commentsManager.GetComment(ctx, request.CommentId)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching comment")
	}

	if comment.BelongsToUser != sessionContextData.GetUserID() {
		return nil, observability.PrepareAndLogGRPCStatus(errors.New("comment does not belong to user"), logger, span, codes.PermissionDenied, "comment does not belong to user")
	}

	if err = s.commentsManager.ArchiveComment(ctx, request.CommentId); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving comment")
	}

	return &commentssvc.ArchiveCommentResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceId: span.SpanContext().TraceID().String(),
		},
	}, nil
}
