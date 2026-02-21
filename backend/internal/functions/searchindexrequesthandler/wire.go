package searchindexrequesthandler

import (
	"github.com/google/wire"
)

var (
	Providers = wire.NewSet(
		NewSearchIndexRequestHandler,
	)
)
