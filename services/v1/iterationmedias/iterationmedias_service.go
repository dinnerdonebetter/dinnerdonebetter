package iterationmedias

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
	// createMiddlewareCtxKey is a string alias we can use for referring to iteration media input data in contexts.
	createMiddlewareCtxKey models.ContextKey = "iteration_media_create_input"
	// updateMiddlewareCtxKey is a string alias we can use for referring to iteration media update data in contexts.
	updateMiddlewareCtxKey models.ContextKey = "iteration_media_update_input"

	counterName        metrics.CounterName = "iterationMedias"
	counterDescription string              = "the number of iterationMedias managed by the iterationMedias service"
	topicName          string              = "iteration_medias"
	serviceName        string              = "iteration_medias_service"
)

var (
	_ models.IterationMediaDataServer = (*Service)(nil)
)

type (
	// Service handles to-do list iteration medias
	Service struct {
		logger                     logging.Logger
		recipeDataManager          models.RecipeDataManager
		recipeIterationDataManager models.RecipeIterationDataManager
		iterationMediaDataManager  models.IterationMediaDataManager
		recipeIDFetcher            RecipeIDFetcher
		recipeIterationIDFetcher   RecipeIterationIDFetcher
		iterationMediaIDFetcher    IterationMediaIDFetcher
		userIDFetcher              UserIDFetcher
		iterationMediaCounter      metrics.UnitCounter
		encoderDecoder             encoding.EncoderDecoder
		reporter                   newsman.Reporter
	}

	// UserIDFetcher is a function that fetches user IDs.
	UserIDFetcher func(*http.Request) uint64

	// RecipeIDFetcher is a function that fetches recipe IDs.
	RecipeIDFetcher func(*http.Request) uint64

	// RecipeIterationIDFetcher is a function that fetches recipe iteration IDs.
	RecipeIterationIDFetcher func(*http.Request) uint64

	// IterationMediaIDFetcher is a function that fetches iteration media IDs.
	IterationMediaIDFetcher func(*http.Request) uint64
)

// ProvideIterationMediasService builds a new IterationMediasService.
func ProvideIterationMediasService(
	logger logging.Logger,
	recipeDataManager models.RecipeDataManager,
	recipeIterationDataManager models.RecipeIterationDataManager,
	iterationMediaDataManager models.IterationMediaDataManager,
	recipeIDFetcher RecipeIDFetcher,
	recipeIterationIDFetcher RecipeIterationIDFetcher,
	iterationMediaIDFetcher IterationMediaIDFetcher,
	userIDFetcher UserIDFetcher,
	encoder encoding.EncoderDecoder,
	iterationMediaCounterProvider metrics.UnitCounterProvider,
	reporter newsman.Reporter,
) (*Service, error) {
	iterationMediaCounter, err := iterationMediaCounterProvider(counterName, counterDescription)
	if err != nil {
		return nil, fmt.Errorf("error initializing counter: %w", err)
	}

	svc := &Service{
		logger:                     logger.WithName(serviceName),
		recipeIDFetcher:            recipeIDFetcher,
		recipeIterationIDFetcher:   recipeIterationIDFetcher,
		iterationMediaIDFetcher:    iterationMediaIDFetcher,
		userIDFetcher:              userIDFetcher,
		recipeDataManager:          recipeDataManager,
		recipeIterationDataManager: recipeIterationDataManager,
		iterationMediaDataManager:  iterationMediaDataManager,
		encoderDecoder:             encoder,
		iterationMediaCounter:      iterationMediaCounter,
		reporter:                   reporter,
	}

	return svc, nil
}
