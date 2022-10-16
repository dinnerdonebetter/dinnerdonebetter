package fakes

import (
	"net/http"

	fake "github.com/brianvoe/gofakeit/v5"

	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/converters"
)

// BuildFakeWebhook builds a faked Webhook.
func BuildFakeWebhook() *types.Webhook {
	return &types.Webhook{
		ID:                 BuildFakeID(),
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
	return converters.ConvertWebhookToWebhookCreationRequestInput(webhook)
}
