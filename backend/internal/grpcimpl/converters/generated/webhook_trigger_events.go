package converters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func ConvertWebhookTriggerEventCreationRequestInputToWebhookTriggerEvent(input *messages.WebhookTriggerEventCreationRequestInput) *messages.WebhookTriggerEvent {

output := &messages.WebhookTriggerEvent{
    TriggerEvent: input.TriggerEvent,
    BelongsToWebhook: input.BelongsToWebhook,
}

return output
}
