package circuitbreaking

import (
	"errors"

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
	name           string
}

type Config struct {
	CircuitBreakerErrorRate  float64 `json:"circuitBreakerErrorPercentage"            toml:"circuitBreaker_error_percentage"`
	CircuitBreakerMinSamples int64   `json:"circuitBreakerMinimumOccurrenceThreshold" toml:"circuitBreaker_minimum_occurrence_threshold"`
}

func (cfg Config) EnsureDefaults() {
	if cfg.CircuitBreakerErrorRate == 0 {
		cfg.CircuitBreakerErrorRate = 200
	}

	if cfg.CircuitBreakerMinSamples == 0 {
		cfg.CircuitBreakerMinSamples = 1_000_000
	}
}

func ProvideCircuitBreaker(cfg *Config) (CircuitBreaker, error) {
	if cfg == nil {
		return nil, ErrNilConfig
	}

	return &BaseCircuitBreaker{
		circuitBreaker: circuit.NewBreakerWithOptions(&circuit.Options{
			ShouldTrip:    circuit.RateTripFunc(cfg.CircuitBreakerErrorRate, cfg.CircuitBreakerMinSamples),
			WindowTime:    circuit.DefaultWindowTime,
			WindowBuckets: circuit.DefaultWindowBuckets,
		}),
	}, nil
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
