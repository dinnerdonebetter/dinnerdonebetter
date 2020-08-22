package validinstruments

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
	// createMiddlewareCtxKey is a string alias we can use for referring to valid instrument input data in contexts.
	createMiddlewareCtxKey models.ContextKey = "valid_instrument_create_input"
	// updateMiddlewareCtxKey is a string alias we can use for referring to valid instrument update data in contexts.
	updateMiddlewareCtxKey models.ContextKey = "valid_instrument_update_input"

	counterName        metrics.CounterName = "validInstruments"
	counterDescription string              = "the number of validInstruments managed by the validInstruments service"
	topicName          string              = "valid_instruments"
	serviceName        string              = "valid_instruments_service"
)

var (
	_ models.ValidInstrumentDataServer = (*Service)(nil)
)

type (
	// SearchIndex is a type alias for dependency injection's sake
	SearchIndex search.IndexManager

	// Service handles to-do list valid instruments
	Service struct {
		logger                     logging.Logger
		validInstrumentDataManager models.ValidInstrumentDataManager
		validInstrumentIDFetcher   ValidInstrumentIDFetcher
		validInstrumentCounter     metrics.UnitCounter
		encoderDecoder             encoding.EncoderDecoder
		reporter                   newsman.Reporter
		search                     SearchIndex
	}

	// ValidInstrumentIDFetcher is a function that fetches valid instrument IDs.
	ValidInstrumentIDFetcher func(*http.Request) uint64
)

// ProvideValidInstrumentsService builds a new ValidInstrumentsService.
func ProvideValidInstrumentsService(
	logger logging.Logger,
	validInstrumentDataManager models.ValidInstrumentDataManager,
	validInstrumentIDFetcher ValidInstrumentIDFetcher,
	encoder encoding.EncoderDecoder,
	validInstrumentCounterProvider metrics.UnitCounterProvider,
	reporter newsman.Reporter,
	searchIndexManager SearchIndex,
) (*Service, error) {
	validInstrumentCounter, err := validInstrumentCounterProvider(counterName, counterDescription)
	if err != nil {
		return nil, fmt.Errorf("error initializing counter: %w", err)
	}

	svc := &Service{
		logger:                     logger.WithName(serviceName),
		validInstrumentIDFetcher:   validInstrumentIDFetcher,
		validInstrumentDataManager: validInstrumentDataManager,
		encoderDecoder:             encoder,
		validInstrumentCounter:     validInstrumentCounter,
		reporter:                   reporter,
		search:                     searchIndexManager,
	}

	return svc, nil
}

// ProvideValidInstrumentsServiceSearchIndex provides a search index for the service
func ProvideValidInstrumentsServiceSearchIndex(
	searchSettings config.SearchSettings,
	indexProvider search.IndexManagerProvider,
	logger logging.Logger,
) (SearchIndex, error) {
	logger.WithValue("index_path", searchSettings.ValidInstrumentsIndexPath).Debug("setting up valid instruments search index")

	searchIndex, indexInitErr := indexProvider(searchSettings.ValidInstrumentsIndexPath, models.ValidInstrumentsSearchIndexName, logger)
	if indexInitErr != nil {
		logger.Error(indexInitErr, "setting up valid instruments search index")
		return nil, indexInitErr
	}

	return searchIndex, nil
}
