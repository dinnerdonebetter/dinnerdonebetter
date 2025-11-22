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

func (s *serviceImpl) ArchiveRecipe(ctx context.Context, request *mealplanning.ArchiveRecipeRequest) (*mealplanning.ArchiveRecipeResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.RecipeIDKey: request.RecipeID,
	}, span, s.logger)

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

func (s *serviceImpl) ArchiveRecipePrepTask(ctx context.Context, request *mealplanning.ArchiveRecipePrepTaskRequest) (*mealplanning.ArchiveRecipePrepTaskResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.RecipeIDKey:         request.RecipeID,
		keys.RecipePrepTaskIDKey: request.RecipePrepTaskID,
	}, span, s.logger)

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

func (s *serviceImpl) ArchiveRecipeRating(ctx context.Context, request *mealplanning.ArchiveRecipeRatingRequest) (*mealplanning.ArchiveRecipeRatingResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.RecipeIDKey:       request.RecipeID,
		keys.RecipeRatingIDKey: request.RecipeRatingID,
	}, span, s.logger)

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

func (s *serviceImpl) ArchiveRecipeStep(ctx context.Context, request *mealplanning.ArchiveRecipeStepRequest) (*mealplanning.ArchiveRecipeStepResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.RecipeIDKey:     request.RecipeID,
		keys.RecipeStepIDKey: request.RecipeStepID,
	}, span, s.logger)

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

func (s *serviceImpl) ArchiveRecipeStepCompletionCondition(ctx context.Context, request *mealplanning.ArchiveRecipeStepCompletionConditionRequest) (*mealplanning.ArchiveRecipeStepCompletionConditionResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.RecipeIDKey:                        request.RecipeID,
		keys.RecipeStepIDKey:                    request.RecipeStepID,
		keys.RecipeStepCompletionConditionIDKey: request.RecipeStepCompletionConditionID,
	}, span, s.logger)

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

func (s *serviceImpl) ArchiveRecipeStepIngredient(ctx context.Context, request *mealplanning.ArchiveRecipeStepIngredientRequest) (*mealplanning.ArchiveRecipeStepIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.RecipeIDKey:               request.RecipeID,
		keys.RecipeStepIDKey:           request.RecipeStepID,
		keys.RecipeStepIngredientIDKey: request.RecipeStepIngredientID,
	}, span, s.logger)

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

func (s *serviceImpl) ArchiveRecipeStepInstrument(ctx context.Context, request *mealplanning.ArchiveRecipeStepInstrumentRequest) (*mealplanning.ArchiveRecipeStepInstrumentResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.RecipeIDKey:               request.RecipeID,
		keys.RecipeStepIDKey:           request.RecipeStepID,
		keys.RecipeStepInstrumentIDKey: request.RecipeStepInstrumentID,
	}, span, s.logger)

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

func (s *serviceImpl) ArchiveRecipeStepProduct(ctx context.Context, request *mealplanning.ArchiveRecipeStepProductRequest) (*mealplanning.ArchiveRecipeStepProductResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.RecipeIDKey:            request.RecipeID,
		keys.RecipeStepIDKey:        request.RecipeStepID,
		keys.RecipeStepProductIDKey: request.RecipeStepProductID,
	}, span, s.logger)

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

func (s *serviceImpl) ArchiveRecipeStepVessel(ctx context.Context, request *mealplanning.ArchiveRecipeStepVesselRequest) (*mealplanning.ArchiveRecipeStepVesselResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.RecipeIDKey:           request.RecipeID,
		keys.RecipeStepIDKey:       request.RecipeStepID,
		keys.RecipeStepVesselIDKey: request.RecipeStepVesselID,
	}, span, s.logger)

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

func (s *serviceImpl) CloneRecipe(ctx context.Context, request *mealplanning.CloneRecipeRequest) (*mealplanning.CloneRecipeResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.RecipeIDKey: request.RecipeID,
	}, span, s.logger)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching session context data")
	}

	cloned, err := s.recipeManager.CloneRecipe(ctx, request.RecipeID, sessionContextData.GetUserID())
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "cloning recipe")
	}

	x := &mealplanning.CloneRecipeResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Cloned: converters.ConvertRecipeToGRPCRecipe(cloned),
	}

	return x, nil
}

func (s *serviceImpl) CreateRecipe(ctx context.Context, request *mealplanning.CreateRecipeRequest) (*mealplanning.CreateRecipeResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching session context data")
	}

	input := converters.ConvertGRPCRecipeCreationRequestInputToRecipeCreationRequestInput(request.Input)

	created, err := s.recipeManager.CreateRecipe(ctx, sessionContextData.GetUserID(), input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating recipe")
	}

	x := &mealplanning.CreateRecipeResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Created: converters.ConvertRecipeToGRPCRecipe(created),
	}

	return x, nil
}

