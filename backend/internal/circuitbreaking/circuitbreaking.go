package circuitbreaking

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"

	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/metrics"

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

type BaseImplementation struct {
	circuitBreaker *circuit.Breaker
}

func (b *BaseImplementation) Failed() {
	b.circuitBreaker.Fail()
}

func (b *BaseImplementation) Succeeded() {
	b.circuitBreaker.Success()
}

func (b *BaseImplementation) CanProceed() bool {
	return b.circuitBreaker.Ready()
}

func (b *BaseImplementation) CannotProceed() bool {
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

func handleCircuitBreakerEvents(cb *circuit.Breaker, logger logging.Logger, brokenCounter metrics.Int64Counter) {
	for be := range <-cb.Subscribe() {
		switch be {
		case circuit.BreakerTripped:
			brokenCounter.Add(context.Background(), 1)
		case circuit.BreakerReset:
			logger.Debug("circuit breaker reset")
		case circuit.BreakerFail:
			logger.Debug("circuit breaker experienced an instance of failure")
		case circuit.BreakerReady:
			logger.Debug("circuit breaker is ready")
		}
	}
}

// ProvideCircuitBreaker provides a CircuitBreaker.
func (cfg *Config) ProvideCircuitBreaker(logger logging.Logger, metricsProvider metrics.Provider) (CircuitBreaker, error) {
	if cfg == nil {
		cfg = &Config{
			Name:                   "UNKNOWN",
			ErrorRate:              1.0,
			MinimumSampleThreshold: math.MaxUint64,
		}
	}
	cfg.EnsureDefaults()

	brokenCounter, err := metricsProvider.NewInt64Counter(fmt.Sprintf("%s_circuit_breaker_tripped", cfg.Name))
	if err != nil {
		return nil, err
	}

	logger = logging.EnsureLogger(logger).WithValue("circuit_breaker", cfg.Name)

	cb := circuit.NewBreakerWithOptions(&circuit.Options{
		ShouldTrip: func(cb *circuit.Breaker) bool {
			return uint64(cb.Failures()+cb.Successes()) >= cfg.MinimumSampleThreshold && cb.ErrorRate() >= cfg.ErrorRate
		},
		WindowTime:    circuit.DefaultWindowTime,
		WindowBuckets: circuit.DefaultWindowBuckets,
	})

	go handleCircuitBreakerEvents(cb, logger, brokenCounter)

	return &BaseImplementation{
		circuitBreaker: cb,
	}, nil
}
