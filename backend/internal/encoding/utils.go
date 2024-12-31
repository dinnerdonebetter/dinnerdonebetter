package encoding

import (
	"bytes"
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
)

// MustEncode encodes a given piece of data to a given encoding.
func MustEncode(data any, ct *contentType) []byte {
	if ct == nil {
		ct = ContentTypeJSON
	}

	var b bytes.Buffer
	if err := ProvideClientEncoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), ct).Encode(context.Background(), &b, data); err != nil {
		panic(err)
	}

	return b.Bytes()
}

// MustEncodeJSON JSON encodes a piece of data.
func MustEncodeJSON(data any) []byte {
	return MustEncode(data, ContentTypeJSON)
}
