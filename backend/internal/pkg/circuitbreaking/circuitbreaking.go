package circuitbreaking

import (
	"errors"
	"log"

	circuit "github.com/rubyist/circuitbreaker"
)

var (
	ErrNilConfig = errors.New("nil circuitbreaker config provided")
)

type CircuitBreaker interface {
	Failed()
	Succeeded()
	CanProceed() bool
	CannotProceed() bool
}

type BaseCircuitBreaker struct {
	circuitBreaker *circuit.Breaker
}

type Config struct {
	ErrorRate              float64 `json:"circuitBreakerErrorPercentage"            toml:"circuitBreaker_error_percentage"             env:"ERROR_RATE"`
	MinimumSampleThreshold int64   `json:"circuitBreakerMinimumOccurrenceThreshold" toml:"circuitBreaker_minimum_occurrence_threshold" env:"MINIMUM_SAMPLE_THRESHOLD"`
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

func ProvideCircuitBreaker(cfg *Config) CircuitBreaker {
	if cfg == nil {
		cfg = &Config{}
	}
	cfg.EnsureDefaults()

	return &BaseCircuitBreaker{
		circuitBreaker: circuit.NewBreakerWithOptions(&circuit.Options{
			ShouldTrip: func(cb *circuit.Breaker) bool {
				samples := cb.Failures() + cb.Successes()
				result := samples >= cfg.MinimumSampleThreshold && cb.ErrorRate() >= cfg.ErrorRate
				return result
			},
			WindowTime:    circuit.DefaultWindowTime,
			WindowBuckets: circuit.DefaultWindowBuckets,
		}),
	}
}

func (b *BaseCircuitBreaker) Failed() {
	b.circuitBreaker.Fail()
}

func (b *BaseCircuitBreaker) Succeeded() {
	b.circuitBreaker.Success()
}

func (b *BaseCircuitBreaker) CanProceed() bool {
	return b.circuitBreaker.Ready()
}

func (b *BaseCircuitBreaker) CannotProceed() bool {
	return !b.circuitBreaker.Ready()
}

// EnsureCircuitBreaker ensures a valid CircuitBreaker is made available.
func EnsureCircuitBreaker(breaker CircuitBreaker) CircuitBreaker {
	if breaker == nil {
		log.Println("NOOP CircuitBreaker implementation in use.")
		return NewNoopCircuitBreaker()
	}

	return breaker
}
