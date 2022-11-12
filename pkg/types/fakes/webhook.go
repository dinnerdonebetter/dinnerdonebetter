package fakes

import (
	"net/http"

	fake "github.com/brianvoe/gofakeit/v5"

	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/converters"
)

// BuildFakeWebhook builds a faked Webhook.
func BuildFakeWebhook() *types.Webhook {
	webhookID := BuildFakeID()

	fakeEvent := BuildFakeWebhookTriggerEvent()
	fakeEvent.BelongsToWebhook = webhookID
	events := []*types.WebhookTriggerEvent{fakeEvent}

	return &types.Webhook{
		ID:                 webhookID,
		Name:               fake.UUID(),
		ContentType:        "application/json",
		URL:                fake.URL(),
		Method:             http.MethodPost,
		Events:             events,
		CreatedAt:          BuildFakeTime(),
		ArchivedAt:         nil,
		BelongsToHousehold: fake.UUID(),
	}
}

// BuildFakeWebhookTriggerEvent builds a faked WebhookTriggerEvent.
func BuildFakeWebhookTriggerEvent() *types.WebhookTriggerEvent {
	return &types.WebhookTriggerEvent{
		ID:               BuildFakeID(),
		TriggerEvent:     string(types.WebhookCreatedCustomerEventType),
		BelongsToWebhook: BuildFakeID(),
		CreatedAt:        BuildFakeTime(),
		ArchivedAt:       nil,
	}
}

// BuildFakeWebhookList builds a faked WebhookList.
func BuildFakeWebhookList() *types.WebhookList {
	var examples []*types.Webhook
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeWebhook())
	}

	return &types.WebhookList{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Webhooks: examples,
	}
}

// BuildFakeWebhookCreationRequestInput builds a faked WebhookCreationRequestInput from a webhook.
func BuildFakeWebhookCreationRequestInput() *types.WebhookCreationRequestInput {
	webhook := BuildFakeWebhook()
	return converters.ConvertWebhookToWebhookCreationRequestInput(webhook)
}
