package payments

import "github.com/google/wire"

var (
	PaymentsRepoProviders = wire.NewSet(
		ProvidePaymentsRepository,
	)
)
