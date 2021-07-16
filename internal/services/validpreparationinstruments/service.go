package validpreparationinstruments

import (
	"net/http"

	"gitlab.com/prixfixe/prixfixe/internal/encoding"
	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
	"gitlab.com/prixfixe/prixfixe/internal/observability/metrics"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	routing "gitlab.com/prixfixe/prixfixe/internal/routing"
	"gitlab.com/prixfixe/prixfixe/internal/search"
	authservice "gitlab.com/prixfixe/prixfixe/internal/services/authentication"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

const (
	counterName        metrics.CounterName = "valid_preparation_instruments"
	counterDescription string              = "the number of valid preparation instruments managed by the valid preparation instruments service"
	serviceName        string              = "valid_preparation_instruments_service"
)

var _ types.ValidPreparationInstrumentDataService = (*service)(nil)

type (
	// SearchIndex is a type alias for dependency injection's sake.
	SearchIndex search.IndexManager

	// service handles valid preparation instruments.
	service struct {
		logger                                logging.Logger
		validPreparationInstrumentDataManager types.ValidPreparationInstrumentDataManager
		validPreparationInstrumentIDFetcher   func(*http.Request) uint64
		sessionContextDataFetcher             func(*http.Request) (*types.SessionContextData, error)
		validPreparationInstrumentCounter     metrics.UnitCounter
		encoderDecoder                        encoding.ServerEncoderDecoder
		tracer                                tracing.Tracer
	}
)

// ProvideService builds a new ValidPreparationInstrumentsService.
func ProvideService(
	logger logging.Logger,
	cfg Config,
	validPreparationInstrumentDataManager types.ValidPreparationInstrumentDataManager,
	encoder encoding.ServerEncoderDecoder,
	counterProvider metrics.UnitCounterProvider,
	routeParamManager routing.RouteParamManager,
) (types.ValidPreparationInstrumentDataService, error) {
	svc := &service{
		logger:                                logging.EnsureLogger(logger).WithName(serviceName),
		validPreparationInstrumentIDFetcher:   routeParamManager.BuildRouteParamIDFetcher(logger, ValidPreparationInstrumentIDURIParamKey, "valid_preparation_instrument"),
		sessionContextDataFetcher:             authservice.FetchContextFromRequest,
		validPreparationInstrumentDataManager: validPreparationInstrumentDataManager,
		encoderDecoder:                        encoder,
		validPreparationInstrumentCounter:     metrics.EnsureUnitCounter(counterProvider, logger, counterName, counterDescription),
		tracer:                                tracing.NewTracer(serviceName),
	}

	return svc, nil
}
