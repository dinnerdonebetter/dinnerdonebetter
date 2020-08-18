package recipestepproducts

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
	// createMiddlewareCtxKey is a string alias we can use for referring to recipe step product input data in contexts.
	createMiddlewareCtxKey models.ContextKey = "recipe_step_product_create_input"
	// updateMiddlewareCtxKey is a string alias we can use for referring to recipe step product update data in contexts.
	updateMiddlewareCtxKey models.ContextKey = "recipe_step_product_update_input"

	counterName        metrics.CounterName = "recipeStepProducts"
	counterDescription string              = "the number of recipeStepProducts managed by the recipeStepProducts service"
	topicName          string              = "recipe_step_products"
	serviceName        string              = "recipe_step_products_service"
)

var (
	_ models.RecipeStepProductDataServer = (*Service)(nil)
)

type (
	// Service handles to-do list recipe step products
	Service struct {
		logger                       logging.Logger
		recipeDataManager            models.RecipeDataManager
		recipeStepDataManager        models.RecipeStepDataManager
		recipeStepProductDataManager models.RecipeStepProductDataManager
		recipeIDFetcher              RecipeIDFetcher
		recipeStepIDFetcher          RecipeStepIDFetcher
		recipeStepProductIDFetcher   RecipeStepProductIDFetcher
		userIDFetcher                UserIDFetcher
		recipeStepProductCounter     metrics.UnitCounter
		encoderDecoder               encoding.EncoderDecoder
		reporter                     newsman.Reporter
	}

	// UserIDFetcher is a function that fetches user IDs.
	UserIDFetcher func(*http.Request) uint64

	// RecipeIDFetcher is a function that fetches recipe IDs.
	RecipeIDFetcher func(*http.Request) uint64

	// RecipeStepIDFetcher is a function that fetches recipe step IDs.
	RecipeStepIDFetcher func(*http.Request) uint64

	// RecipeStepProductIDFetcher is a function that fetches recipe step product IDs.
	RecipeStepProductIDFetcher func(*http.Request) uint64
)

// ProvideRecipeStepProductsService builds a new RecipeStepProductsService.
func ProvideRecipeStepProductsService(
	logger logging.Logger,
	recipeDataManager models.RecipeDataManager,
	recipeStepDataManager models.RecipeStepDataManager,
	recipeStepProductDataManager models.RecipeStepProductDataManager,
	recipeIDFetcher RecipeIDFetcher,
	recipeStepIDFetcher RecipeStepIDFetcher,
	recipeStepProductIDFetcher RecipeStepProductIDFetcher,
	userIDFetcher UserIDFetcher,
	encoder encoding.EncoderDecoder,
	recipeStepProductCounterProvider metrics.UnitCounterProvider,
	reporter newsman.Reporter,
) (*Service, error) {
	recipeStepProductCounter, err := recipeStepProductCounterProvider(counterName, counterDescription)
	if err != nil {
		return nil, fmt.Errorf("error initializing counter: %w", err)
	}

	svc := &Service{
		logger:                       logger.WithName(serviceName),
		recipeIDFetcher:              recipeIDFetcher,
		recipeStepIDFetcher:          recipeStepIDFetcher,
		recipeStepProductIDFetcher:   recipeStepProductIDFetcher,
		userIDFetcher:                userIDFetcher,
		recipeDataManager:            recipeDataManager,
		recipeStepDataManager:        recipeStepDataManager,
		recipeStepProductDataManager: recipeStepProductDataManager,
		encoderDecoder:               encoder,
		recipeStepProductCounter:     recipeStepProductCounter,
		reporter:                     reporter,
	}

	return svc, nil
}
