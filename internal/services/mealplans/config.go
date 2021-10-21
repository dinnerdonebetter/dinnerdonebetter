package mealplans

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
)

// Config configures the service.
type Config struct {
	_                    struct{}
	Logging              logging.Config `json:"logging" mapstructure:"logging" toml:"logging,omitempty"`
	PreWritesTopicName   string         `json:"preWritesTopicName" mapstructure:"pre_writes_topic_name" toml:"pre_writes_topic_name,omitempty"`
	PreUpdatesTopicName  string         `json:"preUpdatesTopicName" mapstructure:"pre_updates_topic_name" toml:"pre_updates_topic_name,omitempty"`
	PreArchivesTopicName string         `json:"preArchivesTopicName" mapstructure:"pre_archives_topic_name" toml:"pre_archives_topic_name,omitempty"`
}

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates a Config struct.
func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		cfg,
		validation.Field(&cfg.Logging, validation.Required),
		validation.Field(&cfg.PreWritesTopicName, validation.Required),
		validation.Field(&cfg.PreUpdatesTopicName, validation.Required),
		validation.Field(&cfg.PreArchivesTopicName, validation.Required),
	)
}