func (s *serviceImpl) CreateRecipePrepTask(ctx context.Context, request *mealplanning.CreateRecipePrepTaskRequest) (*mealplanning.CreateRecipePrepTaskResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.RecipeIDKey: request.RecipeID,
	}, span, s.logger)

	input := converters.ConvertGRPCRecipePrepTaskCreationRequestInputToRecipePrepTaskCreationRequestInput(request.Input)

	created, err := s.recipeManager.CreateRecipePrepTask(ctx, request.RecipeID, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating recipe prep task")
	}

	x := &mealplanning.CreateRecipePrepTaskResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Created: converters.ConvertRecipePrepTaskToGRPCRecipePrepTask(created),
	}

	return x, nil
}

func (s *serviceImpl) CreateRecipeRating(ctx context.Context, request *mealplanning.CreateRecipeRatingRequest) (*mealplanning.CreateRecipeRatingResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.RecipeIDKey: request.RecipeID,
	}, span, s.logger)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to get session context data")
	}

	input := converters.ConvertGRPCRecipeRatingCreationRequestInputToRecipeRatingCreationRequestInput(request.Input)
	input.ByUser = sessionContextData.GetUserID()

	created, err := s.recipeManager.CreateRecipeRating(ctx, request.RecipeID, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating recipe rating")
	}

	x := &mealplanning.CreateRecipeRatingResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Created: converters.ConvertRecipeRatingToGRPCRecipeRating(created),
	}

	return x, nil
}

func (s *serviceImpl) CreateRecipeStep(ctx context.Context, request *mealplanning.CreateRecipeStepRequest) (*mealplanning.CreateRecipeStepResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.RecipeIDKey: request.RecipeID,
	}, span, s.logger)

	recipeStepInput := converters.ConvertGRPCRecipeStepCreationRequestInputToRecipeStepCreationRequestInput(request.Input)

	created, err := s.recipeManager.CreateRecipeStep(ctx, request.RecipeID, recipeStepInput)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating recipe step")
	}

	x := &mealplanning.CreateRecipeStepResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Created: converters.ConvertRecipeStepToGRPCRecipeStep(created),
	}

	return x, nil
}

func (s *serviceImpl) CreateRecipeStepCompletionCondition(ctx context.Context, request *mealplanning.CreateRecipeStepCompletionConditionRequest) (*mealplanning.CreateRecipeStepCompletionConditionResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.RecipeIDKey:     request.RecipeID,
		keys.RecipeStepIDKey: request.RecipeStepID,
	}, span, s.logger)

	creationInput := converters.ConvertGRPCRecipeStepCompletionConditionForExistingRecipeCreationRequestInputToRecipeStepCompletionConditionForExistingRecipeCreationRequestInput(request.Input)

	created, err := s.recipeManager.CreateRecipeStepCompletionCondition(ctx, request.RecipeID, request.RecipeStepID, creationInput)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating recipe step completion condition")
	}

	x := &mealplanning.CreateRecipeStepCompletionConditionResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Created: converters.ConvertRecipeStepCompletionConditionToGRPCRecipeStepCompletionCondition(created),
	}

	return x, nil
}

func (s *serviceImpl) CreateRecipeStepIngredient(ctx context.Context, request *mealplanning.CreateRecipeStepIngredientRequest) (*mealplanning.CreateRecipeStepIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.RecipeIDKey:     request.RecipeID,
		keys.RecipeStepIDKey: request.RecipeStepID,
	}, span, s.logger)

	creationInput := converters.ConvertGRPCRecipeStepIngredientCreationRequestInputToRecipeStepIngredientCreationRequestInput(request.Input)

	created, err := s.recipeManager.CreateRecipeStepIngredient(ctx, request.RecipeID, request.RecipeStepID, creationInput)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating recipe step ingredient")
	}

	x := &mealplanning.CreateRecipeStepIngredientResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Created: converters.ConvertRecipeStepIngredientToGRPCRecipeStepIngredient(created),
	}

	return x, nil
}

