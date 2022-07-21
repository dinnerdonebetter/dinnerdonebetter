package random

import "github.com/google/wire"

// Providers is what we offer to dependency injection.
var Providers = wire.NewSet(
	NewGenerator,
)
