package validinstruments

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Config configures the service.
type Config struct {
	_ struct{}

	DataChangesTopicName string `json:"dataChangesTopicName,omitempty" mapstructure:"data_changes_topic_name" toml:"data_changes_topic_name,omitempty"`
	PreWritesTopicName   string `json:"writesTopicName,omitempty" mapstructure:"pre_writes_topic_name" toml:"pre_writes_topic_name,omitempty"`
	PreUpdatesTopicName  string `json:"updatesTopicName,omitempty" mapstructure:"pre_updates_topic_name" toml:"pre_updates_topic_name,omitempty"`
	PreArchivesTopicName string `json:"archivesTopicName,omitempty" mapstructure:"pre_archives_topic_name" toml:"pre_archives_topic_name,omitempty"`
}

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates a Config struct.
func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		cfg,
		validation.Field(&cfg.PreWritesTopicName, validation.Required),
		validation.Field(&cfg.PreUpdatesTopicName, validation.Required),
		validation.Field(&cfg.PreArchivesTopicName, validation.Required),
	)
}
