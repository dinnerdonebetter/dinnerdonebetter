package grpc

import (
	"github.com/dinnerdonebetter/backend/internal/domain/webhooks"
	grpcconverters "github.com/dinnerdonebetter/backend/internal/grpc/converters"
	configurationsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/configuration"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
)

func ConvertWebhookToGRPCWebhook(webhook *webhooks.Webhook) *configurationsvc.Webhook {
	converted := &configurationsvc.Webhook{
		CreatedAt:        grpcconverters.ConvertTimeToPBTimestamp(webhook.CreatedAt),
		ArchivedAt:       grpcconverters.ConvertTimePointerToPBTimestamp(webhook.ArchivedAt),
		LastUpdatedAt:    grpcconverters.ConvertTimePointerToPBTimestamp(webhook.LastUpdatedAt),
		Name:             webhook.Name,
		URL:              webhook.URL,
		Method:           webhook.Method,
		ID:               webhook.ID,
		BelongsToAccount: webhook.BelongsToAccount,
		ContentType:      webhook.ContentType,
	}

	for _, event := range webhook.Events {
		converted.Events = append(converted.Events, ConvertWebhookTriggerEventToGRPCWebhookTriggerEvent(event))
	}

	return converted
}

func ConvertWebhookTriggerEventToGRPCWebhookTriggerEvent(z *webhooks.WebhookTriggerEvent) *configurationsvc.WebhookTriggerEvent {
	return &configurationsvc.WebhookTriggerEvent{
		CreatedAt:        grpcconverters.ConvertTimeToPBTimestamp(z.CreatedAt),
		ArchivedAt:       grpcconverters.ConvertTimePointerToPBTimestamp(z.ArchivedAt),
		ID:               z.ID,
		BelongsToWebhook: z.BelongsToWebhook,
		TriggerEvent:     z.TriggerEvent,
	}
}

func ConvertGRPCWebhookCreationRequestInputToWebhookDatabaseCreationInput(input *configurationsvc.WebhookCreationRequestInput, accountID string) *webhooks.WebhookDatabaseCreationInput {
	x := &webhooks.WebhookDatabaseCreationInput{
		ID:               identifiers.New(),
		Name:             input.Name,
		ContentType:      input.ContentType,
		URL:              input.URL,
		Method:           input.Method,
		BelongsToAccount: accountID,
	}

	for _, event := range input.Events {
		x.Events = append(x.Events, &webhooks.WebhookTriggerEventDatabaseCreationInput{
			ID:               identifiers.New(),
			BelongsToWebhook: x.ID,
			TriggerEvent:     event,
		})
	}

	return x
}