func (s *serviceImpl) CreateRecipeStepInstrument(ctx context.Context, request *mealplanning.CreateRecipeStepInstrumentRequest) (*mealplanning.CreateRecipeStepInstrumentResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.RecipeIDKey:     request.RecipeID,
		keys.RecipeStepIDKey: request.RecipeStepID,
	}, span, s.logger)

	creationInput := converters.ConvertGRPCRecipeStepInstrumentCreationRequestInputToRecipeStepInstrumentCreationRequestInput(request.Input)

	created, err := s.recipeManager.CreateRecipeStepInstrument(ctx, request.RecipeID, request.RecipeStepID, creationInput)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating recipe step instrument")
	}

	x := &mealplanning.CreateRecipeStepInstrumentResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Created: converters.ConvertRecipeStepInstrumentToGRPCRecipeStepInstrument(created),
	}

	return x, nil
}

func (s *serviceImpl) CreateRecipeStepProduct(ctx context.Context, request *mealplanning.CreateRecipeStepProductRequest) (*mealplanning.CreateRecipeStepProductResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.RecipeIDKey:     request.RecipeID,
		keys.RecipeStepIDKey: request.RecipeStepID,
	}, span, s.logger)

	creationInput := converters.ConvertGRPCRecipeStepProductCreationRequestInputToRecipeStepProductCreationRequestInput(request.Input)

	created, err := s.recipeManager.CreateRecipeStepProduct(ctx, request.RecipeID, request.RecipeStepID, creationInput)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating recipe step product")
	}

	x := &mealplanning.CreateRecipeStepProductResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Created: converters.ConvertRecipeStepProductToGRPCRecipeStepProduct(created),
	}

	return x, nil
}

func (s *serviceImpl) CreateRecipeStepVessel(ctx context.Context, request *mealplanning.CreateRecipeStepVesselRequest) (*mealplanning.CreateRecipeStepVesselResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.RecipeIDKey:     request.RecipeID,
		keys.RecipeStepIDKey: request.RecipeStepID,
	}, span, s.logger)

	creationInput := converters.ConvertGRPCRecipeStepVesselCreationRequestInputToRecipeStepVesselCreationRequestInput(request.Input)

	created, err := s.recipeManager.CreateRecipeStepVessel(ctx, request.RecipeID, request.RecipeStepID, creationInput)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating recipe step vessel")
	}

	x := &mealplanning.CreateRecipeStepVesselResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Created: converters.ConvertRecipeStepVesselToGRPCRecipeStepVessel(created),
	}

	return x, nil
}

func (s *serviceImpl) GetMermaidDiagramForRecipe(ctx context.Context, request *mealplanning.GetMermaidDiagramForRecipeRequest) (*mealplanning.GetMermaidDiagramForRecipeResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.RecipeIDKey: request.RecipeID,
	}, span, s.logger)

	mermaidDiagram, err := s.recipeManager.RecipeMermaid(ctx, request.RecipeID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to generate mermaid diagram")
	}

	x := &mealplanning.GetMermaidDiagramForRecipeResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Response: mermaidDiagram,
	}

	return x, nil
}

func (s *serviceImpl) GetRecipe(ctx context.Context, request *mealplanning.GetRecipeRequest) (*mealplanning.GetRecipeResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.RecipeIDKey: request.RecipeID,
	}, span, s.logger)

	recipe, err := s.recipeManager.ReadRecipe(ctx, request.RecipeID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to get recipe")
	}

	x := &mealplanning.GetRecipeResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Result: converters.ConvertRecipeToGRPCRecipe(recipe),
	}

	return x, nil
}

// TODO: very aware that this basically just emits a list of strings, and that that isn't very useful.
func (s *serviceImpl) EstimateRecipePrepTasks(ctx context.Context, request *mealplanning.EstimateRecipePrepTasksRequest) (*mealplanning.EstimateRecipePrepTasksResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.RecipeIDKey: request.RecipeID,
	}, span, s.logger)

	estimatedPrepSteps, err := s.recipeManager.RecipeEstimatedPrepSteps(ctx, request.RecipeID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to estimate prep steps")
	}

	x := &mealplanning.EstimateRecipePrepTasksResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	for _, step := range estimatedPrepSteps {
		x.Results = append(x.Results, converters.ConvertMealPlanTaskDatabaseCreationEstimateToGRPCMealPlanTask(step))
	}

	return x, nil
}

func (s *serviceImpl) GetRecipePrepTask(ctx context.Context, request *mealplanning.GetRecipePrepTaskRequest) (*mealplanning.GetRecipePrepTaskResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.RecipeIDKey:         request.RecipeID,
		keys.RecipePrepTaskIDKey: request.RecipePrepTaskID,
	}, span, s.logger)

	recipePrepTask, err := s.recipeManager.ReadRecipePrepTask(ctx, request.RecipeID, request.RecipePrepTaskID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to get recipe prep task")
	}

	x := &mealplanning.GetRecipePrepTaskResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Result: converters.ConvertRecipePrepTaskToGRPCRecipePrepTask(recipePrepTask),
	}

	return x, nil
}

