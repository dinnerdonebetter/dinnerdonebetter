package tracing

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/trace"

	"github.com/stretchr/testify/assert"

	"github.com/prixfixeco/api_server/internal/observability/logging"
)

func TestNewInstrumentedSQLTracer(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		assert.NotNil(t, NewInstrumentedSQLTracer(trace.NewNoopTracerProvider(), t.Name()))
	})
}

func Test_instrumentedSQLTracerWrapper_GetSpan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		w := NewInstrumentedSQLTracer(trace.NewNoopTracerProvider(), t.Name())

		assert.NotNil(t, w.GetSpan(ctx))
	})
}

func TestNewInstrumentedSQLLogger(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		assert.NotNil(t, NewInstrumentedSQLLogger(logging.NewNoopLogger()))
	})
}

func Test_instrumentedSQLLoggerWrapper_Log(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		w := NewInstrumentedSQLLogger(logging.NewNoopLogger())

		w.Log(ctx, t.Name())
	})
}
