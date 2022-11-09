package pubsub

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Config configures a PubSub-backed consumer.
type Config struct {
	TopicName string `json:"topicName" mapstructure:"topic_name" toml:"topic_name,omitempty"`
}

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates a Config struct.
func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, cfg,
		validation.Field(&cfg.TopicName, validation.Required),
	)
}