func (s *serviceImpl) GetRecipePrepTasks(ctx context.Context, request *mealplanning.GetRecipePrepTasksRequest) (*mealplanning.GetRecipePrepTasksResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.RecipeIDKey: request.RecipeID,
	}, span, s.logger)

	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)

	recipePrepTasks, err := s.recipeManager.ListRecipePrepTask(ctx, request.RecipeID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to get recipe prep tasks")
	}

	x := &mealplanning.GetRecipePrepTasksResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Filter: grpcconverters.ConvertQueryFilterToGRPCQueryFilter(filter, recipePrepTasks.Pagination),
	}

	for _, recipePrepTask := range recipePrepTasks.Data {
		x.Results = append(x.Results, converters.ConvertRecipePrepTaskToGRPCRecipePrepTask(recipePrepTask))
	}

	return x, nil
}

func (s *serviceImpl) GetRecipeRating(ctx context.Context, request *mealplanning.GetRecipeRatingRequest) (*mealplanning.GetRecipeRatingResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.RecipeIDKey:       request.RecipeID,
		keys.RecipeRatingIDKey: request.RecipeRatingID,
	}, span, s.logger)

	recipeRating, err := s.recipeManager.ReadRecipeRating(ctx, request.RecipeID, request.RecipeRatingID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to get recipe rating")
	}

	x := &mealplanning.GetRecipeRatingResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Result: converters.ConvertRecipeRatingToGRPCRecipeRating(recipeRating),
	}

	return x, nil
}

func (s *serviceImpl) GetRecipeRatingsForRecipe(ctx context.Context, request *mealplanning.GetRecipeRatingsForRecipeRequest) (*mealplanning.GetRecipeRatingsForRecipeResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.RecipeIDKey: request.RecipeID,
	}, span, s.logger)

	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)

	recipeRatings, err := s.recipeManager.ListRecipeRatings(ctx, request.RecipeID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to get recipe ratings")
	}

	x := &mealplanning.GetRecipeRatingsForRecipeResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	for _, recipePrepTask := range recipeRatings.Data {
		x.Results = append(x.Results, converters.ConvertRecipeRatingToGRPCRecipeRating(recipePrepTask))
	}

	return x, nil
}

func (s *serviceImpl) GetRecipeStep(ctx context.Context, request *mealplanning.GetRecipeStepRequest) (*mealplanning.GetRecipeStepResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.RecipeIDKey:     request.RecipeID,
		keys.RecipeStepIDKey: request.RecipeStepID,
	}, span, s.logger)

	recipeStep, err := s.recipeManager.ReadRecipeStep(ctx, request.RecipeID, request.RecipeStepID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to get recipe step")
	}

	x := &mealplanning.GetRecipeStepResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Result: converters.ConvertRecipeStepToGRPCRecipeStep(recipeStep),
	}

	return x, nil
}

func (s *serviceImpl) GetRecipeStepCompletionCondition(ctx context.Context, request *mealplanning.GetRecipeStepCompletionConditionRequest) (*mealplanning.GetRecipeStepCompletionConditionResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.RecipeIDKey:                        request.RecipeID,
		keys.RecipeStepIDKey:                    request.RecipeStepID,
		keys.RecipeStepCompletionConditionIDKey: request.RecipeStepCompletionConditionID,
	}, span, s.logger)

	recipeStepCompletionCondition, err := s.recipeManager.ReadRecipeStepCompletionCondition(ctx, request.RecipeID, request.RecipeStepID, request.RecipeStepCompletionConditionID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to get recipe step completion condition")
	}

	x := &mealplanning.GetRecipeStepCompletionConditionResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Result: converters.ConvertRecipeStepCompletionConditionToGRPCRecipeStepCompletionCondition(recipeStepCompletionCondition),
	}

	return x, nil
}

func (s *serviceImpl) GetRecipeStepCompletionConditions(ctx context.Context, request *mealplanning.GetRecipeStepCompletionConditionsRequest) (*mealplanning.GetRecipeStepCompletionConditionsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.RecipeIDKey:     request.RecipeID,
		keys.RecipeStepIDKey: request.RecipeStepID,
	}, span, s.logger)

	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)

	recipeStepCompletionConditions, err := s.recipeManager.ListRecipeStepCompletionConditions(ctx, request.RecipeID, request.RecipeStepID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to get recipe step completion conditions")
	}

	x := &mealplanning.GetRecipeStepCompletionConditionsResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	for _, recipePrepTask := range recipeStepCompletionConditions.Data {
		x.Results = append(x.Results, converters.ConvertRecipeStepCompletionConditionToGRPCRecipeStepCompletionCondition(recipePrepTask))
	}

	return x, nil
}

