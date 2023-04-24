package servicesettings

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

var _ types.ServiceSettingDataService = (*service)(nil)

type (
	// service handles service settings.
	service struct {
		logger                    logging.Logger
		serviceSettingDataManager types.ServiceSettingDataManager
		serviceSettingIDFetcher   func(*http.Request) string
		sessionContextDataFetcher func(*http.Request) (*types.SessionContextData, error)
		dataChangesPublisher      messagequeue.Publisher
		encoderDecoder            encoding.ServerEncoderDecoder
		tracer                    tracing.Tracer
	}
)

// ProvideService builds a new ServiceSettingsService.
func ProvideService(
	logger logging.Logger,
	cfg *Config,
	serviceSettingDataManager types.ServiceSettingDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
) (types.ServiceSettingDataService, error) {
	dataChangesPublisher, err := publisherProvider.ProvidePublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up service settingss service data changes publisher: %w", err)
	}

	svc := &service{
		logger:                    logging.EnsureLogger(logger).WithName(serviceName),
		serviceSettingIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(ServiceSettingIDURIParamKey),
		sessionContextDataFetcher: authservice.FetchContextFromRequest,
		serviceSettingDataManager: serviceSettingDataManager,
		dataChangesPublisher:      dataChangesPublisher,
		encoderDecoder:            encoder,
		tracer:                    tracing.NewTracer(tracerProvider.Tracer(serviceName)),
	}

	return svc, nil
}
