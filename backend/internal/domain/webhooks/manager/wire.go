package manager

import "github.com/google/wire"

var (
	WebhookManagerProviders = wire.NewSet(
		NewWebhookDataManager,
	)
)