func (s *serviceImpl) GetRecipeStepIngredient(ctx context.Context, request *mealplanning.GetRecipeStepIngredientRequest) (*mealplanning.GetRecipeStepIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.RecipeIDKey:               request.RecipeID,
		keys.RecipeStepIDKey:           request.RecipeStepID,
		keys.RecipeStepIngredientIDKey: request.RecipeStepIngredientID,
	}, span, s.logger)

	recipeStepIngredient, err := s.recipeManager.ReadRecipeStepIngredient(ctx, request.RecipeID, request.RecipeStepID, request.RecipeStepIngredientID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to get recipe step ingredient")
	}

	x := &mealplanning.GetRecipeStepIngredientResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Result: converters.ConvertRecipeStepIngredientToGRPCRecipeStepIngredient(recipeStepIngredient),
	}

	return x, nil
}

func (s *serviceImpl) GetRecipeStepIngredients(ctx context.Context, request *mealplanning.GetRecipeStepIngredientsRequest) (*mealplanning.GetRecipeStepIngredientsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.RecipeIDKey:     request.RecipeID,
		keys.RecipeStepIDKey: request.RecipeStepID,
	}, span, s.logger)

	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)

	recipeStepIngredients, err := s.recipeManager.ListRecipeStepIngredients(ctx, request.RecipeID, request.RecipeStepID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to get recipe step ingredients")
	}

	x := &mealplanning.GetRecipeStepIngredientsResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	for _, recipeStepIngredient := range recipeStepIngredients.Data {
		x.Results = append(x.Results, converters.ConvertRecipeStepIngredientToGRPCRecipeStepIngredient(recipeStepIngredient))
	}

	return x, nil
}

func (s *serviceImpl) GetRecipeStepInstrument(ctx context.Context, request *mealplanning.GetRecipeStepInstrumentRequest) (*mealplanning.GetRecipeStepInstrumentResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.RecipeIDKey:               request.RecipeID,
		keys.RecipeStepIDKey:           request.RecipeStepID,
		keys.RecipeStepInstrumentIDKey: request.RecipeStepInstrumentID,
	}, span, s.logger)

	recipeStepInstrument, err := s.recipeManager.ReadRecipeStepInstrument(ctx, request.RecipeID, request.RecipeStepID, request.RecipeStepInstrumentID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to get recipe step instrument")
	}

	x := &mealplanning.GetRecipeStepInstrumentResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Result: converters.ConvertRecipeStepInstrumentToGRPCRecipeStepInstrument(recipeStepInstrument),
	}

	return x, nil
}

func (s *serviceImpl) GetRecipeStepInstruments(ctx context.Context, request *mealplanning.GetRecipeStepInstrumentsRequest) (*mealplanning.GetRecipeStepInstrumentsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.RecipeIDKey:     request.RecipeID,
		keys.RecipeStepIDKey: request.RecipeStepID,
	}, span, s.logger)

	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)

	recipeStepInstruments, err := s.recipeManager.ListRecipeStepInstruments(ctx, request.RecipeID, request.RecipeStepID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to get recipe step instruments")
	}

	x := &mealplanning.GetRecipeStepInstrumentsResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	for _, recipeStepInstrument := range recipeStepInstruments.Data {
		x.Results = append(x.Results, converters.ConvertRecipeStepInstrumentToGRPCRecipeStepInstrument(recipeStepInstrument))
	}

	return x, nil
}

func (s *serviceImpl) GetRecipeStepProduct(ctx context.Context, request *mealplanning.GetRecipeStepProductRequest) (*mealplanning.GetRecipeStepProductResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.RecipeIDKey:            request.RecipeID,
		keys.RecipeStepIDKey:        request.RecipeStepID,
		keys.RecipeStepProductIDKey: request.RecipeStepProductID,
	}, span, s.logger)

	recipeStepProduct, err := s.recipeManager.ReadRecipeStepProduct(ctx, request.RecipeID, request.RecipeStepID, request.RecipeStepProductID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to get recipe step")
	}

	x := &mealplanning.GetRecipeStepProductResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Result: converters.ConvertRecipeStepProductToGRPCRecipeStepProduct(recipeStepProduct),
	}

	return x, nil
}

