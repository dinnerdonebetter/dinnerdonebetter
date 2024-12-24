package householdinvitations

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Config represents our database configuration.
type Config struct {
	_ struct{} `json:"-"`

	DataChangesTopicName string `env:"DATA_CHANGES_TOPIC_NAME" json:"dataChangesTopicName,omitempty"`
	Debug                bool   `env:"DEBUG"                   json:"debug"`
}

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates a Config struct.
func (cfg Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, &cfg,
		validation.Field(&cfg.DataChangesTopicName, validation.Required),
	)
}
