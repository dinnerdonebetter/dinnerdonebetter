package recipesteppreparations

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
	// CreateMiddlewareCtxKey is a string alias we can use for referring to recipe step preparation input data in contexts.
	CreateMiddlewareCtxKey models.ContextKey = "recipe_step_preparation_create_input"
	// UpdateMiddlewareCtxKey is a string alias we can use for referring to recipe step preparation update data in contexts.
	UpdateMiddlewareCtxKey models.ContextKey = "recipe_step_preparation_update_input"

	counterName        metrics.CounterName = "recipeStepPreparations"
	counterDescription string              = "the number of recipeStepPreparations managed by the recipeStepPreparations service"
	topicName          string              = "recipe_step_preparations"
	serviceName        string              = "recipe_step_preparations_service"
)

var (
	_ models.RecipeStepPreparationDataServer = (*Service)(nil)
)

type (
	// Service handles to-do list recipe step preparations
	Service struct {
		logger                           logging.Logger
		recipeDataManager                models.RecipeDataManager
		recipeStepDataManager            models.RecipeStepDataManager
		recipeStepPreparationDataManager models.RecipeStepPreparationDataManager
		recipeIDFetcher                  RecipeIDFetcher
		recipeStepIDFetcher              RecipeStepIDFetcher
		recipeStepPreparationIDFetcher   RecipeStepPreparationIDFetcher
		userIDFetcher                    UserIDFetcher
		recipeStepPreparationCounter     metrics.UnitCounter
		encoderDecoder                   encoding.EncoderDecoder
		reporter                         newsman.Reporter
	}

	// UserIDFetcher is a function that fetches user IDs.
	UserIDFetcher func(*http.Request) uint64

	// RecipeIDFetcher is a function that fetches recipe IDs.
	RecipeIDFetcher func(*http.Request) uint64

	// RecipeStepIDFetcher is a function that fetches recipe step IDs.
	RecipeStepIDFetcher func(*http.Request) uint64

	// RecipeStepPreparationIDFetcher is a function that fetches recipe step preparation IDs.
	RecipeStepPreparationIDFetcher func(*http.Request) uint64
)

// ProvideRecipeStepPreparationsService builds a new RecipeStepPreparationsService.
func ProvideRecipeStepPreparationsService(
	logger logging.Logger,
	recipeDataManager models.RecipeDataManager,
	recipeStepDataManager models.RecipeStepDataManager,
	recipeStepPreparationDataManager models.RecipeStepPreparationDataManager,
	recipeIDFetcher RecipeIDFetcher,
	recipeStepIDFetcher RecipeStepIDFetcher,
	recipeStepPreparationIDFetcher RecipeStepPreparationIDFetcher,
	userIDFetcher UserIDFetcher,
	encoder encoding.EncoderDecoder,
	recipeStepPreparationCounterProvider metrics.UnitCounterProvider,
	reporter newsman.Reporter,
) (*Service, error) {
	recipeStepPreparationCounter, err := recipeStepPreparationCounterProvider(counterName, counterDescription)
	if err != nil {
		return nil, fmt.Errorf("error initializing counter: %w", err)
	}

	svc := &Service{
		logger:                           logger.WithName(serviceName),
		recipeIDFetcher:                  recipeIDFetcher,
		recipeStepIDFetcher:              recipeStepIDFetcher,
		recipeStepPreparationIDFetcher:   recipeStepPreparationIDFetcher,
		userIDFetcher:                    userIDFetcher,
		recipeDataManager:                recipeDataManager,
		recipeStepDataManager:            recipeStepDataManager,
		recipeStepPreparationDataManager: recipeStepPreparationDataManager,
		encoderDecoder:                   encoder,
		recipeStepPreparationCounter:     recipeStepPreparationCounter,
		reporter:                         reporter,
	}

	return svc, nil
}
