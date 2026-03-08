package authentication

import (
	"context"

	platformerrors "github.com/dinnerdonebetter/backend/internal/platform/errors"
)

var (
	// ErrInvalidTOTPToken indicates that a provided two-factor code is invalid.
	ErrInvalidTOTPToken = platformerrors.New("invalid two factor code")
)

type (
	// Authenticator authenticates users.
	Authenticator interface {
		Hasher

		CredentialsAreValid(ctx context.Context, hash, password, totpSecret, totpCode string) (bool, error)
	}
)
