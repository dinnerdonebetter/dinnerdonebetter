package http

import (
	"context"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// Config describes the settings pertinent to the HTTP serving portion of the service.
	Config struct {
		_ struct{}

		HTTPSCertificateFile    string        `json:"httpsCertificate,omitempty"    mapstructure:"https_certificate"     toml:"https_certificate,omitempty"`
		HTTPSCertificateKeyFile string        `json:"httpsCertificateKey,omitempty" mapstructure:"https_certificate_key" toml:"https_certificate_key,omitempty"`
		StartupDeadline         time.Duration `json:"startupDeadline,omitempty"     mapstructure:"startup_deadline"      toml:"startup_deadline,omitempty"`
		HTTPPort                uint16        `json:"httpPort"                      mapstructure:"http_port"             toml:"http_port,omitempty"`
		Debug                   bool          `json:"debug"                         mapstructure:"debug"                 toml:"debug,omitempty"`
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
