package sessions

import (
	platformerrors "github.com/verygoodsoftwarenotvirus/platform/errors"
)

var (
	ErrAuthenticationNotFound = platformerrors.New("authentication not found")
)
