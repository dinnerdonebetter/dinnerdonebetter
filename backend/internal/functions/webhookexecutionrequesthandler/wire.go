package webhookexecutionrequesthandler

import (
	"github.com/google/wire"
)

var (
	Providers = wire.NewSet(
		NewWebhookExecutionRequestHandler,
	)
)
