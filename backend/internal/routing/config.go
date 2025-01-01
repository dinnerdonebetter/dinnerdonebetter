package routing

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// ChiProvider is the string we use to refer to chi.
	ChiProvider = "chi"
)

// Config configures our router.
type Config struct {
	_ struct{} `json:"-"`

	ServiceName            string   `env:"SERVICE_NAME"              json:"serviceName,omitempty"`
	Provider               string   `env:"PROVIDER"                  json:"provider,omitempty"`
	ValidDomains           []string `env:"VALID_DOMAINS"             json:"validDomains,omitempty"`
	EnableCORSForLocalhost bool     `env:"ENABLE_CORS_FOR_LOCALHOST" json:"enableCORSForLocalhost"`
	SilenceRouteLogging    bool     `env:"SILENCE_ROUTE_LOGGING"     json:"silenceRouteLogging,omitempty"`
}

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates a router config struct.
func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, cfg,
		validation.Field(&cfg.ServiceName, validation.Required),
		validation.Field(&cfg.Provider, validation.In(ChiProvider)),
	)
}
