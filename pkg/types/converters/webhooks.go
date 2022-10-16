package converters

import (
	"github.com/prixfixeco/api_server/pkg/types"
)

// ConvertWebhookToWebhookCreationRequestInput builds a WebhookCreationRequestInput from a Webhook.
func ConvertWebhookToWebhookCreationRequestInput(webhook *types.Webhook) *types.WebhookCreationRequestInput {
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

// ConvertWebhookToWebhookDatabaseCreationInput builds a WebhookCreationRequestInput from a Webhook.
func ConvertWebhookToWebhookDatabaseCreationInput(webhook *types.Webhook) *types.WebhookDatabaseCreationInput {
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

// ConvertWebhookCreationRequestInputToWebhookDatabaseCreationInput creates a WebhookDatabaseCreationInput from a WebhookCreationRequestInput.
func ConvertWebhookCreationRequestInputToWebhookDatabaseCreationInput(input *types.WebhookCreationRequestInput) *types.WebhookDatabaseCreationInput {
	x := &types.WebhookDatabaseCreationInput{}

	x.Name = input.Name
	x.ContentType = input.ContentType
	x.URL = input.URL
	x.Method = input.Method
	x.Events = input.Events
	x.DataTypes = input.DataTypes
	x.Topics = input.Topics

	return x
}
