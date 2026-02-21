package outboundemailhandler

import (
	"github.com/google/wire"
)

var (
	Providers = wire.NewSet(
		NewOutboundEmailHandler,
	)
)
