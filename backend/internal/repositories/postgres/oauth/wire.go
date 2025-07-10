package oauth

import "github.com/google/wire"

var (
	Providers = wire.NewSet(
		ProvideOAuthRepository,
	)
)
