package salsa20

import (
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/cryptography"
)

// salsa20Impl is the standard EncryptorDecryptor implementation.
type salsa20Impl struct {
	tracer tracing.Tracer
	logger logging.Logger
	key    [32]byte
}

func NewEncryptorDecryptor(tracerProvider tracing.TracerProvider, logger logging.Logger, key []byte) (cryptography.EncryptorDecryptor, error) {
	if len(key) != 32 {
		return nil, cryptography.ErrIncorrectKeyLength
	}

	return &salsa20Impl{
		logger: logging.EnsureLogger(logger).WithName("encryptor"),
		tracer: tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer("encryptor")),
	}, nil
}
