package types

import (
	"net/http"
)

type (
	// WebsocketDataService describes a structure capable of serving traffic related to notifications.
	WebsocketDataService interface {
		SubscribeHandler(http.ResponseWriter, *http.Request)
	}
)
