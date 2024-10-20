package usernotifications

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/messagequeue"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/routing"
	authservice "github.com/dinnerdonebetter/backend/internal/services/authentication"
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
		sessionContextDataFetcher   func(*http.Request) (*types.SessionContextData, error)
		userNotificationIDFetcher   func(*http.Request) string
	}
)

// ProvideService builds a new UserNotificationsService.
func ProvideService(
	_ context.Context,
	logger logging.Logger,
	cfg *Config,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
	userNotificationDataManager types.UserNotificationDataManager,
) (types.UserNotificationDataService, error) {
	if cfg == nil {
		return nil, errors.New("nil config provided to user notifications service")
	}

	dataChangesPublisher, err := publisherProvider.ProvidePublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up %s data changes publisher: %w", serviceName, err)
	}

	svc := &service{
		logger:                      logging.EnsureLogger(logger).WithName(serviceName),
		userNotificationIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(UserNotificationIDURIParamKey),
		sessionContextDataFetcher:   authservice.FetchContextFromRequest,
		dataChangesPublisher:        dataChangesPublisher,
		encoderDecoder:              encoder,
		userNotificationDataManager: userNotificationDataManager,
		tracer:                      tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
	}

	return svc, nil
}
