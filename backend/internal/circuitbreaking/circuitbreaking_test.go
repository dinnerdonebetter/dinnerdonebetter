package circuitbreaking

import (
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/metrics"

	"github.com/stretchr/testify/assert"
)

//nolint:paralleltest // race condition in the core circuit breaker library, I think?
func TestProvideCircuitBreaker(t *testing.T) {
	cfg := &Config{}
	cfg.EnsureDefaults()

	cb, err := ProvideCircuitBreaker(cfg, logging.NewNoopLogger(), metrics.NewNoopMetricsProvider())
	assert.NotNil(t, cb)
	assert.NoError(t, err)
}

//nolint:paralleltest // race condition in the core circuit breaker library, I think?
func TestEnsureCircuitBreaker(t *testing.T) {
	assert.NotNil(t, EnsureCircuitBreaker(nil))
}

//nolint:paralleltest // race condition in the core circuit breaker library, I think?
func TestCircuitBreaker_Integration(t *testing.T) {
	// there is a data race bug in the circuit breaker library that prevents this from not tripping the data race detector.
	t.SkipNow()

	cfg := &Config{
		Name:                   t.Name(),
		ErrorRate:              1,
		MinimumSampleThreshold: 1,
	}

	cb, err := ProvideCircuitBreaker(cfg, logging.NewNoopLogger(), metrics.NewNoopMetricsProvider())
	assert.NotNil(t, cb)
	assert.NoError(t, err)

	assert.True(t, cb.CanProceed())
	cb.Failed()
	assert.True(t, cb.CannotProceed())
	cb.Succeeded()
	assert.Eventually(
		t,
		func() bool {
			return cb.CanProceed()
		},
		5*time.Second,
		500*time.Millisecond,
	)
}
