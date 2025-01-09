package routingcfg

import "github.com/google/wire"

var (
	// RoutingConfigProviders are what we provide to the dependency injector.
	RoutingConfigProviders = wire.NewSet(
		ProvideRouter,
		ProvideRouteParamManager,
	)
)
