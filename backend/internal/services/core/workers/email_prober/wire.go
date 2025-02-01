package emailprober

import (
	"github.com/google/wire"
)

var (
	// ProvidersEmailProber are what we provide to dependency injection.
	ProvidersEmailProber = wire.NewSet(
		NewEmailProber,
	)
)
