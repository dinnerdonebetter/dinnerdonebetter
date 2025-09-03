package converters

import (
	"github.com/dinnerdonebetter/backend/internal/domain/notifications"
	grpcconverters "github.com/dinnerdonebetter/backend/internal/grpc/converters"
	notificationssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/notifications"
)

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

func ConvertGRPCUserNotificationToUserNotification(notification *notificationssvc.UserNotification) *notifications.UserNotification {
	return &notifications.UserNotification{
		CreatedAt:     grpcconverters.ConvertPBTimestampToTime(notification.CreatedAt),
		LastUpdatedAt: grpcconverters.ConvertPBTimestampToTimePointer(notification.LastUpdatedAt),
		ID:            notification.ID,
		Content:       notification.Content,
		Status:        notification.Status,
		BelongsToUser: notification.BelongsToUser,
	}
}

func ConvertUserNotificationUpdateRequestInputToGRPCUserNotificationUpdateRequestInput(input *notifications.UserNotificationUpdateRequestInput) *notificationssvc.UserNotificationUpdateRequestInput {
	return &notificationssvc.UserNotificationUpdateRequestInput{
		Status: input.Status,
	}
}
