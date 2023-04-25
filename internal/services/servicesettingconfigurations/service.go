package servicesettingconfigurations

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
	serviceName string = "service_settings_service"
)

var _ types.ServiceSettingConfigurationDataService = (*service)(nil)

type (
	// service handles service setting configurations.
	service struct {
		logger                                 logging.Logger
		serviceSettingConfigurationDataManager types.ServiceSettingConfigurationDataManager
		serviceSettingConfigurationIDFetcher   func(*http.Request) string
		serviceSettingNameFetcher              func(*http.Request) string
		sessionContextDataFetcher              func(*http.Request) (*types.SessionContextData, error)
		dataChangesPublisher                   messagequeue.Publisher
		encoderDecoder                         encoding.ServerEncoderDecoder
		tracer                                 tracing.Tracer
	}
)

// ProvideService builds a new ServiceSettingConfigurationsService.
func ProvideService(
	logger logging.Logger,
	cfg *Config,
	serviceSettingConfigurationDataManager types.ServiceSettingConfigurationDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
) (types.ServiceSettingConfigurationDataService, error) {
	dataChangesPublisher, err := publisherProvider.ProvidePublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up service setting configurationss service data changes publisher: %w", err)
	}

	svc := &service{
		logger:                                 logging.EnsureLogger(logger).WithName(serviceName),
		serviceSettingConfigurationIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(ServiceSettingConfigurationIDURIParamKey),
		serviceSettingNameFetcher:              routeParamManager.BuildRouteParamStringIDFetcher(ServiceSettingConfigurationNameURIParamKey),
		sessionContextDataFetcher:              authservice.FetchContextFromRequest,
		serviceSettingConfigurationDataManager: serviceSettingConfigurationDataManager,
		dataChangesPublisher:                   dataChangesPublisher,
		encoderDecoder:                         encoder,
		tracer:                                 tracing.NewTracer(tracerProvider.Tracer(serviceName)),
	}

	return svc, nil
}
