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
