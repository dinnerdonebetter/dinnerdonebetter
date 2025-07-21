package grpc

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/grpc/generated/types"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"

	"google.golang.org/grpc/codes"
)

func (s *ServiceImpl) ArchiveRecipe(ctx context.Context, request *mealplanning.ArchiveRecipeRequest) (*mealplanning.ArchiveRecipeResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey: request.RecipeID,
	})

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "fetching session context data")
	}

	if err = s.recipeManager.ArchiveRecipe(ctx, request.RecipeID, sessionContextData.GetUserID()); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving recipe")
	}

	x := &mealplanning.ArchiveRecipeResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *ServiceImpl) ArchiveRecipePrepTask(ctx context.Context, request *mealplanning.ArchiveRecipePrepTaskRequest) (*mealplanning.ArchiveRecipePrepTaskResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:         request.RecipeID,
		keys.RecipePrepTaskIDKey: request.RecipePrepTaskID,
	})

	if err := s.recipeManager.ArchiveRecipePrepTask(ctx, request.GetRecipeID(), request.GetRecipePrepTaskID()); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "archiving recipe prep task")
	}

	x := &mealplanning.ArchiveRecipePrepTaskResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *ServiceImpl) ArchiveRecipeRating(ctx context.Context, request *mealplanning.ArchiveRecipeRatingRequest) (*mealplanning.ArchiveRecipeRatingResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:       request.RecipeID,
		keys.RecipeRatingIDKey: request.RecipeRatingID,
	})

	if err := s.recipeManager.ArchiveRecipeRating(ctx, request.RecipeID, request.RecipeRatingID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving recipe rating")
	}

	x := &mealplanning.ArchiveRecipeRatingResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *ServiceImpl) ArchiveRecipeStep(ctx context.Context, request *mealplanning.ArchiveRecipeStepRequest) (*mealplanning.ArchiveRecipeStepResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:     request.RecipeID,
		keys.RecipeStepIDKey: request.RecipeStepID,
	})

	if err := s.recipeManager.ArchiveRecipeStep(ctx, request.RecipeID, request.RecipeStepID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving recipe step")
	}

	x := &mealplanning.ArchiveRecipeStepResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *ServiceImpl) ArchiveRecipeStepCompletionCondition(ctx context.Context, request *mealplanning.ArchiveRecipeStepCompletionConditionRequest) (*mealplanning.ArchiveRecipeStepCompletionConditionResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:                        request.RecipeID,
		keys.RecipeStepIDKey:                    request.RecipeStepID,
		keys.RecipeStepCompletionConditionIDKey: request.RecipeStepCompletionConditionID,
	})

	if err := s.recipeManager.ArchiveRecipeStepCompletionCondition(ctx, request.RecipeID, request.RecipeStepID, request.RecipeStepCompletionConditionID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving recipe step completion condition")
	}

	x := &mealplanning.ArchiveRecipeStepCompletionConditionResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *ServiceImpl) ArchiveRecipeStepIngredient(ctx context.Context, request *mealplanning.ArchiveRecipeStepIngredientRequest) (*mealplanning.ArchiveRecipeStepIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:               request.RecipeID,
		keys.RecipeStepIDKey:           request.RecipeStepID,
		keys.RecipeStepIngredientIDKey: request.RecipeStepIngredientID,
	})

	if err := s.recipeManager.ArchiveRecipeStepIngredient(ctx, request.RecipeID, request.RecipeStepID, request.RecipeStepIngredientID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving recipe step ingredient")
	}

	x := &mealplanning.ArchiveRecipeStepIngredientResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *ServiceImpl) ArchiveRecipeStepInstrument(ctx context.Context, request *mealplanning.ArchiveRecipeStepInstrumentRequest) (*mealplanning.ArchiveRecipeStepInstrumentResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:               request.RecipeID,
		keys.RecipeStepIDKey:           request.RecipeStepID,
		keys.RecipeStepInstrumentIDKey: request.RecipeStepInstrumentID,
	})

	if err := s.recipeManager.ArchiveRecipeStepInstrument(ctx, request.RecipeID, request.RecipeStepID, request.RecipeStepInstrumentID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving recipe step instrument")
	}

	x := &mealplanning.ArchiveRecipeStepInstrumentResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *ServiceImpl) ArchiveRecipeStepProduct(ctx context.Context, request *mealplanning.ArchiveRecipeStepProductRequest) (*mealplanning.ArchiveRecipeStepProductResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:            request.RecipeID,
		keys.RecipeStepIDKey:        request.RecipeStepID,
		keys.RecipeStepProductIDKey: request.RecipeStepProductID,
	})

	if err := s.recipeManager.ArchiveRecipeStepProduct(ctx, request.RecipeID, request.RecipeStepID, request.RecipeStepProductID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving recipe step product")
	}

	x := &mealplanning.ArchiveRecipeStepProductResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *ServiceImpl) ArchiveRecipeStepVessel(ctx context.Context, request *mealplanning.ArchiveRecipeStepVesselRequest) (*mealplanning.ArchiveRecipeStepVesselResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:           request.RecipeID,
		keys.RecipeStepIDKey:       request.RecipeStepID,
		keys.RecipeStepVesselIDKey: request.RecipeStepVesselID,
	})

	if err := s.recipeManager.ArchiveRecipeStepVessel(ctx, request.RecipeID, request.RecipeStepID, request.RecipeStepVesselID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving recipe step vessel")
	}

	x := &mealplanning.ArchiveRecipeStepVesselResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *ServiceImpl) CloneRecipe(ctx context.Context, request *mealplanning.CloneRecipeRequest) (*mealplanning.CloneRecipeResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey: request.RecipeID,
	})

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching session context data")
	}

	_, err = s.recipeManager.CloneRecipe(ctx, request.RecipeID, sessionContextData.GetUserID())
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "cloning recipe")
	}

	x := &mealplanning.CloneRecipeResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Cloned: nil, // TODO: cloned
	}

	return x, nil
}

