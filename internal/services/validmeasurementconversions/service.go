package validmeasurementconversions

import (
	"context"
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
	serviceName string = "valid_measurement_conversion_service"
)

var _ types.ValidMeasurementConversionDataService = (*service)(nil)

type (
	// service handles valid measurement conversions.
	service struct {
		logger                                logging.Logger
		validMeasurementConversionDataManager types.ValidMeasurementConversionDataManager
		validMeasurementConversionIDFetcher   func(*http.Request) string
		sessionContextDataFetcher             func(*http.Request) (*types.SessionContextData, error)
		dataChangesPublisher                  messagequeue.Publisher
		encoderDecoder                        encoding.ServerEncoderDecoder
		tracer                                tracing.Tracer
	}
)

// ProvideService builds a new ValidMeasurementConversionsService.
func ProvideService(
	_ context.Context,
	logger logging.Logger,
	cfg *Config,
	validMeasurementConversionDataManager types.ValidMeasurementConversionDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
) (types.ValidMeasurementConversionDataService, error) {
	dataChangesPublisher, err := publisherProvider.ProviderPublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up valid measurement conversions service data changes publisher: %w", err)
	}

	svc := &service{
		logger:                                logging.EnsureLogger(logger).WithName(serviceName),
		validMeasurementConversionIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(ValidMeasurementConversionIDURIParamKey),
		sessionContextDataFetcher:             authservice.FetchContextFromRequest,
		validMeasurementConversionDataManager: validMeasurementConversionDataManager,
		dataChangesPublisher:                  dataChangesPublisher,
		encoderDecoder:                        encoder,
		tracer:                                tracing.NewTracer(tracerProvider.Tracer(serviceName)),
	}

	return svc, nil
}
