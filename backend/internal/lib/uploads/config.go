package uploads

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/lib/uploads/objectstorage"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Config contains settings regarding search indices.
type Config struct {
	_ struct{} `json:"-"`

	Storage objectstorage.Config `envPrefix:"STORAGE_" json:"storageConfig"`
	Debug   bool                 `env:"DEBUG"          json:"debug"`
}

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates a Config struct.
func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, cfg,
		validation.Field(&cfg.Storage),
	)
}
