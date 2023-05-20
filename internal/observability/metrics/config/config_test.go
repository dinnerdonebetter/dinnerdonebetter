package config

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/metrics/prometheus"

	"github.com/stretchr/testify/assert"
)

func TestConfig_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		cfg := &Config{}

		assert.NoError(t, cfg.ValidateWithContext(ctx))
	})
}

func TestConfig_ProvideMetricsHandler(T *testing.T) {
	T.Parallel()

	T.Run("with prometheus", func(t *testing.T) {
		t.Parallel()

		logger := logging.NewNoopLogger()
		cfg := &Config{
			Provider:   ProviderPrometheus,
			Prometheus: &prometheus.Config{},
		}

		actual, err := cfg.ProvideMetricsHandler(logger)
		assert.NotNil(t, actual)
		assert.NoError(t, err)
	})

	T.Run("with empty provider", func(t *testing.T) {
		t.Parallel()

		logger := logging.NewNoopLogger()
		cfg := &Config{
			Provider: "",
		}

		actual, err := cfg.ProvideMetricsHandler(logger)
		assert.Nil(t, actual)
		assert.NoError(t, err)
	})
}

func TestConfig_ProvideUnitCounterProvider(T *testing.T) {
	T.Parallel()

	T.Run("with prometheus", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		cfg := &Config{
			Provider:   ProviderPrometheus,
			Prometheus: &prometheus.Config{},
		}

		actual, err := cfg.ProvideUnitCounterProvider(ctx, logger)
		assert.NotNil(t, actual)
		assert.NoError(t, err)
	})

	T.Run("no-op", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		cfg := &Config{}

		actual, err := cfg.ProvideUnitCounterProvider(ctx, logger)
		assert.NotNil(t, actual)
		assert.NoError(t, err)
	})

	T.Run("with invalid provider", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		cfg := &Config{
			Provider: t.Name(),
		}

		actual, err := cfg.ProvideUnitCounterProvider(ctx, logger)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
