package rpc

import (
	"github.com/google/wire"
)

// Providers are our wire superset of providers this package offers.
var Providers = wire.NewSet(
	ProvideRPCServer,
)
