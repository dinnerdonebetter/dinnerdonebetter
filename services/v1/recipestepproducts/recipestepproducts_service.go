package recipestepproducts

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
	// CreateMiddlewareCtxKey is a string alias we can use for referring to recipe step product input data in contexts
	CreateMiddlewareCtxKey models.ContextKey = "recipe_step_product_create_input"
	// UpdateMiddlewareCtxKey is a string alias we can use for referring to recipe step product update data in contexts
	UpdateMiddlewareCtxKey models.ContextKey = "recipe_step_product_update_input"

	counterName        metrics.CounterName = "recipeStepProducts"
	counterDescription                     = "the number of recipeStepProducts managed by the recipeStepProducts service"
	topicName          string              = "recipe_step_products"
	serviceName        string              = "recipe_step_products_service"
)

var (
	_ models.RecipeStepProductDataServer = (*Service)(nil)
)

type (
	// Service handles to-do list recipe step products
	Service struct {
		logger                     logging.Logger
		recipeStepProductCounter   metrics.UnitCounter
		recipeStepProductDatabase  models.RecipeStepProductDataManager
		userIDFetcher              UserIDFetcher
		recipeStepProductIDFetcher RecipeStepProductIDFetcher
		encoderDecoder             encoding.EncoderDecoder
		reporter                   newsman.Reporter
	}

	// UserIDFetcher is a function that fetches user IDs
	UserIDFetcher func(*http.Request) uint64

	// RecipeStepProductIDFetcher is a function that fetches recipe step product IDs
	RecipeStepProductIDFetcher func(*http.Request) uint64
)

// ProvideRecipeStepProductsService builds a new RecipeStepProductsService
func ProvideRecipeStepProductsService(
	ctx context.Context,
	logger logging.Logger,
	db models.RecipeStepProductDataManager,
	userIDFetcher UserIDFetcher,
	recipeStepProductIDFetcher RecipeStepProductIDFetcher,
	encoder encoding.EncoderDecoder,
	recipeStepProductCounterProvider metrics.UnitCounterProvider,
	reporter newsman.Reporter,
) (*Service, error) {
	recipeStepProductCounter, err := recipeStepProductCounterProvider(counterName, counterDescription)
	if err != nil {
		return nil, fmt.Errorf("error initializing counter: %w", err)
	}

	svc := &Service{
		logger:                     logger.WithName(serviceName),
		recipeStepProductDatabase:  db,
		encoderDecoder:             encoder,
		recipeStepProductCounter:   recipeStepProductCounter,
		userIDFetcher:              userIDFetcher,
		recipeStepProductIDFetcher: recipeStepProductIDFetcher,
		reporter:                   reporter,
	}

	recipeStepProductCount, err := svc.recipeStepProductDatabase.GetAllRecipeStepProductsCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("setting current recipe step product count: %w", err)
	}
	svc.recipeStepProductCounter.IncrementBy(ctx, recipeStepProductCount)

	return svc, nil
}
