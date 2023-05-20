package prometheus

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/observability/logging"

	"github.com/stretchr/testify/assert"
)

func TestConfig_ProvideMetricsHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{
			RuntimeMetricsCollectionInterval: minimumRuntimeCollectionInterval,
		}

		actual, err := cfg.ProvideMetricsHandler()
		assert.NoError(t, err)
		assert.NotNil(t, actual)
	})
}

func TestConfig_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		cfg := &Config{
			RuntimeMetricsCollectionInterval: minimumRuntimeCollectionInterval,
		}

		assert.NoError(t, cfg.ValidateWithContext(ctx))
	})
}

func Test_initiatePrometheusExporter(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{
			RuntimeMetricsCollectionInterval: minimumRuntimeCollectionInterval,
		}

		_, _, err := cfg.initiateExporter()
		assert.NoError(t, err)
	})
}

func TestConfig_ProvideUnitCounterProvider(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		logger := logging.NewNoopLogger()
		cfg := &Config{}

		actual, err := cfg.ProvideUnitCounterProvider(logger)
		assert.NotNil(t, actual)
		assert.NoError(t, err)

		actual("things", "stuff")
	})
}
