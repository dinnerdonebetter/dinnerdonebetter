package cache

import (
	"context"
	"errors"
)

var (
	ErrNotFound = errors.New("not found")
)

type (
	// Cache is our wrapper interface for a cache.
	Cache[T any] interface {
		Get(ctx context.Context, key string) (*T, error)
		Set(ctx context.Context, key string, value *T) error
		Delete(ctx context.Context, key string) error
	}
)
