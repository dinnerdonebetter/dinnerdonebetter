package tracing

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
)

func TestNewInstrumentedSQLTracer(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		assert.NotNil(t, NewInstrumentedSQLTracer(t.Name()))
	})
}

func Test_instrumentedSQLTracerWrapper_GetSpan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		w := NewInstrumentedSQLTracer(t.Name())

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
