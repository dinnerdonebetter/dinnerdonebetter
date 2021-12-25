package households

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Config configures the service.
type Config struct {
	_ struct{}

	PreWritesTopicName string `json:"writesTopicName,omitempty" mapstructure:"pre_writes_topic_name" toml:"pre_writes_topic_name,omitempty"`
}

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates a Config struct.
func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		cfg,
		validation.Field(&cfg.PreWritesTopicName, validation.Required),
	)
}
