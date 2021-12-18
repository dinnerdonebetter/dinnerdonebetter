package validinstruments

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Config configures the service.
type Config struct {
	_ struct{}

	DataChangesTopicName string `json:"dataChangesTopicName,omitempty" mapstructure:"data_changes_topic_name" toml:"data_changes_topic_name,omitempty"`
	SearchIndexPath      string `json:"searchIndexPath,omitempty" mapstructure:"search_index_path" toml:"search_index_path,omitempty"`
}

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates a Config struct.
func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		cfg,
		validation.Field(&cfg.DataChangesTopicName, validation.Required),
		validation.Field(&cfg.SearchIndexPath, validation.Required),
	)
}
