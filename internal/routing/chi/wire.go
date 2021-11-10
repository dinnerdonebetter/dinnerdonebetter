package chi

import "github.com/google/wire"

var (
	// Providers are what we provide to the dependency injector.
	Providers = wire.NewSet(
		NewRouter,
		NewRouteParamManager,
	)
)
