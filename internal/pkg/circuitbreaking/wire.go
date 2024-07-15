package circuitbreaking

import "github.com/google/wire"

var (
	Providers = wire.NewSet(
		ProvideCircuitBreaker,
	)
)