func (s *ServiceImpl) CreateRecipe(ctx context.Context, request *mealplanning.CreateRecipeRequest) (*mealplanning.CreateRecipeResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateRecipePrepTask(ctx context.Context, request *mealplanning.CreateRecipePrepTaskRequest) (*mealplanning.CreateRecipePrepTaskResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey: request.RecipeID,
	})

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateRecipeRating(ctx context.Context, request *mealplanning.CreateRecipeRatingRequest) (*mealplanning.CreateRecipeRatingResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey: request.RecipeID,
	})

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateRecipeStep(ctx context.Context, request *mealplanning.CreateRecipeStepRequest) (*mealplanning.CreateRecipeStepResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey: request.RecipeID,
	})

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateRecipeStepCompletionCondition(ctx context.Context, request *mealplanning.CreateRecipeStepCompletionConditionRequest) (*mealplanning.CreateRecipeStepCompletionConditionResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:     request.RecipeID,
		keys.RecipeStepIDKey: request.RecipeStepID,
	})

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateRecipeStepIngredient(ctx context.Context, request *mealplanning.CreateRecipeStepIngredientRequest) (*mealplanning.CreateRecipeStepIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:     request.RecipeID,
		keys.RecipeStepIDKey: request.RecipeStepID,
	})

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateRecipeStepInstrument(ctx context.Context, request *mealplanning.CreateRecipeStepInstrumentRequest) (*mealplanning.CreateRecipeStepInstrumentResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:     request.RecipeID,
		keys.RecipeStepIDKey: request.RecipeStepID,
	})

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateRecipeStepProduct(ctx context.Context, request *mealplanning.CreateRecipeStepProductRequest) (*mealplanning.CreateRecipeStepProductResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:     request.RecipeID,
		keys.RecipeStepIDKey: request.RecipeStepID,
	})

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateRecipeStepVessel(ctx context.Context, request *mealplanning.CreateRecipeStepVesselRequest) (*mealplanning.CreateRecipeStepVesselResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:     request.RecipeID,
		keys.RecipeStepIDKey: request.RecipeStepID,
	})

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetMermaidDiagramForRecipe(ctx context.Context, request *mealplanning.GetMermaidDiagramForRecipeRequest) (*mealplanning.GetMermaidDiagramForRecipeResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey: request.RecipeID,
	})

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipe(ctx context.Context, request *mealplanning.GetRecipeRequest) (*mealplanning.GetRecipeResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey: request.RecipeID,
	})

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipeMealPlanTasks(ctx context.Context, request *mealplanning.GetRecipeMealPlanTasksRequest) (*mealplanning.GetRecipeMealPlanTasksResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey: request.RecipeID,
	})

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipePrepTask(ctx context.Context, request *mealplanning.GetRecipePrepTaskRequest) (*mealplanning.GetRecipePrepTaskResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:         request.RecipeID,
		keys.RecipePrepTaskIDKey: request.RecipePrepTaskID,
	})

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipePrepTasks(ctx context.Context, request *mealplanning.GetRecipePrepTasksRequest) (*mealplanning.GetRecipePrepTasksResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey: request.RecipeID,
	})

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipeRating(ctx context.Context, request *mealplanning.GetRecipeRatingRequest) (*mealplanning.GetRecipeRatingResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:       request.RecipeID,
		keys.RecipeRatingIDKey: request.RecipeRatingID,
	})

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipeRatingsForRecipe(ctx context.Context, request *mealplanning.GetRecipeRatingsForRecipeRequest) (*mealplanning.GetRecipeRatingsForRecipeResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey: request.RecipeID,
	})

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipeStep(ctx context.Context, request *mealplanning.GetRecipeStepRequest) (*mealplanning.GetRecipeStepResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:     request.RecipeID,
		keys.RecipeStepIDKey: request.RecipeStepID,
	})

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipeStepCompletionCondition(ctx context.Context, request *mealplanning.GetRecipeStepCompletionConditionRequest) (*mealplanning.GetRecipeStepCompletionConditionResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:                        request.RecipeID,
		keys.RecipeStepIDKey:                    request.RecipeStepID,
		keys.RecipeStepCompletionConditionIDKey: request.RecipeStepCompletionConditionID,
	})

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipeStepCompletionConditions(ctx context.Context, request *mealplanning.GetRecipeStepCompletionConditionsRequest) (*mealplanning.GetRecipeStepCompletionConditionsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:     request.RecipeID,
		keys.RecipeStepIDKey: request.RecipeStepID,
	})

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipeStepIngredient(ctx context.Context, request *mealplanning.GetRecipeStepIngredientRequest) (*mealplanning.GetRecipeStepIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:               request.RecipeID,
		keys.RecipeStepIDKey:           request.RecipeStepID,
		keys.RecipeStepIngredientIDKey: request.RecipeStepIngredientID,
	})

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipeStepIngredients(ctx context.Context, request *mealplanning.GetRecipeStepIngredientsRequest) (*mealplanning.GetRecipeStepIngredientsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:     request.RecipeID,
		keys.RecipeStepIDKey: request.RecipeStepID,
	})

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipeStepInstrument(ctx context.Context, request *mealplanning.GetRecipeStepInstrumentRequest) (*mealplanning.GetRecipeStepInstrumentResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:               request.RecipeID,
		keys.RecipeStepIDKey:           request.RecipeStepID,
		keys.RecipeStepInstrumentIDKey: request.RecipeStepInstrumentID,
	})

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipeStepInstruments(ctx context.Context, request *mealplanning.GetRecipeStepInstrumentsRequest) (*mealplanning.GetRecipeStepInstrumentsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:     request.RecipeID,
		keys.RecipeStepIDKey: request.RecipeStepID,
	})

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipeStepProduct(ctx context.Context, request *mealplanning.GetRecipeStepProductRequest) (*mealplanning.GetRecipeStepProductResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:            request.RecipeID,
		keys.RecipeStepIDKey:        request.RecipeStepID,
		keys.RecipeStepProductIDKey: request.RecipeStepProductID,
	})

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipeStepProducts(ctx context.Context, request *mealplanning.GetRecipeStepProductsRequest) (*mealplanning.GetRecipeStepProductsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:     request.RecipeID,
		keys.RecipeStepIDKey: request.RecipeStepID,
	})

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipeStepVessel(ctx context.Context, request *mealplanning.GetRecipeStepVesselRequest) (*mealplanning.GetRecipeStepVesselResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:           request.RecipeID,
		keys.RecipeStepIDKey:       request.RecipeStepID,
		keys.RecipeStepVesselIDKey: request.RecipeStepVesselID,
	})

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipeStepVessels(ctx context.Context, request *mealplanning.GetRecipeStepVesselsRequest) (*mealplanning.GetRecipeStepVesselsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:     request.RecipeID,
		keys.RecipeStepIDKey: request.RecipeStepID,
	})

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipeSteps(ctx context.Context, request *mealplanning.GetRecipeStepsRequest) (*mealplanning.GetRecipeStepsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey: request.RecipeID,
	})

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipes(ctx context.Context, request *mealplanning.GetRecipesRequest) (*mealplanning.GetRecipesResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) SearchForRecipes(ctx context.Context, request *mealplanning.SearchForRecipesRequest) (*mealplanning.SearchForRecipesResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValues(map[string]any{
		keys.SearchQueryKey: request.Query,
	})

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateRecipe(ctx context.Context, request *mealplanning.UpdateRecipeRequest) (*mealplanning.UpdateRecipeResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey: request.RecipeID,
	})

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateRecipePrepTask(ctx context.Context, request *mealplanning.UpdateRecipePrepTaskRequest) (*mealplanning.UpdateRecipePrepTaskResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:         request.RecipeID,
		keys.RecipePrepTaskIDKey: request.RecipePrepTaskID,
	})

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateRecipeRating(ctx context.Context, request *mealplanning.UpdateRecipeRatingRequest) (*mealplanning.UpdateRecipeRatingResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:       request.RecipeID,
		keys.RecipeRatingIDKey: request.RecipeRatingID,
	})

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateRecipeStep(ctx context.Context, request *mealplanning.UpdateRecipeStepRequest) (*mealplanning.UpdateRecipeStepResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:     request.RecipeID,
		keys.RecipeStepIDKey: request.RecipeStepID,
	})

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateRecipeStepCompletionCondition(ctx context.Context, request *mealplanning.UpdateRecipeStepCompletionConditionRequest) (*mealplanning.UpdateRecipeStepCompletionConditionResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:                        request.RecipeID,
		keys.RecipeStepIDKey:                    request.RecipeStepID,
		keys.RecipeStepCompletionConditionIDKey: request.RecipeStepCompletionConditionID,
	})

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateRecipeStepIngredient(ctx context.Context, request *mealplanning.UpdateRecipeStepIngredientRequest) (*mealplanning.UpdateRecipeStepIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:               request.RecipeID,
		keys.RecipeStepIDKey:           request.RecipeStepID,
		keys.RecipeStepIngredientIDKey: request.RecipeStepIngredientID,
	})

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateRecipeStepInstrument(ctx context.Context, request *mealplanning.UpdateRecipeStepInstrumentRequest) (*mealplanning.UpdateRecipeStepInstrumentResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:               request.RecipeID,
		keys.RecipeStepIDKey:           request.RecipeStepID,
		keys.RecipeStepInstrumentIDKey: request.RecipeStepInstrumentID,
	})

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateRecipeStepProduct(ctx context.Context, request *mealplanning.UpdateRecipeStepProductRequest) (*mealplanning.UpdateRecipeStepProductResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:            request.RecipeID,
		keys.RecipeStepIDKey:        request.RecipeStepID,
		keys.RecipeStepProductIDKey: request.RecipeStepProductID,
	})

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateRecipeStepVessel(ctx context.Context, request *mealplanning.UpdateRecipeStepVesselRequest) (*mealplanning.UpdateRecipeStepVesselResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:           request.RecipeID,
		keys.RecipeStepIDKey:       request.RecipeStepID,
		keys.RecipeStepVesselIDKey: request.RecipeStepVesselID,
	})

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}
