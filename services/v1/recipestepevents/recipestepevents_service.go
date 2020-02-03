package recipestepevents

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
	// CreateMiddlewareCtxKey is a string alias we can use for referring to recipe step event input data in contexts
	CreateMiddlewareCtxKey models.ContextKey = "recipe_step_event_create_input"
	// UpdateMiddlewareCtxKey is a string alias we can use for referring to recipe step event update data in contexts
	UpdateMiddlewareCtxKey models.ContextKey = "recipe_step_event_update_input"

	counterName        metrics.CounterName = "recipeStepEvents"
	counterDescription                     = "the number of recipeStepEvents managed by the recipeStepEvents service"
	topicName          string              = "recipe_step_events"
	serviceName        string              = "recipe_step_events_service"
)

var (
	_ models.RecipeStepEventDataServer = (*Service)(nil)
)

type (
	// Service handles to-do list recipe step events
	Service struct {
		logger                   logging.Logger
		recipeStepEventCounter   metrics.UnitCounter
		recipeStepEventDatabase  models.RecipeStepEventDataManager
		userIDFetcher            UserIDFetcher
		recipeStepEventIDFetcher RecipeStepEventIDFetcher
		encoderDecoder           encoding.EncoderDecoder
		reporter                 newsman.Reporter
	}

	// UserIDFetcher is a function that fetches user IDs
	UserIDFetcher func(*http.Request) uint64

	// RecipeStepEventIDFetcher is a function that fetches recipe step event IDs
	RecipeStepEventIDFetcher func(*http.Request) uint64
)

// ProvideRecipeStepEventsService builds a new RecipeStepEventsService
func ProvideRecipeStepEventsService(
	ctx context.Context,
	logger logging.Logger,
	db models.RecipeStepEventDataManager,
	userIDFetcher UserIDFetcher,
	recipeStepEventIDFetcher RecipeStepEventIDFetcher,
	encoder encoding.EncoderDecoder,
	recipeStepEventCounterProvider metrics.UnitCounterProvider,
	reporter newsman.Reporter,
) (*Service, error) {
	recipeStepEventCounter, err := recipeStepEventCounterProvider(counterName, counterDescription)
	if err != nil {
		return nil, fmt.Errorf("error initializing counter: %w", err)
	}

	svc := &Service{
		logger:                   logger.WithName(serviceName),
		recipeStepEventDatabase:  db,
		encoderDecoder:           encoder,
		recipeStepEventCounter:   recipeStepEventCounter,
		userIDFetcher:            userIDFetcher,
		recipeStepEventIDFetcher: recipeStepEventIDFetcher,
		reporter:                 reporter,
	}

	recipeStepEventCount, err := svc.recipeStepEventDatabase.GetAllRecipeStepEventsCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("setting current recipe step event count: %w", err)
	}
	svc.recipeStepEventCounter.IncrementBy(ctx, recipeStepEventCount)

	return svc, nil
}
