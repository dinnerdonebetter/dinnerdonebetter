package prometheus

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/unit"
)

func Test_unitCounter_Decrement(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		meterProvider, err := initiatePrometheusExporter()
		require.NoError(t, err)

		mustMeter := metric.Must(meterProvider.Meter(defaultNamespace, metric.WithInstrumentationVersion(instrumentationVersion)))

		ctx := context.Background()
		uc := &unitCounter{
			counter: mustMeter.NewInt64Counter(
				t.Name(),
				metric.WithUnit(unit.Dimensionless),
			),
		}

		uc.Decrement(ctx)
	})
}

func Test_unitCounter_Increment(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		meterProvider, err := initiatePrometheusExporter()
		require.NoError(t, err)

		mustMeter := metric.Must(meterProvider.Meter(defaultNamespace, metric.WithInstrumentationVersion(instrumentationVersion)))

		ctx := context.Background()
		uc := &unitCounter{
			counter: mustMeter.NewInt64Counter(
				t.Name(),
				metric.WithUnit(unit.Dimensionless),
			),
		}

		uc.Increment(ctx)
	})
}

func Test_unitCounter_IncrementBy(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		meterProvider, err := initiatePrometheusExporter()
		require.NoError(T, err)

		mustMeter := metric.Must(meterProvider.Meter(defaultNamespace, metric.WithInstrumentationVersion(instrumentationVersion)))

		ctx := context.Background()
		uc := &unitCounter{
			counter: mustMeter.NewInt64Counter(
				t.Name(),
				metric.WithUnit(unit.Dimensionless),
			),
		}

		uc.IncrementBy(ctx, 123)
	})
}
