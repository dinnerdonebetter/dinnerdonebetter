package validinstruments

import (
	"fmt"
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
	counterName        metrics.CounterName = "valid_instruments"
	counterDescription string              = "the number of valid instruments managed by the valid instruments service"
	serviceName        string              = "valid_instruments_service"
)

var _ types.ValidInstrumentDataService = (*service)(nil)

type (
	// SearchIndex is a type alias for dependency injection's sake.
	SearchIndex search.IndexManager

	// service handles valid instruments.
	service struct {
		logger                     logging.Logger
		validInstrumentDataManager types.ValidInstrumentDataManager
		validInstrumentIDFetcher   func(*http.Request) uint64
		sessionContextDataFetcher  func(*http.Request) (*types.SessionContextData, error)
		validInstrumentCounter     metrics.UnitCounter
		encoderDecoder             encoding.ServerEncoderDecoder
		tracer                     tracing.Tracer
		search                     SearchIndex
	}
)

// ProvideService builds a new ValidInstrumentsService.
func ProvideService(
	logger logging.Logger,
	cfg Config,
	validInstrumentDataManager types.ValidInstrumentDataManager,
	encoder encoding.ServerEncoderDecoder,
	counterProvider metrics.UnitCounterProvider,
	searchIndexProvider search.IndexManagerProvider,
	routeParamManager routing.RouteParamManager,
) (types.ValidInstrumentDataService, error) {
	searchIndexManager, err := searchIndexProvider(search.IndexPath(cfg.SearchIndexPath), "valid_instruments", logger)
	if err != nil {
		return nil, fmt.Errorf("setting up search index: %w", err)
	}

	svc := &service{
		logger:                     logging.EnsureLogger(logger).WithName(serviceName),
		validInstrumentIDFetcher:   routeParamManager.BuildRouteParamIDFetcher(logger, ValidInstrumentIDURIParamKey, "valid_instrument"),
		sessionContextDataFetcher:  authservice.FetchContextFromRequest,
		validInstrumentDataManager: validInstrumentDataManager,
		encoderDecoder:             encoder,
		validInstrumentCounter:     metrics.EnsureUnitCounter(counterProvider, logger, counterName, counterDescription),
		search:                     searchIndexManager,
		tracer:                     tracing.NewTracer(serviceName),
	}

	return svc, nil
}