func (s *serviceImpl) GetRecipeStepProducts(ctx context.Context, request *mealplanning.GetRecipeStepProductsRequest) (*mealplanning.GetRecipeStepProductsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.RecipeIDKey:     request.RecipeID,
		keys.RecipeStepIDKey: request.RecipeStepID,
	}, span, s.logger)

	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)

	recipeStepProducts, err := s.recipeManager.ListRecipeStepProducts(ctx, request.RecipeID, request.RecipeStepID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to get recipe step product")
	}

	x := &mealplanning.GetRecipeStepProductsResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	for _, recipeStepProduct := range recipeStepProducts.Data {
		x.Results = append(x.Results, converters.ConvertRecipeStepProductToGRPCRecipeStepProduct(recipeStepProduct))
	}

	return x, nil
}

func (s *serviceImpl) GetRecipeStepVessel(ctx context.Context, request *mealplanning.GetRecipeStepVesselRequest) (*mealplanning.GetRecipeStepVesselResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.RecipeIDKey:           request.RecipeID,
		keys.RecipeStepIDKey:       request.RecipeStepID,
		keys.RecipeStepVesselIDKey: request.RecipeStepVesselID,
	}, span, s.logger)

	recipeStepVessel, err := s.recipeManager.ReadRecipeStepVessel(ctx, request.RecipeID, request.RecipeStepID, request.RecipeStepVesselID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to get recipe step")
	}

	x := &mealplanning.GetRecipeStepVesselResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Result: converters.ConvertRecipeStepVesselToGRPCRecipeStepVessel(recipeStepVessel),
	}

	return x, nil
}

func (s *serviceImpl) GetRecipeStepVessels(ctx context.Context, request *mealplanning.GetRecipeStepVesselsRequest) (*mealplanning.GetRecipeStepVesselsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.RecipeIDKey:     request.RecipeID,
		keys.RecipeStepIDKey: request.RecipeStepID,
	}, span, s.logger)

	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)

	recipeStepVessels, err := s.recipeManager.ListRecipeStepVessels(ctx, request.RecipeID, request.RecipeStepID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to get recipe step vessels")
	}

	x := &mealplanning.GetRecipeStepVesselsResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	for _, recipeStepVessel := range recipeStepVessels.Data {
		x.Results = append(x.Results, converters.ConvertRecipeStepVesselToGRPCRecipeStepVessel(recipeStepVessel))
	}

	return x, nil
}

func (s *serviceImpl) GetRecipeSteps(ctx context.Context, request *mealplanning.GetRecipeStepsRequest) (*mealplanning.GetRecipeStepsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.RecipeIDKey: request.RecipeID,
	}, span, s.logger)

	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)

	recipeSteps, err := s.recipeManager.ListRecipeSteps(ctx, request.RecipeID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to get recipe step vessels")
	}

	x := &mealplanning.GetRecipeStepsResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	for _, recipeStep := range recipeSteps.Data {
		x.Results = append(x.Results, converters.ConvertRecipeStepToGRPCRecipeStep(recipeStep))
	}

	return x, nil
}

func (s *serviceImpl) GetRecipes(ctx context.Context, request *mealplanning.GetRecipesRequest) (*mealplanning.GetRecipesResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	recipes, err := s.recipeManager.ListRecipes(ctx, grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching list of recipes")
	}

	x := &mealplanning.GetRecipesResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	for _, recipe := range recipes.Data {
		x.Results = append(x.Results, converters.ConvertRecipeToGRPCRecipe(recipe))
	}

	return x, nil
}

func (s *serviceImpl) SearchForRecipes(ctx context.Context, request *mealplanning.SearchForRecipesRequest) (*mealplanning.SearchForRecipesResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.SearchQueryKey: request.Query,
	}, span, s.logger)

	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)

	recipes, _, err := s.recipeManager.SearchRecipes(ctx, request.Query, request.UseSearchService, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "searching for recipes")
	}

	x := &mealplanning.SearchForRecipesResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	for _, recipe := range recipes {
		x.Results = append(x.Results, converters.ConvertRecipeToGRPCRecipe(recipe))
	}

	return x, nil
}

func (s *serviceImpl) UpdateRecipe(ctx context.Context, request *mealplanning.UpdateRecipeRequest) (*mealplanning.UpdateRecipeResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.RecipeIDKey: request.RecipeID,
	}, span, s.logger)

	input := converters.ConvertGRPCRecipeUpdateRequestInputToRecipeUpdateRequestInput(request.Input)

	if err := s.recipeManager.UpdateRecipe(ctx, request.RecipeID, input); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed updating recipe")
	}

	updated, err := s.recipeManager.ReadRecipe(ctx, request.RecipeID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed reading recipe")
	}

	x := &mealplanning.UpdateRecipeResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Updated: converters.ConvertRecipeToGRPCRecipe(updated),
	}

	return x, nil
}

