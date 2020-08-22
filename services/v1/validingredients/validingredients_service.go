package validingredients

import (
	"fmt"
	"net/http"

	"gitlab.com/prixfixe/prixfixe/internal/v1/config"
	"gitlab.com/prixfixe/prixfixe/internal/v1/encoding"
	"gitlab.com/prixfixe/prixfixe/internal/v1/metrics"
	"gitlab.com/prixfixe/prixfixe/internal/v1/search"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	"gitlab.com/verygoodsoftwarenotvirus/newsman"
)

const (
	// createMiddlewareCtxKey is a string alias we can use for referring to valid ingredient input data in contexts.
	createMiddlewareCtxKey models.ContextKey = "valid_ingredient_create_input"
	// updateMiddlewareCtxKey is a string alias we can use for referring to valid ingredient update data in contexts.
	updateMiddlewareCtxKey models.ContextKey = "valid_ingredient_update_input"

	counterName        metrics.CounterName = "validIngredients"
	counterDescription string              = "the number of validIngredients managed by the validIngredients service"
	topicName          string              = "valid_ingredients"
	serviceName        string              = "valid_ingredients_service"
)

var (
	_ models.ValidIngredientDataServer = (*Service)(nil)
)

type (
	// SearchIndex is a type alias for dependency injection's sake
	SearchIndex search.IndexManager

	// Service handles to-do list valid ingredients
	Service struct {
		logger                     logging.Logger
		validIngredientDataManager models.ValidIngredientDataManager
		validIngredientIDFetcher   ValidIngredientIDFetcher
		validIngredientCounter     metrics.UnitCounter
		encoderDecoder             encoding.EncoderDecoder
		reporter                   newsman.Reporter
		search                     SearchIndex
	}

	// ValidIngredientIDFetcher is a function that fetches valid ingredient IDs.
	ValidIngredientIDFetcher func(*http.Request) uint64
)

// ProvideValidIngredientsService builds a new ValidIngredientsService.
func ProvideValidIngredientsService(
	logger logging.Logger,
	validIngredientDataManager models.ValidIngredientDataManager,
	validIngredientIDFetcher ValidIngredientIDFetcher,
	encoder encoding.EncoderDecoder,
	validIngredientCounterProvider metrics.UnitCounterProvider,
	reporter newsman.Reporter,
	searchIndexManager SearchIndex,
) (*Service, error) {
	validIngredientCounter, err := validIngredientCounterProvider(counterName, counterDescription)
	if err != nil {
		return nil, fmt.Errorf("error initializing counter: %w", err)
	}

	svc := &Service{
		logger:                     logger.WithName(serviceName),
		validIngredientIDFetcher:   validIngredientIDFetcher,
		validIngredientDataManager: validIngredientDataManager,
		encoderDecoder:             encoder,
		validIngredientCounter:     validIngredientCounter,
		reporter:                   reporter,
		search:                     searchIndexManager,
	}

	return svc, nil
}

// ProvideValidIngredientsServiceSearchIndex provides a search index for the service
func ProvideValidIngredientsServiceSearchIndex(
	searchSettings config.SearchSettings,
	indexProvider search.IndexManagerProvider,
	logger logging.Logger,
) (SearchIndex, error) {
	logger.WithValue("index_path", searchSettings.ValidIngredientsIndexPath).Debug("setting up valid ingredients search index")

	searchIndex, indexInitErr := indexProvider(searchSettings.ValidIngredientsIndexPath, models.ValidIngredientsSearchIndexName, logger)
	if indexInitErr != nil {
		logger.Error(indexInitErr, "setting up valid ingredients search index")
		return nil, indexInitErr
	}

	return searchIndex, nil
}
