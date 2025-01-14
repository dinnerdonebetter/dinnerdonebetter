package users

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/uploads"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// Config configures the users service.
	Config struct {
		PublicMediaURLPrefix string         `env:"PUBLIC_MEDIA_URL_PREFIX" json:"publicMediaURLPrefix"`
		Uploads              uploads.Config `envPrefix:"UPLOADS_"          json:"uploads"`
	}
)

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates a Config struct.
func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		cfg,
		validation.Field(&cfg.Uploads, validation.Required),
	)
}
