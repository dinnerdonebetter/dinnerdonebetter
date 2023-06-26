package hashing

import (
	"context"
)

type (
	Hasher interface {
		Hash(ctx context.Context, content string) (string, error)
	}
)
