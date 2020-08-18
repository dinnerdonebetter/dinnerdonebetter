package recipestepinstruments

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
	// createMiddlewareCtxKey is a string alias we can use for referring to recipe step instrument input data in contexts.
	createMiddlewareCtxKey models.ContextKey = "recipe_step_instrument_create_input"
	// updateMiddlewareCtxKey is a string alias we can use for referring to recipe step instrument update data in contexts.
	updateMiddlewareCtxKey models.ContextKey = "recipe_step_instrument_update_input"

	counterName        metrics.CounterName = "recipeStepInstruments"
	counterDescription string              = "the number of recipeStepInstruments managed by the recipeStepInstruments service"
	topicName          string              = "recipe_step_instruments"
	serviceName        string              = "recipe_step_instruments_service"
)

var (
	_ models.RecipeStepInstrumentDataServer = (*Service)(nil)
)

type (
	// Service handles to-do list recipe step instruments
	Service struct {
		logger                          logging.Logger
		recipeDataManager               models.RecipeDataManager
		recipeStepDataManager           models.RecipeStepDataManager
		recipeStepInstrumentDataManager models.RecipeStepInstrumentDataManager
		recipeIDFetcher                 RecipeIDFetcher
		recipeStepIDFetcher             RecipeStepIDFetcher
		recipeStepInstrumentIDFetcher   RecipeStepInstrumentIDFetcher
		userIDFetcher                   UserIDFetcher
		recipeStepInstrumentCounter     metrics.UnitCounter
		encoderDecoder                  encoding.EncoderDecoder
		reporter                        newsman.Reporter
	}

	// UserIDFetcher is a function that fetches user IDs.
	UserIDFetcher func(*http.Request) uint64

	// RecipeIDFetcher is a function that fetches recipe IDs.
	RecipeIDFetcher func(*http.Request) uint64

	// RecipeStepIDFetcher is a function that fetches recipe step IDs.
	RecipeStepIDFetcher func(*http.Request) uint64

	// RecipeStepInstrumentIDFetcher is a function that fetches recipe step instrument IDs.
	RecipeStepInstrumentIDFetcher func(*http.Request) uint64
)

// ProvideRecipeStepInstrumentsService builds a new RecipeStepInstrumentsService.
func ProvideRecipeStepInstrumentsService(
	logger logging.Logger,
	recipeDataManager models.RecipeDataManager,
	recipeStepDataManager models.RecipeStepDataManager,
	recipeStepInstrumentDataManager models.RecipeStepInstrumentDataManager,
	recipeIDFetcher RecipeIDFetcher,
	recipeStepIDFetcher RecipeStepIDFetcher,
	recipeStepInstrumentIDFetcher RecipeStepInstrumentIDFetcher,
	userIDFetcher UserIDFetcher,
	encoder encoding.EncoderDecoder,
	recipeStepInstrumentCounterProvider metrics.UnitCounterProvider,
	reporter newsman.Reporter,
) (*Service, error) {
	recipeStepInstrumentCounter, err := recipeStepInstrumentCounterProvider(counterName, counterDescription)
	if err != nil {
		return nil, fmt.Errorf("error initializing counter: %w", err)
	}

	svc := &Service{
		logger:                          logger.WithName(serviceName),
		recipeIDFetcher:                 recipeIDFetcher,
		recipeStepIDFetcher:             recipeStepIDFetcher,
		recipeStepInstrumentIDFetcher:   recipeStepInstrumentIDFetcher,
		userIDFetcher:                   userIDFetcher,
		recipeDataManager:               recipeDataManager,
		recipeStepDataManager:           recipeStepDataManager,
		recipeStepInstrumentDataManager: recipeStepInstrumentDataManager,
		encoderDecoder:                  encoder,
		recipeStepInstrumentCounter:     recipeStepInstrumentCounter,
		reporter:                        reporter,
	}

	return svc, nil
}
