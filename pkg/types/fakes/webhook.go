package fakes

import (
	"net/http"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/segmentio/ksuid"

	"github.com/prixfixeco/api_server/pkg/types"
)

// BuildFakeWebhook builds a faked Webhook.
func BuildFakeWebhook() *types.Webhook {
	return &types.Webhook{
		ID:                 ksuid.New().String(),
		Name:               fake.UUID(),
		ContentType:        "application/json",
		URL:                fake.URL(),
		Method:             http.MethodPost,
		Events:             []string{buildUniqueString()},
		DataTypes:          []string{buildUniqueString()},
		Topics:             []string{buildUniqueString()},
		CreatedAt:          fake.Date(),
		ArchivedAt:         nil,
		BelongsToHousehold: fake.UUID(),
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
	return BuildFakeWebhookCreationInputFromWebhook(webhook)
}

// BuildFakeWebhookDatabaseCreationInput builds a faked WebhookCreationRequestInput from a webhook.
func BuildFakeWebhookDatabaseCreationInput() *types.WebhookDatabaseCreationInput {
	webhook := BuildFakeWebhook()
	return BuildFakeWebhookDatabaseCreationInputFromWebhook(webhook)
}

// BuildFakeWebhookCreationInputFromWebhook builds a faked WebhookCreationRequestInput.
func BuildFakeWebhookCreationInputFromWebhook(webhook *types.Webhook) *types.WebhookCreationRequestInput {
	return &types.WebhookCreationRequestInput{
		ID:                 webhook.ID,
		Name:               webhook.Name,
		ContentType:        webhook.ContentType,
		URL:                webhook.URL,
		Method:             webhook.Method,
		Events:             webhook.Events,
		DataTypes:          webhook.DataTypes,
		Topics:             webhook.Topics,
		BelongsToHousehold: webhook.BelongsToHousehold,
	}
}

// BuildFakeWebhookDatabaseCreationInputFromWebhook builds a faked WebhookCreationRequestInput.
func BuildFakeWebhookDatabaseCreationInputFromWebhook(webhook *types.Webhook) *types.WebhookDatabaseCreationInput {
	return &types.WebhookDatabaseCreationInput{
		ID:                 webhook.ID,
		Name:               webhook.Name,
		ContentType:        webhook.ContentType,
		URL:                webhook.URL,
		Method:             webhook.Method,
		Events:             webhook.Events,
		DataTypes:          webhook.DataTypes,
		Topics:             webhook.Topics,
		BelongsToHousehold: webhook.BelongsToHousehold,
	}
}
