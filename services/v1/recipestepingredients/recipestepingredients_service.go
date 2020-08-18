package recipestepingredients

import (
	"fmt"
	"net/http"

	"gitlab.com/prixfixe/prixfixe/internal/v1/encoding"
	"gitlab.com/prixfixe/prixfixe/internal/v1/metrics"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	"gitlab.com/verygoodsoftwarenotvirus/newsman"
)

const (
	// createMiddlewareCtxKey is a string alias we can use for referring to recipe step ingredient input data in contexts.
	createMiddlewareCtxKey models.ContextKey = "recipe_step_ingredient_create_input"
	// updateMiddlewareCtxKey is a string alias we can use for referring to recipe step ingredient update data in contexts.
	updateMiddlewareCtxKey models.ContextKey = "recipe_step_ingredient_update_input"

	counterName        metrics.CounterName = "recipeStepIngredients"
	counterDescription string              = "the number of recipeStepIngredients managed by the recipeStepIngredients service"
	topicName          string              = "recipe_step_ingredients"
	serviceName        string              = "recipe_step_ingredients_service"
)

var (
	_ models.RecipeStepIngredientDataServer = (*Service)(nil)
)

type (
	// Service handles to-do list recipe step ingredients
	Service struct {
		logger                          logging.Logger
		recipeDataManager               models.RecipeDataManager
		recipeStepDataManager           models.RecipeStepDataManager
		recipeStepIngredientDataManager models.RecipeStepIngredientDataManager
		recipeIDFetcher                 RecipeIDFetcher
		recipeStepIDFetcher             RecipeStepIDFetcher
		recipeStepIngredientIDFetcher   RecipeStepIngredientIDFetcher
		userIDFetcher                   UserIDFetcher
		recipeStepIngredientCounter     metrics.UnitCounter
		encoderDecoder                  encoding.EncoderDecoder
		reporter                        newsman.Reporter
	}

	// UserIDFetcher is a function that fetches user IDs.
	UserIDFetcher func(*http.Request) uint64

	// RecipeIDFetcher is a function that fetches recipe IDs.
	RecipeIDFetcher func(*http.Request) uint64

	// RecipeStepIDFetcher is a function that fetches recipe step IDs.
	RecipeStepIDFetcher func(*http.Request) uint64

	// RecipeStepIngredientIDFetcher is a function that fetches recipe step ingredient IDs.
	RecipeStepIngredientIDFetcher func(*http.Request) uint64
)

// ProvideRecipeStepIngredientsService builds a new RecipeStepIngredientsService.
func ProvideRecipeStepIngredientsService(
	logger logging.Logger,
	recipeDataManager models.RecipeDataManager,
	recipeStepDataManager models.RecipeStepDataManager,
	recipeStepIngredientDataManager models.RecipeStepIngredientDataManager,
	recipeIDFetcher RecipeIDFetcher,
	recipeStepIDFetcher RecipeStepIDFetcher,
	recipeStepIngredientIDFetcher RecipeStepIngredientIDFetcher,
	userIDFetcher UserIDFetcher,
	encoder encoding.EncoderDecoder,
	recipeStepIngredientCounterProvider metrics.UnitCounterProvider,
	reporter newsman.Reporter,
) (*Service, error) {
	recipeStepIngredientCounter, err := recipeStepIngredientCounterProvider(counterName, counterDescription)
	if err != nil {
		return nil, fmt.Errorf("error initializing counter: %w", err)
	}

	svc := &Service{
		logger:                          logger.WithName(serviceName),
		recipeIDFetcher:                 recipeIDFetcher,
		recipeStepIDFetcher:             recipeStepIDFetcher,
		recipeStepIngredientIDFetcher:   recipeStepIngredientIDFetcher,
		userIDFetcher:                   userIDFetcher,
		recipeDataManager:               recipeDataManager,
		recipeStepDataManager:           recipeStepDataManager,
		recipeStepIngredientDataManager: recipeStepIngredientDataManager,
		encoderDecoder:                  encoder,
		recipeStepIngredientCounter:     recipeStepIngredientCounter,
		reporter:                        reporter,
	}

	return svc, nil
}
