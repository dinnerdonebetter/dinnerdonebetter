package validingredientmeasurementunits

import (
	"fmt"
	"net/http"

	"github.com/prixfixeco/api_server/internal/messagequeue"

	"github.com/prixfixeco/api_server/internal/encoding"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/internal/routing"
	authservice "github.com/prixfixeco/api_server/internal/services/authentication"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	serviceName string = "valid_ingredient_preparations_service"
)

var _ types.ValidIngredientMeasurementUnitDataService = (*service)(nil)

type (
	// service handles valid ingredient measurement units.
	service struct {
		logger                                    logging.Logger
		validIngredientMeasurementUnitDataManager types.ValidIngredientMeasurementUnitDataManager
		validIngredientMeasurementUnitIDFetcher   func(*http.Request) string
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
	dataChangesPublisher, err := publisherProvider.ProviderPublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up valid ingredient preprarations service data changes publisher: %w", err)
	}

	svc := &service{
		logger:                                    logging.EnsureLogger(logger).WithName(serviceName),
		validIngredientMeasurementUnitIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(ValidIngredientMeasurementUnitIDURIParamKey),
		sessionContextDataFetcher:                 authservice.FetchContextFromRequest,
		validIngredientMeasurementUnitDataManager: validIngredientMeasurementUnitDataManager,
		dataChangesPublisher:                      dataChangesPublisher,
		encoderDecoder:                            encoder,
		tracer:                                    tracing.NewTracer(tracerProvider.Tracer(serviceName)),
	}

	return svc, nil
}
