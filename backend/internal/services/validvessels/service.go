package validvessels

import (
	"context"
	"errors"
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
	"github.com/dinnerdonebetter/backend/internal/search/text"
	textsearchcfg "github.com/dinnerdonebetter/backend/internal/search/text/config"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	serviceName string = "valid_vessels_service"
)

var _ types.ValidVesselDataService = (*service)(nil)

type (
	// service handles valid vessels.
	service struct {
		logger                    logging.Logger
		validVesselDataManager    types.ValidVesselDataManager
		dataChangesPublisher      messagequeue.Publisher
		encoderDecoder            encoding.ServerEncoderDecoder
		tracer                    tracing.Tracer
		validVesselsSearchIndex   textsearch.IndexSearcher[types.ValidVesselSearchSubset]
		validVesselIDFetcher      func(*http.Request) string
		sessionContextDataFetcher func(*http.Request) (*types.SessionContextData, error)
		useSearchService          bool
	}
)

// ProvideService builds a new ValidVesselDataService.
func ProvideService(
	ctx context.Context,
	logger logging.Logger,
	cfg *Config,
	searchConfig *textsearchcfg.Config,
	validVesselDataManager types.ValidVesselDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
	queueConfig *msgconfig.QueuesConfig,
) (types.ValidVesselDataService, error) {
	if cfg == nil {
		return nil, errors.New("nil config provided")
	}

	if queueConfig == nil {
		return nil, fmt.Errorf("nil queue config provided")
	}

	dataChangesPublisher, err := publisherProvider.ProvidePublisher(queueConfig.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up %s data changes publisher: %w", serviceName, err)
	}

	searchIndex, err := textsearchcfg.ProvideIndex[types.ValidVesselSearchSubset](ctx, logger, tracerProvider, searchConfig, textsearch.IndexTypeValidVessels)
	if err != nil {
		return nil, observability.PrepareError(err, nil, "initializing valid vessel index manager")
	}

	svc := &service{
		useSearchService:          cfg.UseSearchService,
		logger:                    logging.EnsureLogger(logger).WithName(serviceName),
		validVesselIDFetcher:      routeParamManager.BuildRouteParamStringIDFetcher(ValidVesselIDURIParamKey),
		sessionContextDataFetcher: authentication.FetchContextFromRequest,
		validVesselDataManager:    validVesselDataManager,
		dataChangesPublisher:      dataChangesPublisher,
		encoderDecoder:            encoder,
		tracer:                    tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
		validVesselsSearchIndex:   searchIndex,
	}

	return svc, nil
}
