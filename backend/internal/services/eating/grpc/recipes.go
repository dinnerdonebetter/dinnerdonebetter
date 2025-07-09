package eatinggrpc

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
)

func (s *ServiceImpl) ArchiveRecipe(ctx context.Context, request *messages.ArchiveRecipeRequest) (*messages.ArchiveRecipeResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) ArchiveRecipePrepTask(ctx context.Context, request *messages.ArchiveRecipePrepTaskRequest) (*messages.ArchiveRecipePrepTaskResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	if err := s.recipeManager.ArchiveRecipePrepTask(ctx, request.GetRecipeID(), request.GetRecipePrepTaskID()); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "archiving recipe prep task")
	}

	return &messages.ArchiveRecipePrepTaskResponse{}, nil
}

func (s *ServiceImpl) ArchiveRecipeRating(ctx context.Context, request *messages.ArchiveRecipeRatingRequest) (*messages.ArchiveRecipeRatingResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) ArchiveRecipeStep(ctx context.Context, request *messages.ArchiveRecipeStepRequest) (*messages.ArchiveRecipeStepResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) ArchiveRecipeStepCompletionCondition(ctx context.Context, request *messages.ArchiveRecipeStepCompletionConditionRequest) (*messages.ArchiveRecipeStepCompletionConditionResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) ArchiveRecipeStepIngredient(ctx context.Context, request *messages.ArchiveRecipeStepIngredientRequest) (*messages.ArchiveRecipeStepIngredientResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) ArchiveRecipeStepInstrument(ctx context.Context, request *messages.ArchiveRecipeStepInstrumentRequest) (*messages.ArchiveRecipeStepInstrumentResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) ArchiveRecipeStepProduct(ctx context.Context, request *messages.ArchiveRecipeStepProductRequest) (*messages.ArchiveRecipeStepProductResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) ArchiveRecipeStepVessel(ctx context.Context, request *messages.ArchiveRecipeStepVesselRequest) (*messages.ArchiveRecipeStepVesselResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CloneRecipe(ctx context.Context, request *messages.CloneRecipeRequest) (*messages.CloneRecipeResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateRecipe(ctx context.Context, request *messages.CreateRecipeRequest) (*messages.CreateRecipeResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateRecipePrepTask(ctx context.Context, request *messages.CreateRecipePrepTaskRequest) (*messages.CreateRecipePrepTaskResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateRecipeRating(ctx context.Context, request *messages.CreateRecipeRatingRequest) (*messages.CreateRecipeRatingResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateRecipeStep(ctx context.Context, request *messages.CreateRecipeStepRequest) (*messages.CreateRecipeStepResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateRecipeStepCompletionCondition(ctx context.Context, request *messages.CreateRecipeStepCompletionConditionRequest) (*messages.CreateRecipeStepCompletionConditionResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateRecipeStepIngredient(ctx context.Context, request *messages.CreateRecipeStepIngredientRequest) (*messages.CreateRecipeStepIngredientResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateRecipeStepInstrument(ctx context.Context, request *messages.CreateRecipeStepInstrumentRequest) (*messages.CreateRecipeStepInstrumentResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateRecipeStepProduct(ctx context.Context, request *messages.CreateRecipeStepProductRequest) (*messages.CreateRecipeStepProductResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateRecipeStepVessel(ctx context.Context, request *messages.CreateRecipeStepVesselRequest) (*messages.CreateRecipeStepVesselResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetMermaidDiagramForRecipe(ctx context.Context, request *messages.GetMermaidDiagramForRecipeRequest) (*messages.GetMermaidDiagramForRecipeResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipe(ctx context.Context, request *messages.GetRecipeRequest) (*messages.GetRecipeResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipeMealPlanTasks(ctx context.Context, request *messages.GetRecipeMealPlanTasksRequest) (*messages.GetRecipeMealPlanTasksResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipePrepTask(ctx context.Context, request *messages.GetRecipePrepTaskRequest) (*messages.GetRecipePrepTaskResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipePrepTasks(ctx context.Context, request *messages.GetRecipePrepTasksRequest) (*messages.GetRecipePrepTasksResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipeRating(ctx context.Context, request *messages.GetRecipeRatingRequest) (*messages.GetRecipeRatingResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipeRatingsForRecipe(ctx context.Context, request *messages.GetRecipeRatingsForRecipeRequest) (*messages.GetRecipeRatingsForRecipeResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipeStep(ctx context.Context, request *messages.GetRecipeStepRequest) (*messages.GetRecipeStepResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipeStepCompletionCondition(ctx context.Context, request *messages.GetRecipeStepCompletionConditionRequest) (*messages.GetRecipeStepCompletionConditionResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipeStepCompletionConditions(ctx context.Context, request *messages.GetRecipeStepCompletionConditionsRequest) (*messages.GetRecipeStepCompletionConditionsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipeStepIngredient(ctx context.Context, request *messages.GetRecipeStepIngredientRequest) (*messages.GetRecipeStepIngredientResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipeStepIngredients(ctx context.Context, request *messages.GetRecipeStepIngredientsRequest) (*messages.GetRecipeStepIngredientsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipeStepInstrument(ctx context.Context, request *messages.GetRecipeStepInstrumentRequest) (*messages.GetRecipeStepInstrumentResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipeStepInstruments(ctx context.Context, request *messages.GetRecipeStepInstrumentsRequest) (*messages.GetRecipeStepInstrumentsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipeStepProduct(ctx context.Context, request *messages.GetRecipeStepProductRequest) (*messages.GetRecipeStepProductResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipeStepProducts(ctx context.Context, request *messages.GetRecipeStepProductsRequest) (*messages.GetRecipeStepProductsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipeStepVessel(ctx context.Context, request *messages.GetRecipeStepVesselRequest) (*messages.GetRecipeStepVesselResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipeStepVessels(ctx context.Context, request *messages.GetRecipeStepVesselsRequest) (*messages.GetRecipeStepVesselsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipeSteps(ctx context.Context, request *messages.GetRecipeStepsRequest) (*messages.GetRecipeStepsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipes(ctx context.Context, request *messages.GetRecipesRequest) (*messages.GetRecipesResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) SearchForRecipes(ctx context.Context, request *messages.SearchForRecipesRequest) (*messages.SearchForRecipesResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateRecipe(ctx context.Context, request *messages.UpdateRecipeRequest) (*messages.UpdateRecipeResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateRecipePrepTask(ctx context.Context, request *messages.UpdateRecipePrepTaskRequest) (*messages.UpdateRecipePrepTaskResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateRecipeRating(ctx context.Context, request *messages.UpdateRecipeRatingRequest) (*messages.UpdateRecipeRatingResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateRecipeStep(ctx context.Context, request *messages.UpdateRecipeStepRequest) (*messages.UpdateRecipeStepResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateRecipeStepCompletionCondition(ctx context.Context, request *messages.UpdateRecipeStepCompletionConditionRequest) (*messages.UpdateRecipeStepCompletionConditionResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateRecipeStepIngredient(ctx context.Context, request *messages.UpdateRecipeStepIngredientRequest) (*messages.UpdateRecipeStepIngredientResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateRecipeStepInstrument(ctx context.Context, request *messages.UpdateRecipeStepInstrumentRequest) (*messages.UpdateRecipeStepInstrumentResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateRecipeStepProduct(ctx context.Context, request *messages.UpdateRecipeStepProductRequest) (*messages.UpdateRecipeStepProductResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateRecipeStepVessel(ctx context.Context, request *messages.UpdateRecipeStepVesselRequest) (*messages.UpdateRecipeStepVesselResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}
