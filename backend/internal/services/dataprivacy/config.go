package dataprivacy

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/uploads"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Config configures the service.
type Config struct {
	_ struct{} `json:"-"`

	UserDataAggregationTopicName string         `env:"USER_DATA_AGGREGATION_TOPIC_NAME" json:"userDataAggregationTopicName,omitempty"`
	DataChangesTopicName         string         `env:"DATA_CHANGES_TOPIC_NAME"          json:"dataChangesTopicName,omitempty"`
	Uploads                      uploads.Config `envPrefix:"UPLOADS_"                   json:"uploads"`
}

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates a Config struct.
func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		cfg,
		validation.Field(&cfg.Uploads, validation.Required),
		validation.Field(&cfg.DataChangesTopicName, validation.Required),
		validation.Field(&cfg.UserDataAggregationTopicName, validation.Required),
	)
}
