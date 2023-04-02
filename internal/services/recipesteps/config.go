package recipesteps

import (
	"context"

	"github.com/prixfixeco/backend/internal/uploads"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Config configures the service.
type Config struct {
	_ struct{}

	PublicMediaURLPrefix string         `json:"mediaUploadPrefix" mapstructure:"media_upload_prefix" toml:"media_upload_prefix"`
	DataChangesTopicName string         `json:"dataChangesTopicName,omitempty" mapstructure:"data_changes_topic_name" toml:"data_changes_topic_name,omitempty"`
	Uploads              uploads.Config `json:"uploads" mapstructure:"uploads" toml:"uploads,omitempty"`
}

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates a Config struct.
func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		cfg,
		validation.Field(&cfg.PublicMediaURLPrefix, validation.Required),
		validation.Field(&cfg.Uploads, validation.Required),
		validation.Field(&cfg.DataChangesTopicName, validation.Required),
	)
}
