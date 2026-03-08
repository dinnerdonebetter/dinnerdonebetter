package sessions

import (
	platformerrors "github.com/dinnerdonebetter/backend/internal/platform/errors"
)

var (
	ErrAuthenticationNotFound = platformerrors.New("authentication not found")
)
