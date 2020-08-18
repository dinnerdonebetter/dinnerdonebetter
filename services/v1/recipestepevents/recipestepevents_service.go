package recipestepevents

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
	// createMiddlewareCtxKey is a string alias we can use for referring to recipe step event input data in contexts.
	createMiddlewareCtxKey models.ContextKey = "recipe_step_event_create_input"
	// updateMiddlewareCtxKey is a string alias we can use for referring to recipe step event update data in contexts.
	updateMiddlewareCtxKey models.ContextKey = "recipe_step_event_update_input"

	counterName        metrics.CounterName = "recipeStepEvents"
	counterDescription string              = "the number of recipeStepEvents managed by the recipeStepEvents service"
	topicName          string              = "recipe_step_events"
	serviceName        string              = "recipe_step_events_service"
)

var (
	_ models.RecipeStepEventDataServer = (*Service)(nil)
)

type (
	// Service handles to-do list recipe step events
	Service struct {
		logger                     logging.Logger
		recipeDataManager          models.RecipeDataManager
		recipeStepDataManager      models.RecipeStepDataManager
		recipeStepEventDataManager models.RecipeStepEventDataManager
		recipeIDFetcher            RecipeIDFetcher
		recipeStepIDFetcher        RecipeStepIDFetcher
		recipeStepEventIDFetcher   RecipeStepEventIDFetcher
		userIDFetcher              UserIDFetcher
		recipeStepEventCounter     metrics.UnitCounter
		encoderDecoder             encoding.EncoderDecoder
		reporter                   newsman.Reporter
	}

	// UserIDFetcher is a function that fetches user IDs.
	UserIDFetcher func(*http.Request) uint64

	// RecipeIDFetcher is a function that fetches recipe IDs.
	RecipeIDFetcher func(*http.Request) uint64

	// RecipeStepIDFetcher is a function that fetches recipe step IDs.
	RecipeStepIDFetcher func(*http.Request) uint64

	// RecipeStepEventIDFetcher is a function that fetches recipe step event IDs.
	RecipeStepEventIDFetcher func(*http.Request) uint64
)

// ProvideRecipeStepEventsService builds a new RecipeStepEventsService.
func ProvideRecipeStepEventsService(
	logger logging.Logger,
	recipeDataManager models.RecipeDataManager,
	recipeStepDataManager models.RecipeStepDataManager,
	recipeStepEventDataManager models.RecipeStepEventDataManager,
	recipeIDFetcher RecipeIDFetcher,
	recipeStepIDFetcher RecipeStepIDFetcher,
	recipeStepEventIDFetcher RecipeStepEventIDFetcher,
	userIDFetcher UserIDFetcher,
	encoder encoding.EncoderDecoder,
	recipeStepEventCounterProvider metrics.UnitCounterProvider,
	reporter newsman.Reporter,
) (*Service, error) {
	recipeStepEventCounter, err := recipeStepEventCounterProvider(counterName, counterDescription)
	if err != nil {
		return nil, fmt.Errorf("error initializing counter: %w", err)
	}

	svc := &Service{
		logger:                     logger.WithName(serviceName),
		recipeIDFetcher:            recipeIDFetcher,
		recipeStepIDFetcher:        recipeStepIDFetcher,
		recipeStepEventIDFetcher:   recipeStepEventIDFetcher,
		userIDFetcher:              userIDFetcher,
		recipeDataManager:          recipeDataManager,
		recipeStepDataManager:      recipeStepDataManager,
		recipeStepEventDataManager: recipeStepEventDataManager,
		encoderDecoder:             encoder,
		recipeStepEventCounter:     recipeStepEventCounter,
		reporter:                   reporter,
	}

	return svc, nil
}
