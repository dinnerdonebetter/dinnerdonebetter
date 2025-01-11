package validingredientstates

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
	"github.com/dinnerdonebetter/backend/internal/search/text"
	textsearchcfg "github.com/dinnerdonebetter/backend/internal/search/text/config"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	serviceName string = "valid_preparations_service"
)

var _ types.ValidIngredientStateDataService = (*service)(nil)

type (
	// service handles valid ingredient states.
	service struct {
		logger                           logging.Logger
		validIngredientStateDataManager  types.ValidIngredientStateDataManager
		dataChangesPublisher             messagequeue.Publisher
		encoderDecoder                   encoding.ServerEncoderDecoder
		tracer                           tracing.Tracer
		validIngredientStatesSearchIndex textsearch.IndexSearcher[types.ValidIngredientStateSearchSubset]
		validIngredientStateIDFetcher    func(*http.Request) string
		sessionContextDataFetcher        func(*http.Request) (*types.SessionContextData, error)
		useSearchService                 bool
	}
)

// ProvideService builds a new ValidIngredientStatesService.
func ProvideService(
	ctx context.Context,
	logger logging.Logger,
	cfg *Config,
	searchConfig *textsearchcfg.Config,
	validIngredientStateDataManager types.ValidIngredientStateDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
	queueConfig *msgconfig.QueuesConfig,
) (types.ValidIngredientStateDataService, error) {
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

	searchIndex, err := textsearchcfg.ProvideIndex[types.ValidIngredientStateSearchSubset](ctx, logger, tracerProvider, searchConfig, textsearch.IndexTypeValidIngredientStates)
	if err != nil {
		return nil, observability.PrepareError(err, nil, "initializing valid ingredient state index manager")
	}

	svc := &service{
		useSearchService:                 cfg.UseSearchService,
		logger:                           logging.EnsureLogger(logger).WithName(serviceName),
		validIngredientStateIDFetcher:    routeParamManager.BuildRouteParamStringIDFetcher(ValidIngredientStateIDURIParamKey),
		sessionContextDataFetcher:        authentication.FetchContextFromRequest,
		validIngredientStateDataManager:  validIngredientStateDataManager,
		dataChangesPublisher:             dataChangesPublisher,
		encoderDecoder:                   encoder,
		tracer:                           tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
		validIngredientStatesSearchIndex: searchIndex,
	}

	return svc, nil
}
