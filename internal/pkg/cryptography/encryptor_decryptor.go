package cryptography

import (
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
)

type (
	EncryptorDecryptor interface {
		Encryptor
		Decryptor
	}
)

// aesImpl is the standard EncryptorDecryptor implementation.
type aesImpl struct {
	tracer tracing.Tracer
	logger logging.Logger
}

func NewAESEncryptorDecryptor(tracerProvider tracing.TracerProvider, logger logging.Logger) EncryptorDecryptor {
	return &aesImpl{
		logger: logging.EnsureLogger(logger).WithName("encryptor"),
		tracer: tracing.NewTracer(tracerProvider.Tracer("encryptor")),
	}
}
