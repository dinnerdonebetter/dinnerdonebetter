package grpc

import "github.com/google/wire"

var (
	PaymentsSvcProviders = wire.NewSet(
		NewService,
		ProvideMethodPermissions,
	)
)
