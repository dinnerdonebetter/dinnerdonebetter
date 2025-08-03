package identity

import "github.com/google/wire"

var (
	Providers = wire.NewSet(
		ProvideIdentityRepository,
	)
)
