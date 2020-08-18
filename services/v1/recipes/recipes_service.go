package recipes

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
	// createMiddlewareCtxKey is a string alias we can use for referring to recipe input data in contexts.
	createMiddlewareCtxKey models.ContextKey = "recipe_create_input"
	// updateMiddlewareCtxKey is a string alias we can use for referring to recipe update data in contexts.
	updateMiddlewareCtxKey models.ContextKey = "recipe_update_input"

	counterName        metrics.CounterName = "recipes"
	counterDescription string              = "the number of recipes managed by the recipes service"
	topicName          string              = "recipes"
	serviceName        string              = "recipes_service"
)

var (
	_ models.RecipeDataServer = (*Service)(nil)
)

type (
	// Service handles to-do list recipes
	Service struct {
		logger            logging.Logger
		recipeDataManager models.RecipeDataManager
		recipeIDFetcher   RecipeIDFetcher
		userIDFetcher     UserIDFetcher
		recipeCounter     metrics.UnitCounter
		encoderDecoder    encoding.EncoderDecoder
		reporter          newsman.Reporter
	}

	// UserIDFetcher is a function that fetches user IDs.
	UserIDFetcher func(*http.Request) uint64

	// RecipeIDFetcher is a function that fetches recipe IDs.
	RecipeIDFetcher func(*http.Request) uint64
)

// ProvideRecipesService builds a new RecipesService.
func ProvideRecipesService(
	logger logging.Logger,
	recipeDataManager models.RecipeDataManager,
	recipeIDFetcher RecipeIDFetcher,
	userIDFetcher UserIDFetcher,
	encoder encoding.EncoderDecoder,
	recipeCounterProvider metrics.UnitCounterProvider,
	reporter newsman.Reporter,
) (*Service, error) {
	recipeCounter, err := recipeCounterProvider(counterName, counterDescription)
	if err != nil {
		return nil, fmt.Errorf("error initializing counter: %w", err)
	}

	svc := &Service{
		logger:            logger.WithName(serviceName),
		recipeIDFetcher:   recipeIDFetcher,
		userIDFetcher:     userIDFetcher,
		recipeDataManager: recipeDataManager,
		encoderDecoder:    encoder,
		recipeCounter:     recipeCounter,
		reporter:          reporter,
	}

	return svc, nil
}
