package validinstruments

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/messagequeue"
	msgconfig "github.com/dinnerdonebetter/backend/internal/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/routing"
	textsearch "github.com/dinnerdonebetter/backend/internal/search/text"
	textsearchcfg "github.com/dinnerdonebetter/backend/internal/search/text/config"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	serviceName string = "valid_instruments_service"
)

var _ types.ValidInstrumentDataService = (*service)(nil)

type (
	// service handles valid instruments.
	service struct {
		logger                     logging.Logger
		validInstrumentDataManager types.ValidInstrumentDataManager
		dataChangesPublisher       messagequeue.Publisher
		encoderDecoder             encoding.ServerEncoderDecoder
		tracer                     tracing.Tracer
		validInstrumentSearchIndex textsearch.IndexSearcher[types.ValidInstrumentSearchSubset]
		validInstrumentIDFetcher   func(*http.Request) string
		sessionContextDataFetcher  func(*http.Request) (*types.SessionContextData, error)
		useSearchService           bool
	}
)

// ProvideService builds a new ValidInstrumentsService.
func ProvideService(
	ctx context.Context,
	logger logging.Logger,
	cfg *Config,
	searchConfig *textsearchcfg.Config,
	validInstrumentDataManager types.ValidInstrumentDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
	queueConfig *msgconfig.QueuesConfig,
) (types.ValidInstrumentDataService, error) {
	if cfg == nil {
		return nil, fmt.Errorf("nil config provided")
	}

	if queueConfig == nil {
		return nil, fmt.Errorf("nil queue config provided")
	}

	dataChangesPublisher, err := publisherProvider.ProvidePublisher(queueConfig.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up %s data changes publisher: %w", serviceName, err)
	}

	searchIndex, err := textsearchcfg.ProvideIndex[types.ValidInstrumentSearchSubset](ctx, logger, tracerProvider, searchConfig, textsearch.IndexTypeValidInstruments)
	if err != nil {
		return nil, observability.PrepareError(err, nil, "initializing valid instrument index manager")
	}

	svc := &service{
		useSearchService:           cfg.UseSearchService,
		logger:                     logging.EnsureLogger(logger).WithName(serviceName),
		validInstrumentIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(ValidInstrumentIDURIParamKey),
		sessionContextDataFetcher:  authentication.FetchContextFromRequest,
		validInstrumentDataManager: validInstrumentDataManager,
		dataChangesPublisher:       dataChangesPublisher,
		encoderDecoder:             encoder,
		tracer:                     tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
		validInstrumentSearchIndex: searchIndex,
	}

	return svc, nil
}
