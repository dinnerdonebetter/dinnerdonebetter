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
	CircuitBreakerErrorRate  float64 `json:"circuitBreakerErrorPercentage"            toml:"circuitBreaker_error_percentage"`
	CircuitBreakerMinSamples int64   `json:"circuitBreakerMinimumOccurrenceThreshold" toml:"circuitBreaker_minimum_occurrence_threshold"`
}

func (cfg *Config) EnsureDefaults() {
	if cfg == nil {
		cfg = &Config{}
	}

	if cfg.CircuitBreakerErrorRate == 0 {
		cfg.CircuitBreakerErrorRate = 200
	}

	if cfg.CircuitBreakerMinSamples == 0 {
		cfg.CircuitBreakerMinSamples = 1_000_000
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
				result := samples >= cfg.CircuitBreakerMinSamples && cb.ErrorRate() >= cfg.CircuitBreakerErrorRate
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
