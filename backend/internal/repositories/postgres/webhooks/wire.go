package webhooks

import "github.com/google/wire"

var (
	WebhookRepoProviders = wire.NewSet(
		ProvideWebhooksRepository,
	)
)
