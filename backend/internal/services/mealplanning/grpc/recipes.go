package grpc

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
)

func (s *ServiceImpl) ArchiveRecipe(ctx context.Context, request *mealplanning.ArchiveRecipeRequest) (*mealplanning.ArchiveRecipeResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) ArchiveRecipePrepTask(ctx context.Context, request *mealplanning.ArchiveRecipePrepTaskRequest) (*mealplanning.ArchiveRecipePrepTaskResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	if err := s.recipeManager.ArchiveRecipePrepTask(ctx, request.GetRecipeID(), request.GetRecipePrepTaskID()); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "archiving recipe prep task")
	}

	return &mealplanning.ArchiveRecipePrepTaskResponse{}, nil
}

func (s *ServiceImpl) ArchiveRecipeRating(ctx context.Context, request *mealplanning.ArchiveRecipeRatingRequest) (*mealplanning.ArchiveRecipeRatingResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) ArchiveRecipeStep(ctx context.Context, request *mealplanning.ArchiveRecipeStepRequest) (*mealplanning.ArchiveRecipeStepResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) ArchiveRecipeStepCompletionCondition(ctx context.Context, request *mealplanning.ArchiveRecipeStepCompletionConditionRequest) (*mealplanning.ArchiveRecipeStepCompletionConditionResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) ArchiveRecipeStepIngredient(ctx context.Context, request *mealplanning.ArchiveRecipeStepIngredientRequest) (*mealplanning.ArchiveRecipeStepIngredientResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) ArchiveRecipeStepInstrument(ctx context.Context, request *mealplanning.ArchiveRecipeStepInstrumentRequest) (*mealplanning.ArchiveRecipeStepInstrumentResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) ArchiveRecipeStepProduct(ctx context.Context, request *mealplanning.ArchiveRecipeStepProductRequest) (*mealplanning.ArchiveRecipeStepProductResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) ArchiveRecipeStepVessel(ctx context.Context, request *mealplanning.ArchiveRecipeStepVesselRequest) (*mealplanning.ArchiveRecipeStepVesselResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CloneRecipe(ctx context.Context, request *mealplanning.CloneRecipeRequest) (*mealplanning.CloneRecipeResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateRecipe(ctx context.Context, request *mealplanning.CreateRecipeRequest) (*mealplanning.CreateRecipeResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateRecipePrepTask(ctx context.Context, request *mealplanning.CreateRecipePrepTaskRequest) (*mealplanning.CreateRecipePrepTaskResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateRecipeRating(ctx context.Context, request *mealplanning.CreateRecipeRatingRequest) (*mealplanning.CreateRecipeRatingResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateRecipeStep(ctx context.Context, request *mealplanning.CreateRecipeStepRequest) (*mealplanning.CreateRecipeStepResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateRecipeStepCompletionCondition(ctx context.Context, request *mealplanning.CreateRecipeStepCompletionConditionRequest) (*mealplanning.CreateRecipeStepCompletionConditionResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateRecipeStepIngredient(ctx context.Context, request *mealplanning.CreateRecipeStepIngredientRequest) (*mealplanning.CreateRecipeStepIngredientResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateRecipeStepInstrument(ctx context.Context, request *mealplanning.CreateRecipeStepInstrumentRequest) (*mealplanning.CreateRecipeStepInstrumentResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateRecipeStepProduct(ctx context.Context, request *mealplanning.CreateRecipeStepProductRequest) (*mealplanning.CreateRecipeStepProductResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateRecipeStepVessel(ctx context.Context, request *mealplanning.CreateRecipeStepVesselRequest) (*mealplanning.CreateRecipeStepVesselResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetMermaidDiagramForRecipe(ctx context.Context, request *mealplanning.GetMermaidDiagramForRecipeRequest) (*mealplanning.GetMermaidDiagramForRecipeResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipe(ctx context.Context, request *mealplanning.GetRecipeRequest) (*mealplanning.GetRecipeResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipeMealPlanTasks(ctx context.Context, request *mealplanning.GetRecipeMealPlanTasksRequest) (*mealplanning.GetRecipeMealPlanTasksResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipePrepTask(ctx context.Context, request *mealplanning.GetRecipePrepTaskRequest) (*mealplanning.GetRecipePrepTaskResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipePrepTasks(ctx context.Context, request *mealplanning.GetRecipePrepTasksRequest) (*mealplanning.GetRecipePrepTasksResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipeRating(ctx context.Context, request *mealplanning.GetRecipeRatingRequest) (*mealplanning.GetRecipeRatingResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipeRatingsForRecipe(ctx context.Context, request *mealplanning.GetRecipeRatingsForRecipeRequest) (*mealplanning.GetRecipeRatingsForRecipeResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipeStep(ctx context.Context, request *mealplanning.GetRecipeStepRequest) (*mealplanning.GetRecipeStepResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipeStepCompletionCondition(ctx context.Context, request *mealplanning.GetRecipeStepCompletionConditionRequest) (*mealplanning.GetRecipeStepCompletionConditionResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipeStepCompletionConditions(ctx context.Context, request *mealplanning.GetRecipeStepCompletionConditionsRequest) (*mealplanning.GetRecipeStepCompletionConditionsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipeStepIngredient(ctx context.Context, request *mealplanning.GetRecipeStepIngredientRequest) (*mealplanning.GetRecipeStepIngredientResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipeStepIngredients(ctx context.Context, request *mealplanning.GetRecipeStepIngredientsRequest) (*mealplanning.GetRecipeStepIngredientsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipeStepInstrument(ctx context.Context, request *mealplanning.GetRecipeStepInstrumentRequest) (*mealplanning.GetRecipeStepInstrumentResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipeStepInstruments(ctx context.Context, request *mealplanning.GetRecipeStepInstrumentsRequest) (*mealplanning.GetRecipeStepInstrumentsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipeStepProduct(ctx context.Context, request *mealplanning.GetRecipeStepProductRequest) (*mealplanning.GetRecipeStepProductResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipeStepProducts(ctx context.Context, request *mealplanning.GetRecipeStepProductsRequest) (*mealplanning.GetRecipeStepProductsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipeStepVessel(ctx context.Context, request *mealplanning.GetRecipeStepVesselRequest) (*mealplanning.GetRecipeStepVesselResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipeStepVessels(ctx context.Context, request *mealplanning.GetRecipeStepVesselsRequest) (*mealplanning.GetRecipeStepVesselsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipeSteps(ctx context.Context, request *mealplanning.GetRecipeStepsRequest) (*mealplanning.GetRecipeStepsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetRecipes(ctx context.Context, request *mealplanning.GetRecipesRequest) (*mealplanning.GetRecipesResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) SearchForRecipes(ctx context.Context, request *mealplanning.SearchForRecipesRequest) (*mealplanning.SearchForRecipesResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateRecipe(ctx context.Context, request *mealplanning.UpdateRecipeRequest) (*mealplanning.UpdateRecipeResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateRecipePrepTask(ctx context.Context, request *mealplanning.UpdateRecipePrepTaskRequest) (*mealplanning.UpdateRecipePrepTaskResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateRecipeRating(ctx context.Context, request *mealplanning.UpdateRecipeRatingRequest) (*mealplanning.UpdateRecipeRatingResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateRecipeStep(ctx context.Context, request *mealplanning.UpdateRecipeStepRequest) (*mealplanning.UpdateRecipeStepResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateRecipeStepCompletionCondition(ctx context.Context, request *mealplanning.UpdateRecipeStepCompletionConditionRequest) (*mealplanning.UpdateRecipeStepCompletionConditionResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateRecipeStepIngredient(ctx context.Context, request *mealplanning.UpdateRecipeStepIngredientRequest) (*mealplanning.UpdateRecipeStepIngredientResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateRecipeStepInstrument(ctx context.Context, request *mealplanning.UpdateRecipeStepInstrumentRequest) (*mealplanning.UpdateRecipeStepInstrumentResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateRecipeStepProduct(ctx context.Context, request *mealplanning.UpdateRecipeStepProductRequest) (*mealplanning.UpdateRecipeStepProductResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateRecipeStepVessel(ctx context.Context, request *mealplanning.UpdateRecipeStepVesselRequest) (*mealplanning.UpdateRecipeStepVesselResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}