func (s *serviceImpl) UpdateRecipePrepTask(ctx context.Context, request *mealplanning.UpdateRecipePrepTaskRequest) (*mealplanning.UpdateRecipePrepTaskResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.RecipeIDKey:         request.RecipeID,
		keys.RecipePrepTaskIDKey: request.RecipePrepTaskID,
	}, span, s.logger)

	input := converters.ConvertGRPCRecipePrepTaskUpdateRequestInputToRecipePrepTaskUpdateRequestInput(request.Input)

	if err := s.recipeManager.UpdateRecipePrepTask(ctx, request.RecipeID, request.RecipePrepTaskID, input); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed updating recipe prep task")
	}

	updated, err := s.recipeManager.ReadRecipePrepTask(ctx, request.RecipeID, request.RecipePrepTaskID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed reading recipe prep task")
	}

	x := &mealplanning.UpdateRecipePrepTaskResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Updated: converters.ConvertRecipePrepTaskToGRPCRecipePrepTask(updated),
	}

	return x, nil
}

func (s *serviceImpl) UpdateRecipeRating(ctx context.Context, request *mealplanning.UpdateRecipeRatingRequest) (*mealplanning.UpdateRecipeRatingResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.RecipeIDKey:       request.RecipeID,
		keys.RecipeRatingIDKey: request.RecipeRatingID,
	}, span, s.logger)

	input := converters.ConvertGRPCRecipeRatingUpdateRequestInputToRecipeRatingUpdateRequestInput(request.Input)

	if err := s.recipeManager.UpdateRecipeRating(ctx, request.RecipeID, request.RecipeRatingID, input); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed updating recipe rating")
	}

	updated, err := s.recipeManager.ReadRecipeRating(ctx, request.RecipeID, request.RecipeRatingID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed reading recipe rating")
	}

	x := &mealplanning.UpdateRecipeRatingResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Updated: converters.ConvertRecipeRatingToGRPCRecipeRating(updated),
	}

	return x, nil
}

func (s *serviceImpl) UpdateRecipeStep(ctx context.Context, request *mealplanning.UpdateRecipeStepRequest) (*mealplanning.UpdateRecipeStepResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.RecipeIDKey:     request.RecipeID,
		keys.RecipeStepIDKey: request.RecipeStepID,
	}, span, s.logger)

	input := converters.ConvertGRPCRecipeStepUpdateRequestInputToRecipeStepUpdateRequestInput(request.Input)

	if err := s.recipeManager.UpdateRecipeStep(ctx, request.RecipeID, request.RecipeStepID, input); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed updating recipe step")
	}

	updated, err := s.recipeManager.ReadRecipeStep(ctx, request.RecipeID, request.RecipeStepID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed reading recipe step")
	}

	x := &mealplanning.UpdateRecipeStepResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Updated: converters.ConvertRecipeStepToGRPCRecipeStep(updated),
	}

	return x, nil
}

func (s *serviceImpl) UpdateRecipeStepCompletionCondition(ctx context.Context, request *mealplanning.UpdateRecipeStepCompletionConditionRequest) (*mealplanning.UpdateRecipeStepCompletionConditionResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.RecipeIDKey:                        request.RecipeID,
		keys.RecipeStepIDKey:                    request.RecipeStepID,
		keys.RecipeStepCompletionConditionIDKey: request.RecipeStepCompletionConditionID,
	}, span, s.logger)

	input := converters.ConvertGRPCRecipeStepCompletionConditionUpdateRequestInputToRecipeStepCompletionConditionUpdateRequestInput(request.Input)

	if err := s.recipeManager.UpdateRecipeStepCompletionCondition(ctx, request.RecipeID, request.RecipeStepID, request.RecipeStepCompletionConditionID, input); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed updating recipe step completion condition")
	}

	updated, err := s.recipeManager.ReadRecipeStepCompletionCondition(ctx, request.RecipeID, request.RecipeStepID, request.RecipeStepCompletionConditionID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed reading recipe step completion condition")
	}

	x := &mealplanning.UpdateRecipeStepCompletionConditionResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Updated: converters.ConvertRecipeStepCompletionConditionToGRPCRecipeStepCompletionCondition(updated),
	}

	return x, nil
}

