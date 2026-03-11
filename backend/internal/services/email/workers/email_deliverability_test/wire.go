package emaildeliverabilitytest

import (
	"github.com/google/wire"
)

var (
	// ProvidersEmailDeliverabilityTest are what we provide to dependency injection.
	ProvidersEmailDeliverabilityTest = wire.NewSet(
		NewJob,
	)
)
