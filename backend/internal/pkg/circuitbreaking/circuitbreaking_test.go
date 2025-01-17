package circuitbreaking

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/metrics"

	"github.com/stretchr/testify/assert"
)

func TestProvideCircuitBreaker(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		cb, err := ProvideCircuitBreaker(&Config{Name: t.Name()}, logging.NewNoopLogger(), metrics.NewNoopMetricsProvider())

		assert.NotNil(t, cb)
		assert.NoError(t, err)
	})
}

func TestEnsureCircuitBreaker(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, EnsureCircuitBreaker(nil))
}
