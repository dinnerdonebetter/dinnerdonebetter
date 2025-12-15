package converters

import (
	"github.com/dinnerdonebetter/backend/internal/domain/webhooks"
	grpcconverters "github.com/dinnerdonebetter/backend/internal/grpc/converters"
	webhookssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/webhooks"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
)

func ConvertStringToWebhookContentType(s string) webhookssvc.WebhookContentType {
	value, ok := webhookssvc.WebhookContentType_value[s]
	if !ok {
		return webhookssvc.WebhookContentType_WEBHOOK_CONTENT_TYPE_JSON
	}
	return webhookssvc.WebhookContentType(value)
}

func ConvertStringToWebhookMethod(s string) webhookssvc.WebhookMethod {
	value, ok := webhookssvc.WebhookMethod_value[s]
	if !ok {
		return webhookssvc.WebhookMethod_WEBHOOK_METHOD_POST
	}
	return webhookssvc.WebhookMethod(value)
}

func ConvertWebhookToGRPCWebhook(webhook *webhooks.Webhook) *webhookssvc.Webhook {
	converted := &webhookssvc.Webhook{
		CreatedAt:        grpcconverters.ConvertTimeToPBTimestamp(webhook.CreatedAt),
		ArchivedAt:       grpcconverters.ConvertTimePointerToPBTimestamp(webhook.ArchivedAt),
		LastUpdatedAt:    grpcconverters.ConvertTimePointerToPBTimestamp(webhook.LastUpdatedAt),
		Name:             webhook.Name,
		Url:              webhook.URL,
		Method:           ConvertStringToWebhookMethod(webhook.Method),
		Id:               webhook.ID,
		BelongsToAccount: webhook.BelongsToAccount,
		ContentType:      ConvertStringToWebhookContentType(webhook.ContentType),
	}

	for _, event := range webhook.Events {
		converted.Events = append(converted.Events, ConvertWebhookTriggerEventToGRPCWebhookTriggerEvent(event))
	}

	return converted
}

func ConvertWebhookTriggerEventToGRPCWebhookTriggerEvent(z *webhooks.WebhookTriggerEvent) *webhookssvc.WebhookTriggerEvent {
	return &webhookssvc.WebhookTriggerEvent{
		CreatedAt:        grpcconverters.ConvertTimeToPBTimestamp(z.CreatedAt),
		ArchivedAt:       grpcconverters.ConvertTimePointerToPBTimestamp(z.ArchivedAt),
		Id:               z.ID,
		BelongsToWebhook: z.BelongsToWebhook,
		TriggerEvent:     z.TriggerEvent,
	}
}

func ConvertGRPCWebhookToWebhook(webhook *webhookssvc.Webhook) *webhooks.Webhook {
	converted := &webhooks.Webhook{
		CreatedAt:        grpcconverters.ConvertPBTimestampToTime(webhook.CreatedAt),
		ArchivedAt:       grpcconverters.ConvertPBTimestampToTimePointer(webhook.ArchivedAt),
		LastUpdatedAt:    grpcconverters.ConvertPBTimestampToTimePointer(webhook.LastUpdatedAt),
		Name:             webhook.Name,
		URL:              webhook.Url,
		Method:           webhook.Method.String(),
		ID:               webhook.Id,
		BelongsToAccount: webhook.BelongsToAccount,
		ContentType:      webhook.ContentType.String(),
	}

	for _, event := range webhook.Events {
		converted.Events = append(converted.Events, ConvertGRPCWebhookTriggerEventToWebhookTriggerEvent(event))
	}

	return converted
}

func ConvertGRPCWebhookTriggerEventToWebhookTriggerEvent(z *webhookssvc.WebhookTriggerEvent) *webhooks.WebhookTriggerEvent {
	return &webhooks.WebhookTriggerEvent{
		CreatedAt:        grpcconverters.ConvertPBTimestampToTime(z.CreatedAt),
		ArchivedAt:       grpcconverters.ConvertPBTimestampToTimePointer(z.ArchivedAt),
		ID:               z.Id,
		BelongsToWebhook: z.BelongsToWebhook,
		TriggerEvent:     z.TriggerEvent,
	}
}

func ConvertGRPCWebhookCreationRequestInputToWebhookDatabaseCreationInput(input *webhookssvc.WebhookCreationRequestInput, accountID string) *webhooks.WebhookDatabaseCreationInput {
	webhookID := identifiers.New()

	var events []*webhooks.WebhookTriggerEventDatabaseCreationInput
	for _, event := range input.Events {
		events = append(events, &webhooks.WebhookTriggerEventDatabaseCreationInput{
			ID:               identifiers.New(),
			BelongsToWebhook: webhookID,
			TriggerEvent:     event,
		})
	}

	x := &webhooks.WebhookDatabaseCreationInput{
		ID:               webhookID,
		Name:             input.Name,
		ContentType:      input.ContentType.String(),
		URL:              input.Url,
		Method:           input.Method.String(),
		BelongsToAccount: accountID,
		Events:           events,
	}

	return x
}

func ConvertWebhookCreationRequestInputToGRPCWebhookCreationRequestInput(input *webhooks.WebhookCreationRequestInput) *webhookssvc.WebhookCreationRequestInput {
	return &webhookssvc.WebhookCreationRequestInput{
		Name:        input.Name,
		ContentType: ConvertStringToWebhookContentType(input.ContentType),
		Url:         input.URL,
		Method:      ConvertStringToWebhookMethod(input.Method),
		Events:      input.Events,
	}
}

func ConvertGRPCWebhookTriggerEventDatabaseCreationInputToWebhookTriggerEventDatabaseCreationInput(input *webhookssvc.WebhookTriggerEventCreationRequestInput) *webhooks.WebhookTriggerEventDatabaseCreationInput {
	return &webhooks.WebhookTriggerEventDatabaseCreationInput{
		ID:               identifiers.New(),
		BelongsToWebhook: input.BelongsToWebhook,
		TriggerEvent:     input.TriggerEvent,
	}
}
