package cloudtrace

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// Config contains settings related to tracing.
	Config struct {
		_ struct{}

		ProjectID                 string  `json:"projectID,omitempty" mapstructure:"project_id" toml:"project_id,omitempty"`
		ServiceName               string  `json:"service_name,omitempty" mapstructure:"service_name" toml:"service_name,omitempty"`
		SpanCollectionProbability float64 `json:"spanCollectionProbability,omitempty" mapstructure:"span_collection_probability" toml:"span_collection_probability,omitempty"`
	}
)

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates the config struct.
func (c *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, c,
		validation.Field(&c.ProjectID, validation.Required),
		validation.Field(&c.ServiceName, validation.Required),
	)
}
