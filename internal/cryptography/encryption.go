package cryptography

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	"github.com/prixfixeco/backend/internal/observability"
)

var errSecretNotTheRightLength = errors.New("secret is not the right length")

type (
	Encryptor interface {
		Encrypt(ctx context.Context, text, secret string) (string, error)
	}
)

func (e *aesImpl) Encrypt(ctx context.Context, text, secret string) (string, error) {
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
