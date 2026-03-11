package queuetest

import (
	"github.com/google/wire"
)

var (
	// ProvidersQueueTest are what we provide to dependency injection.
	ProvidersQueueTest = wire.NewSet(
		NewJob,
	)
)
