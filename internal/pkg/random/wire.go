package random

import "github.com/google/wire"

// Providers are what we offer to dependency injection.
var Providers = wire.NewSet(
	NewGenerator,
)
