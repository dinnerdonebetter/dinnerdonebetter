package eatinggrpc

import (
	"github.com/google/wire"
)

var (
	ProvidersGRPCImpl = wire.NewSet(
		NewService,
	)
)
