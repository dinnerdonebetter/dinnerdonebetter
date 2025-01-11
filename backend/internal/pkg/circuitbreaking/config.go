package circuitbreaking

type Config struct {
	ErrorRate              float64 `env:"ERROR_RATE"               json:"circuitBreakerErrorPercentage"`
	MinimumSampleThreshold int64   `env:"MINIMUM_SAMPLE_THRESHOLD" json:"circuitBreakerMinimumOccurrenceThreshold"`
}

func (cfg *Config) EnsureDefaults() {
	if cfg == nil {
		cfg = &Config{}
	}

	if cfg.ErrorRate == 0 {
		cfg.ErrorRate = 200
	}

	if cfg.MinimumSampleThreshold == 0 {
		cfg.MinimumSampleThreshold = 1_000_000
	}
}
