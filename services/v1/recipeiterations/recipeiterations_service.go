package recipeiterations

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
	// CreateMiddlewareCtxKey is a string alias we can use for referring to recipe iteration input data in contexts.
	CreateMiddlewareCtxKey models.ContextKey = "recipe_iteration_create_input"
	// UpdateMiddlewareCtxKey is a string alias we can use for referring to recipe iteration update data in contexts.
	UpdateMiddlewareCtxKey models.ContextKey = "recipe_iteration_update_input"

	counterName        metrics.CounterName = "recipeIterations"
	counterDescription string              = "the number of recipeIterations managed by the recipeIterations service"
	topicName          string              = "recipe_iterations"
	serviceName        string              = "recipe_iterations_service"
)

var (
	_ models.RecipeIterationDataServer = (*Service)(nil)
)

type (
	// Service handles to-do list recipe iterations
	Service struct {
		logger                     logging.Logger
		recipeDataManager          models.RecipeDataManager
		recipeIterationDataManager models.RecipeIterationDataManager
		recipeIDFetcher            RecipeIDFetcher
		recipeIterationIDFetcher   RecipeIterationIDFetcher
		userIDFetcher              UserIDFetcher
		recipeIterationCounter     metrics.UnitCounter
		encoderDecoder             encoding.EncoderDecoder
		reporter                   newsman.Reporter
	}

	// UserIDFetcher is a function that fetches user IDs.
	UserIDFetcher func(*http.Request) uint64

	// RecipeIDFetcher is a function that fetches recipe IDs.
	RecipeIDFetcher func(*http.Request) uint64

	// RecipeIterationIDFetcher is a function that fetches recipe iteration IDs.
	RecipeIterationIDFetcher func(*http.Request) uint64
)

// ProvideRecipeIterationsService builds a new RecipeIterationsService.
func ProvideRecipeIterationsService(
	logger logging.Logger,
	recipeDataManager models.RecipeDataManager,
	recipeIterationDataManager models.RecipeIterationDataManager,
	recipeIDFetcher RecipeIDFetcher,
	recipeIterationIDFetcher RecipeIterationIDFetcher,
	userIDFetcher UserIDFetcher,
	encoder encoding.EncoderDecoder,
	recipeIterationCounterProvider metrics.UnitCounterProvider,
	reporter newsman.Reporter,
) (*Service, error) {
	recipeIterationCounter, err := recipeIterationCounterProvider(counterName, counterDescription)
	if err != nil {
		return nil, fmt.Errorf("error initializing counter: %w", err)
	}

	svc := &Service{
		logger:                     logger.WithName(serviceName),
		recipeIDFetcher:            recipeIDFetcher,
		recipeIterationIDFetcher:   recipeIterationIDFetcher,
		userIDFetcher:              userIDFetcher,
		recipeDataManager:          recipeDataManager,
		recipeIterationDataManager: recipeIterationDataManager,
		encoderDecoder:             encoder,
		recipeIterationCounter:     recipeIterationCounter,
		reporter:                   reporter,
	}

	return svc, nil
}
