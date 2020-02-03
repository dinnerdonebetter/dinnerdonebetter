package instruments

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
	// CreateMiddlewareCtxKey is a string alias we can use for referring to instrument input data in contexts
	CreateMiddlewareCtxKey models.ContextKey = "instrument_create_input"
	// UpdateMiddlewareCtxKey is a string alias we can use for referring to instrument update data in contexts
	UpdateMiddlewareCtxKey models.ContextKey = "instrument_update_input"

	counterName        metrics.CounterName = "instruments"
	counterDescription                     = "the number of instruments managed by the instruments service"
	topicName          string              = "instruments"
	serviceName        string              = "instruments_service"
)

var (
	_ models.InstrumentDataServer = (*Service)(nil)
)

type (
	// Service handles to-do list instruments
	Service struct {
		logger              logging.Logger
		instrumentCounter   metrics.UnitCounter
		instrumentDatabase  models.InstrumentDataManager
		userIDFetcher       UserIDFetcher
		instrumentIDFetcher InstrumentIDFetcher
		encoderDecoder      encoding.EncoderDecoder
		reporter            newsman.Reporter
	}

	// UserIDFetcher is a function that fetches user IDs
	UserIDFetcher func(*http.Request) uint64

	// InstrumentIDFetcher is a function that fetches instrument IDs
	InstrumentIDFetcher func(*http.Request) uint64
)

// ProvideInstrumentsService builds a new InstrumentsService
func ProvideInstrumentsService(
	ctx context.Context,
	logger logging.Logger,
	db models.InstrumentDataManager,
	userIDFetcher UserIDFetcher,
	instrumentIDFetcher InstrumentIDFetcher,
	encoder encoding.EncoderDecoder,
	instrumentCounterProvider metrics.UnitCounterProvider,
	reporter newsman.Reporter,
) (*Service, error) {
	instrumentCounter, err := instrumentCounterProvider(counterName, counterDescription)
	if err != nil {
		return nil, fmt.Errorf("error initializing counter: %w", err)
	}

	svc := &Service{
		logger:              logger.WithName(serviceName),
		instrumentDatabase:  db,
		encoderDecoder:      encoder,
		instrumentCounter:   instrumentCounter,
		userIDFetcher:       userIDFetcher,
		instrumentIDFetcher: instrumentIDFetcher,
		reporter:            reporter,
	}

	instrumentCount, err := svc.instrumentDatabase.GetAllInstrumentsCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("setting current instrument count: %w", err)
	}
	svc.instrumentCounter.IncrementBy(ctx, instrumentCount)

	return svc, nil
}
