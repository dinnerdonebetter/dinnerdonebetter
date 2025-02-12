package converters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func ConvertWebhookCreationRequestInputToWebhook(input *messages.WebhookCreationRequestInput) *messages.Webhook {
convertedevents := make([]*messages.WebhookTriggerEvent, 0, len(input.Events))
for _, item := range input.Events {
    convertedevents = append(convertedevents, ConvertstringToWebhookTriggerEvent(item))
}

output := &messages.Webhook{
    Name: input.Name,
    URL: input.URL,
    Method: input.Method,
    ContentType: input.ContentType,
    Events: convertedevents,
}

return output
}
