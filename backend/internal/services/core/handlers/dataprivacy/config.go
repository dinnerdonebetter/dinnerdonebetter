package dataprivacy

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/lib/uploads"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Config configures the service.
type Config struct {
	_ struct{} `json:"-"`

	Uploads uploads.Config `envPrefix:"UPLOADS_" json:"uploads"`
}

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates a Config struct.
func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		cfg,
		validation.Field(&cfg.Uploads, validation.Required),
	)
}
