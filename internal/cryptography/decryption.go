package cryptography

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"

	"github.com/prixfixeco/backend/internal/observability"
)

type (
	Decryptor interface {
		Decrypt(ctx context.Context, text, secret string) (string, error)
	}
)

func (s *aesImpl) Decrypt(ctx context.Context, text, secret string) (string, error) {
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
