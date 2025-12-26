package grpc

import (
	"github.com/google/wire"
)

var (
	MPSvcProviders = wire.NewSet(
		NewService,
		ProvideMethodPermissions,
	)
)
