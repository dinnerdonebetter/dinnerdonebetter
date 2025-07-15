package grpc

import (
	"context"
	grpcconverters "github.com/dinnerdonebetter/backend/internal/grpc/converters"
	grpctypes "github.com/dinnerdonebetter/backend/internal/grpc/generated/types"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"google.golang.org/grpc/codes"

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

func ConvertUserNotificationToGRPCUserNotification(notification *notifications.UserNotification) *notificationssvc.UserNotification {
	return &notificationssvc.UserNotification{
		CreatedAt:     grpcconverters.ConvertTimeToPBTimestamp(notification.CreatedAt),
		LastUpdatedAt: grpcconverters.ConvertTimePointerToPBTimestamp(notification.LastUpdatedAt),
		ID:            notification.ID,
		Content:       notification.Content,
		Status:        notification.Status,
		BelongsToUser: notification.BelongsToUser,
	}
}

func (s *ServiceImpl) GetUserNotification(ctx context.Context, request *notificationssvc.GetUserNotificationRequest) (*notificationssvc.GetUserNotificationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	userID := "TODO"
	logger := s.logger.WithValue(keys.UserNotificationIDKey, request.UserNotificationID)

	notification, err := s.notificationsRepository.GetUserNotification(ctx, userID, request.UserNotificationID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "")
	}

	x := &notificationssvc.GetUserNotificationResponse{
		ResponseDetails: &grpctypes.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Result: ConvertUserNotificationToGRPCUserNotification(notification),
	}

	return x, nil
}

func (s *ServiceImpl) GetUserNotifications(ctx context.Context, request *notificationssvc.GetUserNotificationsRequest) (*notificationssvc.GetUserNotificationsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) UpdateUserNotification(ctx context.Context, request *notificationssvc.UpdateUserNotificationRequest) (*notificationssvc.UpdateUserNotificationResponse, error) {
	//TODO implement me
	panic("implement me")
}
