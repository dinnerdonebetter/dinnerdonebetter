package validingredients

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
	"github.com/dinnerdonebetter/backend/internal/search"
	searchcfg "github.com/dinnerdonebetter/backend/internal/search/config"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	serviceName string = "valid_ingredients_service"
)

var _ types.ValidIngredientDataService = (*service)(nil)

type (
	// service handles valid ingredients.
	service struct {
		cfg                        *Config
		logger                     logging.Logger
		validIngredientDataManager types.ValidIngredientDataManager
		validIngredientIDFetcher   func(*http.Request) string
		validPreparationIDFetcher  func(*http.Request) string
		sessionContextDataFetcher  func(*http.Request) (*types.SessionContextData, error)
		dataChangesPublisher       messagequeue.Publisher
		encoderDecoder             encoding.ServerEncoderDecoder
		tracer                     tracing.Tracer
		validIngredientSearchIndex search.IndexSearcher[types.ValidIngredientSearchSubset]
	}
)

// ProvideService builds a new ValidIngredientsService.
func ProvideService(
	ctx context.Context,
	logger logging.Logger,
	cfg *Config,
	searchConfig *searchcfg.Config,
	validIngredientDataManager types.ValidIngredientDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
) (types.ValidIngredientDataService, error) {
	dataChangesPublisher, err := publisherProvider.ProvidePublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up %s data changes publisher: %w", serviceName, err)
	}

	searchIndex, err := searchcfg.ProvideIndex[types.ValidIngredientSearchSubset](ctx, logger, tracerProvider, searchConfig, search.IndexTypeValidIngredients)
	if err != nil {
		return nil, observability.PrepareError(err, nil, "initializing valid ingredient index manager")
	}

	svc := &service{
		cfg:                        cfg,
		logger:                     logging.EnsureLogger(logger).WithName(serviceName),
		validIngredientIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(ValidIngredientIDURIParamKey),
		validPreparationIDFetcher:  routeParamManager.BuildRouteParamStringIDFetcher(ValidPreparationIDURIParamKey),
		sessionContextDataFetcher:  authentication.FetchContextFromRequest,
		validIngredientDataManager: validIngredientDataManager,
		dataChangesPublisher:       dataChangesPublisher,
		encoderDecoder:             encoder,
		tracer:                     tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
		validIngredientSearchIndex: searchIndex,
	}

	return svc, nil
}
