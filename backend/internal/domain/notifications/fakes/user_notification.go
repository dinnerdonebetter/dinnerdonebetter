package fakes

import (
	types "github.com/dinnerdonebetter/backend/internal/domain/notifications"
	"github.com/dinnerdonebetter/backend/internal/domain/notifications/converters"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
)

// BuildFakeUserNotification builds a faked valid ingredient.
func BuildFakeUserNotification() *types.UserNotification {
	return &types.UserNotification{
		CreatedAt:     BuildFakeTime(),
		ID:            BuildFakeID(),
		Content:       buildUniqueString(),
		Status:        types.UserNotificationStatusTypeUnread,
		BelongsToUser: BuildFakeID(),
	}
}

// BuildFakeUserNotificationsList builds a faked UserNotificationList.
func BuildFakeUserNotificationsList() *filtering.QueryFilteredResult[types.UserNotification] {
	var notifications []*types.UserNotification
	for i := 0; i < exampleQuantity; i++ {
		notifications = append(notifications, BuildFakeUserNotification())
	}

	return &filtering.QueryFilteredResult[types.UserNotification]{
		Pagination: filtering.Pagination{
			Page:          1,
			Limit:         50,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: notifications,
	}
}

// BuildFakeUserNotificationUpdateRequestInput builds a faked UserNotificationUpdateRequestInput.
func BuildFakeUserNotificationUpdateRequestInput() *types.UserNotificationUpdateRequestInput {
	userNotification := BuildFakeUserNotification()
	return converters.ConvertUserNotificationToUserNotificationUpdateRequestInput(userNotification)
}

// BuildFakeUserNotificationCreationRequestInput builds a faked UserNotificationCreationRequestInput.
func BuildFakeUserNotificationCreationRequestInput() *types.UserNotificationCreationRequestInput {
	userNotification := BuildFakeUserNotification()
	return converters.ConvertUserNotificationToUserNotificationCreationRequestInput(userNotification)
}
