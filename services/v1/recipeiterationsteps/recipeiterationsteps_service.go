package recipeiterationsteps

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
	// CreateMiddlewareCtxKey is a string alias we can use for referring to recipe iteration step input data in contexts.
	CreateMiddlewareCtxKey models.ContextKey = "recipe_iteration_step_create_input"
	// UpdateMiddlewareCtxKey is a string alias we can use for referring to recipe iteration step update data in contexts.
	UpdateMiddlewareCtxKey models.ContextKey = "recipe_iteration_step_update_input"

	counterName        metrics.CounterName = "recipeIterationSteps"
	counterDescription string              = "the number of recipeIterationSteps managed by the recipeIterationSteps service"
	topicName          string              = "recipe_iteration_steps"
	serviceName        string              = "recipe_iteration_steps_service"
)

var (
	_ models.RecipeIterationStepDataServer = (*Service)(nil)
)

type (
	// Service handles to-do list recipe iteration steps
	Service struct {
		logger                         logging.Logger
		recipeDataManager              models.RecipeDataManager
		recipeIterationStepDataManager models.RecipeIterationStepDataManager
		recipeIDFetcher                RecipeIDFetcher
		recipeIterationStepIDFetcher   RecipeIterationStepIDFetcher
		userIDFetcher                  UserIDFetcher
		recipeIterationStepCounter     metrics.UnitCounter
		encoderDecoder                 encoding.EncoderDecoder
		reporter                       newsman.Reporter
	}

	// UserIDFetcher is a function that fetches user IDs.
	UserIDFetcher func(*http.Request) uint64

	// RecipeIDFetcher is a function that fetches recipe IDs.
	RecipeIDFetcher func(*http.Request) uint64

	// RecipeIterationStepIDFetcher is a function that fetches recipe iteration step IDs.
	RecipeIterationStepIDFetcher func(*http.Request) uint64
)

// ProvideRecipeIterationStepsService builds a new RecipeIterationStepsService.
func ProvideRecipeIterationStepsService(
	logger logging.Logger,
	recipeDataManager models.RecipeDataManager,
	recipeIterationStepDataManager models.RecipeIterationStepDataManager,
	recipeIDFetcher RecipeIDFetcher,
	recipeIterationStepIDFetcher RecipeIterationStepIDFetcher,
	userIDFetcher UserIDFetcher,
	encoder encoding.EncoderDecoder,
	recipeIterationStepCounterProvider metrics.UnitCounterProvider,
	reporter newsman.Reporter,
) (*Service, error) {
	recipeIterationStepCounter, err := recipeIterationStepCounterProvider(counterName, counterDescription)
	if err != nil {
		return nil, fmt.Errorf("error initializing counter: %w", err)
	}

	svc := &Service{
		logger:                         logger.WithName(serviceName),
		recipeIDFetcher:                recipeIDFetcher,
		recipeIterationStepIDFetcher:   recipeIterationStepIDFetcher,
		userIDFetcher:                  userIDFetcher,
		recipeDataManager:              recipeDataManager,
		recipeIterationStepDataManager: recipeIterationStepDataManager,
		encoderDecoder:                 encoder,
		recipeIterationStepCounter:     recipeIterationStepCounter,
		reporter:                       reporter,
	}

	return svc, nil
}