func (s *serviceImpl) UpdateRecipeStepIngredient(ctx context.Context, request *mealplanning.UpdateRecipeStepIngredientRequest) (*mealplanning.UpdateRecipeStepIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.RecipeIDKey:               request.RecipeID,
		keys.RecipeStepIDKey:           request.RecipeStepID,
		keys.RecipeStepIngredientIDKey: request.RecipeStepIngredientID,
	}, span, s.logger)

	input := converters.ConvertGRPCRecipeStepIngredientUpdateRequestInputToRecipeStepIngredientUpdateRequestInput(request.Input)

	if err := s.recipeManager.UpdateRecipeStepIngredient(ctx, request.RecipeID, request.RecipeStepID, request.RecipeStepIngredientID, input); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed updating recipe step ingredient")
	}

	updated, err := s.recipeManager.ReadRecipeStepIngredient(ctx, request.RecipeID, request.RecipeStepID, request.RecipeStepIngredientID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed reading recipe step ingredient")
	}

	x := &mealplanning.UpdateRecipeStepIngredientResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Updated: converters.ConvertRecipeStepIngredientToGRPCRecipeStepIngredient(updated),
	}

	return x, nil
}

func (s *serviceImpl) UpdateRecipeStepInstrument(ctx context.Context, request *mealplanning.UpdateRecipeStepInstrumentRequest) (*mealplanning.UpdateRecipeStepInstrumentResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.RecipeIDKey:               request.RecipeID,
		keys.RecipeStepIDKey:           request.RecipeStepID,
		keys.RecipeStepInstrumentIDKey: request.RecipeStepInstrumentID,
	}, span, s.logger)

	input := converters.ConvertGRPCRecipeStepInstrumentUpdateRequestInputToRecipeStepInstrumentUpdateRequestInput(request.Input)

	if err := s.recipeManager.UpdateRecipeStepInstrument(ctx, request.RecipeID, request.RecipeStepID, request.RecipeStepInstrumentID, input); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed updating recipe step instrument")
	}

	updated, err := s.recipeManager.ReadRecipeStepInstrument(ctx, request.RecipeID, request.RecipeStepID, request.RecipeStepInstrumentID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed reading recipe step instrument")
	}

	x := &mealplanning.UpdateRecipeStepInstrumentResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Updated: converters.ConvertRecipeStepInstrumentToGRPCRecipeStepInstrument(updated),
	}

	return x, nil
}

func (s *serviceImpl) UpdateRecipeStepProduct(ctx context.Context, request *mealplanning.UpdateRecipeStepProductRequest) (*mealplanning.UpdateRecipeStepProductResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.RecipeIDKey:            request.RecipeID,
		keys.RecipeStepIDKey:        request.RecipeStepID,
		keys.RecipeStepProductIDKey: request.RecipeStepProductID,
	}, span, s.logger)

	input := converters.ConvertGRPCRecipeStepProductUpdateRequestInputToRecipeStepProductUpdateRequestInput(request.Input)

	if err := s.recipeManager.UpdateRecipeStepProduct(ctx, request.RecipeID, request.RecipeStepID, request.RecipeStepProductID, input); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed updating recipe step product")
	}

	updated, err := s.recipeManager.ReadRecipeStepProduct(ctx, request.RecipeID, request.RecipeStepID, request.RecipeStepProductID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed reading recipe step product")
	}

	x := &mealplanning.UpdateRecipeStepProductResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Updated: converters.ConvertRecipeStepProductToGRPCRecipeStepProduct(updated),
	}

	return x, nil
}

func (s *serviceImpl) UpdateRecipeStepVessel(ctx context.Context, request *mealplanning.UpdateRecipeStepVesselRequest) (*mealplanning.UpdateRecipeStepVesselResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.RecipeIDKey:           request.RecipeID,
		keys.RecipeStepIDKey:       request.RecipeStepID,
		keys.RecipeStepVesselIDKey: request.RecipeStepVesselID,
	}, span, s.logger)

	input := converters.ConvertGRPCRecipeStepVesselUpdateRequestInputToRecipeStepVesselUpdateRequestInput(request.Input)

	if err := s.recipeManager.UpdateRecipeStepVessel(ctx, request.RecipeID, request.RecipeStepID, request.RecipeStepVesselID, input); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed updating recipe step vessel")
	}

	updated, err := s.recipeManager.ReadRecipeStepVessel(ctx, request.RecipeID, request.RecipeStepID, request.RecipeStepVesselID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed reading recipe step vessel")
	}

	x := &mealplanning.UpdateRecipeStepVesselResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Updated: converters.ConvertRecipeStepVesselToGRPCRecipeStepVessel(updated),
	}

	return x, nil
}
