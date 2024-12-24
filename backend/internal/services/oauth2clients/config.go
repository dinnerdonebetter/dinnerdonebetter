package oauth2clients

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Config manages our body validation.
type Config struct {
	DataChangesTopicName string `env:"DATA_CHANGES_TOPIC_NAME" json:"dataChangesTopicName,omitempty"`
	CreationEnabled      bool   `env:"CREATION_ENABLED"        json:"creationEnabled"`
}

func (c *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, c,
		validation.Field(&c.DataChangesTopicName, validation.Required),
	)
}
