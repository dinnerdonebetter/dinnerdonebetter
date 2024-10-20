package types

import (
	"net/http"
)

type (
	// DataDeletionResponse is returned when a user requests their data be deleted.
	DataDeletionResponse struct {
		Successful bool `json:"Successful"`
	}

	// DataPrivacyService describes a structure capable of serving CCPA/GRPC-related requests.
	DataPrivacyService interface {
		DataDeletionHandler(http.ResponseWriter, *http.Request)
	}
)
