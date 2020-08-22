package validpreparations

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
	// createMiddlewareCtxKey is a string alias we can use for referring to valid preparation input data in contexts.
	createMiddlewareCtxKey models.ContextKey = "valid_preparation_create_input"
	// updateMiddlewareCtxKey is a string alias we can use for referring to valid preparation update data in contexts.
	updateMiddlewareCtxKey models.ContextKey = "valid_preparation_update_input"

	counterName        metrics.CounterName = "validPreparations"
	counterDescription string              = "the number of validPreparations managed by the validPreparations service"
	topicName          string              = "valid_preparations"
	serviceName        string              = "valid_preparations_service"
)

var (
	_ models.ValidPreparationDataServer = (*Service)(nil)
)

type (
	// SearchIndex is a type alias for dependency injection's sake
	SearchIndex search.IndexManager

	// Service handles to-do list valid preparations
	Service struct {
		logger                      logging.Logger
		validPreparationDataManager models.ValidPreparationDataManager
		validPreparationIDFetcher   ValidPreparationIDFetcher
		validPreparationCounter     metrics.UnitCounter
		encoderDecoder              encoding.EncoderDecoder
		reporter                    newsman.Reporter
		search                      SearchIndex
	}

	// ValidPreparationIDFetcher is a function that fetches valid preparation IDs.
	ValidPreparationIDFetcher func(*http.Request) uint64
)

// ProvideValidPreparationsService builds a new ValidPreparationsService.
func ProvideValidPreparationsService(
	logger logging.Logger,
	validPreparationDataManager models.ValidPreparationDataManager,
	validPreparationIDFetcher ValidPreparationIDFetcher,
	encoder encoding.EncoderDecoder,
	validPreparationCounterProvider metrics.UnitCounterProvider,
	reporter newsman.Reporter,
	searchIndexManager SearchIndex,
) (*Service, error) {
	validPreparationCounter, err := validPreparationCounterProvider(counterName, counterDescription)
	if err != nil {
		return nil, fmt.Errorf("error initializing counter: %w", err)
	}

	svc := &Service{
		logger:                      logger.WithName(serviceName),
		validPreparationIDFetcher:   validPreparationIDFetcher,
		validPreparationDataManager: validPreparationDataManager,
		encoderDecoder:              encoder,
		validPreparationCounter:     validPreparationCounter,
		reporter:                    reporter,
		search:                      searchIndexManager,
	}

	return svc, nil
}

// ProvideValidPreparationsServiceSearchIndex provides a search index for the service
func ProvideValidPreparationsServiceSearchIndex(
	searchSettings config.SearchSettings,
	indexProvider search.IndexManagerProvider,
	logger logging.Logger,
) (SearchIndex, error) {
	logger.WithValue("index_path", searchSettings.ValidPreparationsIndexPath).Debug("setting up valid preparations search index")

	searchIndex, indexInitErr := indexProvider(searchSettings.ValidPreparationsIndexPath, models.ValidPreparationsSearchIndexName, logger)
	if indexInitErr != nil {
		logger.Error(indexInitErr, "setting up valid preparations search index")
		return nil, indexInitErr
	}

	return searchIndex, nil
}
