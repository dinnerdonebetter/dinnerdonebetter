package encoding

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"

	"github.com/BurntSushi/toml"
	"github.com/go-yaml/yaml"
	"github.com/keith-turner/ecoji/v2"
)

type (
	// ClientEncoder is an encoder for a service client.
	ClientEncoder interface {
		ContentType() string
		Unmarshal(ctx context.Context, data []byte, v any) error
		Encode(ctx context.Context, dest io.Writer, v any) error
		EncodeReader(ctx context.Context, data any) (io.Reader, error)
	}

	// clientEncoder is our concrete implementation of ClientEncoder.
	clientEncoder struct {
		logger      logging.Logger
		tracer      tracing.Tracer
		contentType *contentType
	}
)

func (e *clientEncoder) Unmarshal(ctx context.Context, data []byte, v any) error {
	_, span := e.tracer.StartSpan(ctx)
	defer span.End()

	logger := e.logger.WithValue("data_length", len(data))
	var unmarshalFunc func(data []byte, v any) error

	switch e.contentType {
	case ContentTypeXML:
		unmarshalFunc = xml.Unmarshal
	case ContentTypeEmoji:
		unmarshalFunc = unmarshalEmoji
	case ContentTypeTOML:
		unmarshalFunc = toml.Unmarshal
	case ContentTypeYAML:
		unmarshalFunc = yaml.Unmarshal
	default:
		unmarshalFunc = json.Unmarshal
	}

	if err := unmarshalFunc(data, v); err != nil {
		return observability.PrepareError(err, span, "unmarshaling JSON content")
	}

	logger.Debug("unmarshalled")

	return nil
}

func (e *clientEncoder) Encode(ctx context.Context, dest io.Writer, data any) error {
	_, span := e.tracer.StartSpan(ctx)
	defer span.End()

	var err error

	switch e.contentType {
	case ContentTypeXML:
		err = xml.NewEncoder(dest).Encode(data)
	case ContentTypeTOML:
		err = toml.NewEncoder(dest).Encode(data)
	case ContentTypeYAML:
		err = yaml.NewEncoder(dest).Encode(data)
	case ContentTypeEmoji:
		emojiEncoded, emojiEncodeErr := marshalEmoji(data)
		if emojiEncodeErr != nil {
			return emojiEncodeErr
		}

		_, err = dest.Write(emojiEncoded)
	default:
		err = json.NewEncoder(dest).Encode(data)
	}

	if err != nil {
		return observability.PrepareError(err, span, "encoding content")
	}

	return nil
}

func marshalEmoji(v any) ([]byte, error) {
	var gobWriter bytes.Buffer
	if err := gob.NewEncoder(&gobWriter).Encode(v); err != nil {
		return nil, fmt.Errorf("encoding to gob: %w", err)
	}

	r := bytes.NewBuffer(gobWriter.Bytes())
	w := bytes.NewBuffer([]byte{})

	// lord help me, I don't know why it's 76 here
	if err := ecoji.EncodeV2(r, w, 76); err != nil {
		return nil, fmt.Errorf("encoding to emoji: %w", err)
	}

	return w.Bytes(), nil
}

func unmarshalEmoji(data []byte, v any) error {
	w := bytes.NewBuffer([]byte{})

	if err := ecoji.Decode(bytes.NewReader(data), w); err != nil {
		return fmt.Errorf("decoding emoji: %w", err)
	}

	return gob.NewDecoder(w).Decode(v)
}

func tomlMarshalFunc(v any) ([]byte, error) {
	var b bytes.Buffer
	if err := toml.NewEncoder(&b).Encode(v); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (e *clientEncoder) EncodeReader(ctx context.Context, data any) (io.Reader, error) {
	_, span := e.tracer.StartSpan(ctx)
	defer span.End()

	var marshalFunc func(v any) ([]byte, error)

	switch e.contentType {
	case ContentTypeXML:
		marshalFunc = xml.Marshal
	case ContentTypeTOML:
		marshalFunc = tomlMarshalFunc
	case ContentTypeYAML:
		marshalFunc = yaml.Marshal
	case ContentTypeEmoji:
		marshalFunc = marshalEmoji
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
		tracer:      tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer("client_encoder")),
		contentType: encoding,
	}
}
