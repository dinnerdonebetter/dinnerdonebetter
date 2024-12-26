package config

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// QueuesConfig is primarily used for development.
type QueuesConfig struct {
	_ struct{} `json:"-"`

	DataChangesTopicName              string `env:"DATA_CHANGES_TOPIC_NAME"               json:"dataChangesTopicName"`
	OutboundEmailsTopicName           string `env:"OUTBOUND_EMAILS_TOPIC_NAME"            json:"outboundEmailsTopicName"`
	SearchIndexRequestsTopicName      string `env:"SEARCH_INDEX_REQUESTS_TOPIC_NAME"      json:"searchIndexRequestsTopicName"`
	UserDataAggregationTopicName      string `env:"USER_DATA_AGGREGATION_TOPIC_NAME"      json:"userDataAggregationTopicName"`
	WebhookExecutionRequestsTopicName string `env:"WEBHOOK_EXECUTION_REQUESTS_TOPIC_NAME" json:"webhookExecutionRequestsTopicName"`
}

var _ validation.ValidatableWithContext = (*QueuesConfig)(nil)

// ValidateWithContext validates a QueuesConfig struct.
func (c *QueuesConfig) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, c,
		validation.Field(&c.DataChangesTopicName, validation.Required),
		validation.Field(&c.OutboundEmailsTopicName, validation.Required),
		validation.Field(&c.SearchIndexRequestsTopicName, validation.Required),
		validation.Field(&c.UserDataAggregationTopicName, validation.Required),
		validation.Field(&c.WebhookExecutionRequestsTopicName, validation.Required),
	)
}
