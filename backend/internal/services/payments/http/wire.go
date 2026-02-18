package http

import "github.com/google/wire"

var (
	PaymentsHTTPProviders = wire.NewSet(
		NewWebhookHandler,
		wire.Value(WebhookSignatureHeader(StripeSignatureHeader)),
	)
)
