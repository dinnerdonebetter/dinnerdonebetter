package converters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func ConvertUserNotificationCreationRequestInputToUserNotification(input *messages.UserNotificationCreationRequestInput) *messages.UserNotification {

output := &messages.UserNotification{
    Content: input.Content,
    Status: input.Status,
    BelongsToUser: input.BelongsToUser,
}

return output
}
func ConvertUserNotificationUpdateRequestInputToUserNotification(input *messages.UserNotificationUpdateRequestInput) *messages.UserNotification {

output := &messages.UserNotification{
    Status: input.Status,
}

return output
}
