package recipestepinstruments

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
	// CreateMiddlewareCtxKey is a string alias we can use for referring to recipe step instrument input data in contexts
	CreateMiddlewareCtxKey models.ContextKey = "recipe_step_instrument_create_input"
	// UpdateMiddlewareCtxKey is a string alias we can use for referring to recipe step instrument update data in contexts
	UpdateMiddlewareCtxKey models.ContextKey = "recipe_step_instrument_update_input"

	counterName        metrics.CounterName = "recipeStepInstruments"
	counterDescription                     = "the number of recipeStepInstruments managed by the recipeStepInstruments service"
	topicName          string              = "recipe_step_instruments"
	serviceName        string              = "recipe_step_instruments_service"
)

var (
	_ models.RecipeStepInstrumentDataServer = (*Service)(nil)
)

type (
	// Service handles to-do list recipe step instruments
	Service struct {
		logger                        logging.Logger
		recipeStepInstrumentCounter   metrics.UnitCounter
		recipeStepInstrumentDatabase  models.RecipeStepInstrumentDataManager
		userIDFetcher                 UserIDFetcher
		recipeStepInstrumentIDFetcher RecipeStepInstrumentIDFetcher
		encoderDecoder                encoding.EncoderDecoder
		reporter                      newsman.Reporter
	}

	// UserIDFetcher is a function that fetches user IDs
	UserIDFetcher func(*http.Request) uint64

	// RecipeStepInstrumentIDFetcher is a function that fetches recipe step instrument IDs
	RecipeStepInstrumentIDFetcher func(*http.Request) uint64
)

// ProvideRecipeStepInstrumentsService builds a new RecipeStepInstrumentsService
func ProvideRecipeStepInstrumentsService(
	ctx context.Context,
	logger logging.Logger,
	db models.RecipeStepInstrumentDataManager,
	userIDFetcher UserIDFetcher,
	recipeStepInstrumentIDFetcher RecipeStepInstrumentIDFetcher,
	encoder encoding.EncoderDecoder,
	recipeStepInstrumentCounterProvider metrics.UnitCounterProvider,
	reporter newsman.Reporter,
) (*Service, error) {
	recipeStepInstrumentCounter, err := recipeStepInstrumentCounterProvider(counterName, counterDescription)
	if err != nil {
		return nil, fmt.Errorf("error initializing counter: %w", err)
	}

	svc := &Service{
		logger:                        logger.WithName(serviceName),
		recipeStepInstrumentDatabase:  db,
		encoderDecoder:                encoder,
		recipeStepInstrumentCounter:   recipeStepInstrumentCounter,
		userIDFetcher:                 userIDFetcher,
		recipeStepInstrumentIDFetcher: recipeStepInstrumentIDFetcher,
		reporter:                      reporter,
	}

	recipeStepInstrumentCount, err := svc.recipeStepInstrumentDatabase.GetAllRecipeStepInstrumentsCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("setting current recipe step instrument count: %w", err)
	}
	svc.recipeStepInstrumentCounter.IncrementBy(ctx, recipeStepInstrumentCount)

	return svc, nil
}
