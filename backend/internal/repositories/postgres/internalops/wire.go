package internalops

import "github.com/google/wire"

var (
	Providers = wire.NewSet(
		ProvideInternalOpsRepository,
	)
)
