package grpc

import "github.com/google/wire"

var (
	CommentsSvcProviders = wire.NewSet(
		NewService,
		ProvideMethodPermissions,
	)
)
