package requiredpreparationinstruments

import (
	"context"
	"fmt"
	"net/http"

	"gitlab.com/prixfixe/prixfixe/internal/v1/encoding"
	"gitlab.com/prixfixe/prixfixe/internal/v1/metrics"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	"gitlab.com/verygoodsoftwarenotvirus/newsman"
)

const (
	// CreateMiddlewareCtxKey is a string alias we can use for referring to required preparation instrument input data in contexts
	CreateMiddlewareCtxKey models.ContextKey = "required_preparation_instrument_create_input"
	// UpdateMiddlewareCtxKey is a string alias we can use for referring to required preparation instrument update data in contexts
	UpdateMiddlewareCtxKey models.ContextKey = "required_preparation_instrument_update_input"

	counterName        metrics.CounterName = "requiredPreparationInstruments"
	counterDescription                     = "the number of requiredPreparationInstruments managed by the requiredPreparationInstruments service"
	topicName          string              = "required_preparation_instruments"
	serviceName        string              = "required_preparation_instruments_service"
)

var (
	_ models.RequiredPreparationInstrumentDataServer = (*Service)(nil)
)

type (
	// Service handles to-do list required preparation instruments
	Service struct {
		logger                                 logging.Logger
		requiredPreparationInstrumentCounter   metrics.UnitCounter
		requiredPreparationInstrumentDatabase  models.RequiredPreparationInstrumentDataManager
		userIDFetcher                          UserIDFetcher
		requiredPreparationInstrumentIDFetcher RequiredPreparationInstrumentIDFetcher
		encoderDecoder                         encoding.EncoderDecoder
		reporter                               newsman.Reporter
	}

	// UserIDFetcher is a function that fetches user IDs
	UserIDFetcher func(*http.Request) uint64

	// RequiredPreparationInstrumentIDFetcher is a function that fetches required preparation instrument IDs
	RequiredPreparationInstrumentIDFetcher func(*http.Request) uint64
)

// ProvideRequiredPreparationInstrumentsService builds a new RequiredPreparationInstrumentsService
func ProvideRequiredPreparationInstrumentsService(
	ctx context.Context,
	logger logging.Logger,
	db models.RequiredPreparationInstrumentDataManager,
	userIDFetcher UserIDFetcher,
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
		logger:                                 logger.WithName(serviceName),
		requiredPreparationInstrumentDatabase:  db,
		encoderDecoder:                         encoder,
		requiredPreparationInstrumentCounter:   requiredPreparationInstrumentCounter,
		userIDFetcher:                          userIDFetcher,
		requiredPreparationInstrumentIDFetcher: requiredPreparationInstrumentIDFetcher,
		reporter:                               reporter,
	}

	requiredPreparationInstrumentCount, err := svc.requiredPreparationInstrumentDatabase.GetAllRequiredPreparationInstrumentsCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("setting current required preparation instrument count: %w", err)
	}
	svc.requiredPreparationInstrumentCounter.IncrementBy(ctx, requiredPreparationInstrumentCount)

	return svc, nil
}
