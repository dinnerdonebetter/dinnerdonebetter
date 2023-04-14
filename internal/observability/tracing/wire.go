package tracing

import "github.com/google/wire"

var (
	Providers = wire.NewSet(
		BuildTracedHTTPClient,
	)
)
