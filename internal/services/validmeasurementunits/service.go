package validmeasurementunits

import (
	"fmt"
	"net/http"

	"github.com/prixfixeco/backend/internal/encoding"
	"github.com/prixfixeco/backend/internal/messagequeue"
	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/internal/routing"
	authservice "github.com/prixfixeco/backend/internal/services/authentication"
	"github.com/prixfixeco/backend/pkg/types"
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
		validMeasurementUnitIDFetcher   func(*http.Request) string
		validIngredientIDFetcher        func(*http.Request) string
		sessionContextDataFetcher       func(*http.Request) (*types.SessionContextData, error)
		dataChangesPublisher            messagequeue.Publisher
		encoderDecoder                  encoding.ServerEncoderDecoder
		tracer                          tracing.Tracer
	}
)

// ProvideService builds a new ValidMeasurementUnitsService.
func ProvideService(
	logger logging.Logger,
	cfg *Config,
	validMeasurementUnitDataManager types.ValidMeasurementUnitDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
) (types.ValidMeasurementUnitDataService, error) {
	dataChangesPublisher, err := publisherProvider.ProviderPublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up valid ingredients service data changes publisher: %w", err)
	}

	svc := &service{
		logger:                          logging.EnsureLogger(logger).WithName(serviceName),
		validMeasurementUnitIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(ValidMeasurementUnitIDURIParamKey),
		validIngredientIDFetcher:        routeParamManager.BuildRouteParamStringIDFetcher(ValidIngredientIDURIParamKey),
		sessionContextDataFetcher:       authservice.FetchContextFromRequest,
		validMeasurementUnitDataManager: validMeasurementUnitDataManager,
		dataChangesPublisher:            dataChangesPublisher,
		encoderDecoder:                  encoder,
		tracer:                          tracing.NewTracer(tracerProvider.Tracer(serviceName)),
	}

	return svc, nil
}
