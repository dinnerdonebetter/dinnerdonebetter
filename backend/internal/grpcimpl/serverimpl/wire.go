package serverimpl

import (
	"github.com/google/wire"
)

var (
	Providers = wire.NewSet(
		NewServer,
	)
)
