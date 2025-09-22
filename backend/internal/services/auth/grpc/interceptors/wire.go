package interceptors

import "github.com/google/wire"

var (
	Providers = wire.NewSet(
		ProvideAuthInterceptor,
	)
)
