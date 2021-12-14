package encoding

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"go.opentelemetry.io/otel/trace"

	"github.com/stretchr/testify/assert"

	"github.com/prixfixeco/api_server/internal/observability/logging"
)

func TestProvideClientEncoder(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		assert.NotNil(t, ProvideClientEncoder(logging.NewNoopLogger(), trace.NewNoopTracerProvider(), ContentTypeJSON))
	})
}

func Test_clientEncoder_Unmarshal(T *testing.T) {
	T.Parallel()

	T.Run("with JSON", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		e := ProvideClientEncoder(logging.NewNoopLogger(), trace.NewNoopTracerProvider(), ContentTypeJSON)

		expected := &example{Name: "name"}
		actual := &example{}

		assert.NoError(t, e.Unmarshal(ctx, []byte(`{"name": "name"}`), &actual))
		assert.Equal(t, expected, actual)
	})

	T.Run("with XML", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		e := ProvideClientEncoder(logging.NewNoopLogger(), trace.NewNoopTracerProvider(), ContentTypeXML)

		expected := &example{Name: "name"}
		actual := &example{}

		assert.NoError(t, e.Unmarshal(ctx, []byte(`<example><name>name</name></example>`), &actual))
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid data", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		e := ProvideClientEncoder(logging.NewNoopLogger(), trace.NewNoopTracerProvider(), ContentTypeJSON)

		actual := &example{}

		assert.Error(t, e.Unmarshal(ctx, []byte(`{"name"   `), &actual))
		assert.Empty(t, actual.Name)
	})
}

func Test_clientEncoder_Encode(T *testing.T) {
	T.Parallel()

	T.Run("with JSON", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		e := ProvideClientEncoder(logging.NewNoopLogger(), trace.NewNoopTracerProvider(), ContentTypeJSON)

		res := httptest.NewRecorder()

		assert.NoError(t, e.Encode(ctx, res, &example{Name: t.Name()}))
	})

	T.Run("with XML", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		e := ProvideClientEncoder(logging.NewNoopLogger(), trace.NewNoopTracerProvider(), ContentTypeXML)

		res := httptest.NewRecorder()

		assert.NoError(t, e.Encode(ctx, res, &example{Name: t.Name()}))
	})

	T.Run("with invalid data", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		e := ProvideClientEncoder(logging.NewNoopLogger(), trace.NewNoopTracerProvider(), ContentTypeJSON)

		assert.Error(t, e.Encode(ctx, nil, &broken{Name: json.Number(t.Name())}))
	})
}

func Test_clientEncoder_EncodeReader(T *testing.T) {
	T.Parallel()

	T.Run("with JSON", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		e := ProvideClientEncoder(logging.NewNoopLogger(), trace.NewNoopTracerProvider(), ContentTypeJSON)

		actual, err := e.EncodeReader(ctx, &example{Name: t.Name()})
		assert.NoError(t, err)
		assert.NotNil(t, actual)
	})

	T.Run("with XML", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		e := ProvideClientEncoder(logging.NewNoopLogger(), trace.NewNoopTracerProvider(), ContentTypeXML)

		actual, err := e.EncodeReader(ctx, &example{Name: t.Name()})
		assert.NoError(t, err)
		assert.NotNil(t, actual)
	})

	T.Run("with invalid data", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		e := ProvideClientEncoder(logging.NewNoopLogger(), trace.NewNoopTracerProvider(), ContentTypeJSON)

		actual, err := e.EncodeReader(ctx, &broken{Name: json.Number(t.Name())})
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}
