package validinstruments

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
	// CreateMiddlewareCtxKey is a string alias we can use for referring to valid instrument input data in contexts.
	CreateMiddlewareCtxKey models.ContextKey = "valid_instrument_create_input"
	// UpdateMiddlewareCtxKey is a string alias we can use for referring to valid instrument update data in contexts.
	UpdateMiddlewareCtxKey models.ContextKey = "valid_instrument_update_input"

	counterName        metrics.CounterName = "validInstruments"
	counterDescription string              = "the number of validInstruments managed by the validInstruments service"
	topicName          string              = "valid_instruments"
	serviceName        string              = "valid_instruments_service"
)

var (
	_ models.ValidInstrumentDataServer = (*Service)(nil)
)

type (
	// Service handles to-do list valid instruments
	Service struct {
		logger                     logging.Logger
		validInstrumentDataManager models.ValidInstrumentDataManager
		validInstrumentIDFetcher   ValidInstrumentIDFetcher
		validInstrumentCounter     metrics.UnitCounter
		encoderDecoder             encoding.EncoderDecoder
		reporter                   newsman.Reporter
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
	}

	return svc, nil
}
