package validpreparationvessels

import (
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/messagequeue"
	msgconfig "github.com/dinnerdonebetter/backend/internal/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/routing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	serviceName string = "valid_preparation_vessels_service"
)

var _ types.ValidPreparationVesselDataService = (*service)(nil)

type (
	// service handles valid preparation vessels.
	service struct {
		logger                            logging.Logger
		validPreparationVesselDataManager types.ValidPreparationVesselDataManager
		validPreparationVesselIDFetcher   func(*http.Request) string
		validPreparationIDFetcher         func(*http.Request) string
		validVesselIDFetcher              func(*http.Request) string
		sessionContextDataFetcher         func(*http.Request) (*types.SessionContextData, error)
		dataChangesPublisher              messagequeue.Publisher
		encoderDecoder                    encoding.ServerEncoderDecoder
		tracer                            tracing.Tracer
	}
)

// ProvideService builds a new ValidPreparationVesselsService.
func ProvideService(
	logger logging.Logger,
	validPreparationVesselDataManager types.ValidPreparationVesselDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
	queueConfig *msgconfig.QueuesConfig,
) (types.ValidPreparationVesselDataService, error) {
	if queueConfig == nil {
		return nil, fmt.Errorf("nil queue config provided")
	}

	dataChangesPublisher, err := publisherProvider.ProvidePublisher(queueConfig.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up %s data changes publisher: %w", serviceName, err)
	}

	svc := &service{
		logger:                            logging.EnsureLogger(logger).WithName(serviceName),
		validPreparationVesselIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(ValidPreparationVesselIDURIParamKey),
		validPreparationIDFetcher:         routeParamManager.BuildRouteParamStringIDFetcher(ValidPreparationIDURIParamKey),
		validVesselIDFetcher:              routeParamManager.BuildRouteParamStringIDFetcher(ValidVesselIDURIParamKey),
		sessionContextDataFetcher:         authentication.FetchContextFromRequest,
		validPreparationVesselDataManager: validPreparationVesselDataManager,
		dataChangesPublisher:              dataChangesPublisher,
		encoderDecoder:                    encoder,
		tracer:                            tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
	}

	return svc, nil
}
