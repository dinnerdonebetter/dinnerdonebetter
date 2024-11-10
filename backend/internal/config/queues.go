package config

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// QueueSettings is primarily used for development.
type QueueSettings struct {
	_ struct{} `json:"-"`

	DataChangesTopicName              string `json:"dataChangesTopicName"              toml:"dataChangesTopicName,omitempty"`
	OutboundEmailsTopicName           string `json:"outboundEmailsTopicName"           toml:"outboundEmailsTopicName,omitempty"`
	SearchIndexRequestsTopicName      string `json:"searchIndexRequestsTopicName"      toml:"searchIndexRequestsTopicName,omitempty"`
	UserDataAggregationTopicName      string `json:"userDataAggregationTopicName"      toml:"userDataAggregationTopicName,omitempty"`
	WebhookExecutionRequestsTopicName string `json:"webhookExecutionRequestsTopicName" toml:"webhookExecutionRequestsTopicName,omitempty"`
}

var _ validation.ValidatableWithContext = (*QueueSettings)(nil)

// ValidateWithContext validates an MetaSettings struct.
func (s QueueSettings) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, &s,
		validation.Field(&s.DataChangesTopicName, validation.Required),
		validation.Field(&s.OutboundEmailsTopicName, validation.Required),
		validation.Field(&s.SearchIndexRequestsTopicName, validation.Required),
		validation.Field(&s.UserDataAggregationTopicName, validation.Required),
		validation.Field(&s.WebhookExecutionRequestsTopicName, validation.Required),
	)
}
