package random

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"io"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
)

const (
	arbitrarySize uint16 = 128
)

var (
	_ Generator = (*standardGenerator)(nil)

	defaultGenerator = NewGenerator(logging.NewNoopLogger(), tracing.NewNoopTracerProvider())
)

func init() {
	if _, err := rand.Read(make([]byte, arbitrarySize)); err != nil {
		panic(err)
	}
}

type (
	// Generator should generate random strings securely.
	Generator interface {
		GenerateHexEncodedString(ctx context.Context, length int) (string, error)
		GenerateBase32EncodedString(context.Context, int) (string, error)
		GenerateBase64EncodedString(context.Context, int) (string, error)
		GenerateRawBytes(context.Context, int) ([]byte, error)
	}

	standardGenerator struct {
		logger     logging.Logger
		tracer     tracing.Tracer
		randReader io.Reader
	}
)

// NewGenerator builds a new Generator.
func NewGenerator(logger logging.Logger, tracerProvider tracing.TracerProvider) Generator {
	return &standardGenerator{
		logger:     logging.EnsureLogger(logger).WithName("random_generator"),
		tracer:     tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer("secret_generator")),
		randReader: rand.Reader,
	}
}

// GenerateHexEncodedString generates a one-off value with an anonymous Generator.
func GenerateHexEncodedString(ctx context.Context, length int) (string, error) {
	return defaultGenerator.GenerateHexEncodedString(ctx, length)
}

// GenerateBase32EncodedString generates a one-off value with an anonymous Generator.
func GenerateBase32EncodedString(ctx context.Context, length int) (string, error) {
	return defaultGenerator.GenerateBase32EncodedString(ctx, length)
}

// GenerateBase64EncodedString generates a one-off value with an anonymous Generator.
func GenerateBase64EncodedString(ctx context.Context, length int) (string, error) {
	return defaultGenerator.GenerateBase64EncodedString(ctx, length)
}

// GenerateRawBytes generates a one-off value with an anonymous Generator.
func GenerateRawBytes(ctx context.Context, length int) ([]byte, error) {
	return defaultGenerator.GenerateRawBytes(ctx, length)
}

// generateSecret generates a securely random byte array of a given length.
func (g *standardGenerator) generateSecret(ctx context.Context, length int) ([]byte, error) {
	_, span := g.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachToSpan(span, "generated_string.requested_length", length)

	b := make([]byte, length)
	if _, err := g.randReader.Read(b); err != nil {
		return nil, observability.PrepareError(err, span, "reading from secure random source")
	}

	return b, nil
}

// GenerateRawBytes generates a securely random byte array.
func (g *standardGenerator) GenerateRawBytes(ctx context.Context, length int) ([]byte, error) {
	_, span := g.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachToSpan(span, "rand_gen.requested_length", length)

	return g.generateSecret(ctx, length)
}

// GenerateHexEncodedString generates a base64-encoded string of a securely random byte array of a given length.
func (g *standardGenerator) GenerateHexEncodedString(ctx context.Context, length int) (string, error) {
	_, span := g.tracer.StartSpan(ctx)
	defer span.End()

	logger := g.logger.WithValue("length", length)
	tracing.AttachToSpan(span, "rand_gen.requested_length", length)

	b, err := g.GenerateRawBytes(ctx, length)
	if err != nil {
		return "", observability.PrepareAndLogError(err, logger, span, "reading from secure random source")
	}

	return hex.EncodeToString(b), nil
}

// GenerateBase32EncodedString generates a base64-encoded string of a securely random byte array of a given length.
func (g *standardGenerator) GenerateBase32EncodedString(ctx context.Context, length int) (string, error) {
	_, span := g.tracer.StartSpan(ctx)
	defer span.End()

	logger := g.logger.WithValue("length", length)
	tracing.AttachToSpan(span, "rand_gen.requested_length", length)

	b, err := g.GenerateRawBytes(ctx, length)
	if err != nil {
		return "", observability.PrepareAndLogError(err, logger, span, "reading from secure random source")
	}

	return base32.StdEncoding.EncodeToString(b), nil
}

// GenerateBase64EncodedString generates a base64-encoded string of a securely random byte array of a given length.
func (g *standardGenerator) GenerateBase64EncodedString(ctx context.Context, length int) (string, error) {
	_, span := g.tracer.StartSpan(ctx)
	defer span.End()

	logger := g.logger.WithValue("length", length)
	tracing.AttachToSpan(span, "rand_gen.requested_length", length)

	b, err := g.GenerateRawBytes(ctx, length)
	if err != nil {
		return "", observability.PrepareAndLogError(err, logger, span, "reading from secure random source")
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}
