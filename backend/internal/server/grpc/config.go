package grpc

import (
	"context"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// Config describes the settings pertinent to the HTTP serving portion of the service.
	Config struct {
		_ struct{} `json:"-"`

		HTTPSCertificateFile    string        `json:"httpsCertificate,omitempty"    toml:"https_certificate,omitempty"`
		HTTPSCertificateKeyFile string        `json:"httpsCertificateKey,omitempty" toml:"https_certificate_key,omitempty"`
		StartupDeadline         time.Duration `json:"startupDeadline,omitempty"     toml:"startup_deadline,omitempty"`
		Port                    uint16        `json:"port"                          toml:"port,omitempty"`
		Debug                   bool          `json:"debug"                         toml:"debug,omitempty"`
	}
)

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates a Config struct.
func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, cfg,
		validation.Field(&cfg.Port, validation.Required),
		validation.Field(&cfg.StartupDeadline, validation.Required),
	)
}
