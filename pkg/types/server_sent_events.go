package types

import "net/http"

type (
	// ServerSentEventsService describes a service capable of serving SSE streams.
	ServerSentEventsService interface {
		StreamSubscriptionHandler(res http.ResponseWriter, req *http.Request)
	}
)
