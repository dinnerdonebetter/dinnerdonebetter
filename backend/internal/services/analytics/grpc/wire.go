package grpc

import "github.com/google/wire"

var (
	AnalyticsSvcProviders = wire.NewSet(
		NewService,
		ProvideMethodPermissions,
	)
)
