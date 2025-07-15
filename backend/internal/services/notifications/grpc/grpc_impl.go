package grpc

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/domain/notifications"
	notificationssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/notifications"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

var _ notificationssvc.UserNotificationsServiceServer = (*ServiceImpl)(nil)

type (
	ServiceImpl struct {
		notificationssvc.UnimplementedUserNotificationsServiceServer
		tracer                  tracing.Tracer
		logger                  logging.Logger
		notificationsRepository notifications.Repository
	}
)

func NewService(
	notificationsRepository notifications.Repository,
) notificationssvc.UserNotificationsServiceServer {
	return &ServiceImpl{
		notificationsRepository: notificationsRepository,
	}
}

func (s *ServiceImpl) GetUserNotification(ctx context.Context, request *notificationssvc.GetUserNotificationRequest) (*notificationssvc.GetUserNotificationResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) GetUserNotifications(ctx context.Context, request *notificationssvc.GetUserNotificationsRequest) (*notificationssvc.GetUserNotificationsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) UpdateUserNotification(ctx context.Context, request *notificationssvc.UpdateUserNotificationRequest) (*notificationssvc.UpdateUserNotificationResponse, error) {
	//TODO implement me
	panic("implement me")
}
