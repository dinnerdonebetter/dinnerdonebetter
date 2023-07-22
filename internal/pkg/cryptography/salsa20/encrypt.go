package salsa20

import (
	"context"
	"encoding/base64"

	"golang.org/x/crypto/salsa20"
)

var nonce = [8]byte{}

func (e *salsa20Impl) Encrypt(ctx context.Context, content string) (string, error) {
	_, span := e.tracer.StartSpan(ctx)
	defer span.End()

	out := make([]byte, len([]byte(content)))
	salsa20.XORKeyStream(out, []byte(content), nonce[:], &e.key)

	return base64.URLEncoding.EncodeToString(out), nil
}
