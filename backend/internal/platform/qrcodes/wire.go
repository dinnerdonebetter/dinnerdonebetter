package qrcodes

import "github.com/google/wire"

var (
	Providers = wire.NewSet(
		NewBuilder,
	)
)
