package grpc

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/authentication/sessions"
	notificationsmanager "github.com/dinnerdonebetter/backend/internal/domain/notifications/manager"
	notificationssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/notifications"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

const (
	o11yName = "notifications_service"
)

var _ notificationssvc.UserNotificationsServiceServer = (*serviceImpl)(nil)

type (
	serviceImpl struct {
		notificationssvc.UnimplementedUserNotificationsServiceServer
		tracer                    tracing.Tracer
		logger                    logging.Logger
		sessionContextDataFetcher func(context.Context) (*sessions.ContextData, error)
		notificationsManager      notificationsmanager.NotificationsDataManager
	}
)

func NewService(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	notificationsManager notificationsmanager.NotificationsDataManager,
) notificationssvc.UserNotificationsServiceServer {
	return &serviceImpl{
		logger:                    logging.EnsureLogger(logger).WithName(o11yName),
		tracer:                    tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		notificationsManager:      notificationsManager,
		sessionContextDataFetcher: sessions.FetchContextDataFromContext,
	}
}
