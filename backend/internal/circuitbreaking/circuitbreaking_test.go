package circuitbreaking

import (
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/metrics"

	"github.com/stretchr/testify/assert"
)

func TestProvideCircuitBreaker(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		cb, err := ProvideCircuitBreaker(nil, logging.NewNoopLogger(), metrics.NewNoopMetricsProvider())
		assert.NotNil(t, cb)
		assert.NoError(t, err)
	})
}

func TestEnsureCircuitBreaker(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, EnsureCircuitBreaker(nil))
}

func TestCircuitBreaker_Integration(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

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
	})
}
