package validingredientmeasurementunits

import (
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/messagequeue"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/routing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	serviceName string = "valid_ingredient_measurement_units_service"
)

var _ types.ValidIngredientMeasurementUnitDataService = (*service)(nil)

type (
	// service handles valid ingredient measurement units.
	service struct {
		logger                                    logging.Logger
		validIngredientMeasurementUnitDataManager types.ValidIngredientMeasurementUnitDataManager
		validIngredientMeasurementUnitIDFetcher   func(*http.Request) string
		validIngredientIDFetcher                  func(*http.Request) string
		validMeasurementUnitIDFetcher             func(*http.Request) string
		sessionContextDataFetcher                 func(*http.Request) (*types.SessionContextData, error)
		dataChangesPublisher                      messagequeue.Publisher
		encoderDecoder                            encoding.ServerEncoderDecoder
		tracer                                    tracing.Tracer
	}
)

// ProvideService builds a new ValidIngredientMeasurementUnitsService.
func ProvideService(
	logger logging.Logger,
	cfg *Config,
	validIngredientMeasurementUnitDataManager types.ValidIngredientMeasurementUnitDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
) (types.ValidIngredientMeasurementUnitDataService, error) {
	dataChangesPublisher, err := publisherProvider.ProvidePublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up %s data changes publisher: %w", serviceName, err)
	}

	svc := &service{
		logger:                                    logging.EnsureLogger(logger).WithName(serviceName),
		validIngredientMeasurementUnitIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(ValidIngredientMeasurementUnitIDURIParamKey),
		validIngredientIDFetcher:                  routeParamManager.BuildRouteParamStringIDFetcher(ValidIngredientIDURIParamKey),
		validMeasurementUnitIDFetcher:             routeParamManager.BuildRouteParamStringIDFetcher(ValidMeasurementUnitIDURIParamKey),
		sessionContextDataFetcher:                 authentication.FetchContextFromRequest,
		validIngredientMeasurementUnitDataManager: validIngredientMeasurementUnitDataManager,
		dataChangesPublisher:                      dataChangesPublisher,
		encoderDecoder:                            encoder,
		tracer:                                    tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
	}

	return svc, nil
}
