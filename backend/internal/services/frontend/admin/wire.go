package admin

import (
	"github.com/google/wire"
)

var (
	ProvidersAdminWebapp = wire.NewSet(
		NewServer,
	)
)
