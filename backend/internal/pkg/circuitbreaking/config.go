package circuitbreaking

import (
	"context"
	"math"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Config struct {
	Name                   string  `env:"NAME"                     json:"name"`
	ErrorRate              float64 `env:"ERROR_RATE"               json:"circuitBreakerErrorPercentage"`
	MinimumSampleThreshold uint64  `env:"MINIMUM_SAMPLE_THRESHOLD" json:"circuitBreakerMinimumOccurrenceThreshold"`
}

func (cfg *Config) EnsureDefaults() {
	if cfg == nil {
		cfg = &Config{
			Name:                   "UNKNOWN",
			ErrorRate:              1.0,
			MinimumSampleThreshold: math.MaxUint64,
		}
	}

	if cfg.ErrorRate == 0 {
		cfg.ErrorRate = 200
	}

	if cfg.MinimumSampleThreshold == 0 {
		cfg.MinimumSampleThreshold = 1_000_000
	}
}

func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, cfg,
		validation.Field(&cfg.Name, validation.Required),
		validation.Field(&cfg.ErrorRate, validation.Min(0.01), validation.Max(0.99)),
		validation.Field(&cfg.MinimumSampleThreshold, validation.Min(0.01), validation.Max(0.99)),
	)
}
