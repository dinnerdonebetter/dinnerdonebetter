package cloudtrace

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// Config contains settings related to tracing.
	Config struct {
		_ struct{} `json:"-"`

		ProjectID string `json:"projectID,omitempty" toml:"project_id,omitempty" env:"GOOGLE_CLOUD_TRACE_PROJECT_ID"`
	}
)

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates the config struct.
func (c *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, c,
		validation.Field(&c.ProjectID, validation.Required),
	)
}
