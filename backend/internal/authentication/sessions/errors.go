package sessions

import (
	platformerrors "github.com/verygoodsoftwarenotvirus/platform/v5/errors"
)

var (
	ErrAuthenticationNotFound = platformerrors.New("authentication not found")
)
