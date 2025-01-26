package servicesettingconfigurations

import (
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/lib/authentication/sessioncontext"
	"github.com/dinnerdonebetter/backend/internal/lib/encoding"
	"github.com/dinnerdonebetter/backend/internal/lib/messagequeue"
	msgconfig "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/lib/routing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	serviceName string = "service_setting_configurations_service"
)

var _ types.ServiceSettingConfigurationDataService = (*service)(nil)

type (
	// service handles service setting configurations.
	service struct {
		logger                                 logging.Logger
		serviceSettingConfigurationDataManager types.ServiceSettingConfigurationDataManager
		serviceSettingConfigurationIDFetcher   func(*http.Request) string
		serviceSettingNameFetcher              func(*http.Request) string
		sessionContextDataFetcher              func(*http.Request) (*sessioncontext.SessionContextData, error)
		dataChangesPublisher                   messagequeue.Publisher
		encoderDecoder                         encoding.ServerEncoderDecoder
		tracer                                 tracing.Tracer
	}
)

// ProvideService builds a new ServiceSettingConfigurationsService.
func ProvideService(
	logger logging.Logger,
	serviceSettingConfigurationDataManager types.ServiceSettingConfigurationDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
	queueConfig *msgconfig.QueuesConfig,
) (types.ServiceSettingConfigurationDataService, error) {
	if queueConfig == nil {
		return nil, fmt.Errorf("nil queue config provided")
	}

	dataChangesPublisher, err := publisherProvider.ProvidePublisher(queueConfig.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up %s data changes publisher: %w", serviceName, err)
	}

	svc := &service{
		logger:                                 logging.EnsureLogger(logger).WithName(serviceName),
		serviceSettingConfigurationIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(ServiceSettingConfigurationIDURIParamKey),
		serviceSettingNameFetcher:              routeParamManager.BuildRouteParamStringIDFetcher(ServiceSettingConfigurationNameURIParamKey),
		sessionContextDataFetcher:              sessioncontext.FetchContextFromRequest,
		serviceSettingConfigurationDataManager: serviceSettingConfigurationDataManager,
		dataChangesPublisher:                   dataChangesPublisher,
		encoderDecoder:                         encoder,
		tracer:                                 tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
	}

	return svc, nil
}
