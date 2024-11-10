package config

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// QueuesConfig is primarily used for development.
type QueuesConfig struct {
	_ struct{} `json:"-"`

	DataChangesTopicName              string `json:"dataChangesTopicName"              toml:"dataChangesTopicName,omitempty"`
	OutboundEmailsTopicName           string `json:"outboundEmailsTopicName"           toml:"outboundEmailsTopicName,omitempty"`
	SearchIndexRequestsTopicName      string `json:"searchIndexRequestsTopicName"      toml:"searchIndexRequestsTopicName,omitempty"`
	UserDataAggregationTopicName      string `json:"userDataAggregationTopicName"      toml:"userDataAggregationTopicName,omitempty"`
	WebhookExecutionRequestsTopicName string `json:"webhookExecutionRequestsTopicName" toml:"webhookExecutionRequestsTopicName,omitempty"`
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
