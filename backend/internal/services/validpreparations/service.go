package validpreparations

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/messagequeue"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/routing"
	"github.com/dinnerdonebetter/backend/internal/search/text"
	textsearchcfg "github.com/dinnerdonebetter/backend/internal/search/text/config"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	serviceName string = "valid_preparations_service"
)

var _ types.ValidPreparationDataService = (*service)(nil)

type (
	// service handles valid preparations.
	service struct {
		cfg                          *Config
		logger                       logging.Logger
		validPreparationDataManager  types.ValidPreparationDataManager
		validPreparationIDFetcher    func(*http.Request) string
		sessionContextDataFetcher    func(*http.Request) (*types.SessionContextData, error)
		dataChangesPublisher         messagequeue.Publisher
		encoderDecoder               encoding.ServerEncoderDecoder
		tracer                       tracing.Tracer
		validPreparationsSearchIndex textsearch.IndexSearcher[types.ValidPreparationSearchSubset]
	}
)

// ProvideService builds a new ValidPreparationsService.
func ProvideService(
	ctx context.Context,
	logger logging.Logger,
	cfg *Config,
	searchConfig *textsearchcfg.Config,
	validPreparationDataManager types.ValidPreparationDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
) (types.ValidPreparationDataService, error) {
	dataChangesPublisher, err := publisherProvider.ProvidePublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up %s data changes publisher: %w", serviceName, err)
	}

	searchIndex, err := textsearchcfg.ProvideIndex[types.ValidPreparationSearchSubset](ctx, logger, tracerProvider, searchConfig, textsearch.IndexTypeValidPreparations)
	if err != nil {
		return nil, observability.PrepareError(err, nil, "initializing valid preparation index manager")
	}

	svc := &service{
		logger:                       logging.EnsureLogger(logger).WithName(serviceName),
		validPreparationIDFetcher:    routeParamManager.BuildRouteParamStringIDFetcher(ValidPreparationIDURIParamKey),
		sessionContextDataFetcher:    authentication.FetchContextFromRequest,
		validPreparationDataManager:  validPreparationDataManager,
		dataChangesPublisher:         dataChangesPublisher,
		cfg:                          cfg,
		encoderDecoder:               encoder,
		tracer:                       tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
		validPreparationsSearchIndex: searchIndex,
	}

	return svc, nil
}
