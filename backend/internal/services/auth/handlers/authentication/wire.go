package authentication

import (
	"github.com/google/wire"
)

// Providers are our collection of what we provide to other services.
var Providers = wire.NewSet(
	ProvideService,
	ProvideOAuth2ClientManager,
	ProvideOAuth2ServerImplementation,
)
