package authentication

import (
	"context"

	platformerrors "github.com/dinnerdonebetter/backend/internal/platform/errors"
)

var (
	// ErrPasswordDoesNotMatch indicates that a provided passwords does not match.
	ErrPasswordDoesNotMatch = platformerrors.New("password does not match")
)

type (
	// Hasher hashes passwords.
	Hasher interface {
		HashPassword(ctx context.Context, password string) (string, error)
	}
)
