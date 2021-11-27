package householdinvitations

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/prixfixeco/api_server/internal/observability/logging"
)

// Config represents our database configuration.
type Config struct {
	_                  struct{}
	Logging            *logging.Config `json:"logging,omitempty" mapstructure:"logging" toml:"logging,omitempty"`
	PreWritesTopicName string          `json:"writesTopicName,omitempty" mapstructure:"pre_writes_topic_name" toml:"pre_writes_topic_name,omitempty"`
	Debug              bool            `json:"debug" mapstructure:"debug" toml:"debug,omitempty"`
}

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates a Config struct.
func (cfg Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, &cfg)
}
