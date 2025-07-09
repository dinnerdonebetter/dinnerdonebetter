package coregrpc

import "github.com/google/wire"

var (
	ProvidersCoreGRPC = wire.NewSet(
		NewCoreService,
	)
)
