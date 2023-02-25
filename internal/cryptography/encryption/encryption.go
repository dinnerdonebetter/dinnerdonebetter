package encryption

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	"github.com/prixfixeco/backend/internal/observability"
	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/tracing"
)

var errSecretNotTheRightLength = errors.New("secret is not the right length")

type Encryptor interface {
	Encrypt(ctx context.Context, text, secret string) (string, error)
	Decrypt(ctx context.Context, text, secret string) (string, error)
}

// StandardEncryptor is the standard Encryptor implementation.
type StandardEncryptor struct {
	tracer tracing.Tracer
	logger logging.Logger
}

func NewStandardEncryptor(tracerProvider tracing.TracerProvider, logger logging.Logger) Encryptor {
	return &StandardEncryptor{
		logger: logging.EnsureLogger(logger).WithName("encryptor"),
		tracer: tracing.NewTracer(tracerProvider.Tracer("encryptor")),
	}
}

func (e *StandardEncryptor) Encrypt(ctx context.Context, text, secret string) (string, error) {
	_, span := e.tracer.StartSpan(ctx)
	defer span.End()

	logger := e.logger.WithValue("text", text)

	if len(secret) != 32 {
		return "", observability.PrepareAndLogError(errSecretNotTheRightLength, logger, span, "secret is too small")
	}

	aesBlock, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return "", observability.PrepareAndLogError(err, logger, span, "creating aes cipher")
	}

	gcmInstance, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return "", observability.PrepareAndLogError(err, logger, span, "creating gcm instance")
	}

	nonce := make([]byte, gcmInstance.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", observability.PrepareAndLogError(err, logger, span, "generating nonce")
	}

	cipheredText := gcmInstance.Seal(nonce, nonce, []byte(text), nil)

	return base64.URLEncoding.EncodeToString(cipheredText), nil
}

func (s *StandardEncryptor) Decrypt(ctx context.Context, text, secret string) (string, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithValue("text", text)

	if len(secret) != 32 {
		return "", observability.PrepareAndLogError(errSecretNotTheRightLength, logger, span, "secret is too small")
	}

	ciphered, err := base64.URLEncoding.DecodeString(text)
	if err != nil {
		return "", observability.PrepareAndLogError(err, logger, span, "decoding ciphered text")
	}

	aesBlock, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return "", observability.PrepareAndLogError(err, logger, span, "creating aes cipher")
	}

	gcmInstance, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return "", observability.PrepareAndLogError(err, logger, span, "creating gcm instance")
	}

	nonceSize := gcmInstance.NonceSize()
	nonce, cipheredText := ciphered[:nonceSize], ciphered[nonceSize:]

	originalText, err := gcmInstance.Open(nil, nonce, cipheredText, nil)
	if err != nil {
		return "", observability.PrepareAndLogError(err, logger, span, "decrypting ciphered text")
	}

	return string(originalText), nil
}
