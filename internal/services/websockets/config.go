package websockets

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/prixfixeco/api_server/internal/observability/logging"
)

// Config configures the service.
type Config struct {
	_ struct{}

	Logging              *logging.Config `json:"logging,omitempty" mapstructure:"logging" toml:"logging,omitempty"`
	DataChangesTopicName string          `json:"dataChangesTopicName,omitempty" mapstructure:"pre_data_changes_topic_name" toml:"pre_data_changes_topic_name,omitempty"`
}

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates a Config struct.
func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		cfg,
	)
}
