package types

import "net/http"

type (
	// CapitalismDataService describes a structure capable of serving worker-oriented requests.
	CapitalismDataService interface {
		IncomingWebhookHandler(res http.ResponseWriter, req *http.Request)
	}
)
