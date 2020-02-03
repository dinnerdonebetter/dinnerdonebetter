package recipes

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
	// CreateMiddlewareCtxKey is a string alias we can use for referring to recipe input data in contexts
	CreateMiddlewareCtxKey models.ContextKey = "recipe_create_input"
	// UpdateMiddlewareCtxKey is a string alias we can use for referring to recipe update data in contexts
	UpdateMiddlewareCtxKey models.ContextKey = "recipe_update_input"

	counterName        metrics.CounterName = "recipes"
	counterDescription                     = "the number of recipes managed by the recipes service"
	topicName          string              = "recipes"
	serviceName        string              = "recipes_service"
)

var (
	_ models.RecipeDataServer = (*Service)(nil)
)

type (
	// Service handles to-do list recipes
	Service struct {
		logger          logging.Logger
		recipeCounter   metrics.UnitCounter
		recipeDatabase  models.RecipeDataManager
		userIDFetcher   UserIDFetcher
		recipeIDFetcher RecipeIDFetcher
		encoderDecoder  encoding.EncoderDecoder
		reporter        newsman.Reporter
	}

	// UserIDFetcher is a function that fetches user IDs
	UserIDFetcher func(*http.Request) uint64

	// RecipeIDFetcher is a function that fetches recipe IDs
	RecipeIDFetcher func(*http.Request) uint64
)

// ProvideRecipesService builds a new RecipesService
func ProvideRecipesService(
	ctx context.Context,
	logger logging.Logger,
	db models.RecipeDataManager,
	userIDFetcher UserIDFetcher,
	recipeIDFetcher RecipeIDFetcher,
	encoder encoding.EncoderDecoder,
	recipeCounterProvider metrics.UnitCounterProvider,
	reporter newsman.Reporter,
) (*Service, error) {
	recipeCounter, err := recipeCounterProvider(counterName, counterDescription)
	if err != nil {
		return nil, fmt.Errorf("error initializing counter: %w", err)
	}

	svc := &Service{
		logger:          logger.WithName(serviceName),
		recipeDatabase:  db,
		encoderDecoder:  encoder,
		recipeCounter:   recipeCounter,
		userIDFetcher:   userIDFetcher,
		recipeIDFetcher: recipeIDFetcher,
		reporter:        reporter,
	}

	recipeCount, err := svc.recipeDatabase.GetAllRecipesCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("setting current recipe count: %w", err)
	}
	svc.recipeCounter.IncrementBy(ctx, recipeCount)

	return svc, nil
}
