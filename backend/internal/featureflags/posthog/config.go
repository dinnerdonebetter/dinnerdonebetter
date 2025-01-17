package posthog

import (
	"github.com/dinnerdonebetter/backend/internal/circuitbreaking"
)

type (
	Config struct {
		ProjectAPIKey        string                 `env:"PROJECT_API_KEY"         json:"projectAPIKey"`
		PersonalAPIKey       string                 `env:"PERSONAL_API_KEY"        json:"personalAPIKey"`
		CircuitBreakerConfig circuitbreaking.Config `envPrefix:"CIRCUIT_BREAKING_" json:"circuitBreakerConfig"`
	}
)
