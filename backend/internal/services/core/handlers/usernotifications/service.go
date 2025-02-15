package usernotifications

import (
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/lib/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/lib/encoding"
	"github.com/dinnerdonebetter/backend/internal/lib/messagequeue"
	msgconfig "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/lib/routing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	serviceName string = "user_notifications_service"
)

var _ types.UserNotificationDataService = (*service)(nil)

type (
	// service handles user notifications.
	service struct {
		logger                      logging.Logger
		dataChangesPublisher        messagequeue.Publisher
		tracer                      tracing.Tracer
		encoderDecoder              encoding.ServerEncoderDecoder
		userNotificationDataManager types.UserNotificationDataManager
		sessionContextDataFetcher   func(*http.Request) (*sessions.ContextData, error)
		userNotificationIDFetcher   func(*http.Request) string
	}
)

// ProvideService builds a new UserNotificationsService.
func ProvideService(
	logger logging.Logger,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
	userNotificationDataManager types.UserNotificationDataManager,
	queueConfig *msgconfig.QueuesConfig,
) (types.UserNotificationDataService, error) {
	if queueConfig == nil {
		return nil, fmt.Errorf("nil queue config provided")
	}

	dataChangesPublisher, err := publisherProvider.ProvidePublisher(queueConfig.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up %s data changes publisher: %w", serviceName, err)
	}

	svc := &service{
		logger:                      logging.EnsureLogger(logger).WithName(serviceName),
		userNotificationIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(UserNotificationIDURIParamKey),
		sessionContextDataFetcher:   sessions.FetchContextFromRequest,
		dataChangesPublisher:        dataChangesPublisher,
		encoderDecoder:              encoder,
		userNotificationDataManager: userNotificationDataManager,
		tracer:                      tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
	}

	return svc, nil
}
