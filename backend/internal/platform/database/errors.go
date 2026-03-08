package database

import (
	"github.com/dinnerdonebetter/backend/internal/platform/errors"
)

// ErrUserAlreadyExists indicates that a user with that username has already been created.
var ErrUserAlreadyExists = errors.New("user already exists")

// Re-exports from platform/errors for backward compatibility.
// Prefer importing platform/errors directly in new code.
var (
	ErrNilInputProvided   = errors.ErrNilInputProvided
	ErrInvalidIDProvided  = errors.ErrInvalidIDProvided
	ErrEmptyInputProvided = errors.ErrEmptyInputProvided
)
