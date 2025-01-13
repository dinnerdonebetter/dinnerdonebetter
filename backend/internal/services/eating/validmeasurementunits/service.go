package validmeasurementunits

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
	serviceName string = "valid_ingredients_service"
)

var _ types.ValidMeasurementUnitDataService = (*service)(nil)

type (
	// service handles valid ingredients.
	service struct {
		logger                          logging.Logger
		validMeasurementUnitDataManager types.ValidMeasurementUnitDataManager
		dataChangesPublisher            messagequeue.Publisher
		encoderDecoder                  encoding.ServerEncoderDecoder
		tracer                          tracing.Tracer
		validMeasurementUnitSearchIndex textsearch.IndexSearcher[types.ValidMeasurementUnitSearchSubset]
		validMeasurementUnitIDFetcher   func(*http.Request) string
		validIngredientIDFetcher        func(*http.Request) string
		sessionContextDataFetcher       func(*http.Request) (*types.SessionContextData, error)
		useSearchService                bool
	}
)

// ProvideService builds a new ValidMeasurementUnitsService.
func ProvideService(
	ctx context.Context,
	logger logging.Logger,
	cfg *Config,
	searchConfig *textsearchcfg.Config,
	validMeasurementUnitDataManager types.ValidMeasurementUnitDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
	queueConfig *msgconfig.QueuesConfig,
) (types.ValidMeasurementUnitDataService, error) {
	if queueConfig == nil {
		return nil, fmt.Errorf("nil queue config provided")
	}

	dataChangesPublisher, err := publisherProvider.ProvidePublisher(queueConfig.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up %s data changes publisher: %w", serviceName, err)
	}

	searchIndex, err := textsearchcfg.ProvideIndex[types.ValidMeasurementUnitSearchSubset](ctx, logger, tracerProvider, searchConfig, textsearch.IndexTypeValidMeasurementUnits)
	if err != nil {
		return nil, observability.PrepareError(err, nil, "initializing valid preparation index manager")
	}

	svc := &service{
		useSearchService:                cfg.UseSearchService,
		logger:                          logging.EnsureLogger(logger).WithName(serviceName),
		validMeasurementUnitIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(ValidMeasurementUnitIDURIParamKey),
		validIngredientIDFetcher:        routeParamManager.BuildRouteParamStringIDFetcher(ValidIngredientIDURIParamKey),
		sessionContextDataFetcher:       authentication.FetchContextFromRequest,
		validMeasurementUnitDataManager: validMeasurementUnitDataManager,
		dataChangesPublisher:            dataChangesPublisher,
		encoderDecoder:                  encoder,
		tracer:                          tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
		validMeasurementUnitSearchIndex: searchIndex,
	}

	return svc, nil
}
