package webhooks

import "github.com/google/wire"

var (
	Providers = wire.NewSet(
		ProvideWebhooksRepository,
	)
)
