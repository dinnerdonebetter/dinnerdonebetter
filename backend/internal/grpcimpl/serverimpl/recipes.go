package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func (s *Server) CloneRecipe(ctx context.Context, request *messages.CloneRecipeRequest) (*messages.Recipe, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveRecipe(ctx context.Context, request *messages.ArchiveRecipeRequest) (*messages.Recipe, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveRecipePrepTask(ctx context.Context, request *messages.ArchiveRecipePrepTaskRequest) (*messages.RecipePrepTask, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveRecipeRating(ctx context.Context, request *messages.ArchiveRecipeRatingRequest) (*messages.RecipeRating, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveRecipeStep(ctx context.Context, request *messages.ArchiveRecipeStepRequest) (*messages.RecipeStep, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveRecipeStepCompletionCondition(ctx context.Context, request *messages.ArchiveRecipeStepCompletionConditionRequest) (*messages.RecipeStepCompletionCondition, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveRecipeStepIngredient(ctx context.Context, request *messages.ArchiveRecipeStepIngredientRequest) (*messages.RecipeStepIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveRecipeStepInstrument(ctx context.Context, request *messages.ArchiveRecipeStepInstrumentRequest) (*messages.RecipeStepInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveRecipeStepProduct(ctx context.Context, request *messages.ArchiveRecipeStepProductRequest) (*messages.RecipeStepProduct, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveRecipeStepVessel(ctx context.Context, request *messages.ArchiveRecipeStepVesselRequest) (*messages.RecipeStepVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateRecipe(ctx context.Context, input *messages.RecipeCreationRequestInput) (*messages.Recipe, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateRecipePrepTask(ctx context.Context, request *messages.CreateRecipePrepTaskRequest) (*messages.RecipePrepTask, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateRecipeRating(ctx context.Context, request *messages.CreateRecipeRatingRequest) (*messages.RecipeRating, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateRecipeStep(ctx context.Context, request *messages.CreateRecipeStepRequest) (*messages.RecipeStep, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateRecipeStepCompletionCondition(ctx context.Context, request *messages.CreateRecipeStepCompletionConditionRequest) (*messages.RecipeStepCompletionCondition, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateRecipeStepIngredient(ctx context.Context, request *messages.CreateRecipeStepIngredientRequest) (*messages.RecipeStepIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateRecipeStepInstrument(ctx context.Context, request *messages.CreateRecipeStepInstrumentRequest) (*messages.RecipeStepInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateRecipeStepProduct(ctx context.Context, request *messages.CreateRecipeStepProductRequest) (*messages.RecipeStepProduct, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateRecipeStepVessel(ctx context.Context, request *messages.CreateRecipeStepVesselRequest) (*messages.RecipeStepVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetMermaidDiagramForRecipe(ctx context.Context, request *messages.GetMermaidDiagramForRecipeRequest) (*messages.GetMermaidDiagramForRecipeResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipe(ctx context.Context, request *messages.GetRecipeRequest) (*messages.Recipe, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipeMealPlanTasks(ctx context.Context, request *messages.GetRecipeMealPlanTasksRequest) (*messages.RecipePrepTaskStep, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipePrepTask(ctx context.Context, request *messages.GetRecipePrepTaskRequest) (*messages.RecipePrepTask, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipePrepTasks(ctx context.Context, request *messages.GetRecipePrepTasksRequest) (*messages.RecipePrepTask, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipeRating(ctx context.Context, request *messages.GetRecipeRatingRequest) (*messages.RecipeRating, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipeRatingsForRecipe(ctx context.Context, request *messages.GetRecipeRatingsForRecipeRequest) (*messages.RecipeRating, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipeStep(ctx context.Context, request *messages.GetRecipeStepRequest) (*messages.RecipeStep, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipeStepCompletionCondition(ctx context.Context, request *messages.GetRecipeStepCompletionConditionRequest) (*messages.RecipeStepCompletionCondition, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipeStepCompletionConditions(ctx context.Context, request *messages.GetRecipeStepCompletionConditionsRequest) (*messages.RecipeStepCompletionCondition, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipeStepIngredient(ctx context.Context, request *messages.GetRecipeStepIngredientRequest) (*messages.RecipeStepIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipeStepIngredients(ctx context.Context, request *messages.GetRecipeStepIngredientsRequest) (*messages.RecipeStepIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipeStepInstrument(ctx context.Context, request *messages.GetRecipeStepInstrumentRequest) (*messages.RecipeStepInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipeStepInstruments(ctx context.Context, request *messages.GetRecipeStepInstrumentsRequest) (*messages.RecipeStepInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipeStepProduct(ctx context.Context, request *messages.GetRecipeStepProductRequest) (*messages.RecipeStepProduct, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipeStepProducts(ctx context.Context, request *messages.GetRecipeStepProductsRequest) (*messages.RecipeStepProduct, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipeStepVessel(ctx context.Context, request *messages.GetRecipeStepVesselRequest) (*messages.RecipeStepVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipeStepVessels(ctx context.Context, request *messages.GetRecipeStepVesselsRequest) (*messages.RecipeStepVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipeSteps(ctx context.Context, request *messages.GetRecipeStepsRequest) (*messages.RecipeStep, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipes(ctx context.Context, request *messages.GetRecipesRequest) (*messages.Recipe, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) SearchForRecipes(ctx context.Context, request *messages.SearchForRecipesRequest) (*messages.Recipe, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateRecipe(ctx context.Context, request *messages.UpdateRecipeRequest) (*messages.Recipe, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateRecipePrepTask(ctx context.Context, request *messages.UpdateRecipePrepTaskRequest) (*messages.RecipePrepTask, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateRecipeRating(ctx context.Context, request *messages.UpdateRecipeRatingRequest) (*messages.RecipeRating, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateRecipeStep(ctx context.Context, request *messages.UpdateRecipeStepRequest) (*messages.RecipeStep, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateRecipeStepCompletionCondition(ctx context.Context, request *messages.UpdateRecipeStepCompletionConditionRequest) (*messages.RecipeStepCompletionCondition, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateRecipeStepIngredient(ctx context.Context, request *messages.UpdateRecipeStepIngredientRequest) (*messages.RecipeStepIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateRecipeStepInstrument(ctx context.Context, request *messages.UpdateRecipeStepInstrumentRequest) (*messages.RecipeStepInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateRecipeStepProduct(ctx context.Context, request *messages.UpdateRecipeStepProductRequest) (*messages.RecipeStepProduct, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateRecipeStepVessel(ctx context.Context, request *messages.UpdateRecipeStepVesselRequest) (*messages.RecipeStepVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}
