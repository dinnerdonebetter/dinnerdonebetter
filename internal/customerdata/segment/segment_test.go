package segment

import (
	"context"
	"testing"

	"github.com/prixfixeco/backend/pkg/types"

	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/pkg/types/fakes"
)

func TestNewSegmentCustomerDataCollector(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		logger := logging.NewNoopLogger()

		collector, err := NewSegmentCustomerDataCollector(logger, tracing.NewNoopTracerProvider(), t.Name())
		require.NoError(t, err)
		require.NotNil(t, collector)
	})

	T.Run("with empty API key", func(t *testing.T) {
		t.Parallel()

		logger := logging.NewNoopLogger()

		collector, err := NewSegmentCustomerDataCollector(logger, tracing.NewNoopTracerProvider(), "")
		require.Error(t, err)
		require.Nil(t, collector)
	})
}

func TestCustomerDataCollector_Close(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		logger := logging.NewNoopLogger()

		collector, err := NewSegmentCustomerDataCollector(logger, tracing.NewNoopTracerProvider(), t.Name())
		require.NoError(t, err)
		require.NotNil(t, collector)

		require.NoError(t, collector.Close())
	})
}

func TestCustomerDataCollector_AddUser(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		exampleUserID := fakes.BuildFakeID()
		properties := map[string]any{
			"test.name": t.Name(),
		}

		collector, err := NewSegmentCustomerDataCollector(logger, tracing.NewNoopTracerProvider(), t.Name())
		require.NoError(t, err)
		require.NotNil(t, collector)

		require.NoError(t, collector.AddUser(ctx, exampleUserID, properties))
	})
}

func TestCustomerDataCollector_EventOccurred(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		exampleUserID := fakes.BuildFakeID()
		properties := map[string]any{
			"test.name": t.Name(),
		}

		collector, err := NewSegmentCustomerDataCollector(logger, tracing.NewNoopTracerProvider(), t.Name())
		require.NoError(t, err)
		require.NotNil(t, collector)

		require.NoError(t, collector.EventOccurred(ctx, types.CustomerEventType(t.Name()), exampleUserID, properties))
	})
}
