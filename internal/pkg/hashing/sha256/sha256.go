package sha256

import (
	"context"
	"crypto/sha256"

	"github.com/dinnerdonebetter/backend/internal/pkg/hashing"
)

var _ hashing.Hasher = (*sha256Hasher)(nil)

type (
	sha256Hasher struct{}
)

func NewSHA256Hasher() hashing.Hasher {
	return &sha256Hasher{}
}

func (s *sha256Hasher) Hash(_ context.Context, content string) (string, error) {
	return string(sha256.New().Sum([]byte(content))), nil
}
