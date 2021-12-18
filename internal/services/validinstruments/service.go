package validinstruments

import (
	"context"
	"fmt"
	"net/http"

	"github.com/prixfixeco/api_server/internal/encoding"
	"github.com/prixfixeco/api_server/internal/messagequeue/publishers"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/internal/routing"
	"github.com/prixfixeco/api_server/internal/search"
	authservice "github.com/prixfixeco/api_server/internal/services/authentication"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	serviceName string = "valid_instruments_service"
)

var _ types.ValidInstrumentDataService = (*service)(nil)

type (
	// SearchIndex is a type alias for dependency injection's sake.
	SearchIndex search.IndexManager

	// service handles valid instruments.
	service struct {
		logger                     logging.Logger
		validInstrumentDataManager types.ValidInstrumentDataManager
		validInstrumentIDFetcher   func(*http.Request) string
		sessionContextDataFetcher  func(*http.Request) (*types.SessionContextData, error)
		dataChangesPublisher       publishers.Publisher
		encoderDecoder             encoding.ServerEncoderDecoder
		tracer                     tracing.Tracer
		search                     SearchIndex
	}
)

// ProvideService builds a new ValidInstrumentsService.
func ProvideService(
	ctx context.Context,
	logger logging.Logger,
	cfg *Config,
	validInstrumentDataManager types.ValidInstrumentDataManager,
	encoder encoding.ServerEncoderDecoder,
	searchIndexProvider search.IndexManagerProvider,
	routeParamManager routing.RouteParamManager,
	publisherProvider publishers.PublisherProvider,
	tracerProvider tracing.TracerProvider,
) (types.ValidInstrumentDataService, error) {
	searchIndexManager, err := searchIndexProvider.ProvideIndexManager(ctx, logger, "valid_instruments", "name", "variant", "description")
	if err != nil {
		return nil, fmt.Errorf("setting up valid instrument search index: %w", err)
	}

	dataChangesPublisher, err := publisherProvider.ProviderPublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up valid instrument queue pre-archives publisher: %w", err)
	}

	svc := &service{
		logger:                     logging.EnsureLogger(logger).WithName(serviceName),
		validInstrumentIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(ValidInstrumentIDURIParamKey),
		sessionContextDataFetcher:  authservice.FetchContextFromRequest,
		validInstrumentDataManager: validInstrumentDataManager,
		dataChangesPublisher:       dataChangesPublisher,
		encoderDecoder:             encoder,
		search:                     searchIndexManager,
		tracer:                     tracing.NewTracer(tracerProvider.Tracer(serviceName)),
	}

	return svc, nil
}
