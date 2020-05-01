package recipesteps

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
	// CreateMiddlewareCtxKey is a string alias we can use for referring to recipe step input data in contexts.
	CreateMiddlewareCtxKey models.ContextKey = "recipe_step_create_input"
	// UpdateMiddlewareCtxKey is a string alias we can use for referring to recipe step update data in contexts.
	UpdateMiddlewareCtxKey models.ContextKey = "recipe_step_update_input"

	counterName        metrics.CounterName = "recipeSteps"
	counterDescription string              = "the number of recipeSteps managed by the recipeSteps service"
	topicName          string              = "recipe_steps"
	serviceName        string              = "recipe_steps_service"
)

var (
	_ models.RecipeStepDataServer = (*Service)(nil)
)

type (
	// Service handles to-do list recipe steps
	Service struct {
		logger                logging.Logger
		recipeDataManager     models.RecipeDataManager
		recipeStepDataManager models.RecipeStepDataManager
		recipeIDFetcher       RecipeIDFetcher
		recipeStepIDFetcher   RecipeStepIDFetcher
		userIDFetcher         UserIDFetcher
		recipeStepCounter     metrics.UnitCounter
		encoderDecoder        encoding.EncoderDecoder
		reporter              newsman.Reporter
	}

	// UserIDFetcher is a function that fetches user IDs.
	UserIDFetcher func(*http.Request) uint64

	// RecipeIDFetcher is a function that fetches recipe IDs.
	RecipeIDFetcher func(*http.Request) uint64

	// RecipeStepIDFetcher is a function that fetches recipe step IDs.
	RecipeStepIDFetcher func(*http.Request) uint64
)

// ProvideRecipeStepsService builds a new RecipeStepsService.
func ProvideRecipeStepsService(
	logger logging.Logger,
	recipeDataManager models.RecipeDataManager,
	recipeStepDataManager models.RecipeStepDataManager,
	recipeIDFetcher RecipeIDFetcher,
	recipeStepIDFetcher RecipeStepIDFetcher,
	userIDFetcher UserIDFetcher,
	encoder encoding.EncoderDecoder,
	recipeStepCounterProvider metrics.UnitCounterProvider,
	reporter newsman.Reporter,
) (*Service, error) {
	recipeStepCounter, err := recipeStepCounterProvider(counterName, counterDescription)
	if err != nil {
		return nil, fmt.Errorf("error initializing counter: %w", err)
	}

	svc := &Service{
		logger:                logger.WithName(serviceName),
		recipeIDFetcher:       recipeIDFetcher,
		recipeStepIDFetcher:   recipeStepIDFetcher,
		userIDFetcher:         userIDFetcher,
		recipeDataManager:     recipeDataManager,
		recipeStepDataManager: recipeStepDataManager,
		encoderDecoder:        encoder,
		recipeStepCounter:     recipeStepCounter,
		reporter:              reporter,
	}

	return svc, nil
}
