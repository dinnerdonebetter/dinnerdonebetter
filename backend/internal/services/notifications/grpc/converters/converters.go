package converters

import (
	"github.com/dinnerdonebetter/backend/internal/domain/notifications"
	grpcconverters "github.com/dinnerdonebetter/backend/internal/grpc/converters"
	notificationssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/notifications"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
)

func ConvertStringToUserNotificationStatus(s string) notificationssvc.UserNotificationStatus {
	value, ok := notificationssvc.UserNotificationStatus_value[s]
	if !ok {
		return notificationssvc.UserNotificationStatus_USER_NOTIFICATION_STATUS_UNREAD
	}
	return notificationssvc.UserNotificationStatus(value)
}

func ConvertUserNotificationToGRPCUserNotification(notification *notifications.UserNotification) *notificationssvc.UserNotification {
	return &notificationssvc.UserNotification{
		CreatedAt:     grpcconverters.ConvertTimeToPBTimestamp(notification.CreatedAt),
		LastUpdatedAt: grpcconverters.ConvertTimePointerToPBTimestamp(notification.LastUpdatedAt),
		Id:            notification.ID,
		Content:       notification.Content,
		Status:        ConvertStringToUserNotificationStatus(notification.Status),
		BelongsToUser: notification.BelongsToUser,
	}
}

func ConvertGRPCUserNotificationToUserNotification(notification *notificationssvc.UserNotification) *notifications.UserNotification {
	return &notifications.UserNotification{
		CreatedAt:     grpcconverters.ConvertPBTimestampToTime(notification.CreatedAt),
		LastUpdatedAt: grpcconverters.ConvertPBTimestampToTimePointer(notification.LastUpdatedAt),
		ID:            notification.Id,
		Content:       notification.Content,
		Status:        notification.Status.String(),
		BelongsToUser: notification.BelongsToUser,
	}
}

func ConvertUserNotificationUpdateRequestInputToGRPCUserNotificationUpdateRequestInput(input *notifications.UserNotificationUpdateRequestInput) *notificationssvc.UserNotificationUpdateRequestInput {
	var newStatus *notificationssvc.UserNotificationStatus
	if input.Status != nil {
		newStatus = pointer.To(ConvertStringToUserNotificationStatus(*input.Status))
	}

	return &notificationssvc.UserNotificationUpdateRequestInput{
		Status: newStatus,
	}
}
