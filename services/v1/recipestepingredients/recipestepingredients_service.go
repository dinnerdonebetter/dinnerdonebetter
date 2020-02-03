package recipestepingredients

import (
	"context"
	"fmt"
	"net/http"

	"gitlab.com/prixfixe/prixfixe/internal/v1/encoding"
	"gitlab.com/prixfixe/prixfixe/internal/v1/metrics"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	"gitlab.com/verygoodsoftwarenotvirus/newsman"
)

const (
	// CreateMiddlewareCtxKey is a string alias we can use for referring to recipe step ingredient input data in contexts
	CreateMiddlewareCtxKey models.ContextKey = "recipe_step_ingredient_create_input"
	// UpdateMiddlewareCtxKey is a string alias we can use for referring to recipe step ingredient update data in contexts
	UpdateMiddlewareCtxKey models.ContextKey = "recipe_step_ingredient_update_input"

	counterName        metrics.CounterName = "recipeStepIngredients"
	counterDescription                     = "the number of recipeStepIngredients managed by the recipeStepIngredients service"
	topicName          string              = "recipe_step_ingredients"
	serviceName        string              = "recipe_step_ingredients_service"
)

var (
	_ models.RecipeStepIngredientDataServer = (*Service)(nil)
)

type (
	// Service handles to-do list recipe step ingredients
	Service struct {
		logger                        logging.Logger
		recipeStepIngredientCounter   metrics.UnitCounter
		recipeStepIngredientDatabase  models.RecipeStepIngredientDataManager
		userIDFetcher                 UserIDFetcher
		recipeStepIngredientIDFetcher RecipeStepIngredientIDFetcher
		encoderDecoder                encoding.EncoderDecoder
		reporter                      newsman.Reporter
	}

	// UserIDFetcher is a function that fetches user IDs
	UserIDFetcher func(*http.Request) uint64

	// RecipeStepIngredientIDFetcher is a function that fetches recipe step ingredient IDs
	RecipeStepIngredientIDFetcher func(*http.Request) uint64
)

// ProvideRecipeStepIngredientsService builds a new RecipeStepIngredientsService
func ProvideRecipeStepIngredientsService(
	ctx context.Context,
	logger logging.Logger,
	db models.RecipeStepIngredientDataManager,
	userIDFetcher UserIDFetcher,
	recipeStepIngredientIDFetcher RecipeStepIngredientIDFetcher,
	encoder encoding.EncoderDecoder,
	recipeStepIngredientCounterProvider metrics.UnitCounterProvider,
	reporter newsman.Reporter,
) (*Service, error) {
	recipeStepIngredientCounter, err := recipeStepIngredientCounterProvider(counterName, counterDescription)
	if err != nil {
		return nil, fmt.Errorf("error initializing counter: %w", err)
	}

	svc := &Service{
		logger:                        logger.WithName(serviceName),
		recipeStepIngredientDatabase:  db,
		encoderDecoder:                encoder,
		recipeStepIngredientCounter:   recipeStepIngredientCounter,
		userIDFetcher:                 userIDFetcher,
		recipeStepIngredientIDFetcher: recipeStepIngredientIDFetcher,
		reporter:                      reporter,
	}

	recipeStepIngredientCount, err := svc.recipeStepIngredientDatabase.GetAllRecipeStepIngredientsCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("setting current recipe step ingredient count: %w", err)
	}
	svc.recipeStepIngredientCounter.IncrementBy(ctx, recipeStepIngredientCount)

	return svc, nil
}
