package grpc

import (
	"context"
	"errors"

	"github.com/dinnerdonebetter/backend/internal/domain/comments"
	grpcconverters "github.com/dinnerdonebetter/backend/internal/grpc/converters"
	mealplanningsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/grpc/generated/types"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	converters "github.com/dinnerdonebetter/backend/internal/services/mealplanning/grpc/converters"

	"google.golang.org/grpc/codes"
)

func (s *serviceImpl) AddCommentToRecipe(ctx context.Context, request *mealplanningsvc.AddCommentToRecipeRequest) (*mealplanningsvc.AddCommentToRecipeResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.RecipeIDKey: request.RecipeId,
	}, span, s.logger)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "fetching session context data")
	}

	input := &comments.CommentCreationRequestInput{
		Content:         request.Content,
		TargetType:      comments.CommentTargetTypeRecipes,
		ReferencedID:    request.RecipeId,
		ParentCommentID: request.ParentCommentId,
		BelongsToUser:   sessionContextData.GetUserID(),
	}

	comment, err := s.commentsManager.CreateComment(ctx, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating comment")
	}

	return &mealplanningsvc.AddCommentToRecipeResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceId: span.SpanContext().TraceID().String(),
		},
		Comment: converters.ConvertCommentToGRPCComment(comment),
	}, nil
}

func (s *serviceImpl) AddCommentToMeal(ctx context.Context, request *mealplanningsvc.AddCommentToMealRequest) (*mealplanningsvc.AddCommentToMealResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.MealIDKey: request.MealId,
	}, span, s.logger)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "fetching session context data")
	}

	input := &comments.CommentCreationRequestInput{
		Content:         request.Content,
		TargetType:      comments.CommentTargetTypeMeals,
		ReferencedID:    request.MealId,
		ParentCommentID: request.ParentCommentId,
		BelongsToUser:   sessionContextData.GetUserID(),
	}

	comment, err := s.commentsManager.CreateComment(ctx, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating comment")
	}

	return &mealplanningsvc.AddCommentToMealResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceId: span.SpanContext().TraceID().String(),
		},
		Comment: converters.ConvertCommentToGRPCComment(comment),
	}, nil
}

func (s *serviceImpl) AddCommentToMealPlan(ctx context.Context, request *mealplanningsvc.AddCommentToMealPlanRequest) (*mealplanningsvc.AddCommentToMealPlanResponse, error) {
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

	input := &comments.CommentCreationRequestInput{
		Content:         request.Content,
		TargetType:      comments.CommentTargetTypeMealPlans,
		ReferencedID:    request.MealPlanId,
		ParentCommentID: request.ParentCommentId,
		BelongsToUser:   sessionContextData.GetUserID(),
	}

	comment, err := s.commentsManager.CreateComment(ctx, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating comment")
	}

	return &mealplanningsvc.AddCommentToMealPlanResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceId: span.SpanContext().TraceID().String(),
		},
		Comment: converters.ConvertCommentToGRPCComment(comment),
	}, nil
}

func (s *serviceImpl) GetCommentsForReference(ctx context.Context, request *mealplanningsvc.GetCommentsForReferenceRequest) (*mealplanningsvc.GetCommentsForReferenceResponse, error) {
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
	// everybody can read recipes and meals so there's no need to check access.
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

	x := &mealplanningsvc.GetCommentsForReferenceResponse{
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

func (s *serviceImpl) UpdateComment(ctx context.Context, request *mealplanningsvc.UpdateCommentRequest) (*mealplanningsvc.UpdateCommentResponse, error) {
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

	if err = s.commentsManager.UpdateComment(ctx, request.CommentId, sessionContextData.GetUserID(), &comments.CommentUpdateRequestInput{Content: request.Content}); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "updating comment")
	}

	updated, err := s.commentsManager.GetComment(ctx, request.CommentId)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching updated comment")
	}

	return &mealplanningsvc.UpdateCommentResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceId: span.SpanContext().TraceID().String(),
		},
		Comment: converters.ConvertCommentToGRPCComment(updated),
	}, nil
}

func (s *serviceImpl) ArchiveComment(ctx context.Context, request *mealplanningsvc.ArchiveCommentRequest) (*mealplanningsvc.ArchiveCommentResponse, error) {
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

	return &mealplanningsvc.ArchiveCommentResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceId: span.SpanContext().TraceID().String(),
		},
	}, nil
}
