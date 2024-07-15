package posthog

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/circuitbreaking"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/require"
)

func TestNewSegmentEventReporter(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		logger := logging.NewNoopLogger()
		cfg := &Config{APIKey: t.Name()}

		collector, err := NewPostHogEventReporter(logger, tracing.NewNoopTracerProvider(), cfg.APIKey, circuitbreaking.NewNoopCircuitBreaker())
		require.NoError(t, err)
		require.NotNil(t, collector)
	})

	T.Run("with empty API key", func(t *testing.T) {
		t.Parallel()

		logger := logging.NewNoopLogger()
		cfg := &Config{}

		collector, err := NewPostHogEventReporter(logger, tracing.NewNoopTracerProvider(), cfg.APIKey, circuitbreaking.NewNoopCircuitBreaker())
		require.Error(t, err)
		require.Nil(t, collector)
	})
}

func TestSegmentEventReporter_Close(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		logger := logging.NewNoopLogger()
		cfg := &Config{APIKey: t.Name()}

		collector, err := NewPostHogEventReporter(logger, tracing.NewNoopTracerProvider(), cfg.APIKey, circuitbreaking.NewNoopCircuitBreaker())
		require.NoError(t, err)
		require.NotNil(t, collector)

		collector.Close()
	})
}

func TestSegmentEventReporter_AddUser(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		cfg := &Config{APIKey: t.Name()}
		exampleUserID := fakes.BuildFakeID()
		properties := map[string]any{
			"test.name": t.Name(),
		}

		collector, err := NewPostHogEventReporter(logger, tracing.NewNoopTracerProvider(), cfg.APIKey, circuitbreaking.NewNoopCircuitBreaker())
		require.NoError(t, err)
		require.NotNil(t, collector)

		require.NoError(t, collector.AddUser(ctx, exampleUserID, properties))
	})
}

func TestSegmentEventReporter_EventOccurred(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		cfg := &Config{APIKey: t.Name()}
		exampleUserID := fakes.BuildFakeID()
		properties := map[string]any{
			"test.name": t.Name(),
		}

		collector, err := NewPostHogEventReporter(logger, tracing.NewNoopTracerProvider(), cfg.APIKey, circuitbreaking.NewNoopCircuitBreaker())
		require.NoError(t, err)
		require.NotNil(t, collector)

		require.NoError(t, collector.EventOccurred(ctx, types.ServiceEventType(t.Name()), exampleUserID, properties))
	})
}
