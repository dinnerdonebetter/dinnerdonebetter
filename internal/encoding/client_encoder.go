package encoding

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"io"

	"github.com/prixfixeco/backend/internal/observability"
	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/tracing"
)

type (
	// ClientEncoder is an encoder for a service client.
	ClientEncoder interface {
		ContentType() string
		Unmarshal(ctx context.Context, data []byte, v interface{}) error
		Encode(ctx context.Context, dest io.Writer, v interface{}) error
		EncodeReader(ctx context.Context, data interface{}) (io.Reader, error)
	}

	// clientEncoder is our concrete implementation of ClientEncoder.
	clientEncoder struct {
		logger      logging.Logger
		tracer      tracing.Tracer
		contentType *contentType
	}
)

func (e *clientEncoder) Unmarshal(ctx context.Context, data []byte, v interface{}) error {
	_, span := e.tracer.StartSpan(ctx)
	defer span.End()

	logger := e.logger.WithValue("data_length", len(data))
	var unmarshalFunc func(data []byte, v interface{}) error

	switch e.contentType {
	case ContentTypeXML:
		unmarshalFunc = xml.Unmarshal
	default:
		unmarshalFunc = json.Unmarshal
	}

	if err := unmarshalFunc(data, v); err != nil {
		return observability.PrepareError(err, span, "unmarshalling JSON content")
	}

	logger.Debug("unmarshalled")

	return nil
}

func (e *clientEncoder) Encode(ctx context.Context, dest io.Writer, data interface{}) error {
	_, span := e.tracer.StartSpan(ctx)
	defer span.End()

	var err error

	switch e.contentType {
	case ContentTypeXML:
		err = xml.NewEncoder(dest).Encode(data)
	default:
		err = json.NewEncoder(dest).Encode(data)
	}

	if err != nil {
		return observability.PrepareError(err, span, "encoding JSON content")
	}

	return nil
}

func (e *clientEncoder) EncodeReader(ctx context.Context, data interface{}) (io.Reader, error) {
	_, span := e.tracer.StartSpan(ctx)
	defer span.End()

	var marshalFunc func(v interface{}) ([]byte, error)

	switch e.contentType {
	case ContentTypeXML:
		marshalFunc = xml.Marshal
	default:
		marshalFunc = json.Marshal
	}

	out, err := marshalFunc(data)
	if err != nil {
		return nil, observability.PrepareError(err, span, "marshaling to XML")
	}

	return bytes.NewReader(out), nil
}

// ProvideClientEncoder provides a ClientEncoder.
func ProvideClientEncoder(logger logging.Logger, tracerProvider tracing.TracerProvider, encoding *contentType) ClientEncoder {
	return &clientEncoder{
		logger:      logging.EnsureLogger(logger).WithName("client_encoder"),
		tracer:      tracing.NewTracer(tracerProvider.Tracer("client_encoder")),
		contentType: encoding,
	}
}
