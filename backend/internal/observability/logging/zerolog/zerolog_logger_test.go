package zerolog

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/observability/logging"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
)

func Test_buildZerologger(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		assert.NotNil(t, buildZerologger(logging.DebugLevel))
	})
}

func TestNewLogger(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		assert.NotNil(t, NewZerologLogger(logging.DebugLevel))
	})
}

func Test_zerologLogger_WithName(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		l := NewZerologLogger(logging.DebugLevel)

		assert.NotNil(t, l.WithName(t.Name()))
	})
}

func Test_zerologLogger_SetRequestIDFunc(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		l := NewZerologLogger(logging.DebugLevel)

		l.SetRequestIDFunc(func(*http.Request) string {
			return ""
		})
	})
}

func Test_zerologLogger_Info(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		l := NewZerologLogger(logging.DebugLevel)

		l.Info(t.Name())
	})
}

func Test_zerologLogger_Debug(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		l := NewZerologLogger(logging.DebugLevel)

		l.Debug(t.Name())
	})
}

func Test_zerologLogger_Error(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		l := NewZerologLogger(logging.DebugLevel)

		l.Error(errors.New("blah"), t.Name())
	})
}

func Test_zerologLogger_Clone(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		l := NewZerologLogger(logging.DebugLevel)

		assert.NotNil(t, l.Clone())
	})
}

func Test_zerologLogger_WithValue(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		l := NewZerologLogger(logging.DebugLevel)

		assert.NotNil(t, l.WithValue("name", t.Name()))
	})
}

func Test_zerologLogger_WithValues(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		l := NewZerologLogger(logging.DebugLevel)

		assert.NotNil(t, l.WithValues(map[string]any{"name": t.Name()}))
	})
}

func Test_zerologLogger_WithError(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		l := NewZerologLogger(logging.DebugLevel)

		assert.NotNil(t, l.WithError(errors.New("blah")))
	})
}

func Test_zerologLogger_WithSpan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		l := NewZerologLogger(logging.DebugLevel)

		span := trace.SpanFromContext(ctx)

		assert.NotNil(t, l.WithSpan(span))
	})
}

func Test_zerologLogger_WithRequest(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		l, ok := NewZerologLogger(logging.DebugLevel).(*zerologLogger)
		require.True(t, ok)

		l.requestIDFunc = func(*http.Request) string {
			return t.Name()
		}

		u, err := url.ParseRequestURI("https://whatever.whocares.gov?things=stuff")
		require.NoError(t, err)

		assert.NotNil(t, l.WithRequest(&http.Request{
			URL: u,
		}))
	})
}

func Test_zerologLogger_WithResponse(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		l := NewZerologLogger(logging.DebugLevel)

		assert.NotNil(t, l.WithResponse(&http.Response{}))
	})
}
