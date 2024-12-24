package http

import (
	"context"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// Config describes the settings pertinent to the HTTP serving portion of the service.
	Config struct {
		_ struct{} `json:"-"`

		HTTPSCertificateFile    string        `env:"HTTPS_CERTIFICATE_FILEPATH"     json:"httpsCertificate,omitempty"`
		HTTPSCertificateKeyFile string        `env:"HTTPS_CERTIFICATE_KEY_FILEPATH" json:"httpsCertificateKey,omitempty"`
		StartupDeadline         time.Duration `env:"STARTUP_DEADLINE"               json:"startupDeadline,omitempty"`
		HTTPPort                uint16        `env:"HTTP_PORT"                      json:"httpPort"`
		Debug                   bool          `env:"DEBUG"                          json:"debug"`
	}
)

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates a Config struct.
func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, cfg,
		validation.Field(&cfg.HTTPPort, validation.Required),
		validation.Field(&cfg.StartupDeadline, validation.Required),
	)
}
