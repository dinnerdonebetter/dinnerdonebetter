package sessions

import (
	platformerrors "github.com/verygoodsoftwarenotvirus/platform/v3/errors"
)

var (
	ErrAuthenticationNotFound = platformerrors.New("authentication not found")
)
