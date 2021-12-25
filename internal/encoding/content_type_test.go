package encoding

import (
	"testing"

	"go.opentelemetry.io/otel/trace"

	"github.com/stretchr/testify/assert"

	"github.com/prixfixeco/api_server/internal/observability/logging"
)

func Test_clientEncoder_ContentType(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		e := ProvideClientEncoder(logging.NewNoopLogger(), trace.NewNoopTracerProvider(), ContentTypeJSON)

		assert.NotEmpty(t, e.ContentType())
	})
}

func Test_buildContentType(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		//
	})
}

func Test_contentTypeToString(T *testing.T) {
	T.Parallel()

	T.Run("with JSON", func(t *testing.T) {
		t.Parallel()

		assert.NotEmpty(t, contentTypeToString(ContentTypeJSON))
	})

	T.Run("with XML", func(t *testing.T) {
		t.Parallel()

		assert.NotEmpty(t, contentTypeToString(ContentTypeXML))
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		assert.Empty(t, contentTypeToString(nil))
	})
}

func Test_contentTypeFromString(T *testing.T) {
	T.Parallel()

	T.Run("with JSON", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, ContentTypeJSON, contentTypeFromString(contentTypeJSON))
	})

	T.Run("with XML", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, ContentTypeXML, contentTypeFromString(contentTypeXML))
	})
}
