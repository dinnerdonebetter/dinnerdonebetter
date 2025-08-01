package converters

import (
	"github.com/dinnerdonebetter/backend/internal/domain/notifications"
	grpcconverters "github.com/dinnerdonebetter/backend/internal/grpc/converters"
	notifications2 "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/notifications"
)

func ConvertUserNotificationToGRPCUserNotification(notification *notifications.UserNotification) *notifications2.UserNotification {
	return &notifications2.UserNotification{
		CreatedAt:     grpcconverters.ConvertTimeToPBTimestamp(notification.CreatedAt),
		LastUpdatedAt: grpcconverters.ConvertTimePointerToPBTimestamp(notification.LastUpdatedAt),
		ID:            notification.ID,
		Content:       notification.Content,
		Status:        notification.Status,
		BelongsToUser: notification.BelongsToUser,
	}
}
