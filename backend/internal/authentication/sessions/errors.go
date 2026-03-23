package sessions

import (
	platformerrors "github.com/verygoodsoftwarenotvirus/platform/v2/errors"
)

var (
	ErrAuthenticationNotFound = platformerrors.New("authentication not found")
)
