package requiredpreparationinstruments

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
	// createMiddlewareCtxKey is a string alias we can use for referring to required preparation instrument input data in contexts.
	createMiddlewareCtxKey models.ContextKey = "required_preparation_instrument_create_input"
	// updateMiddlewareCtxKey is a string alias we can use for referring to required preparation instrument update data in contexts.
	updateMiddlewareCtxKey models.ContextKey = "required_preparation_instrument_update_input"

	counterName        metrics.CounterName = "requiredPreparationInstruments"
	counterDescription string              = "the number of requiredPreparationInstruments managed by the requiredPreparationInstruments service"
	topicName          string              = "required_preparation_instruments"
	serviceName        string              = "required_preparation_instruments_service"
)

var (
	_ models.RequiredPreparationInstrumentDataServer = (*Service)(nil)
)

type (
	// Service handles to-do list required preparation instruments
	Service struct {
		logger                                   logging.Logger
		requiredPreparationInstrumentDataManager models.RequiredPreparationInstrumentDataManager
		requiredPreparationInstrumentIDFetcher   RequiredPreparationInstrumentIDFetcher
		requiredPreparationInstrumentCounter     metrics.UnitCounter
		encoderDecoder                           encoding.EncoderDecoder
		reporter                                 newsman.Reporter
	}

	// RequiredPreparationInstrumentIDFetcher is a function that fetches required preparation instrument IDs.
	RequiredPreparationInstrumentIDFetcher func(*http.Request) uint64
)

// ProvideRequiredPreparationInstrumentsService builds a new RequiredPreparationInstrumentsService.
func ProvideRequiredPreparationInstrumentsService(
	logger logging.Logger,
	requiredPreparationInstrumentDataManager models.RequiredPreparationInstrumentDataManager,
	requiredPreparationInstrumentIDFetcher RequiredPreparationInstrumentIDFetcher,
	encoder encoding.EncoderDecoder,
	requiredPreparationInstrumentCounterProvider metrics.UnitCounterProvider,
	reporter newsman.Reporter,
) (*Service, error) {
	requiredPreparationInstrumentCounter, err := requiredPreparationInstrumentCounterProvider(counterName, counterDescription)
	if err != nil {
		return nil, fmt.Errorf("error initializing counter: %w", err)
	}

	svc := &Service{
		logger:                                   logger.WithName(serviceName),
		requiredPreparationInstrumentIDFetcher:   requiredPreparationInstrumentIDFetcher,
		requiredPreparationInstrumentDataManager: requiredPreparationInstrumentDataManager,
		encoderDecoder:                           encoder,
		requiredPreparationInstrumentCounter:     requiredPreparationInstrumentCounter,
		reporter:                                 reporter,
	}

	return svc, nil
}
