package manager

import "github.com/google/wire"

var (
	PaymentsManagerProviders = wire.NewSet(
		NewPaymentsDataManager,
	)
)
