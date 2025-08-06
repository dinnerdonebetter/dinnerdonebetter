package manager

import (
	"github.com/google/wire"
)

var (
	Providers = wire.NewSet(
		NewOAuth2Manager,
	)
)
